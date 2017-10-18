# trade-addr-tx-service
trade service for public address and transaction creation

Example:

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"sort"

	bip39 "github.com/GameLeLe/trade-addr-tx-service/bip39"
	"github.com/GameLeLe/trade-addr-tx-service/btc"
	"github.com/GameLeLe/trade-addr-tx-service/hdwallet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func main() {
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	//fmt.Println(mnemonic)
	mnemonic = "duty capital transfer goose segment trap good kite ramp before amused fiber alter awful into chair smile erupt burger scare culture quote visit dragon"

	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, "")

	// Create a master private key
	masterprv := hdwallet.MasterKey(seed)

	// Convert a private key to public key
	masterpub := masterprv.Pub()
	// type HDWallet struct {
	// 	Vbytes      []byte //4 bytes
	// 	Depth       uint16 //1 byte
	// 	Fingerprint []byte //4 bytes
	// 	I           []byte //4 bytes
	// 	Chaincode   []byte //32 bytes
	// 	Key         []byte //33 bytes
	// }

	// vBytes := masterpub.Vbytes
	// depth := masterpub.Depth
	// fingerPrint := masterpub.Fingerprint
	// i := masterpub.I
	// chainCode := masterpub.Chaincode
	// key := masterpub.Key

	buf := new(bytes.Buffer)
	tmpBytes := make([]byte, 2)
	//write vBytes
	binary.BigEndian.PutUint16(tmpBytes, uint16(len(masterpub.Vbytes)))
	buf.Write(tmpBytes)
	buf.Write(masterpub.Vbytes)
	//write depth
	//binary.BigEndian.PutUint16(tmpBytes, uint16(masterpub.Depth))
	//write Fingerprint
	binary.BigEndian.PutUint16(tmpBytes, uint16(len(masterpub.Fingerprint)))
	buf.Write(tmpBytes)
	buf.Write(masterpub.Fingerprint)
	//write I
	binary.BigEndian.PutUint16(tmpBytes, uint16(len(masterpub.I)))
	buf.Write(tmpBytes)
	buf.Write(masterpub.I)
	//write Chaincode
	binary.BigEndian.PutUint16(tmpBytes, uint16(len(masterpub.Chaincode)))
	buf.Write(tmpBytes)
	buf.Write(masterpub.Chaincode)
	//write Key
	binary.BigEndian.PutUint16(tmpBytes, uint16(len(masterpub.Key)))
	buf.Write(tmpBytes)
	buf.Write(masterpub.Key)
	//write to file
	ioutil.WriteFile("./master_pubkey", buf.Bytes(), 0666)

	//read from file
	fi, _ := os.Open("./master_pubkey")
	defer fi.Close()
	masterPubInFile := &hdwallet.HDWallet{}
	r := bufio.NewReader(fi)
	//read vbytes
	r.Read(tmpBytes)
	length := binary.BigEndian.Uint16(tmpBytes)
	masterPubInFile.Vbytes = make([]byte, length)
	r.Read(masterPubInFile.Vbytes)
	//read depth
	//r.Read(tmpBytes)
	//masterPubInFile.Depth = binary.BigEndian.Uint16(tmpBytes)
	//read Fingerprint
	r.Read(tmpBytes)
	length = binary.BigEndian.Uint16(tmpBytes)
	masterPubInFile.Fingerprint = make([]byte, length)
	r.Read(masterPubInFile.Fingerprint)
	//read I
	r.Read(tmpBytes)
	length = binary.BigEndian.Uint16(tmpBytes)
	masterPubInFile.I = make([]byte, length)
	r.Read(masterPubInFile.I)
	//read Chaincode
	r.Read(tmpBytes)
	length = binary.BigEndian.Uint16(tmpBytes)
	masterPubInFile.Chaincode = make([]byte, length)
	r.Read(masterPubInFile.Chaincode)
	//read Chaincode
	r.Read(tmpBytes)
	length = binary.BigEndian.Uint16(tmpBytes)
	masterPubInFile.Key = make([]byte, length)
	r.Read(masterPubInFile.Key)

	// fmt.Println("==========master private key===========")
	// fmt.Println(masterprv.String())

	fmt.Println("==========master public key===========")
	fmt.Println(masterpub.String())
	fmt.Println("==========master public key in file===========")
	fmt.Println(masterPubInFile.String())

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

