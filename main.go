package main

import (
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/whoisnian/glb/ansi"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/glb/util/netutil"
	"github.com/whoisnian/share-Go/internal/config"
	"github.com/whoisnian/share-Go/internal/router"
	"github.com/whoisnian/share-Go/pkg/ftpd"
)

func main() {
	config.Init()
	logger.SetDebug(config.Debug)

	if _, port, err := net.SplitHostPort(config.HTTPListenAddr); err == nil {
		if ip, err := netutil.GetOutBoundIP(); err == nil {
			logger.Info("Try visiting ", ansi.Green, "http://", net.JoinHostPort(ip.String(), port), ansi.Reset, " in your browser.")
		}
	}

	go func() {
		mux := router.Init()
		logger.Info("Service httpd started: <http://", config.HTTPListenAddr, ">")
		if err := http.ListenAndServe(config.HTTPListenAddr, logger.Req(mux)); err != nil {
			logger.Fatal(err)
		}
	}()

	go func() {
		logger.Info("Service ftpd started: <ftp://", config.FTPListenAddr, ">")
		ftpd.Start(config.FTPListenAddr, config.RootPath)
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
