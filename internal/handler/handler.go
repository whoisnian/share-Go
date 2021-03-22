package handler

import (
	"github.com/whoisnian/share-Go/internal/config"
	"github.com/whoisnian/share-Go/pkg/storage"
)

type jsonMap map[string]interface{}

var fsStore *storage.Store

func Init() {
	fsStore = storage.New(config.RootPath)
}
