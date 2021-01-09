package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

// HTTPListenAddr example: `127.0.0.1:9000`.
var HTTPListenAddr string

// FTPListenAddr example: `127.0.0.1:2121`.
var FTPListenAddr string

var configInstance = &struct {
	HTTPListenAddr *string
	FTPListenAddr  *string
}{
	&HTTPListenAddr,
	&FTPListenAddr,
}

var configFilePath = flag.String("c", "config.json", "Specify a path to a custom config file")

func init() {
	flag.Parse()
	loadFromJSON()
}

func loadFromJSON() {
	fi, err := os.Open(*configFilePath)
	if err != nil {
		log.Panicln(err)
	}
	defer fi.Close()

	content, err := ioutil.ReadAll(fi)
	if err != nil {
		log.Panicln(err)
	}

	err = json.Unmarshal(content, configInstance)
	if err != nil {
		log.Panicln(err)
	}
}

func saveAsJSON() {
	fi, err := os.OpenFile(*configFilePath, os.O_WRONLY, 0)
	if err != nil {
		log.Panicln(err)
	}
	defer fi.Close()

	content, err := json.MarshalIndent(configInstance, "", "  ")
	if err != nil {
		log.Panicln(err)
	}

	_, err = fi.Write(content)
	if err != nil {
		log.Panicln(err)
	}
}
