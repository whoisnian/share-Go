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

		{"/api/raw/*", "GET", handler.RawHandler},
		{"/api/download/*", "GET", handler.DownloadHandler},
		{"/api/upload/*", "POST", handler.UploadHandler},

		{"/view/*", "GET", handler.ViewHander},
		{"/*", "GET", handler.IndexHandler},
	}

	for _, route := range routes {
		httpd.Handle(route.pattern, route.method, route.handler)
	}
}
