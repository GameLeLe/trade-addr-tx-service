package main

import (
	"fmt"
	"os"

	"git.apache.org/thrift.git/lib/go/thrift"
	addrtx "github.com/GameLeLe/trade-addr-tx-service/thrift/addrtx"
	"github.com/surge/glog"
)

func main0() {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket("127.0.0.1:8099")
	if err != nil {
		glog.Error("error resolving address:", err)
		os.Exit(1)
	}
	useTransport := transportFactory.GetTransport(transport)
	client0 := addrtx.NewAddrTXServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		glog.Error("Error opening socket", err)
		os.Exit(1)
	}
	defer transport.Close()
	msg := &addrtx.GetAddrMsg{}
	msg.UID = 10
	msg.CoinType = "BTC"
	ret, _ := client0.GetAddr(msg)
	fmt.Println("addr:", ret)
}
