package handler

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	fe "github.com/whoisnian/share-Go/fe/dist"
	"github.com/whoisnian/share-Go/pkg/httpd"
	"github.com/whoisnian/share-Go/pkg/logger"
)

func serveFileFromFE(store httpd.Store, path string) {
	file, err := fe.FS.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			store.WriteHeader(http.StatusNotFound)
			return
		}
		logger.Panic(err)
	}
	defer file.Close()

	if fileInfo, _ := file.Stat(); fileInfo.IsDir() {
		store.WriteHeader(http.StatusForbidden)
		return
	}

	ctype := mime.TypeByExtension(filepath.Ext(path))
	if ctype == "" {
		ctype = "application/octet-stream"
	} else if strings.Contains(ctype, "text/css") || strings.Contains(ctype, "application/javascript") {
		// nginx expires max
		// https://nginx.org/en/docs/http/ngx_http_headers_module.html#expires
		store.ResponseHeader().Add("cache-control", "max-age:315360000, public")
		store.ResponseHeader().Add("expires", "Thu, 31 Dec 2037 23:55:55 GMT")
	}
	store.ResponseHeader().Add("content-type", ctype)

	writer := store.Writer()
	if _, err := io.Copy(writer, file); err != nil {
		logger.Panic(err)
	}
}

func ViewHander(store httpd.Store) {
	serveFileFromFE(store, "static/index.html")
}

func IndexHandler(store httpd.Store) {
	path := store.RouteAny()
	if path == "" {
		store.Redirect("/view/", http.StatusFound)
		return
	}

	serveFileFromFE(store, path)
}
