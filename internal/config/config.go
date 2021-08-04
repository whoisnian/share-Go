package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/whoisnian/glb/logger"
)

// Debug example: `true`.
var Debug bool

// ReadOnly example: `false`.
var ReadOnly bool

// HTTPListenAddr example: `127.0.0.1:9000`.
var HTTPListenAddr string

// RootPath example: `/srv/share`.
var RootPath string

var configInstance = &struct {
	Debug          *bool
	ReadOnly       *bool
	HTTPListenAddr *string
	RootPath       *string
}{
	&Debug,
	&ReadOnly,
	&HTTPListenAddr,
	&RootPath,
}

var configFilePath = flag.String("c", "config.json", "Specify a path to a custom config file")

// Init load config options from specified json file.
func Init() {
	flag.Parse()

	fi, err := os.Open(*configFilePath)
	if err != nil {
		logger.Fatal(err)
	}
	defer fi.Close()
	err = json.NewDecoder(fi).Decode(configInstance)
	if err != nil {
		logger.Fatal(err)
	}
}
