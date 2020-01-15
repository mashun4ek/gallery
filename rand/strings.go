package rand

import (
	"crypto/rand"
	"encoding/base64"
)

var RememberTokenBytes = 32

// Bytes generate n bytes or return error - helper func
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// NBytes returns the number of bytes used in the base64 URL encoded string
func NBytes(base64String string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(base64String)
	if err != nil {
		return -1, err
	}
	return len(b), nil
}

// String generates a byte slice of size nBytes and then return a string
// that is the base64 URL encoded version of that byte slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	// it's not encrypting, just encode to string a slice of bytes
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken a helper func designed to generate remember tokens
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
