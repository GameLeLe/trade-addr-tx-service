package main

import (
	"log"
	"strconv"

	"git.apache.org/thrift.git/lib/go/thrift"
	bip39 "github.com/GameLeLe/trade-addr-tx-service/bip39"
	hdwallet "github.com/GameLeLe/trade-addr-tx-service/hdwallet"
	addrtx "github.com/GameLeLe/trade-addr-tx-service/thrift/addrtx"
)

type rpcServer struct {
	addr         string
	thriftServer *thrift.TSimpleServer
}

func newRPCServer(port int) *rpcServer {
	server := &rpcServer{}
	server.addr = "0.0.0.0:" + strconv.Itoa(port)
	return server
}

type rpcThrift struct {
}

func (rpcT *rpcThrift) GetAddr(msg *addrtx.GetAddrMsg) (string, error) {
	coinType := msg.CoinType
	uid := msg.UID
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
	childpubUID, _ := childpub0.Child(uint32(uid))
	switch coinType {
	case "BTC":
		addr := genBTCAddr(childpubUID.Pub().Key, false)
		return addr, nil
	case "ETH":
		addr := genETHAddr(childpubUID.Pub().Key)
		return addr, nil
	default:
		return "", nil
	}
}

func (server *rpcServer) start() {
	// wg.Add(1)
	// defer wg.Add(-1)
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	serverTransport, err := thrift.NewTServerSocket(server.addr)
	if err != nil {
		log.Fatal(err)
	}

	handler := &rpcThrift{}
	processor := addrtx.NewAddrTXServiceProcessor(handler)

	server.thriftServer = thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	log.Println("thrift server start...")
	server.thriftServer.Serve()
}

func (server *rpcServer) stop() {
	server.thriftServer.Stop()
}
