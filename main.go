package main

import (
	"log"
	"net/http"
	"time"

	"github.com/whoisnian/share-Go/pkg/server"
)

func indexHander(store server.Store) {
	url := store.URL()
	if url.Path != "/" {
		store.Respond404()
		return
	}
	time.Sleep(time.Second * 1)
	store.Respond200(nil)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", server.MakeHander(indexHander))

	log.Printf("Server started: <http://127.0.0.1:9000>\n")
	http.ListenAndServe("127.0.0.1:9000", nil)
}
