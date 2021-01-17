package handler

import (
	"io"
	"os"
	"path"

	"github.com/whoisnian/share-Go/pkg/httpd"
	"github.com/whoisnian/share-Go/pkg/logger"
)

// UploadHandler saves received files to specified directory.
func UploadHandler(store httpd.Store) {
	reader, err := store.MultipartReader()
	if err != nil {
		logger.Panic(err)
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Panic(err)
		}
		if part.FormName() != "fileList" {
			continue
		}

		file, err := os.Create(path.Join("./uploads", part.FileName()))
		if err != nil {
			logger.Panic(err)
		}
		defer file.Close()

		io.Copy(file, part)
	}
	store.Respond200(nil)
}
