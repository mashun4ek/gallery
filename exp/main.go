package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func main() {
	toHash := []byte("this is my string to hash")
	h := hmac.New(sha256.New, []byte("my-secret-key"))
	h.Write(toHash)
	b := h.Sum(nil)
	fmt.Println(b)
	h.Reset()

	h = hmac.New(sha256.New, []byte("my-secret-key"))
	h.Write(toHash)
	b = h.Sum(nil)
	fmt.Println(b)
}
