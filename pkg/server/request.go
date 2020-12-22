package server

import (
	"io"
	"net/url"
	"os"
	"path"
)

// CookieValue ...
func (store Store) CookieValue(name string) string {
	cookie, err := store.r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// URL ...
func (store Store) URL() *url.URL {
	return store.r.URL
}

// Save ...
func (store Store) Save(dir string) {
	reader, err := store.r.MultipartReader()
	if err != nil {
		return
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			return
		}
		if err != nil {
			return
		}
		file, err := os.Create(path.Join(dir, part.FileName()))
		if err != nil {
			return
		}
		io.Copy(file, part)
		file.Close()
	}
}
