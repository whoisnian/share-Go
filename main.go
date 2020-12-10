package main

import (
	logger "log"
	"net/http"
	"time"

	"github.com/whoisnian/share-Go/pkg/log"
	"github.com/whoisnian/share-Go/pkg/state"
)

func indexHander(store state.Store) {
	time.Sleep(time.Second * 2)
	store.Respond200(nil)
}

func main() {
	http.Handle("/", log.MakeHander(indexHander))

	logger.Printf("Server started: <http://127.0.0.1:9000>\n")
	http.ListenAndServe("127.0.0.1:9000", nil)
}
