package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/GameLeLe/trade-addr-tx-service/hdwallet"
)

var wg sync.WaitGroup
var daConfig *DigitalAssetsConfig
var daRPCServer *rpcServer
var configFile string
var startFlag bool
var stopFlag bool
var cc chan struct{}

func main() {
	flag.Parse()

	//according to the param stop to kill the pid == running trade deamon
	if stopFlag {
		pid, err := getPID()
		if err != nil {
			log.Fatalln(err)
		} else if pid != 0 {
			syscall.Kill(pid, syscall.SIGINT)
			err = os.Remove(".pid")
			if err != nil {
				log.Fatalln("remove pid file failed!")
			}
			log.Printf("kill {PID:%d} success!", pid)
		}
		return
	}
	if !startFlag {
		log.Fatalln("no start flag found!")
		return
	}

	//init pid file
	err := initPID()
	if err != nil {
		log.Fatalln("init PID:", err)
		return
	}
	//start service
	daConfig, err := ParseConfig(configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	port := daConfig.RPCConfig.Port
	daRPCServer = newRPCServer(port, &wg)
	ethPubKey, _ := hdwallet.ReadWalletFromFile(daConfig.ETHMasterPubKeyFile)
	btcPubKey, _ := hdwallet.ReadWalletFromFile(daConfig.BTCMasterPubKeyFile)
	go daRPCServer.start(ethPubKey, btcPubKey)

	cc = make(chan struct{})
	//listening the signal, Ctrl+C eg.
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	for s := range c {
		if s == syscall.SIGINT {
			daRPCServer.stop()
			close(c)
			close(cc)
		}
	}
	wg.Wait()
}

func getPID() (int, error) {
	fileName := ".pid"
	_, err := os.Stat(fileName)
	if err != nil {
		return 0, err
	}
	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		return 0, err
	}
	line, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(string(line))
	if err != nil {
		return 0, err
	}
	return pid, nil
}

func init() {
	//start or stop
	flag.BoolVar(&startFlag, "start", false, "whether start the program")
	flag.BoolVar(&stopFlag, "stop", false, "whether stop the program")
	flag.StringVar(&configFile, "config", "config.toml", "config file path")
}

func initPID() error {
	var err error
	var f *os.File
	pid := os.Getpid()
	filename := ".pid"
	if checkFileIsExist(filename) {
		f, err = os.OpenFile(filename, os.O_WRONLY, 0666)
	} else {
		f, err = os.Create(filename)
	}
	_, err = io.WriteString(f, strconv.Itoa(pid))
	return err
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
