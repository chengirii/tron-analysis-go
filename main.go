package main

import (
	"tron/core"
	"tron/global"
	"tron/initialize"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

func main() {

	global.TRON_VP = core.Viper()       // 初始化Viper
	global.TRON_DB = initialize.Neo4j() // 连接数据库
	global.TRON_LOG = core.Zap()        // 初始zap日志库
	global.TRON_MongoDB = initialize.MongoDB()
	core.RunServer()

}
