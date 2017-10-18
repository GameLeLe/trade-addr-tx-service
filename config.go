package main

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

//DigitalAssetsConfig config
type DigitalAssetsConfig struct {
	Title               string      `toml:"title"`
	BTCMasterPubKeyFile string      `toml:"btc_master_pub_key_file"`
	ETHMasterPubKeyFile string      `toml:"eth_master_pub_key_file"`
	RPCConfig           rpcConfig   `toml:"rpc"`
	DBConfig            mysqlConfig `toml:"mysql"`
	RedisConfig         redisConfig `toml:"redis"`
}

type mysqlConfig struct {
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	DBName      string `toml:"db"`
	User        string `toml:"user"`
	Pwd         string `toml:"password"`
	MaxIdleConn int    `toml:"max_idle_conns"`
}

type redisConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
	DB   int    `toml:"db"`
	User string `toml:"user"`
	Pwd  string `toml:"password"`
}

type rpcConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

//ParseConfig parse config file in TOML format
func ParseConfig(filename string) (*DigitalAssetsConfig, error) {
	var config DigitalAssetsConfig
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	_, err = toml.Decode(string(data), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
