package main

import (
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
	childpub44_0_0 := childprv44_0_0.Pub()
	if err != nil {
		t.Errorf("get master pub key child 44/60/0/0 error: %v", err)
	}
	childpub44_0_0_0, err := childpub44_0_0.Child(0)
	if err != nil {
		t.Errorf("get master pub key child 44/60/0/0 error: %v", err)
	}
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
	//m/44'/0'/0'/0/*
	seed := getSeed()
	masterprv := hdwallet.MasterKey(seed)
	childprv44, err := masterprv.Child(2147483692)
	if err != nil {
		t.Errorf("get master pub key child 44 error: %v", err)
	}
	childprv44_0, err := childprv44.Child(2147483708)
	if err != nil {
		t.Errorf("get master pub key child 44/60 error: %v", err)
	}
	childprv44_0_0, err := childprv44_0.Child(2147483648)
	if err != nil {
		t.Errorf("get master pub key child 44/60/0 error: %v", err)
	}
	childpub44_0_0 := childprv44_0_0.Pub()
	if err != nil {
		t.Errorf("get master pub key child 44/60/0/0 error: %v", err)
	}
	childpub44_0_0_0, err := childpub44_0_0.Child(0)
	if err != nil {
		t.Errorf("get master pub key child 44/60/0/0 error: %v", err)
	}
	if err != nil {
		t.Errorf("get master pub key child 44/60/0/0 error: %v", err)
	}
	cases := []struct {
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
	for _, c := range cases {
		uid := c.uid
		expectedAddr := strings.ToLower(c.expected)
		childpubUID, err := childpub44_0_0_0.Child(uint32(uid))
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
