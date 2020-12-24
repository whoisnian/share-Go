package server

import (
	"mime/multipart"
	"net/url"
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

// MultipartReader ...
func (store Store) MultipartReader() (*multipart.Reader, error) {
	return store.r.MultipartReader()
}
