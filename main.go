package main

import (
	"os"
	"os/signal"

	"github.com/whoisnian/share-Go/internal/config"
	"github.com/whoisnian/share-Go/internal/router"
	"github.com/whoisnian/share-Go/pkg/ftpd"
	"github.com/whoisnian/share-Go/pkg/httpd"
	"github.com/whoisnian/share-Go/pkg/logger"
)

func main() {
	config.Init()
	logger.SetDebug(config.Debug)

	router.Init()
	go httpd.Start(config.HTTPListenAddr)
	go ftpd.Start(config.FTPListenAddr, config.RootPath)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
