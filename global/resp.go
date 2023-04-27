package global

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// GetTronData 通过指定块高获取数据
func GetTronData(runBlockNum uint64) map[string]interface{} {
	data := make(map[string]interface{})
	for {
		response, err := http.Get(TRON_CONFIG.Tron.Dns() + "/wallet/getblockbynum?num=" + strconv.FormatUint(runBlockNum, 10))
		if err != nil {
			time.Sleep(5 * time.Second)
			TRON_LOG.Warn(fmt.Sprint("request wallet getblock err:", err))
			continue
		}
		_ = json.NewDecoder(response.Body).Decode(&data)
		_ = response.Body.Close()
		break
	}
	return data
}

// GetTronNowBlock 获取当前钱包最新块高
func GetTronNowBlock() (blockNum uint64) {
	data := make(map[string]interface{})
	for {
		response, err := http.Get(TRON_CONFIG.Tron.Dns() + "/wallet/getnowblock")
		if err != nil {
			time.Sleep(5 * time.Second)
			TRON_LOG.Warn(fmt.Sprint("request wallet getnowblock err:", err))
			continue
		}
		_ = json.NewDecoder(response.Body).Decode(&data)
		_ = response.Body.Close()
		break
	}
	blockNum = uint64(data["block_header"].(map[string]interface{})["raw_data"].(map[string]interface{})["number"].(float64))
	return blockNum
}
