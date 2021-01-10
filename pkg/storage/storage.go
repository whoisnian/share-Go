package storage

import (
	"log"
	"os"
	"path/filepath"
)

// Store ...
type Store struct {
	base string
}

// New ...
func New(path string) *Store {
	base, err := filepath.Abs(path)
	if err != nil {
		log.Panicln(err)
	}
	if _, err := os.Lstat(base); os.IsNotExist(err) {
		log.Panicln(err)
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
