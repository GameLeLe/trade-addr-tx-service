package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup
var daConfig *DigitalAssetsConfig

func main() {
	daConfig, err := ParseConfig("./config.toml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(daConfig.Title)
}
