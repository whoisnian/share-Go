package main

import (
	"net"
	"os"
	"os/signal"

	"github.com/whoisnian/share-Go/internal/config"
	"github.com/whoisnian/share-Go/internal/router"
	"github.com/whoisnian/share-Go/pkg/ftpd"
	"github.com/whoisnian/share-Go/pkg/httpd"
	"github.com/whoisnian/share-Go/pkg/logger"
	"github.com/whoisnian/share-Go/pkg/util"
)

func main() {
	config.Init()
	logger.SetDebug(config.Debug)

	if _, port, err := net.SplitHostPort(config.HTTPListenAddr); err == nil {
		if ip, err := util.GetOutBoundIP(); err == nil {
			logger.Info("Try visiting \x1b[32mhttp://", net.JoinHostPort(ip.String(), port), "\x1b[0m in your browser.")
		}
	}

	router.Init()
	go httpd.Start(config.HTTPListenAddr)
	go ftpd.Start(config.FTPListenAddr, config.RootPath)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
