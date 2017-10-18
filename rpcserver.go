package main

import (
	"errors"
	"log"
	"strconv"
	"sync"

	"git.apache.org/thrift.git/lib/go/thrift"
	bip39 "github.com/GameLeLe/trade-addr-tx-service/bip39"
	hdwallet "github.com/GameLeLe/trade-addr-tx-service/hdwallet"
	addrtx "github.com/GameLeLe/trade-addr-tx-service/thrift/addrtx"
)

type rpcServer struct {
	started      bool
	wg           *sync.WaitGroup
	addr         string
	thriftServer *thrift.TSimpleServer
}

func newRPCServer(port int, wg *sync.WaitGroup) *rpcServer {
	server := &rpcServer{}
	server.addr = "0.0.0.0:" + strconv.Itoa(port)
	server.wg = wg
	return server
}

type rpcThrift struct {
	ethPubKey *hdwallet.HDWallet
	btcPubKey *hdwallet.HDWallet
}

func (rpcT *rpcThrift) GetTX(msg *addrtx.GetTXMsg) (string, error) {
	coinType := msg.CoinType
	fromUID := msg.FromUID
	totalAmount := msg.FromAmount
	toUID := msg.ToUID
	mnemonic := "duty capital transfer goose segment trap good kite ramp before amused fiber alter awful into chair smile erupt burger scare culture quote visit dragon"
	password := "222222"
	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, password)
	// Create a master private key
	masterprv := hdwallet.MasterKey(seed)
	// Convert a private key to public key
	masterpub := masterprv.Pub()
	// Generate new child key based on private or public key
	//childprv, err := masterprv.Child(0)
	childpub0, _ := masterpub.Child(0)
	childpubFrom, _ := childpub0.Child(uint32(fromUID))
	childpubTO, _ := childpub0.Child(uint32(toUID))
	switch coinType {
	case "BTC":
		return "", nil
	case "ETH":
		txJSONStr := getETHTX(childpubFrom.Pub().Key, childpubTO.Pub().Key, totalAmount)
		return txJSONStr, nil
	default:
		return "", nil
	}
}

func (rpcT *rpcThrift) GetAddr(msg *addrtx.GetAddrMsg) (string, error) {
	coinType := msg.CoinType
	uid := msg.UID

	switch coinType {
	case "BTC":
		childpubUID, _ := rpcT.btcPubKey.Child(uint32(uid))
		addr := genBTCAddr(childpubUID.Pub().Key, false)
		return addr, nil
	case "ETH":
		childpubUID, _ := rpcT.ethPubKey.Child(uint32(uid))
		addr := genETHAddr(childpubUID.Pub().Key)
		return addr, nil
	default:
		return "", errors.New("coin type not supported")
	}
}

func (server *rpcServer) start(ethPubKey, btcPubKey *hdwallet.HDWallet) {
	server.wg.Add(1)
	defer server.wg.Done()
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	serverTransport, err := thrift.NewTServerSocket(server.addr)
	if err != nil {
		log.Fatal(err)
	}

	handler := &rpcThrift{}
	handler.ethPubKey = ethPubKey
	handler.btcPubKey = btcPubKey
	processor := addrtx.NewAddrTXServiceProcessor(handler)

	server.thriftServer = thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	log.Println("thrift server start...")
	server.thriftServer.Serve()
}

func (server *rpcServer) stop() {
	server.thriftServer.Stop()
}
