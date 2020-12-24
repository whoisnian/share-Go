package main

import (
	"log"
	"net/http"

	"github.com/whoisnian/share-Go/internal/handler"
	"github.com/whoisnian/share-Go/pkg/server"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/upload", server.MakeHander(handler.UploadHandler))
	http.Handle("/", server.MakeHander(handler.IndexHander))

	log.Printf("Server started: <http://127.0.0.1:9000>\n")
	http.ListenAndServe("127.0.0.1:9000", nil)
}
