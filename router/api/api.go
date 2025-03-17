package api

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"syscall"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/glb/util/fsutil"
	"github.com/whoisnian/glb/util/osutil"
	"github.com/whoisnian/share-Go/global"
)

func FileInfoHandler(store *httpd.Store) {
	fpath := fsutil.ResolveUrlPath(global.CFG.RootPath, store.RouteParamAny())
	info, err := os.Stat(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			store.W.WriteHeader(http.StatusNotFound)
			return
		}
		global.LOG.Error(store.R.Context(), "os.Stat failed", logger.Error(err))
		store.W.WriteHeader(http.StatusInternalServerError)
		return
	}
	store.RespondJson(prepareRespFileInfo(info))
}

func NewFileHandler(store *httpd.Store) {
	fpath := fsutil.ResolveUrlPath(global.CFG.RootPath, store.RouteParamAny())
	if err := writeNewFile(fpath, store.R.Body); err != nil {
		global.LOG.Error(store.R.Context(), "writeNewFile failed", logger.Error(err))
		store.W.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func DeleteFileHandler(store *httpd.Store) {
	fpath := fsutil.ResolveUrlPath(global.CFG.RootPath, store.RouteParamAny())
	info, err := os.Stat(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		global.LOG.Error(store.R.Context(), "os.Stat failed", logger.Error(err))
		store.W.WriteHeader(http.StatusInternalServerError)
		return
	}
	if info.IsDir() {
		store.W.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := os.Remove(fpath); err != nil {
		global.LOG.Error(store.R.Context(), "os.Remove failed", logger.Error(err))
		store.W.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ListDirHandler(store *httpd.Store) {
	fpath := fsutil.ResolveUrlPath(global.CFG.RootPath, store.RouteParamAny())
	dir, err := os.Open(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			store.W.WriteHeader(http.StatusNotFound)
			return
		}
		global.LOG.Error(store.R.Context(), "os.Open failed", logger.Error(err))
		store.W.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dir.Close()

	infos, err := dir.Readdir(-1)
	if err != nil {
		if errors.Is(err, syscall.ENOTDIR) {
			store.W.WriteHeader(http.StatusBadRequest)
			return
		}
		global.LOG.Error(store.R.Context(), "dir.Readdir failed", logger.Error(err))
		store.W.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := make([]respFileInfo, len(infos))
	for i := 0; i < len(infos); i++ {
		result[i] = prepareRespFileInfo(infos[i])
	}
	store.RespondJson(respFileInfos{result})
}

func NewDirHandler(store *httpd.Store) {
	fpath := fsutil.ResolveUrlPath(global.CFG.RootPath, store.RouteParamAny())
	if err := os.MkdirAll(fpath, osutil.DefaultDirMode); err != nil {
		global.LOG.Error(store.R.Context(), "os.MkdirAll failed", logger.Error(err))
		store.W.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func DeleteDirHandler(store *httpd.Store) {
	fpath := fsutil.ResolveUrlPath(global.CFG.RootPath, store.RouteParamAny())
	dirInfo, _ := os.Stat(fpath)
	rootInfo, _ := os.Stat(global.CFG.RootPath)
	if os.SameFile(rootInfo, dirInfo) {
		store.W.WriteHeader(http.StatusForbidden)
		return
	}

	if err := os.RemoveAll(fpath); err != nil {
		global.LOG.Error(store.R.Context(), "os.RemoveAll failed", logger.Error(err))
		store.W.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func RenameHandler(store *httpd.Store) {
	fromPath := fsutil.ResolveUrlPath(global.CFG.RootPath, store.RouteParamAny())
	toPath := fsutil.ResolveUrlPath(global.CFG.RootPath, store.R.URL.Query().Get("to"))

	fromInfo, err := os.Stat(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			store.W.WriteHeader(http.StatusNotFound)
			store.RespondJson(respMessage{"source file or folder not found"})
			return
		}
		global.LOG.Error(store.R.Context(), "os.Stat failed", logger.Error(err))
		store.W.WriteHeader(http.StatusInternalServerError)
		return
	}
	rootInfo, _ := os.Stat(global.CFG.RootPath)
	if os.SameFile(rootInfo, fromInfo) {
		store.W.WriteHeader(http.StatusForbidden)
		store.RespondJson(respMessage{"forbidden to rename the root folder"})
		return
	}

	if _, err := os.Stat(toPath); err == nil {
		store.W.WriteHeader(http.StatusConflict)
		store.RespondJson(respMessage{"destination file or folder already exists"})
		return
	}

	toPathParent := filepath.Dir(toPath)
	if err := os.MkdirAll(toPathParent, osutil.DefaultDirMode); err != nil {
		global.LOG.Error(store.R.Context(), "os.MkdirAll failed", logger.Error(err))
		store.W.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := os.Rename(fromPath, toPath); err != nil {
		global.LOG.Error(store.R.Context(), "os.Rename failed", logger.Error(err))
		store.W.WriteHeader(http.StatusInternalServerError)
		return
	}
	store.RespondJson(respMessage{"success"})
}

func UploadHandler(store *httpd.Store) {
	fpath := fsutil.ResolveUrlPath(global.CFG.RootPath, store.RouteParamAny())
	reader, err := store.R.MultipartReader()
	if err != nil {
		global.LOG.Error(store.R.Context(), "store.R.MultipartReader failed", logger.Error(err))
		store.W.WriteHeader(http.StatusInternalServerError)
		return
	}
	shortestQueue := _taskLane.ShortestQueueIndex()
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			global.LOG.Error(store.R.Context(), "reader.NextPart failed", logger.Error(err))
			store.W.WriteHeader(http.StatusInternalServerError)
			return
		}
		if part.FormName() == "urlList" {
			url, err := io.ReadAll(part)
			if err != nil {
				global.LOG.Error(store.R.Context(), "io.ReadAll failed", logger.Error(err))
				store.W.WriteHeader(http.StatusInternalServerError)
				return
			}
			task := newDownloadTask(string(url), fpath)
			if err := _taskLane.PushTask(task, shortestQueue); err != nil {
				global.LOG.Error(store.R.Context(), "taskLane.PushTask failed", logger.Error(err))
				store.W.WriteHeader(http.StatusInternalServerError)
				return
			}
			global.LOG.Infof(store.R.Context(), "task(%s) pushed with %s %s", task.tid, task.url, task.dir)
		} else if part.FormName() == "fileList" {
			if err := writeNewFile(filepath.Join(fpath, part.FileName()), part); err != nil {
				global.LOG.Error(store.R.Context(), "writeNewFile failed", logger.Error(err))
				store.W.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}

func RawHandler(store *httpd.Store) {
	fpath := fsutil.ResolveUrlPath(global.CFG.RootPath, store.RouteParamAny())
	info, err := os.Stat(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			store.W.WriteHeader(http.StatusNotFound)
			return
		}
		global.LOG.Error(store.R.Context(), "os.Stat failed", logger.Error(err))
		store.W.WriteHeader(http.StatusInternalServerError)
		return
	}

	if info.Mode().IsRegular() {
		file, err := os.Open(fpath)
		if err != nil {
			global.LOG.Error(store.R.Context(), "os.Open failed", logger.Error(err))
			store.W.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		http.ServeFile(store.W, store.R, fpath)
	} else {
		store.W.WriteHeader(http.StatusUnprocessableEntity)
	}
}

func archiveDirAsZip(dirPath string, zipWriter *zip.Writer) error {
	walkFunc := func(fullPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(dirPath, fullPath)
		if err != nil || relativePath == "" || relativePath == "." {
			return err
		}
		if info.IsDir() {
			// https://unix.stackexchange.com/q/743511
			// https://github.com/python/cpython/pull/9419
			fh := &zip.FileHeader{
				Name:     relativePath,
				Modified: info.ModTime(),
				Method:   zip.Store,
			}
			fh.SetMode(info.Mode())
			if fh.Name[len(fh.Name)-1] != '/' {
				fh.Name += "/"
			}
			_, err = zipWriter.CreateHeader(fh)
			return err
		}

		file, err := os.Open(fullPath)
		if err != nil {
			return err
		}
		defer file.Close()

		fh := &zip.FileHeader{
			Name:               relativePath,
			UncompressedSize64: uint64(info.Size()),
			Modified:           info.ModTime(),
			Method:             zip.Deflate,
		}
		fh.SetMode(info.Mode())
		zipFile, err := zipWriter.CreateHeader(fh)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipFile, file)
		return err
	}
	if err := filepath.Walk(dirPath, walkFunc); err != nil {
		return err
	}
	return zipWriter.Close()
}

func DownloadHandler(store *httpd.Store) {
	fpath := fsutil.ResolveUrlPath(global.CFG.RootPath, store.RouteParamAny())
	info, err := os.Stat(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			store.W.WriteHeader(http.StatusNotFound)
			return
		}
		global.LOG.Error(store.R.Context(), "os.Stat failed", logger.Error(err))
		store.W.WriteHeader(http.StatusInternalServerError)
		return
	}

	if info.Mode().IsRegular() {
		file, err := os.Open(fpath)
		if err != nil {
			global.LOG.Error(store.R.Context(), "os.Open failed", logger.Error(err))
			store.W.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		filename := url.PathEscape(filepath.Base(fpath))
		store.W.Header().Set("content-disposition", "attachment; filename*=UTF-8''"+filename+"; filename=\""+filename+"\"")

		http.ServeFile(store.W, store.R, fpath)
	} else if info.Mode().IsDir() {
		filename := filepath.Base(fpath)
		if filename == "/" {
			filename = "root"
		}
		filename = url.PathEscape(filename)
		store.W.Header().Set("content-disposition", "attachment; filename*=UTF-8''"+filename+".zip; filename=\""+filename+".zip\"")
		zipWriter := zip.NewWriter(store.W)
		if err := archiveDirAsZip(fpath, zipWriter); err != nil {
			store.W.Header().Del("content-disposition")
			global.LOG.Error(store.R.Context(), "archiveDirAsZip failed", logger.Error(err))
			store.W.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		store.W.WriteHeader(http.StatusUnprocessableEntity)
	}
}
