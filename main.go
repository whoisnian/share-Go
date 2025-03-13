package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/glb/util/netutil"
	"github.com/whoisnian/glb/util/osutil"
	"github.com/whoisnian/share-Go/global"
	"github.com/whoisnian/share-Go/router"
)

func main() {
	ctx := context.Background()
	global.SetupConfig(ctx)
	global.SetupLogger(ctx)
	global.LOG.Debugf(ctx, "use config: %+v", global.CFG)

	if global.CFG.Version {
		fmt.Printf("%s version %s built with %s at %s\n", global.AppName, global.Version, runtime.Version(), global.BuildTime)
		return
	}

	predictScheme := "http"
	if global.CFG.TlsCert != "" && global.CFG.TlsKey != "" {
		predictScheme = "https"
	}
	predictAddr := global.CFG.ListenAddr
	if host, port, err := net.SplitHostPort(global.CFG.ListenAddr); err == nil && (host == "" || host == "0.0.0.0") {
		if ip, err := netutil.GetOutBoundIP(); err == nil {
			predictAddr = net.JoinHostPort(ip.String(), port)
		}
	}
	global.LOG.Infof(ctx, "try visiting %s://%s in your browser.", predictScheme, predictAddr)

	server := &http.Server{Addr: global.CFG.ListenAddr, Handler: router.Setup(ctx)}
	go func() {
		var serverErr error
		if global.CFG.TlsCert != "" && global.CFG.TlsKey != "" {
			global.LOG.Infof(ctx, "service httpd started: https://%s", global.CFG.ListenAddr)
			serverErr = server.ListenAndServeTLS(global.CFG.TlsCert, global.CFG.TlsKey)
		} else {
			global.LOG.Infof(ctx, "service httpd started: http://%s", global.CFG.ListenAddr)
			serverErr = server.ListenAndServe()
		}
		if errors.Is(serverErr, http.ErrServerClosed) {
			global.LOG.Warn(ctx, "service shutting down")
		} else if serverErr != nil {
			global.LOG.Fatal(ctx, "service start", logger.Error(serverErr))
		}
	}()

	osutil.WaitForStop()

	shutdownCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		global.LOG.Warn(ctx, "service stop", logger.Error(err))
	}
}
