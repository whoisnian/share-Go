package router

import (
	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/util/fsutil"
	"github.com/whoisnian/share-Go/internal/config"
	"github.com/whoisnian/share-Go/pkg/tasklane"
	"golang.org/x/net/webdav"
)

type jsonMap map[string]interface{}

var lockedFS *fsutil.LockedFS
var downloadTaskLane *tasklane.TaskLane

func Init() *httpd.Mux {
	lockedFS = fsutil.NewLockedFS()

	webdavHander := httpd.CreateHandler((&webdav.Handler{
		Prefix:     "/webdav",
		FileSystem: webdav.Dir(config.RootPath),
		LockSystem: webdav.NewMemLS(),
	}).ServeHTTP)

	mux := httpd.NewMux()
	mux.Handle("/api/file/*", "GET", FileInfoHandler)
	mux.Handle("/api/file/*", "POST", NewFileHandler)
	mux.Handle("/api/file/*", "DELETE", DeleteFileHandler)
	mux.Handle("/api/dir/*", "GET", ListDirHandler)
	mux.Handle("/api/dir/*", "POST", NewDirHandler)
	mux.Handle("/api/dir/*", "DELETE", DeleteDirHandler)

	mux.Handle("/api/raw/*", "GET", RawHandler)
	mux.Handle("/api/download/*", "GET", DownloadHandler)
	mux.Handle("/api/upload/*", "POST", UploadHandler)

	mux.Handle("/webdav/*", "*", webdavHander)
	mux.Handle("/view/*", "GET", ViewHandler)
	mux.Handle("/*", "GET", IndexHandler)
	return mux
}
