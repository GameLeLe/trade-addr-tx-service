package main

import (
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	hdwallet "github.com/GameLeLe/trade-addr-tx-service/hdwallet"
	addrtx "github.com/GameLeLe/trade-addr-tx-service/thrift/addrtx"
	"github.com/stretchr/testify/assert"
)

func TestGetAddr(t *testing.T) {
	var wg sync.WaitGroup
	daConfig, err := ParseConfig("config.toml")
	if err != nil {
		t.Errorf("parse config file error: %v", err)
	}

	port := 8095
	server := newRPCServer(port, &wg)
	ethPubKey, _ := hdwallet.ReadWalletFromFile(daConfig.ETHMasterPubKeyFile)
	btcPubKey, _ := hdwallet.ReadWalletFromFile(daConfig.BTCMasterPubKeyFile)
	go server.start(ethPubKey, btcPubKey)
	time.Sleep(100 * time.Millisecond)

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket("127.0.0.1:" + strconv.Itoa(port))
	if err != nil {
		t.Errorf("error resolving address: %v", err)
	}
	useTransport := transportFactory.GetTransport(transport)
	client0 := addrtx.NewAddrTXServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		t.Errorf("Error opening socket: %v", err)
	}

	btcCases := []struct {
		uid      int
		expected string
	}{
		{0, "13tBtZwgZ7usfEfbf7bKcErY9AimBzNNUq"},
		{2, "16kWBh377tPqZmXu4XTz3dGzgxFLKUL7ey"},
		{6, "1Be99wrYzyztKGVPkxAuJcLyZGArcUcdVg"},
		{9, "1HJnu3GeY16UjdzFVsqqdKTKxkJ2u25ZGN"},
		{12, "19Xpf841sPbPW2j8WsEvPKzsQVA1JAAant"},
		{33, "1Ka8BG5eEnbKkXqsa5F2QZfGYVxpLnNWo1"},
		{41, "14WPkc3VF2BxQbNdhsjuVUnzYrzjJhF7mU"},
		{45, "1NsfBLTCb3QpKgVwLwJtuqRuUEYV8qwrMo"},
		{58, "19CXZnEr9XiTdKmvQWEyrHLAPDLBnt8v3r"},
	}
	for _, c := range btcCases {
		msg := &addrtx.GetAddrMsg{}
		msg.UID = int64(c.uid)
		msg.CoinType = "BTC"
		ret, err := client0.GetAddr(msg)
		if err != nil {
			t.Errorf("get message from server error return: %v", err)
		}
		expected := strings.ToLower(c.expected)
		ret = strings.ToLower(ret)
		assert.Equal(t, expected, ret, "btc address not matched: %s|%s", expected, ret)
	}

	ethCases := []struct {
		uid      int
		expected string
	}{
		{0, "0x20EBEf408Dd557df2C1F3c1Ff3c541655f30c68D"},
		{2, "0x3726Ad8dF5C8eBBDac562e32d08Ad50a1F6aC4cD"},
		{6, "0x81b0ba94a11632E478C12E974b027E54f3e9226a"},
		{9, "0xaEA3d615e5822F38931083AF70d5Cfb0A7eD0d73"},
		{12, "0xd3b212FF285FF13ffad6f09543bE235d0241E272"},
		{33, "0x071971219EFbd8177d353D13335365fF9d800eC7"},
		{41, "0xdf6EC73F12252B0682b6d17961B1d0587FFbE8e4"},
		{45, "0xaE8b5Fa895EdfF1D5Fdd11140e2D9F0EDf51Ef81"},
		{58, "0x9a5fe3ABcEae5472a553266a35A3a02ae7Df23ad"},
	}
	for _, c := range ethCases {
		msg := &addrtx.GetAddrMsg{}
		msg.UID = int64(c.uid)
		msg.CoinType = "ETH"
		ret, err := client0.GetAddr(msg)
		if err != nil {
			t.Errorf("get message from server error return: %v", err)
		}
		expected := strings.ToLower(c.expected)
		ret = strings.ToLower("0x" + ret)
		assert.Equal(t, expected, ret, "eth address not matched: %s|%s", expected, ret)
	}

	transport.Close()
	server.stop()
	wg.Wait()
}
