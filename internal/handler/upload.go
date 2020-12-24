package handler

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/whoisnian/share-Go/pkg/server"
)

// UploadHandler ...
func UploadHandler(store server.Store) {
	reader, err := store.MultipartReader()
	if err != nil {
		log.Panic(err)
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}
		if part.FormName() != "fileList" {
			continue
		}

		file, err := os.Create(path.Join("./uploads", part.FileName()))
		if err != nil {
			log.Panic(err)
		}
		defer file.Close()

		io.Copy(file, part)
	}
	store.Respond200(nil)
}