func main() {
	rpcServe := newRPCServer(8099)
	rpcServe.start()
}

func main() {
	mnemonic := "duty capital transfer goose segment trap good kite ramp before amused fiber alter awful into chair smile erupt burger scare culture quote visit dragon"
	password := "222222"
	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, password)
	// Create a master private key
	masterprv := hdwallet.MasterKey(seed)
	// Convert a private key to public key
	masterpub := masterprv.Pub()
	childpub0, _ := masterpub.Child(0)
	// Generate new child key based on private or public key
	//childprv, err := masterprv.Child(0)
	fromUID := 100
	//fromAmount := 1
	toUID := 101
	//toAmount := 1
	childpubFromUID, _ := childpub0.Child(uint32(fromUID))
	childpubToUID, _ := childpub0.Child(uint32(toUID))

	// secp256k1 := btcec.S256()
	// key, err := btcec.ParsePubKey(pubKeyByte, secp256k1)
	// if err != nil {
	// 	return nil, err
	// }

	var totalAmount uint64
	totalAmount = 1

	fromAddr := genBTCAddr(childpubFromUID.Pub().Key, false)
	fromKey := &btc.Key{}
	fromPub, _ := btc.GetPublicKey(childpubFromUID.Pub().Key, false)
	fromKey.Pub = fromPub

	toAddr := genBTCAddr(childpubToUID.Pub().Key, false)
	toKey := &btc.Key{}
	toPub, _ := btc.GetPublicKey(childpubFromUID.Pub().Key, false)
	toKey.Pub = toPub

	service, _ := btc.NewBlockrService()
	utxos, _ := service.GetUTXO(fromAddr, fromKey)

	sort.Sort(utxos)

	tx := btc.TX{}

	txins := make([]*btc.TXin, 0, 10)
	var amount uint64
	for i := range utxos {
		utxo := utxos[len(utxos)-1-i]
		txin := btc.TXin{}
		txin.Hash = utxo.Hash
		txin.Index = utxo.Index
		txin.Sequence = uint32(0xffffffff)
		txin.PrevScriptPubkey = utxo.Script
		txin.CreateScriptSig = nil
		txins = append(txins, &txin)
		if amount += utxo.Amount; amount >= totalAmount {
			break
		}
	}

	txouts := make([]*btc.TXout, 0, 10)
	txout := &btc.TXout{}
	txout.Value = totalAmount
	txout.ScriptPubkey, _ = btc.CreateP2PKHScriptPubkey(toAddr)
	txouts = append(txouts, txout)
	if amount-totalAmount > 0 {
		txout := &btc.TXout{}
		txout.Value = amount - totalAmount
		txout.ScriptPubkey, _ = btc.CreateP2PKHScriptPubkey(fromAddr)
		txouts = append(txouts, txout)
	}
	rawtx, err := tx.MakeTX()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rawtx)

	// //get utxo
	// service, _ := btc.NewBlockrService()
	// service.GetUTXO(genBTCAddr(childpubFromUID.Pub().Key))

	// var rawtx []byte
	// tx := btc.TX{}
	// tx.Locktime = 0
	// var kb uint = 1

	// fmt.Println(childpubFromUID)
	// fmt.Println(childpubToUID)
}

func main() {
	//fromAddrStr := "0xbAc66419aC1DCD91c5CBAbEAAd865d70F48DD7ac"
	toAddrStr := "0x34B02FdF4De0048AE3e3e159268486A3A43C6940"
	//fromAddr := common.HexToAddress(fromAddrStr)
	toAddr := common.HexToAddress(toAddrStr)
	var totalAmount *big.Int
	var nonce uint64
	totalAmount = new(big.Int)
	totalAmount.SetInt64(1)
	nonce = 1
	tx := types.NewTransaction(nonce, toAddr, totalAmount, nil, nil, nil)
	//hex.EncodeToString(tx.MarshalJSON())
	fmt.Println(tx.String())
}

