package core

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"tron/core/internal"
	"tron/global"
	"tron/util"
)

func Zap() (logger *zap.Logger) {
	if ok, _ := util.PathExists(global.TRON_CONFIG.Zap.Director); !ok {
		fmt.Printf("create %v directory\n", global.TRON_CONFIG.Zap.Director)
		_ = os.Mkdir(global.TRON_CONFIG.Zap.Director, os.ModePerm)
	}

	cores := internal.Zap.GetZapCores()
	logger = zap.New(zapcore.NewTee(cores...))

	if global.TRON_CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}
