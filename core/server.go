package core

import (
	"fmt"
	"time"
	"tron/global"
	"tron/server"
	"tron/server/usdt"
)

func RunServer() {
	global.TRON_CONFIG.Tron.VerifyConnectivity()
	missBlock := server.CheckMissingBlocks()
	if len(missBlock) > 0 {
		for _, v := range missBlock {
			for key, value := range v {
				global.TRON_LOG.Warn(fmt.Sprintf("block num %d  err: %s ", key, value))
				fmt.Println(fmt.Sprintf("block num %d  err: %s ", key, value))
			}
		}
		panic("missing block")
	}
	runBlockNum := server.ReadEndBlockNum() + 1 // 从上次结束的块高开始
	global.TRON_LOG.Info(fmt.Sprintf("server run success start block num %d", runBlockNum))
	fmt.Printf(`
	欢迎使用 波场区块链浏览器分析
	钱包节点: %s
	代币: trc20-%s
	运行块高: %d
	启动时间: %s`, global.TRON_CONFIG.Tron.Dns(), global.TRON_CONFIG.Tron.Tokens, runBlockNum, time.Now().Format("2006-01-02 15:04:05"))
	switch {
	case global.TRON_CONFIG.Tron.Tokens == "usdt":
		usdt.Run(runBlockNum)
	}

}
