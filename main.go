package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/whoisnian/glb/util/netutil"
	"github.com/whoisnian/glb/util/osutil"
	"github.com/whoisnian/share-Go/internal/global"
	"github.com/whoisnian/share-Go/internal/router"
)

func main() {
	global.Init()
	if global.CFG.Version {
		fmt.Printf("share-Go %s(%s)\n", global.Version, global.BuildTime)
		return
	}

	predictAddr := global.CFG.HTTPListenAddr
	if host, port, err := net.SplitHostPort(global.CFG.HTTPListenAddr); err == nil && (host == "" || host == "0.0.0.0") {
		if ip, err := netutil.GetOutBoundIP(); err == nil {
			predictAddr = net.JoinHostPort(ip.String(), port)
		}
	}
	global.LOG.Infof("Try visiting http://%s in your browser.", predictAddr)

	server := &http.Server{Addr: global.CFG.HTTPListenAddr, Handler: router.Init()}
	go func() {
		global.LOG.Infof("Service httpd started: <http://%s>", global.CFG.HTTPListenAddr)
		if err := server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			global.LOG.Warn("Service shutting down")
		} else if err != nil {
			global.LOG.Fatal(err.Error())
		}
	}()

	osutil.WaitForInterrupt()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		global.LOG.Warn(err.Error())
	}
}
