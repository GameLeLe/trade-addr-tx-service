package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/GameLeLe/trade-addr-tx-service/hdwallet"
)

var wg sync.WaitGroup
var daConfig *DigitalAssetsConfig

func main() {
	daConfig, err := ParseConfig("./config.toml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	port := 8095
	server := newRPCServer(port, &wg)
	ethPubKey, _ := hdwallet.ReadWalletFromFile(daConfig.ETHMasterPubKeyFile)
	btcPubKey, _ := hdwallet.ReadWalletFromFile(daConfig.BTCMasterPubKeyFile)
	server.start(ethPubKey, btcPubKey)
	wg.Wait()
}
