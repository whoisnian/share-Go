package handler

import (
	"net/http"
	"strings"

	"github.com/whoisnian/share-Go/internal/config"
	"github.com/whoisnian/share-Go/pkg/logger"
	"github.com/whoisnian/share-Go/pkg/storage"
	"github.com/whoisnian/share-Go/pkg/tasklane"
	"golang.org/x/net/webdav"
)

type jsonMap map[string]interface{}

var fsStore *storage.Store
var downloadTaskLane *tasklane.TaskLane
var webdavHander *webdav.Handler

func Init() {
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
		Logger: func(r *http.Request, err error) {
			logger.Info(
				r.RemoteAddr[0:strings.IndexByte(r.RemoteAddr, ':')],
				" [webdav] ",
				r.Method, " ",
				r.URL.Path, " ",
				r.UserAgent(), " ",
				err,
			)
		},
	}
}
