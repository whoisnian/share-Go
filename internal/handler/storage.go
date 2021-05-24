package handler

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
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
			store.WriteHeader(http.StatusNotFound)
			store.RespondJson(nil)
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
			store.WriteHeader(http.StatusNotFound)
			store.RespondJson(nil)
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

func RawHandler(store httpd.Store) {
	path := filepath.Join("/", store.RouteAny())
	info, err := fsStore.FileInfo(path)
	if err != nil {
		if os.IsNotExist(err) {
			store.WriteHeader(http.StatusNotFound)
			store.RespondJson(nil)
			return
		}
		logger.Panic(err)
	}

	if info.Mode().IsRegular() {
		file, err := fsStore.GetFile(path)
		if err != nil {
			logger.Panic(err)
		}
		defer file.Close()

		writer := store.Writer()
		if _, err := io.Copy(writer, file); err != nil {
			logger.Panic(err)
		}
		return
	} else {
		store.WriteHeader(http.StatusUnprocessableEntity)
		store.RespondJson(nil)
		return
	}
}

func DownloadHandler(store httpd.Store) {
	path := filepath.Join("/", store.RouteAny())
	info, err := fsStore.FileInfo(path)
	if err != nil {
		if os.IsNotExist(err) {
			store.WriteHeader(http.StatusNotFound)
			store.RespondJson(nil)
			return
		}
		logger.Panic(err)
	}

	if info.Mode().IsRegular() {
		file, err := fsStore.GetFile(path)
		if err != nil {
			logger.Panic(err)
		}
		defer file.Close()

		filename := url.PathEscape(filepath.Base(path))
		store.ResponseHeader().Add("content-disposition", "attachment; filename*=UTF-8''"+filename+"; filename=\""+filename+"\"")

		writer := store.Writer()
		if _, err := io.Copy(writer, file); err != nil {
			store.ResponseHeader().Del("content-disposition")
			logger.Panic(err)
		}
		return
	} else if info.Mode().IsDir() {
		filename := filepath.Base(path)
		if filename == "/" {
			filename = "root"
		}
		filename = url.PathEscape(filename)
		store.ResponseHeader().Add("content-disposition", "attachment; filename*=UTF-8''"+filename+".zip; filename=\""+filename+".zip\"")
		if err := fsStore.GetDirAsZip(path, store.Writer()); err != nil {
			store.ResponseHeader().Del("content-disposition")
			logger.Panic(err)
		}
		return
	} else {
		store.WriteHeader(http.StatusUnprocessableEntity)
		store.RespondJson(nil)
		return
	}
}

func createDownloadTask(url string, dir string) func() {
	return func() {
		file, err := fsStore.CreateFile(filepath.Join(dir, path.Base(url)))
		if err != nil {
			return
		}
		defer file.Close()

		if resp, err := http.Get(url); err == nil {
			defer resp.Body.Close()
			io.Copy(file, resp.Body)
		}
	}
}

func UploadHandler(store httpd.Store) {
	path := filepath.Join("/", store.RouteAny())
	reader, err := store.MultipartReader()
	if err != nil {
		logger.Panic(err)
	}
	shortestQueue := downloadTaskLane.ShortestQueueIndex()
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Panic(err)
		}
		if part.FormName() == "urlList" {
			url, err := io.ReadAll(part)
			if err != nil {
				logger.Panic(err)
			}
			downloadTaskLane.PushTask(createDownloadTask(string(url), path), shortestQueue)
		} else if part.FormName() == "fileList" {
			file, err := fsStore.CreateFile(filepath.Join(path, part.FileName()))
			if err != nil {
				logger.Panic(err)
			}
			defer file.Close()
			io.Copy(file, part)
		}
	}
	store.Respond200(nil)
}
