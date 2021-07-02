package util

import (
	"net"
)

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
	return ch == ' ' ||
		ch == '\n' ||
		ch == '\r' ||
		ch == '\t' ||
		ch == '\v' ||
		ch == '\f'
}

// GetOutBoundIP get preferred outbound ip of current process.
// https://stackoverflow.com/a/37382208
func GetOutBoundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
