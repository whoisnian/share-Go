package router

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/tasklane"
	"github.com/whoisnian/share-Go/internal/global"
	"golang.org/x/net/webdav"
)

var lockerMap *sync.Map
var downloadTaskLane *tasklane.TaskLane

type jsonMap map[string]interface{}

func checkReadOnly(handler httpd.HandlerFunc) httpd.HandlerFunc {
	if global.CFG.ReadOnly {
		return func(store *httpd.Store) { store.W.WriteHeader(http.StatusForbidden) }
	} else {
		return handler
	}
}

func Init() *httpd.Mux {
	lockerMap = new(sync.Map)

	downloadTaskLane = tasklane.New(context.Background(), 2, 16)

	webdavFS := webdav.Dir(global.CFG.RootPath)
	webdavLS := webdav.NewMemLS()
	webdavHander := func(store *httpd.Store) {
		if !global.CFG.ReadOnly ||
			store.R.Method == "PROPFIND" ||
			store.R.Method == "GET" ||
			store.R.Method == "HEAD" ||
			store.R.Method == "OPTIONS" {
			(&webdav.Handler{
				Prefix:     "/webdav",
				FileSystem: webdavFS,
				LockSystem: webdavLS,
				Logger: func(_ *http.Request, err error) {
					global.LOG.Error("webdav.ServeHTTP failed", slog.Any("error", err), slog.String("tid", store.GetID()))
				},
			}).ServeHTTP(store.W, store.R)
		} else {
			store.W.WriteHeader(http.StatusForbidden)
		}
	}

	mux := httpd.NewMux()
	mux.HandleRelay(global.LOG.Relay)
	mux.Handle("/api/file/*", "GET", fileInfoHandler)
	mux.Handle("/api/file/*", "POST", checkReadOnly(newFileHandler))
	mux.Handle("/api/file/*", "DELETE", checkReadOnly(deleteFileHandler))
	mux.Handle("/api/dir/*", "GET", listDirHandler)
	mux.Handle("/api/dir/*", "POST", checkReadOnly(newDirHandler))
	mux.Handle("/api/dir/*", "DELETE", checkReadOnly(deleteDirHandler))

	mux.Handle("/api/raw/*", "GET", rawHandler)
	mux.Handle("/api/download/*", "GET", downloadHandler)
	mux.Handle("/api/rename/*", "POST", checkReadOnly(renameHandler))
	mux.Handle("/api/upload/*", "POST", checkReadOnly(uploadHandler))

	mux.Handle("/webdav/*", "*", webdavHander)
	mux.Handle("/view/*", "GET", viewHandler)
	mux.Handle("/*", "GET", indexHandler)
	return mux
}
