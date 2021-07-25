package util

// IsSpace determines whether a byte belongs to space.
func IsSpace(ch byte) bool {
	return ch == ' ' ||
		ch == '\n' ||
		ch == '\r' ||
		ch == '\t' ||
		ch == '\v' ||
		ch == '\f'
}
