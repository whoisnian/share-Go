package server

import (
	"mime/multipart"
	"net/http"
	"net/url"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Store consists of both request and responseWriter
type Store struct {
	w *statusResponseWriter
	r *http.Request
}

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

// Respond200 ...
func (store Store) Respond200(content []byte) error {
	store.w.WriteHeader(http.StatusOK)
	if len(content) > 0 {
		_, err := store.w.Write(content)
		return err
	}
	return nil
}

// Redirect ...
func (store Store) Redirect(url string, code int) {
	http.Redirect(store.w, store.r, url, code)
}

// Redirect301 ...
func (store Store) Redirect301(url string) {
	store.Redirect(url, http.StatusMovedPermanently)
}

// Redirect302 ...
func (store Store) Redirect302(url string) {
	store.Redirect(url, http.StatusFound)
}

// Respond404 ...
func (store Store) Respond404() {
	http.NotFound(store.w, store.r)
}

// Error500 ...
func (store Store) Error500(err string) {
	http.Error(store.w, err, http.StatusInternalServerError)
}
