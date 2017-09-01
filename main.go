package main

import (
	"encoding/hex"
	"fmt"

	"crypto/sha256"

	"github.com/GameLeLe/trade-addr-tx-service/base58check"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ripemd160"

	bip39 "github.com/GameLeLe/trade-addr-tx-service/bip39"
	hdwallet "github.com/GameLeLe/trade-addr-tx-service/hdwallet"
)

func main() {
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	fmt.Println(mnemonic)
	//mnemonic = "duty capital transfer goose segment trap good kite ramp before amused fiber alter awful into chair smile erupt burger scare culture quote visit dragon"

	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, "222222")

	// Create a master private key
	masterprv := hdwallet.MasterKey(seed)

	// Convert a private key to public key
	masterpub := masterprv.Pub()

	fmt.Println("==========master private key===========")
	fmt.Println(masterprv.String())

	fmt.Println("==========master public key===========")
	fmt.Println(masterpub.String())

	// Generate new child key based on private or public key
	//childprv, err := masterprv.Child(0)
	childpub0, _ := masterpub.Child(0)
	childprv0, _ := masterprv.Child(0)

	fmt.Println("==========m/0 public key===========")
	fmt.Println(childpub0.String())
	fmt.Println("==========m/0 private key===========")
	fmt.Println(childprv0.String())

	//childstring0, _ := hdwallet.StringChild(walletstring0, 0)
	// childstring1, _ := hdwallet.StringChild(walletstring1, 0)
	// childstring2, _ := hdwallet.StringChild(walletstring2, 0)
	childpub00, _ := childpub0.Child(0)
	childpub01, _ := childpub0.Child(1)
	childpub02, _ := childpub0.Child(2)

	childprv00, _ := childprv0.Child(0)
	childprv01, _ := childprv0.Child(1)
	childprv02, _ := childprv0.Child(2)
	//HDWallet
	// childpub00.Pub().
	fmt.Println("==========address============")
	fmt.Println(childpub00.Pub().Address())
	fmt.Println(childpub01.Pub().Address())
	fmt.Println(childpub02.Pub().Address())
	fmt.Println("==========pub key============")
	fmt.Println(hex.EncodeToString(childpub00.Pub().Key))
	fmt.Println(hex.EncodeToString(childpub01.Pub().Key))
	fmt.Println(hex.EncodeToString(childpub02.Pub().Key))
	fmt.Println("==========prv key============")
	fmt.Println(hex.EncodeToString(childprv00.Key))
	fmt.Println(hex.EncodeToString(childprv01.Key))
	fmt.Println(hex.EncodeToString(childprv02.Key))
	//address
	fmt.Println("==========eth address============")
	fmt.Println(genETHAddr(childpub00.Pub().Key))
	fmt.Println(genETHAddr(childpub01.Pub().Key))
	fmt.Println(genETHAddr(childpub02.Pub().Key))
	fmt.Println("==========bitcoin address============")
	fmt.Println(genBTCAddr(childpub00.Pub().Key, false))
	fmt.Println(genBTCAddr(childpub01.Pub().Key, false))
	fmt.Println(genBTCAddr(childpub02.Pub().Key, false))
}

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
