package main

import (
	"os"
	"os/signal"

	"github.com/whoisnian/share-Go/internal/config"
	_ "github.com/whoisnian/share-Go/internal/router"
	"github.com/whoisnian/share-Go/pkg/ftpd"
	"github.com/whoisnian/share-Go/pkg/httpd"
)

func main() {
	go ftpd.Start(config.FTPListenAddr)
	go httpd.Start(config.ListenAddr)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
