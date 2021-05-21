package handler

import (
	"github.com/whoisnian/share-Go/internal/config"
	"github.com/whoisnian/share-Go/pkg/logger"
	"github.com/whoisnian/share-Go/pkg/storage"
	"github.com/whoisnian/share-Go/pkg/tasklane"
)

type jsonMap map[string]interface{}

var fsStore *storage.Store
var downloadTaskLane *tasklane.TaskLane

func Init() {
	var err error
	if fsStore, err = storage.New(config.RootPath); err != nil {
		logger.Fatal(err)
	}
	// runtime.GOMAXPROCS(0)
	downloadTaskLane = tasklane.New(2, 16)
}
