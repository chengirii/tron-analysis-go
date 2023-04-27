package global

import (
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"tron/config"
)

var (
	TRON_DB      *config.Neo4jConn
	TRON_VP      *viper.Viper
	TRON_CONFIG  config.Server
	TRON_LOG     *zap.Logger
	TRON_MongoDB *mongo.Client
)
