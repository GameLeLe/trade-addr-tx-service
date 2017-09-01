package main

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/GameLeLe/trade-addr-tx-service/base58check"
	"github.com/GameLeLe/trade-addr-tx-service/hdwallet"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ripemd160"
)

func genETHAddr(compressedKey []byte) string {
	x, y := hdwallet.Expand(compressedKey)
	four, _ := hex.DecodeString("04")
	paddedKey := append(four, append(x.Bytes(), y.Bytes()...)...)
	pubKey := crypto.ToECDSAPub(paddedKey)
	addr := crypto.PubkeyToAddress(*pubKey)
	return hex.EncodeToString(addr[:])
}
func genBTCAddr(compressedKey []byte, isTestnet bool) string {
	var publicKeyPrefix byte
	if isTestnet {
		publicKeyPrefix = 0x6F
	} else {
		publicKeyPrefix = 0x00
	}
	shadPublicKeyBytes := sha256.Sum256(compressedKey)
	ripeHash := ripemd160.New()
	ripeHash.Write(shadPublicKeyBytes[:])
	ripeHashedBytes := ripeHash.Sum(nil)

	address := base58check.Encode(publicKeyPrefix, ripeHashedBytes)
	return address
}
