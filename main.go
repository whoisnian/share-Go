package main

import (
	"github.com/whoisnian/share-Go/internal/config"
	_ "github.com/whoisnian/share-Go/internal/router"
	"github.com/whoisnian/share-Go/pkg/ftp"
	"github.com/whoisnian/share-Go/pkg/server"
)

func main() {
	go ftp.Start(config.FTPListenAddr)
	server.Start(config.ListenAddr)
}
