package core

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"tron/global"
)

func Viper(path ...string) *viper.Viper {
	var config string
	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "choose config file")
		flag.Parse()
	}
	// 测试
	if config == "" {
		panic("no config file specified")
	}
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if err = v.Unmarshal(&global.TRON_CONFIG); err != nil {
		fmt.Println(err)
	}

	return v
}
