package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

// config file structure
type ConfigT struct {
	Client ClientConfig
}

type ClientConfig struct {
	Servers []string
}

// end config file

// global config
var Conf ConfigT

func ReadConf() {
	// read config
	if _, err := toml.DecodeFile("pafimi.conf", &Conf); err != nil {
		// handle error
		log.Print("error in reading pafimi.conf:")
		log.Fatal(err)
	}
}

// Request RPC Argument type
type Request struct {
	User string
	Src  string
	Dst  string
}
