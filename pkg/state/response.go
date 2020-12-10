package state

import "net/http"

// Respond200 ...
func (store Store) Respond200(content []byte) error {
	store.Code = http.StatusOK
	store.W.WriteHeader(http.StatusOK)
	if len(content) > 0 {
		_, err := store.W.Write(content)
		return err
	}
	return nil
}

// Redirect ...
func (store Store) Redirect(url string, code int) {
	store.Code = code
	http.Redirect(store.W, store.R, url, code)
}

// Redirect301 ...
func (store Store) Redirect301(url string) {
	store.Redirect(url, http.StatusMovedPermanently)
}

// Redirect302 ...
func (store Store) Redirect302(url string) {
	store.Redirect(url, http.StatusFound)
}

// Error500 ...
func (store Store) Error500(err string) {
	store.Code = http.StatusInternalServerError
	http.Error(store.W, err, http.StatusInternalServerError)
}
