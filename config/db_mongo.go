package config

import (
	"fmt"
)

type MongoDB struct {
	Db       string `mapstructure:"db" json:"db"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
}

func (m *MongoDB) Dns() string {
	if m.Username != "" && m.Password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%d", m.Username, m.Password, m.Host, m.Port)
	} else {

		return fmt.Sprintf("mongodb://%s:%d", m.Host, m.Port)
	}
}
