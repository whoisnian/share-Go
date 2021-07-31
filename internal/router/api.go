package router

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/glb/util/fsutil"
	"github.com/whoisnian/share-Go/internal/config"
)

const (
	typeRegular   int64 = 0
	typeDirectory int64 = 1
)

type respFileInfo struct {
	Type  int64
	Name  string
	Size  int64
	MTime int64
}

func parseFileInfo(info os.FileInfo) respFileInfo {
	t := typeRegular
	if info.Mode().IsDir() {
		t = typeDirectory
	}
	return respFileInfo{
		Type:  t,
		Name:  info.Name(),
		Size:  info.Size(),
		MTime: info.ModTime().Unix(),
	}
}

func FileInfoHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(config.RootPath, store.RouteParamAny())
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			store.W.WriteHeader(http.StatusNotFound)
			return
		}
		logger.Panic(err)
	}
	store.RespondJson(parseFileInfo(info))
}

func NewFileHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(config.RootPath, store.RouteParamAny())
	file, err := lockedFS.Create(path)
	if err != nil {
		logger.Panic(err)
	}
	defer file.Close()

	body := store.R.Body
	defer body.Close()

	io.Copy(file, body)
}

func DeleteFileHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(config.RootPath, store.RouteParamAny())
	if err := os.Remove(path); err != nil {
		if errors.Is(err, syscall.ENOTEMPTY) {
			store.W.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Panic(err)
	}
}

func ListDirHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(config.RootPath, store.RouteParamAny())
	dir, err := lockedFS.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			store.W.WriteHeader(http.StatusNotFound)
			return
		}
		logger.Panic(err)
	}
	defer dir.Close()

	infos, err := dir.Readdir(-1)
	if err != nil {
		if errors.Is(err, syscall.ENOTDIR) {
			store.W.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Panic(err)
	}

	result := make([]respFileInfo, len(infos))
	for i := 0; i < len(infos); i++ {
		result[i] = parseFileInfo(infos[i])
	}

	store.RespondJson(jsonMap{"FileInfos": result})
}

func NewDirHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(config.RootPath, store.RouteParamAny())
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		logger.Panic(err)
	}
}

func DeleteDirHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(config.RootPath, store.RouteParamAny())
	if err := os.RemoveAll(path); err != nil {
		logger.Panic(err)
	}
}

func RawHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(config.RootPath, store.RouteParamAny())
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			store.W.WriteHeader(http.StatusNotFound)
			return
		}
		logger.Panic(err)
	}

	if info.Mode().IsRegular() {
		file, err := lockedFS.Open(path)
		if err != nil {
			logger.Panic(err)
		}
		defer file.Close()

		if _, err := io.Copy(store.W, file); err != nil {
			logger.Panic(err)
		}
	} else {
		store.W.WriteHeader(http.StatusUnprocessableEntity)
	}
}

func archiveDirAsZip(dirPath string, zipWriter *zip.Writer) error {
	walkFunc := func(fullPath string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		file, err := lockedFS.Open(fullPath)
		if err != nil {
			return err
		}
		defer file.Close()

		relativePath, err := filepath.Rel(dirPath, fullPath)
		if err != nil {
			return err
		}
		zipFile, err := zipWriter.CreateHeader(&zip.FileHeader{
			Name:   relativePath,
			Method: zip.Store,
		})
		if err != nil {
			return err
		}

		if _, err := io.Copy(zipFile, file); err != nil {
			return err
		}

		return nil
	}
	if err := filepath.WalkDir(dirPath, walkFunc); err != nil {
		return err
	}
	return zipWriter.Close()
}

func DownloadHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(config.RootPath, store.RouteParamAny())
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			store.W.WriteHeader(http.StatusNotFound)
			return
		}
		logger.Panic(err)
	}

	if info.Mode().IsRegular() {
		file, err := lockedFS.Open(path)
		if err != nil {
			logger.Panic(err)
		}
		defer file.Close()

		filename := url.PathEscape(filepath.Base(path))
		store.W.Header().Add("content-disposition", "attachment; filename*=UTF-8''"+filename+"; filename=\""+filename+"\"")

		if _, err := io.Copy(store.W, file); err != nil {
			store.W.Header().Del("content-disposition")
			logger.Panic(err)
		}
	} else if info.Mode().IsDir() {
		filename := filepath.Base(path)
		if filename == "/" {
			filename = "root"
		}
		filename = url.PathEscape(filename)
		store.W.Header().Add("content-disposition", "attachment; filename*=UTF-8''"+filename+".zip; filename=\""+filename+".zip\"")
		zipWriter := zip.NewWriter(store.W)
		if err := archiveDirAsZip(path, zipWriter); err != nil {
			store.W.Header().Del("content-disposition")
			logger.Panic(err)
		}
	} else {
		store.W.WriteHeader(http.StatusUnprocessableEntity)
	}
}

func createDownloadTask(url string, dir string) func() {
	return func() {
		file, err := lockedFS.Create(filepath.Join(dir, path.Base(url)))
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

func UploadHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(config.RootPath, store.RouteParamAny())
	reader, err := store.R.MultipartReader()
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
			file, err := lockedFS.Create(filepath.Join(path, part.FileName()))
			if err != nil {
				logger.Panic(err)
			}
			defer file.Close()
			io.Copy(file, part)
		}
	}
}
