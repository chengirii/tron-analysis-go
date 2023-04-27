package initialize

import (
	"tron/config"
	"tron/global"
)

func Neo4j() *config.Neo4jConn {
	return global.TRON_CONFIG.Neo4j.NewConnect()
}
