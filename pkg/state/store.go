package state

import "net/http"

// Store ...
type Store struct {
	W    http.ResponseWriter
	R    *http.Request
	Code int
}

// NewStore ...
func NewStore(w http.ResponseWriter, r *http.Request) Store {
	return Store{
		W:    w,
		R:    r,
		Code: http.StatusOK,
	}
}

// MakeHander ...
func MakeHander(fn func(store Store)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(NewStore(w, r))
	}
}
