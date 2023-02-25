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
	"sync"
	"syscall"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/glb/util/fsutil"
	"github.com/whoisnian/share-Go/internal/global"
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

type lockedFile struct {
	*os.File
	locker sync.Locker
}

func (f *lockedFile) Close() error {
	f.locker.Unlock()
	return f.File.Close()
}

func createFile(name string) (*lockedFile, error) {
	v, _ := lockerMap.LoadOrStore(name, new(sync.RWMutex))
	locker := v.(*sync.RWMutex)
	locker.Lock()

	file, err := os.Create(name)
	if err != nil {
		locker.Unlock()
		return nil, err
	}
	return &lockedFile{file, locker}, nil
}

func openFile(name string) (*lockedFile, error) {
	v, _ := lockerMap.LoadOrStore(name, new(sync.RWMutex))
	locker := v.(*sync.RWMutex)
	locker.RLock()

	file, err := os.Open(name)
	if err != nil {
		locker.RUnlock()
		return nil, err
	}
	return &lockedFile{file, locker.RLocker()}, nil
}

func fileInfoHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(global.CFG.RootPath, store.RouteParamAny())
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

func newFileHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(global.CFG.RootPath, store.RouteParamAny())
	file, err := createFile(path)
	if err != nil {
		logger.Panic(err)
	}
	defer file.Close()

	body := store.R.Body
	defer body.Close()

	io.Copy(file, body)
}

func deleteFileHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(global.CFG.RootPath, store.RouteParamAny())

	fileInfo, _ := os.Stat(path)
	rootInfo, _ := os.Stat(global.CFG.RootPath)
	if os.SameFile(rootInfo, fileInfo) {
		store.W.WriteHeader(http.StatusForbidden)
		return
	}

	if err := os.Remove(path); err != nil {
		if errors.Is(err, syscall.ENOTEMPTY) {
			store.W.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Panic(err)
	}
}

func listDirHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(global.CFG.RootPath, store.RouteParamAny())
	dir, err := openFile(path)
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

func newDirHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(global.CFG.RootPath, store.RouteParamAny())
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		logger.Panic(err)
	}
}

func deleteDirHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(global.CFG.RootPath, store.RouteParamAny())

	dirInfo, _ := os.Stat(path)
	rootInfo, _ := os.Stat(global.CFG.RootPath)
	if os.SameFile(rootInfo, dirInfo) {
		store.W.WriteHeader(http.StatusForbidden)
		return
	}

	if err := os.RemoveAll(path); err != nil {
		logger.Panic(err)
	}
}

func rawHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(global.CFG.RootPath, store.RouteParamAny())
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			store.W.WriteHeader(http.StatusNotFound)
			return
		}
		logger.Panic(err)
	}

	if info.Mode().IsRegular() {
		file, err := openFile(path)
		if err != nil {
			logger.Panic(err)
		}
		defer file.Close()

		http.ServeFile(store.W, store.R, path)
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

		file, err := openFile(fullPath)
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

func downloadHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(global.CFG.RootPath, store.RouteParamAny())
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			store.W.WriteHeader(http.StatusNotFound)
			return
		}
		logger.Panic(err)
	}

	if info.Mode().IsRegular() {
		file, err := openFile(path)
		if err != nil {
			logger.Panic(err)
		}
		defer file.Close()

		filename := url.PathEscape(filepath.Base(path))
		store.W.Header().Set("content-disposition", "attachment; filename*=UTF-8''"+filename+"; filename=\""+filename+"\"")

		http.ServeFile(store.W, store.R, path)
	} else if info.Mode().IsDir() {
		filename := filepath.Base(path)
		if filename == "/" {
			filename = "root"
		}
		filename = url.PathEscape(filename)
		store.W.Header().Set("content-disposition", "attachment; filename*=UTF-8''"+filename+".zip; filename=\""+filename+".zip\"")
		zipWriter := zip.NewWriter(store.W)
		if err := archiveDirAsZip(path, zipWriter); err != nil {
			store.W.Header().Del("content-disposition")
			logger.Panic(err)
		}
	} else {
		store.W.WriteHeader(http.StatusUnprocessableEntity)
	}
}

func renameHandler(store *httpd.Store) {
	fromPath := fsutil.ResolveBase(global.CFG.RootPath, store.RouteParamAny())
	toPath := fsutil.ResolveBase(global.CFG.RootPath, store.R.URL.Query().Get("to"))

	fromInfo, err := os.Stat(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			store.W.WriteHeader(http.StatusNotFound)
			store.RespondJson(jsonMap{"Message": "Source file or folder not found"})
			return
		}
		logger.Panic(err)
	}
	rootInfo, _ := os.Stat(global.CFG.RootPath)
	if os.SameFile(rootInfo, fromInfo) {
		store.W.WriteHeader(http.StatusForbidden)
		store.RespondJson(jsonMap{"Message": "Forbidden to rename the root folder"})
		return
	}

	if _, err := os.Stat(toPath); err == nil {
		store.W.WriteHeader(http.StatusConflict)
		store.RespondJson(jsonMap{"Message": "Destination file or folder already exists"})
		return
	}

	toPathParent := filepath.Dir(toPath)
	if err := os.MkdirAll(toPathParent, os.ModePerm); err != nil {
		logger.Panic(err)
	}
	if err := os.Rename(fromPath, toPath); err != nil {
		logger.Panic(err)
	}
	store.RespondJson(jsonMap{"Message": "Success"})
}

func uploadHandler(store *httpd.Store) {
	path := fsutil.ResolveBase(global.CFG.RootPath, store.RouteParamAny())
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
			downloadTaskLane.PushTask(newDownloadTask(string(url), path), shortestQueue)
		} else if part.FormName() == "fileList" {
			file, err := createFile(filepath.Join(path, part.FileName()))
			if err != nil {
				logger.Panic(err)
			}
			defer file.Close()
			io.Copy(file, part)
		}
	}
}

type downloadTask struct {
	url string
	dir string
	err error
}

func newDownloadTask(url string, dir string) *downloadTask {
	return &downloadTask{url, dir, nil}
}

func (t *downloadTask) Start() {
	file, err := createFile(filepath.Join(t.dir, path.Base(t.url)))
	if err != nil {
		t.Done(err)
		return
	}
	defer file.Close()

	resp, err := http.Get(t.url)
	if err != nil {
		t.Done(err)
		return
	}
	defer resp.Body.Close()
	io.Copy(file, resp.Body)
	t.Done(nil)
}

func (t *downloadTask) Done(err error) {
	if err != nil {
		logger.Error(err)
	}
}
