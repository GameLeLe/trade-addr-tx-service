package main

import (
	"fmt"
	"net"
	"strings"
	"testing"

	bip39 "github.com/GameLeLe/trade-addr-tx-service/bip39"
	hdwallet "github.com/GameLeLe/trade-addr-tx-service/hdwallet"
)

func TestBTCAddrBIP32(t *testing.T) {
	seed := getSeed()
	masterprv := hdwallet.MasterKey(seed)
	masterpub := masterprv.Pub()
	childpub0, err := masterpub.Child(0)
	if err != nil {
		t.Errorf("get master pub key child error: %v", err)
	}
	cases := []struct {
		uid      int
		expected string
	}{
		{0, "1GgVcAD7jPgcf9QcHoru7SHr5oKB7rZTXN"},
		{2, "1wuTBRg16T9W3MxkqiADaPnhNN87nYSuq"},
		{6, "15AwR5TAMKoKUygqrRhjyVQK6ExFPBe2Rb"},
		{9, "1HupN5rQd5hkJuHrpYRRve9LJcDX33JVU5"},
		{12, "1Jtg7RTLSEvLzCbicRx31vr1b1NFACsrfm"},
		{33, "1FhVU7ajqiXQJDKaMZaSB2bgHdfS4dGZbV"},
		{41, "1MJChbDtn3z3KDD1pcgMfLirvJPeqdpfmn"},
		{45, "1Np99xgGb8RCdHSywfX3cLGUc7MsRNaou"},
		{58, "147wTsyhhCoPpdUYpEKHCArXrfnLjyRtQz"},
	}
	for _, c := range cases {
		uid := c.uid
		expectedAddr := c.expected
		childpubUID, err := childpub0.Child(uint32(uid))
		if err != nil {
			t.Errorf("get uid related pub key error: %v", err)
		}
		addr := genBTCAddr(childpubUID.Pub().Key, false)
		if addr != expectedAddr {
			t.Errorf("BTC addr not matched: %s|%s", addr, expectedAddr)
		}
	}
}

func TestETHAddrBIP32(t *testing.T) {
	seed := getSeed()
	masterprv := hdwallet.MasterKey(seed)
	masterpub := masterprv.Pub()
	childpub0, err := masterpub.Child(0)
	if err != nil {
		t.Errorf("get master pub key child error: %v", err)
	}
	cases := []struct {
		uid      int
		expected string
	}{
		{0, "0x7cA63F24d0B9e052459c58037028ad6A3dF35410"},
		{2, "0x646bA6727ce2708839B0DA6FFacAb0878bbB8862"},
		{6, "0x4A569ae33b39a06D0c7B7f287149925dFf583FF2"},
		{9, "0x82f32A3749a494Dfc3cE4cea5f086a2157Db8eAb"},
		{12, "0x9D76b584e41B6a32a3908Cd38B5E3a11F81f4dbd"},
		{33, "0xAE9d9f072c0b9ABBc9187d8e6789815299B555b7"},
		{41, "0xa846057EDa04F3B9332D19caaD502328F4F880Ab"},
		{45, "0x51605b60fE8Ad3418F8F9Cd26647596fD27521cf"},
		{58, "0xE8773172E36E18FF18BEF7A21d09287242D81f8C"},
	}
	for _, c := range cases {
		uid := c.uid
		expectedAddr := strings.ToLower(c.expected)
		childpubUID, err := childpub0.Child(uint32(uid))
		if err != nil {
			t.Errorf("get uid related pub key error: %v", err)
		}
		addr := strings.ToLower("0x" + genETHAddr(childpubUID.Pub().Key))
		if addr != expectedAddr {
			t.Errorf("ETH addr not matched: %s|%s", addr, expectedAddr)
		}
	}
}

func TestBTCAddrBIP44(t *testing.T) {
	//m/44'/0'/0'/0/*
	seed := getSeed()
	masterprv := hdwallet.MasterKey(seed)
	childprv44, err := masterprv.Child(2147483692)
	if err != nil {
		t.Errorf("get master pub key child 44 error: %v", err)
	}
	childprv44_0, err := childprv44.Child(2147483648)
	if err != nil {
		t.Errorf("get master pub key child 44/60 error: %v", err)
	}
	childprv44_0_0, err := childprv44_0.Child(2147483648)
	if err != nil {
		t.Errorf("get master pub key child 44/60/0 error: %v", err)
	}
	childprv44_0_0_0, err := childprv44_0_0.Child(0)
	if err != nil {
		t.Errorf("get master pub key child 44/60/0/0 error: %v", err)
	}
	childpub44_0_0_0 := childprv44_0_0_0.Pub()
	if err != nil {
		t.Errorf("get master pub key child 44/60/0/0 error: %v", err)
	}
	cases := []struct {
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
	for _, c := range cases {
		uid := c.uid
		expectedAddr := c.expected
		childpubUID, err := childpub44_0_0_0.Child(uint32(uid))
		if err != nil {
			t.Errorf("get uid related pub key error: %v", err)
		}
		addr := genBTCAddr(childpubUID.Pub().Key, false)
		if addr != expectedAddr {
			t.Errorf("BTC addr not matched: %s|%s", addr, expectedAddr)
		}
	}
}

func TestETHAddrBIP44(t *testing.T) {
	//m/44'/60'/0'/*
	seed := getSeed()
	masterprv := hdwallet.MasterKey(seed)
	childprv44, err := masterprv.Child(2147483692)
	if err != nil {
		t.Errorf("get master pub key child 44 error: %v", err)
	}
	childprv44_60, err := childprv44.Child(2147483708)
	if err != nil {
		t.Errorf("get master pub key child 44/60 error: %v", err)
	}
	childprv44_60_0, err := childprv44_60.Child(2147483648)
	if err != nil {
		t.Errorf("get master pub key child 44/60/0 error: %v", err)
	}
	// Convert a private key to public key
	masterpub := childprv44_60_0.Pub()

	cases := []struct {
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
	for _, c := range cases {
		uid := c.uid
		expectedAddr := strings.ToLower(c.expected)
		childpubUID, err := masterpub.Child(uint32(uid))
		if err != nil {
			t.Errorf("get uid related pub key error: %v", err)
		}
		addr := strings.ToLower("0x" + genETHAddr(childpubUID.Pub().Key))
		if addr != expectedAddr {
			t.Errorf("ETH addr not matched: %s|%s", addr, expectedAddr)
		}
	}
}

func getSeed() []byte {
	mnemonic := "bird march express devote nature tone rich shadow invest husband table chicken input pull zero shove stage typical color chimney fat entire split aware"
	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	return bip39.NewSeed(mnemonic, "")
}

func TestSocket(t *testing.T) {
	conn, err := net.Dial("tcp", "192.168.1.25:3250")
	defer conn.Close()
	if err != nil {
		t.Errorf("conn error:%v", err)
	}
	for {
		tmp := make([]byte, 1)
		_, err := conn.Read(tmp)
		if err != nil {
			t.Errorf("read error:%v", err)
			break
		}
	}
	fmt.Println("break loop")
}
