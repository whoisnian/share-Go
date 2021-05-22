package handler

import (
	"net/http"

	"github.com/whoisnian/share-Go/pkg/httpd"
)

func IndexHandler(store httpd.Store) {
	store.Redirect("/view/", http.StatusFound)
}

func ViewHander(store httpd.Store) {
	store.Respond200([]byte("TODO"))
}
