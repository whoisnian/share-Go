package router

import (
	"github.com/whoisnian/share-Go/internal/handler"
	"github.com/whoisnian/share-Go/pkg/server"
)

func init() {
	routes := []struct {
		pattern string
		method  string
		handler func(server.Store)
	}{
		{"/upload", "POST", handler.UploadHandler},
		{"/", "GET", handler.IndexHander},
	}

	for _, route := range routes {
		server.Handle(route.pattern, route.method, route.handler)
	}
}
