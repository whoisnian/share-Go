package router

import (
	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/share-Go/internal/config"
	"github.com/whoisnian/share-Go/pkg/storage"
	"github.com/whoisnian/share-Go/pkg/tasklane"
	"golang.org/x/net/webdav"
)

type jsonMap map[string]interface{}

var fsStore *storage.Store
var downloadTaskLane *tasklane.TaskLane
var webdavHander *webdav.Handler

func Init() *httpd.Mux {
	var err error
	if fsStore, err = storage.New(config.RootPath); err != nil {
		logger.Fatal(err)
	}

	// runtime.GOMAXPROCS(0)
	downloadTaskLane = tasklane.New(2, 16)

	webdavHander = &webdav.Handler{
		Prefix:     "/webdav",
		FileSystem: webdav.Dir(config.RootPath),
		LockSystem: webdav.NewMemLS(),
	}

	mux := httpd.NewMux()
	mux.Handle("/api/file/*", "GET", FileInfoHandler)
	mux.Handle("/api/file/*", "POST", CreateFileHandler)
	mux.Handle("/api/file/*", "DELETE", DeleteFileHandler)
	mux.Handle("/api/dir/*", "GET", ListDirHandler)
	mux.Handle("/api/dir/*", "POST", CreateDirHandler)
	mux.Handle("/api/dir/*", "DELETE", DeleteDirHandler)

	mux.Handle("/api/raw/*", "GET", RawHandler)
	mux.Handle("/api/download/*", "GET", DownloadHandler)
	mux.Handle("/api/upload/*", "POST", UploadHandler)

	mux.Handle("/webdav/*", "*", httpd.CreateHandler(webdavHander.ServeHTTP))
	mux.Handle("/view/*", "GET", ViewHandler)
	mux.Handle("/*", "GET", IndexHandler)
	return mux
}
