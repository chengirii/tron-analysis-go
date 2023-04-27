package config

import (
	"fmt"
	"net/http"
	"strconv"
)

type Tron struct {
	Tokens     string `mapstructure:"tokens" json:"tokens"`
	WalletHost string `mapstructure:"wallet-host" json:"wallet-host"`
	WalletPort int    `mapstructure:"wallet-port" json:"wallet-port"`
}

func (t *Tron) Dns() string {
	return "http://" + t.WalletHost + ":" + strconv.Itoa(t.WalletPort)
}

func (t *Tron) VerifyConnectivity() {
	if t.WalletHost == "" || t.WalletPort == 0 {
		panic("Wallet address or port cannot be empty")
	}
	hostname := t.Dns() + "/wallet/getnowblock"
	resp, err := http.Get(hostname)
	if err != nil {
		panic(fmt.Errorf("connet fail check tron config", err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("Failed to access URL %s. Status code: %d", hostname, resp.StatusCode))
	}
}
