package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/whoisnian/share-Go/internal/config"
	"github.com/whoisnian/share-Go/internal/handler"
	"github.com/whoisnian/share-Go/pkg/server"
)

// CONFIG ...
var CONFIG *config.Config
var configFilePath = flag.String("config", "config.json", "Specify a path to a custom config file")

func main() {
	flag.Parse()
	CONFIG = config.Load(*configFilePath)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/upload", server.MakeHander(handler.UploadHandler))
	http.Handle("/", server.MakeHander(handler.IndexHander))

	log.Printf("Server started: <http://%s>\n", CONFIG.ListenAddr)
	http.ListenAndServe(CONFIG.ListenAddr, nil)
}
