package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const configFile = "wallet-config.json"

var conf Config

type Config struct {
	Option      WalletOption `json:"wallet_option"`
	Version     string       `json:"version"`
	CultureInfo string       `json:"culture_info"`
	Theme       string       `json:"theme"`
}

type WalletOption struct {
	RpcEnabled           bool    `json:"rpc_enabled"`
	RpcPort              *int    `json:"rpc_port"`
	SpecifiedPoolAddress *string `json:"specified_pool_address"`
	IsTestNet            bool    `json:"is_test_net"`
	DisableMining        bool    `json:"disable_mining"`
	PoolAddress          string  `json:"pool_address"`
}

func initConfig() {
	newConf := false

	conf.Version = "0.4.0"
	conf.Theme = "Dark"
	conf.CultureInfo = "en-US"
	conf.Option.DisableMining = true
	conf.Option.PoolAddress = "equal.xdag.org:13656"

	data, err := ioutil.ReadFile(configFile)
	if err == nil {
		err = json.Unmarshal(data, &conf)
		if err != nil {
			log.Println("fail to Unmarshal configure,conf.json")
			newConf = true
		}
	} else {
		newConf = true
	}
	if conf.CultureInfo == "" {
		conf.CultureInfo = "en-US"
	}

	if newConf {
		data, _ = json.MarshalIndent(conf, "", "  ")
		ioutil.WriteFile(configFile, data, 666)
	}
}

func SaveConfig() error {
	data, _ := json.MarshalIndent(conf, "", "  ")
	ioutil.WriteFile(configFile, data, 666)
	return nil
}

func GetConfig() *Config {
	return &conf
}
