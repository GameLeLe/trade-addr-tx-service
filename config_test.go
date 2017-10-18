package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDAConfig(t *testing.T) {
	fileStr := `
	title = "digital assets service"
	
	[rpc]
	host = "0.0.0.0"
	port = 8090
	
	[mysql]
	host = "127.0.0.1"
	port = 3306
	db = "dump_test"
	user = "root"
	password = ""
	max_idle_conns = 10
	
	[redis]
	host = "127.0.0.1"
	port = 6379
	user = "root"
	password = ""
	db = 0
	`

	tmpFileName := "./config_tmp.toml"
	ioutil.WriteFile(tmpFileName, []byte(fileStr), 0666)
	defer os.Remove(tmpFileName)
	config, err := ParseConfig(tmpFileName)
	if err != nil {
		t.Errorf("parse config error: %v", err)
	}
	//check title
	assert.Equal(t, "digital assets service", config.Title, "config title not matched")
	//check rpc config
	assert.Equal(t, "0.0.0.0", config.RPCConfig.Host, "rpc host not matched")
	assert.Equal(t, 8090, config.RPCConfig.Port, "rpc port not matched")
	//check mysql config
	assert.Equal(t, "127.0.0.1", config.DBConfig.Host, "mysql host not matched")
	assert.Equal(t, 3306, config.DBConfig.Port, "mysql port not matched")
	assert.Equal(t, "dump_test", config.DBConfig.DBName, "mysql db not matched")
	assert.Equal(t, "root", config.DBConfig.User, "mysql user not matched")
	assert.Equal(t, "", config.DBConfig.Pwd, "mysql password not matched")
	assert.Equal(t, 10, config.DBConfig.MaxIdleConn, "mysql max idle conns not matched")
	//check redis config
	assert.Equal(t, "127.0.0.1", config.RedisConfig.Host, "redis host not matched")
	assert.Equal(t, 6379, config.RedisConfig.Port, "redis port not matched")
	assert.Equal(t, 0, config.RedisConfig.DB, "redis db not matched")
	assert.Equal(t, "root", config.RedisConfig.User, "redis user not matched")
	assert.Equal(t, "", config.RedisConfig.Pwd, "redis password not matched")
}
