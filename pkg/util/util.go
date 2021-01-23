package util

// Contain determines whether a string slice includes a certain value.
func Contain(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// IsSpace determines whether a byte belongs to space.
func IsSpace(ch byte) bool {
	return ' ' == ch ||
		'\n' == ch ||
		'\r' == ch ||
		'\t' == ch ||
		'\v' == ch ||
		'\f' == ch
}
