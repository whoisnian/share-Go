package storage

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/whoisnian/share-Go/pkg/logger"
)

// Store ...
type Store struct {
	base    string
	lockMap *sync.Map
}

func (store *Store) getLocker(path string) *sync.RWMutex {
	locker, _ := store.lockMap.LoadOrStore(path, new(sync.RWMutex))
	return locker.(*sync.RWMutex)
}

type readCloser struct {
	*os.File
	locker *sync.RWMutex
}

func (r readCloser) Close() error {
	if err := r.File.Close(); err != nil {
		return err
	}
	r.locker.RUnlock()
	return nil
}

type writeCloser struct {
	*os.File
	locker *sync.RWMutex
}

func (w writeCloser) Close() error {
	if err := w.File.Close(); err != nil {
		return err
	}
	w.locker.Unlock()
	return nil
}

// New ...
func New(path string) *Store {
	base, err := filepath.Abs(path)
	if err != nil {
		logger.Panic(err)
	}
	if _, err := os.Stat(base); os.IsNotExist(err) {
		logger.Panic(err)
	}
	return &Store{base, new(sync.Map)}
}

// IsDir ...
func (store *Store) IsDir(path string) bool {
	if fileInfo, err := os.Stat(filepath.Join(store.base, path)); err == nil {
		return fileInfo.Mode().IsDir()
	}
	return false
}

// IsFile ...
func (store *Store) IsFile(path string) bool {
	if fileInfo, err := os.Stat(filepath.Join(store.base, path)); err == nil {
		return fileInfo.Mode().IsRegular()
	}
	return false
}

// FileInfo ...
func (store *Store) FileInfo(path string) (os.FileInfo, error) {
	return os.Stat(filepath.Join(store.base, path))
}

// ListDir ...
func (store *Store) ListDir(path string) ([]os.FileInfo, error) {
	dir, err := os.Open(filepath.Join(store.base, path))
	if err != nil {
		return nil, err
	}
	return dir.Readdir(-1)
}

// Delete ...
func (store *Store) Delete(path string) error {
	return os.Remove(filepath.Join(store.base, path))
}

// DeleteAll ...
func (store *Store) DeleteAll(path string) error {
	return os.RemoveAll(filepath.Join(store.base, path))
}

// CreateDir ...
func (store *Store) CreateDir(path string) error {
	return os.MkdirAll(filepath.Join(store.base, path), os.ModePerm)
}

// GetDirAsZip
func (store *Store) GetDirAsZip(path string, writer io.Writer) error {
	dirPath := filepath.Join(store.base, path)
	zipWriter := zip.NewWriter(writer)
	walkFunc := func(fullPath string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		lock := store.getLocker(fullPath)
		lock.RLock()

		file, err := os.Open(fullPath)
		if err != nil {
			lock.RUnlock()
			return err
		}
		defer lock.RUnlock()
		defer file.Close()

		relativePath, err := filepath.Rel(dirPath, fullPath)
		if err != nil {
			return err
		}
		zipFile, err := zipWriter.Create(relativePath)
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

// GetFile ...
func (store *Store) GetFile(path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(store.base, path)
	lock := store.getLocker(fullPath)
	lock.RLock()

	file, err := os.Open(fullPath)
	if err != nil {
		lock.RUnlock()
		return nil, err
	}
	return &readCloser{file, lock}, nil
}

// CreateFile ...
func (store *Store) CreateFile(path string) (io.WriteCloser, error) {
	fullPath := filepath.Join(store.base, path)
	lock := store.getLocker(fullPath)
	lock.Lock()

	file, err := os.Create(fullPath)
	if err != nil {
		lock.Unlock()
		return nil, err
	}
	return &writeCloser{file, lock}, nil
}
