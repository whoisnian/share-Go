package storage

import (
	"io"
	"os"
	"path/filepath"

	"github.com/whoisnian/share-Go/pkg/logger"
)

// Store ...
type Store struct {
	base string
}

// New ...
func New(path string) *Store {
	base, err := filepath.Abs(path)
	if err != nil {
		logger.Panic(err)
	}
	if _, err := os.Lstat(base); os.IsNotExist(err) {
		logger.Panic(err)
	}
	return &Store{base}
}

// IsDir ...
func (store *Store) IsDir(path string) bool {
	if fileInfo, err := os.Lstat(filepath.Join(store.base, path)); err == nil {
		return fileInfo.Mode().IsDir()
	}
	return false
}

// IsFile ...
func (store *Store) IsFile(path string) bool {
	if fileInfo, err := os.Lstat(filepath.Join(store.base, path)); err == nil {
		return fileInfo.Mode().IsRegular()
	}
	return false
}

// FileInfo ...
func (store *Store) FileInfo(path string) (os.FileInfo, error) {
	return os.Lstat(filepath.Join(store.base, path))
}

// ListDir ...
func (store *Store) ListDir(path string) ([]os.FileInfo, error) {
	dir, err := os.Open(filepath.Join(store.base, path))
	if err != nil {
		return nil, err
	}
	return dir.Readdir(-1)
}

// DeleteDir ...
func (store *Store) DeleteDir(path string) error {
	return os.RemoveAll(filepath.Join(store.base, path))
}

// CreateDir ...
func (store *Store) CreateDir(path string) error {
	return os.MkdirAll(filepath.Join(store.base, path), os.ModePerm)
}

// GetFile ...
func (store *Store) GetFile(path string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(store.base, path))
}

// CreateFile ...
func (store *Store) CreateFile(path string) (io.WriteCloser, error) {
	return os.Create(filepath.Join(store.base, path))
}
