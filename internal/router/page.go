package router

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/logger"
	fe "github.com/whoisnian/share-Go/fe/dist"
)

func serveFileFromFE(store *httpd.Store, path string) {
	file, err := fe.FS.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			store.W.WriteHeader(http.StatusNotFound)
			return
		}
		logger.Panic(err)
	}
	defer file.Close()

	if info, err := file.Stat(); err != nil || info.IsDir() {
		store.W.WriteHeader(http.StatusForbidden)
		return
	}

	ctype := mime.TypeByExtension(filepath.Ext(path))
	if ctype == "" {
		ctype = "application/octet-stream"
	} else if strings.Contains(ctype, "text/css") || strings.Contains(ctype, "application/javascript") {
		// nginx expires max
		// https://nginx.org/en/docs/http/ngx_http_headers_module.html#expires
		store.W.Header().Add("cache-control", "max-age:315360000, public")
		store.W.Header().Add("expires", "Thu, 31 Dec 2037 23:55:55 GMT")
	}
	store.W.Header().Add("content-type", ctype)

	if _, err := io.Copy(store.W, file); err != nil {
		logger.Panic(err)
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
