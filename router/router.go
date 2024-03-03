package router

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/tasklane"
	"github.com/whoisnian/glb/util/strutil"
	"github.com/whoisnian/share-Go/global"
	"golang.org/x/net/webdav"
)

type jsonMap map[string]any

var (
	lockerMap *sync.Map
	dTaskLane *tasklane.TaskLane
	webdavFS  webdav.FileSystem
	webdavLS  webdav.LockSystem
)

func Setup() *httpd.Mux {
	lockerMap = new(sync.Map)
	dTaskLane = tasklane.New(context.Background(), 2, 16)
	webdavFS = webdav.Dir(global.CFG.RootPath)
	webdavLS = webdav.NewMemLS()

	mux := httpd.NewMux()
	mux.HandleRelay(global.LOG.Relay)
	muxHandleCheck := func(p string, m string, h httpd.HandlerFunc, readOnlyCheck bool) {
		if readOnlyCheck && global.CFG.ReadOnly {
			mux.Handle(p, m, readOnlyHander)
			return
		}
		mux.Handle(p, m, h)
	}

	muxHandleCheck("/api/file/*", http.MethodGet, fileInfoHandler, false)
	muxHandleCheck("/api/file/*", http.MethodPost, newFileHandler, true)
	muxHandleCheck("/api/file/*", http.MethodDelete, deleteFileHandler, true)
	muxHandleCheck("/api/dir/*", http.MethodGet, listDirHandler, false)
	muxHandleCheck("/api/dir/*", http.MethodPost, newDirHandler, true)
	muxHandleCheck("/api/dir/*", http.MethodDelete, deleteDirHandler, true)

	muxHandleCheck("/api/raw/*", http.MethodGet, rawHandler, false)
	muxHandleCheck("/api/download/*", http.MethodGet, downloadHandler, false)
	muxHandleCheck("/api/rename/*", http.MethodPost, renameHandler, true)
	muxHandleCheck("/api/upload/*", http.MethodPost, uploadHandler, true)

	mux.Handle("/webdav/*", httpd.MethodAll, webdavHander)
	mux.Handle("/view/*", http.MethodGet, viewHandler)
	mux.Handle("/*", http.MethodGet, indexHandler)
	return mux
}

func webdavHander(store *httpd.Store) {
	if global.CFG.ReadOnly && !strutil.SliceContain([]string{"GET", "HEAD", "OPTIONS", "PROPFIND"}, store.R.Method) {
		readOnlyHander(store)
		return
	}
	(&webdav.Handler{
		Prefix:     "/webdav",
		FileSystem: webdavFS,
		LockSystem: webdavLS,
		Logger: func(_ *http.Request, err error) {
			global.LOG.Error("webdav.ServeHTTP failed", slog.Any("error", err), slog.String("tid", store.GetID()))
		},
	}).ServeHTTP(store.W, store.R)
}

func readOnlyHander(store *httpd.Store) {
	store.W.WriteHeader(http.StatusForbidden)
	store.W.Write([]byte("Forbidden on read-only mode"))
}
