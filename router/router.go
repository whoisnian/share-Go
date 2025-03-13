package router

import (
	"context"
	"net/http"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/share-Go/global"
	"github.com/whoisnian/share-Go/router/api"
	"github.com/whoisnian/share-Go/router/page"
)

func Setup(ctx context.Context) *httpd.Mux {
	api.Setup(ctx)

	mux := httpd.NewMux()
	mux.HandleMiddleware(global.LOG.NewMiddleware())
	mux.HandleMiddleware(readOnlyMiddleware(global.CFG.ReadOnly))

	mux.Handle("/api/file/*", http.MethodGet, api.FileInfoHandler)
	mux.Handle("/api/file/*", http.MethodPost, api.NewFileHandler)
	mux.Handle("/api/file/*", http.MethodDelete, api.DeleteFileHandler)
	mux.Handle("/api/dir/*", http.MethodGet, api.ListDirHandler)
	mux.Handle("/api/dir/*", http.MethodPost, api.NewDirHandler)
	mux.Handle("/api/dir/*", http.MethodDelete, api.DeleteDirHandler)

	mux.Handle("/api/raw/*", http.MethodGet, api.RawHandler)
	mux.Handle("/api/download/*", http.MethodGet, api.DownloadHandler)
	mux.Handle("/api/rename/*", http.MethodPost, api.RenameHandler)
	mux.Handle("/api/upload/*", http.MethodPost, api.UploadHandler)

	mux.Handle("/view/*", http.MethodGet, page.ViewHandler)
	mux.Handle("/*", http.MethodGet, page.IndexHandler)
	return mux
}

func readOnlyMiddleware(enabled bool) httpd.HandlerFunc {
	if enabled {
		return func(store *httpd.Store) {
			if store.R.Method != http.MethodGet {
				store.W.WriteHeader(http.StatusForbidden)
				store.W.Write([]byte("forbidden on read-only mode"))
			} else {
				store.Next()
			}
		}
	} else {
		return func(store *httpd.Store) { store.Next() }
	}
}
