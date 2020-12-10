package state

// CookieValue ...
func (store Store) CookieValue(name string) string {
	cookie, err := store.R.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}
