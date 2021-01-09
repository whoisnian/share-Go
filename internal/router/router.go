package router

import (
	"github.com/whoisnian/share-Go/internal/handler"
	"github.com/whoisnian/share-Go/pkg/httpd"
)

func init() {
	routes := []struct {
		pattern string
		method  string
		handler func(httpd.Store)
	}{
		{"/upload", "POST", handler.UploadHandler},
		{"/", "GET", handler.IndexHander},
	}

	for _, route := range routes {
		httpd.Handle(route.pattern, route.method, route.handler)
	}
}
