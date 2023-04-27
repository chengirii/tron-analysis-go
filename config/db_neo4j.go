package config

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4j struct {
	Username    string `mapstructure:"username" json:"username"`
	Password    string `mapstructure:"password" json:"password"`
	Host        string `mapstructure:"host" json:"host"`
	Port        int    `mapstructure:"port" json:"port"`
	MaxPoolSize int    `mapstructure:"max_pool_size" json:"max_pool_size"`
}
type Neo4jConn struct {
	driver neo4j.Driver
}

func (config *Neo4j) NewConnect() *Neo4jConn {
	driver, _ := neo4j.NewDriver(
		fmt.Sprintf("bolt://%s:%d", config.Host, config.Port),
		neo4j.BasicAuth(config.Username, config.Password, ""),
	)
	if err := driver.VerifyConnectivity(); err != nil {
		panic(fmt.Errorf("connet fail check neo4j config : %s\n", err))
	}
	// session := driver.NewSession(neo4j.SessionConfig{})

	return &Neo4jConn{driver: driver}
}

func (config *Neo4jConn) RunCypher(query string, params map[string]interface{}) (result neo4j.Result, err error) {
	session := config.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		if sessionErr := session.Close(); sessionErr != nil && err == nil {
			err = sessionErr
		}
	}()
	result, err = session.Run(query, params)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		record := result.Record()
		_ = record.Values
		fmt.Println(record)
	}
	return result, nil
}
func (config *Neo4jConn) Close() {
	_ = config.driver.Close()
}
