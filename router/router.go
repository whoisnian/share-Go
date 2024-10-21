package router

import (
	"context"
	"net/http"
	"sync"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/tasklane"
	"github.com/whoisnian/share-Go/global"
)

type jsonMap map[string]any

var (
	lockerMap *sync.Map
	dTaskLane *tasklane.TaskLane
)

func Setup() *httpd.Mux {
	lockerMap = new(sync.Map)
	dTaskLane = tasklane.New(context.Background(), 2, 16)

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

	mux.Handle("/view/*", http.MethodGet, viewHandler)
	mux.Handle("/*", http.MethodGet, indexHandler)
	return mux
}

func readOnlyHander(store *httpd.Store) {
	store.W.WriteHeader(http.StatusForbidden)
	store.W.Write([]byte("Forbidden on read-only mode"))
}
