package server

import (
	"log"
	"net/http"
	"strings"
	"time"
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

// MakeHander ...
func MakeHander(fn func(store Store)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		store := Store{&statusResponseWriter{w, http.StatusOK}, r}

		fn(store)

		log.Printf("%s [%d] %s %s %s %d",
			r.RemoteAddr[0:strings.IndexByte(r.RemoteAddr, ':')],
			store.w.status,
			r.Method,
			r.URL.Path,
			r.UserAgent(),
			time.Now().Sub(start).Milliseconds())
	}
}
