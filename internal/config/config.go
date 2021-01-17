package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	"github.com/whoisnian/share-Go/pkg/logger"
)

// Debug example: `true`.
var Debug bool

// HTTPListenAddr example: `127.0.0.1:9000`.
var HTTPListenAddr string

// FTPListenAddr example: `127.0.0.1:2121`.
var FTPListenAddr string

// RootPath example: `/srv/share`.
var RootPath string

var configInstance = &struct {
	Debug          *bool
	HTTPListenAddr *string
	FTPListenAddr  *string
	RootPath       *string
}{
	&Debug,
	&HTTPListenAddr,
	&FTPListenAddr,
	&RootPath,
}

var configFilePath = flag.String("c", "config.json", "Specify a path to a custom config file")

// Init load config options from specified json file.
func Init() {
	flag.Parse()
	loadFromJSON()
}

func loadFromJSON() {
	fi, err := os.Open(*configFilePath)
	if err != nil {
		logger.Fatal(err)
	}
	defer fi.Close()

	content, err := ioutil.ReadAll(fi)
	if err != nil {
		logger.Fatal(err)
	}

	err = json.Unmarshal(content, configInstance)
	if err != nil {
		logger.Fatal(err)
	}
}

func saveAsJSON() {
	fi, err := os.OpenFile(*configFilePath, os.O_WRONLY, 0)
	if err != nil {
		logger.Panic(err)
	}
	defer fi.Close()

	content, err := json.MarshalIndent(configInstance, "", "  ")
	if err != nil {
		logger.Panic(err)
	}

	_, err = fi.Write(content)
	if err != nil {
		logger.Panic(err)
	}
}
