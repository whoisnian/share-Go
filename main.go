package main

import (
	"flag"

	"github.com/whoisnian/share-Go/internal/config"
	_ "github.com/whoisnian/share-Go/internal/router"
	"github.com/whoisnian/share-Go/pkg/server"
)

// CONFIG ...
var CONFIG *config.Config
var configFilePath = flag.String("config", "config.json", "Specify a path to a custom config file")

func main() {
	flag.Parse()
	CONFIG = config.Load(*configFilePath)

	server.Start(CONFIG.ListenAddr)
}
