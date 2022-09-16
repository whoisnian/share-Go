package main

import (
	"net"
	"net/http"

	"github.com/whoisnian/glb/ansi"
	"github.com/whoisnian/glb/config"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/glb/util/netutil"
	"github.com/whoisnian/glb/util/osutil"
	"github.com/whoisnian/share-Go/internal/global"
	"github.com/whoisnian/share-Go/internal/router"
)

func main() {
	err := config.FromCommandLine(&global.CFG)
	if err != nil {
		logger.Fatal(err)
	}
	logger.SetDebug(global.CFG.Debug)
	logger.SetColorful(global.CFG.Debug)

	predictAddr := global.CFG.HTTPListenAddr
	if host, port, err := net.SplitHostPort(global.CFG.HTTPListenAddr); err == nil && host == "0.0.0.0" {
		if ip, err := netutil.GetOutBoundIP(); err == nil {
			predictAddr = net.JoinHostPort(ip.String(), port)
		}
	}
	logger.Info("Try visiting ", ansi.Green, "http://", predictAddr, ansi.Reset, " in your browser.")

	go func() {
		mux := router.Init()
		logger.Info("Service httpd started: <http://", global.CFG.HTTPListenAddr, ">")
		if err := http.ListenAndServe(global.CFG.HTTPListenAddr, logger.Req(logger.Recovery(mux))); err != nil {
			logger.Fatal(err)
		}
	}()

	osutil.WaitForInterrupt()
}
