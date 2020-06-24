package main

import (
	"fmt"

	"github.com/mashun4ek/webdevcalhoun/gallery/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = ""
	password = ""
	dbname   = ""
)

func main() {
	// toHash := []byte("this is my string to hash")
	// h := hmac.New(sha256.New, []byte("my-secret-key"))
	// h.Write(toHash)
	// b := h.Sum(nil)
	// fmt.Println(b)
	// h.Reset()

	// h = hmac.New(sha256.New, []byte("my-secret-key"))
	// h.Write(toHash)
	// b = h.Sum(nil)
	// fmt.Println(b)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.DestructiveReset()

	user := models.User{
		Name:     "Maria Ker",
		Email:    "mashun4ek@gmail.com",
		Password: "slon",
		Remember: "hey123",
	}
	err = us.Create(&user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", user)
	user2, err := us.ByRemember("hey123")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *user2)
}
