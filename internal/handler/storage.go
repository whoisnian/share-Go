package handler

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"syscall"

	"github.com/whoisnian/share-Go/pkg/httpd"
	"github.com/whoisnian/share-Go/pkg/logger"
)

type fileType int

const (
	typeRegular fileType = iota
	typeDirectory
)

type fileInfo struct {
	Type  fileType
	Name  string
	Size  int64
	MTime int64
}

func parseFileInfo(info os.FileInfo) fileInfo {
	t := typeRegular
	if info.Mode().IsDir() {
		t = typeDirectory
	}
	return fileInfo{
		Type:  t,
		Name:  info.Name(),
		Size:  info.Size(),
		MTime: info.ModTime().Unix(),
	}
}

func FileInfoHandler(store httpd.Store) {
	path := filepath.Join("/", store.RouteAny())
	info, err := fsStore.FileInfo(path)
	if err != nil {
		if os.IsNotExist(err) {
			store.Respond404()
			return
		}
		logger.Panic(err)
	}

	store.RespondJson(parseFileInfo(info))
}

func CreateFileHandler(store httpd.Store) {
	path := filepath.Join("/", store.RouteAny())
	file, err := fsStore.CreateFile(path)
	if err != nil {
		logger.Panic(err)
	}
	defer file.Close()

	bodyReader := store.Body()
	defer bodyReader.Close()

	io.Copy(file, bodyReader)

	store.Respond200(nil)
}

func DeleteFileHandler(store httpd.Store) {
	path := filepath.Join("/", store.RouteAny())
	if err := fsStore.Delete(path); err != nil {
		if errors.Is(err, syscall.ENOTEMPTY) {
			store.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Panic(err)
	}
	store.Respond200(nil)
}

func ListDirHandler(store httpd.Store) {
	path := filepath.Join("/", store.RouteAny())
	infos, err := fsStore.ListDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			store.Respond404()
			return
		} else if errors.Is(err, syscall.ENOTDIR) {
			store.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Panic(err)
	}

	result := make([]fileInfo, len(infos))
	for i := 0; i < len(infos); i++ {
		result[i] = parseFileInfo(infos[i])
	}

	store.RespondJson(jsonMap{
		"FileInfos": result,
	})
}

func CreateDirHandler(store httpd.Store) {
	path := filepath.Join("/", store.RouteAny())
	if err := fsStore.CreateDir(path); err != nil {
		logger.Panic(err)
	}
	store.Respond200(nil)
}

func DeleteDirHandler(store httpd.Store) {
	path := filepath.Join("/", store.RouteAny())
	if err := fsStore.DeleteAll(path); err != nil {
		logger.Panic(err)
	}
	store.Respond200(nil)
}
