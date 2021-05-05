package router

import (
	"github.com/whoisnian/share-Go/internal/handler"
	"github.com/whoisnian/share-Go/pkg/httpd"
)

func Init() {
	handler.Init()

	routes := []struct {
		pattern string
		method  string
		handler func(httpd.Store)
	}{
		{"/api/file/*", "GET", handler.FileInfoHandler},
		{"/api/file/*", "POST", handler.CreateFileHandler},
		{"/api/file/*", "DELETE", handler.DeleteFileHandler},
		{"/api/dir/*", "GET", handler.ListDirHandler},
		{"/api/dir/*", "POST", handler.CreateDirHandler},
		{"/api/dir/*", "DELETE", handler.DeleteDirHandler},
		{"/download/*", "GET", handler.DownloadHandler},

		{"/upload", "POST", handler.UploadHandler}, // TODO
		{"/", "GET", handler.IndexHander},          // TODO
	}

	for _, route := range routes {
		httpd.Handle(route.pattern, route.method, route.handler)
	}
}
