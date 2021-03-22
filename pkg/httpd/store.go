package httpd

import (
	"encoding/json"
	"io"
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

// Store consists of request, routeParams and responseWriter.
type Store struct {
	w *statusResponseWriter
	r *http.Request
	m map[string]string
}

// RouteParam returns the value of specified route param, or empty string if param not found.
func (store Store) RouteParam(name string) string {
	if param, ok := store.m[name]; ok {
		return param
	}
	return ""
}

// RouteAny returns the value of route param "/*".
func (store Store) RouteAny() string {
	return store.RouteParam(routeAny)
}

// CookieValue returns the value of specified cookie, or empty string if cookie not found.
func (store Store) CookieValue(name string) string {
	if cookie, err := store.r.Cookie(name); err == nil {
		return cookie.Value
	}
	return ""
}

// URL equals to `http.Request.URL`.
func (store Store) URL() *url.URL {
	return store.r.URL
}

// Body equals to `http.Request.Body`.
func (store Store) Body() io.ReadCloser {
	return store.r.Body
}

// MultipartReader equals to `http.Request.MultipartReader()`.
func (store Store) MultipartReader() (*multipart.Reader, error) {
	return store.r.MultipartReader()
}

// WriteHeader equals to `http.ResponseWriter.WriteHeader()`.
func (store Store) WriteHeader(code int) {
	store.w.WriteHeader(code)
}

// Respond200 replies 200 to client request with optional body.
func (store Store) Respond200(content []byte) error {
	store.w.WriteHeader(http.StatusOK)
	if len(content) > 0 {
		_, err := store.w.Write(content)
		return err
	}
	return nil
}

// RespondJson replies 200 to client request with json body.
func (store Store) RespondJson(v interface{}) error {
	store.w.Header().Add("content-type", "application/json; charset=UTF-8")
	return json.NewEncoder(store.w).Encode(v)
}

// Redirect is similar to `http.Redirect()`.
func (store Store) Redirect(url string, code int) {
	http.Redirect(store.w, store.r, url, code)
}

// Respond404 is similar to `http.NotFound()`.
func (store Store) Respond404() {
	http.NotFound(store.w, store.r)
}

// Error500 is similar to `http.Error()`.
func (store Store) Error500(err string) {
	http.Error(store.w, err, http.StatusInternalServerError)
}
