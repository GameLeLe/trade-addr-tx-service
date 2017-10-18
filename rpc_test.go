package main

import (
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	addrtx "github.com/GameLeLe/trade-addr-tx-service/thrift/addrtx"
	"github.com/stretchr/testify/assert"
)

func TestGetAddr(t *testing.T) {
	var wg sync.WaitGroup
	port := 8095
	server := newRPCServer(port, &wg)
	go server.start()
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
		{0, "0x5688cC7C4e79B6E65cffbE4e53A588a094324467"},
		{2, "0xDDa4D8D48Eff4355cC1375c54E79c563f0284f46"},
		{6, "0x471d34c3080aBE8Cd989706E64Af3D938D58D185"},
		{9, "0xddA8AE84a3d6bA0AABf6e5B35331E78a5735321f"},
		{12, "0x6af7DC6322cF03373FF39193c97a0CED888BF21a"},
		{33, "0x7E3105EF8F1A9eDec032AB230b30E9e2A7eaD3Ee"},
		{41, "0x150fF0b9B4C2aC7162e0009FE8244dbceCd1352A"},
		{45, "0x7A66f224ed65dBdD96c538Ede50E54fe80082D23"},
		{58, "0x1a461056433b93cfefc569aaB95f7bC40E5FC5e1"},
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
