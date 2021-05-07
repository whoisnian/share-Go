package handler

import (
	"github.com/whoisnian/share-Go/internal/config"
	"github.com/whoisnian/share-Go/pkg/logger"
	"github.com/whoisnian/share-Go/pkg/storage"
)

type jsonMap map[string]interface{}

var fsStore *storage.Store

func Init() {
	var err error
	if fsStore, err = storage.New(config.RootPath); err != nil {
		logger.Fatal(err)
	}
}
