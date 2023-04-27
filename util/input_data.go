package util

import (
	"math/big"
	"strings"
)

func TransferParseData(data string) (string, uint64) {
	//address
	toAddress := data[30:72]
	b58ToAddress := addressHexToB58("41" + toAddress[2:])
	// amount
	hexAmount := strings.TrimPrefix(data[72:136], "0000000000000000000000000000000000000000000000000000000")
	if hexAmount == "" {
		return b58ToAddress, 0
	}
	amount := hexToUint64(hexAmount)
	return b58ToAddress, amount
}
func TransferFromParseData(data string) (string, string, uint64) {
	fromAddress := data[30:72]
	b58FromAddress := addressHexToB58("41" + fromAddress[2:])
	toAddress := data[94:136]
	b58ToAddress := addressHexToB58("41" + toAddress[2:])
	hexAmount := strings.TrimPrefix(data[136:200], "0000000000000000000000000000000000000000000000000000000")
	if hexAmount == "" {
		return b58FromAddress, b58ToAddress, 0
	}
	amount := hexToUint64(hexAmount)
	return b58FromAddress, b58ToAddress, amount
}

// hexToUint64
func hexToUint64(hex string) uint64 {
	//hex = strings.Replace(hex, "0x", "", -1)
	n := new(big.Int)
	n, _ = n.SetString(hex, 16)
	return n.Uint64()
}

// addressHexToB58
func addressHexToB58(hexAddress string) string {
	address, _ := FromHexAddress(hexAddress)
	return address
}
