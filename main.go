package main

import (
	"github.com/whoisnian/share-Go/internal/config"
	_ "github.com/whoisnian/share-Go/internal/router"
	"github.com/whoisnian/share-Go/pkg/server"
)

func main() {
	server.Start(config.ListenAddr)
}
