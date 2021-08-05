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
)

func main() {
	config.Init()
	logger.SetDebug(config.Debug)

	predictAddr := config.HTTPListenAddr
	if host, port, err := net.SplitHostPort(config.HTTPListenAddr); err == nil && host == "0.0.0.0" {
		if ip, err := netutil.GetOutBoundIP(); err == nil {
			predictAddr = net.JoinHostPort(ip.String(), port)
		}
	}
	logger.Info("Try visiting ", ansi.Green, "http://", predictAddr, ansi.Reset, " in your browser.")

	go func() {
		mux := router.Init()
		logger.Info("Service httpd started: <http://", config.HTTPListenAddr, ">")
		if err := http.ListenAndServe(config.HTTPListenAddr, logger.Req(mux)); err != nil {
			logger.Fatal(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
