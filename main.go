package main

import (
	"log"
	"net/http"

	"github.com/whoisnian/share-Go/pkg/server"
)

func indexHander(store server.Store) {
	page := `<!DOCTYPE html>
<html>
	<head>
	  <meta charset="utf-8">
		<title>Index Page</title>
	</head>
	<body>
	  <h2>Upload a file</h2>
	  <form action="/upload" method="post" enctype="multipart/form-data">
		  <input type="file" name="file" multiple>
		  <br>
		  <input type="submit" name="submit" value="Submit">
	  </form>
	</body>
</html>`
	store.Respond200([]byte(page))
}

// https://golang.org/pkg/mime/multipart/
func uploadHandler(store server.Store) {
	store.Save("uploads")
	store.Respond200(nil)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/upload", server.MakeHander(uploadHandler))
	http.Handle("/", server.MakeHander(indexHander))

	log.Printf("Server started: <http://127.0.0.1:9000>\n")
	http.ListenAndServe("127.0.0.1:9000", nil)
}
