package usdt

import (
	"encoding/json"
	"fmt"
	"time"
	"tron/global"
	"tron/server"
	"tron/util"
)

const (
	Transfer        = "a9059cbb"                           // MethodID
	TransferFrom    = "23b872dd"                           // MethodID
	ContractAddress = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" //usdt合约地址
	Type            = "TriggerSmartContract"
	ContractRet     = "SUCCESS"
)

type TransactionData struct {
	Ret []struct {
		ContractRet string `json:"contractRet"`
	} `json:"ret"`
	Signature []string `json:"signature"`
	TxID      string   `json:"txID"`
	RawData   struct {
		Contract []struct {
			Parameter struct {
				Value struct {
					Data            string `json:"data"`
					OwnerAddress    string `json:"owner_address"`
					ContractAddress string `json:"contract_address"`
				} `json:"value"`
				TypeURL string `json:"type_url"`
			} `json:"parameter"`
			Type string `json:"type"`
		} `json:"contract"`
		RefBlockBytes string `json:"ref_block_bytes"`
		RefBlockHash  string `json:"ref_block_hash"`
		Expiration    int64  `json:"expiration"`
		FeeLimit      int    `json:"fee_limit"`
		Timestamp     int64  `json:"timestamp"`
	} `json:"raw_data"`
	RawDataHex string `json:"raw_data_hex"`
}

// IsUSDTTransfer 验证是否为usdt数据
func (usdtData *TransactionData) IsUSDTTransfer() bool {
	// 成功
	if len(usdtData.Ret) == 0 || usdtData.Ret[0].ContractRet != ContractRet {
		return false
	}
	// 合约类型
	if len(usdtData.RawData.Contract) != 1 || usdtData.RawData.Contract[0].Type != Type {
		return false
	}
	// 合约地址
	param := usdtData.RawData.Contract[0].Parameter.Value
	if param.ContractAddress != util.ToHexAddress(ContractAddress) {
		return false
	}
	// 方法ID
	methID := param.Data[0:8]
	if methID != Transfer && methID != TransferFrom {
		return false
	}
	return true
}

// InsertData 插入数据库
func InsertData(dataChan chan []map[string]interface{}) {
	for data := range dataChan {
		cypher := fmt.Sprintf(`
			unwind $data as row
			merge (f:Address{address:row.from_address})
			merge (t:Address{address:row.to_address})
			with f,t,row
			create (f)-[trade: 正常交易 {hash:row.hash}]->(t)
			set trade.timestamp = row.timestamp,trade.value = row.value,trade.from = row.from_address,trade.to= row.to_address,trade.blockNum = row.block_num
		`)
		_, err := global.TRON_DB.RunCypher(cypher, map[string]interface{}{"data": data})
		if err != nil {
			// 如果重复失败会导致chan阻塞，请检查块高的信息进行排查
			global.TRON_LOG.Error(fmt.Sprintf("blockNum :%d InsertData err: %s", data[0]["block_num"].(uint64), err))
			dataChan <- data
		}
		server.WriteLog(data[0]["block_num"].(uint64), err)
	}
}

func Run(runBlockNum uint64) {
	nowBlockNum := global.GetTronNowBlock()
	dataChan := make(chan []map[string]interface{}, 10000)
	go InsertData(dataChan)
	for {
		if nowBlockNum == runBlockNum {
			nowBlockNum = global.GetTronNowBlock()
			time.Sleep(5 * time.Second)
			continue
		}
		usdtTradeData := make([]map[string]interface{}, 0)
		data := global.GetTronData(runBlockNum)

		if _, ok := data["transactions"]; !ok {
			server.WriteLog(runBlockNum, nil)
			runBlockNum++
			continue
		}
		for _, trx := range data["transactions"].([]interface{}) {
			jsonData, _ := json.Marshal(trx)
			var transactionData TransactionData
			_ = json.Unmarshal(jsonData, &transactionData)
			if ok := transactionData.IsUSDTTransfer(); ok {
				usdtData := AnalysisData(&transactionData, &runBlockNum)
				usdtTradeData = append(usdtTradeData, usdtData)
			}

		}
		if len(usdtTradeData) == 0 {
			server.WriteLog(runBlockNum, nil)
			runBlockNum++
			continue
		}
		dataChan <- usdtTradeData
		runBlockNum++
	}

}

// AnalysisData 解析数据
func AnalysisData(usdtData *TransactionData, runBlockNum *uint64) map[string]interface{} {
	var fromAddress, toAddress string
	var value uint64
	methID := usdtData.RawData.Contract[0].Parameter.Value.Data[0:8]
	switch methID {
	case Transfer:
		toAddress, value = util.TransferParseData(usdtData.RawData.Contract[0].Parameter.Value.Data)
		fromAddress, _ = util.FromHexAddress(usdtData.RawData.Contract[0].Parameter.Value.OwnerAddress)
	case TransferFrom:
		fromAddress, toAddress, value = util.TransferFromParseData(usdtData.RawData.Contract[0].Parameter.Value.Data)
	}
	data := map[string]interface{}{
		"hash":         usdtData.TxID,
		"from_address": fromAddress,
		"to_address":   toAddress,
		"value":        value,
		"block_num":    *runBlockNum,
	}
	return data

}
