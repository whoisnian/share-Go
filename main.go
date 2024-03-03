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
	"github.com/whoisnian/share-Go/global"
	"github.com/whoisnian/share-Go/router"
)

func main() {
	global.SetupConfig()
	global.SetupLogger()
	global.LOG.Debugf("use config: %+v", global.CFG)

	if global.CFG.Version {
		fmt.Printf("%s %s(%s)\n", global.AppName, global.Version, global.BuildTime)
		return
	}

	predictAddr := global.CFG.ListenAddr
	if host, port, err := net.SplitHostPort(global.CFG.ListenAddr); err == nil && (host == "" || host == "0.0.0.0") {
		if ip, err := netutil.GetOutBoundIP(); err == nil {
			predictAddr = net.JoinHostPort(ip.String(), port)
		}
	}
	global.LOG.Infof("Try visiting http://%s in your browser.", predictAddr)

	server := &http.Server{Addr: global.CFG.ListenAddr, Handler: router.Setup()}
	go func() {
		global.LOG.Infof("Service httpd started: <http://%s>", global.CFG.ListenAddr)
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
