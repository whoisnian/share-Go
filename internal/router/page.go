package router

import (
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/whoisnian/glb/httpd"
	fe "github.com/whoisnian/share-Go/fe/dist"
	"github.com/whoisnian/share-Go/internal/global"
)

func serveFileFromFE(store *httpd.Store, path string) {
	file, err := fe.FS.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			store.W.WriteHeader(http.StatusNotFound)
			return
		}
		global.LOG.Panic("serveFileFromFE failed", slog.Any("error", err), slog.String("tid", store.GetID()))
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil || info.IsDir() {
		store.W.WriteHeader(http.StatusForbidden)
		return
	}

	ctype := mime.TypeByExtension(filepath.Ext(path))
	if ctype == "" {
		ctype = "application/octet-stream"
	} else if strings.Contains(ctype, "text/css") || strings.Contains(ctype, "application/javascript") {
		// nginx expires max
		// https://nginx.org/en/docs/http/ngx_http_headers_module.html#expires
		store.W.Header().Set("cache-control", "max-age:315360000, public")
		store.W.Header().Set("expires", "Thu, 31 Dec 2037 23:55:55 GMT")
	}
	store.W.Header().Set("content-type", ctype)

	if store.W.Header().Get("content-encoding") == "" {
		store.W.Header().Set("content-length", strconv.FormatInt(info.Size(), 10))
	}
	if _, err := io.CopyN(store.W, file, info.Size()); err != nil {
		global.LOG.Panic("io.CopyN failed", slog.Any("error", err), slog.String("tid", store.GetID()))
	}
}

func viewHandler(store *httpd.Store) {
	serveFileFromFE(store, "static/index.html")
}

func indexHandler(store *httpd.Store) {
	path := store.RouteParamAny()
	if path == "" {
		store.Redirect("/view/", http.StatusFound)
		return
	}

	serveFileFromFE(store, path)
}
