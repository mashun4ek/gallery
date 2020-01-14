package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mashun4ek/webdevcalhoun/gallery/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "mariaker"
	password = "nopassword"
	dbname   = "gallery"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"unique;not null"`
	Color string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()
	user := models.User{
		Name:  "Michael Scott",
		Email: "DM",
	}
	if err := us.Create(&user); err != nil {
		panic(err)
	}
	// user.Email = "newEmail"
	// if err := us.Update(&user); err != nil {
	// 	panic(err)
	// }
	// userByEmail, err := us.ByEmail("newEmail")
	// if err != nil {
	// 	panic(err)
	// }
	if err := us.Delete(user.ID); err != nil {
		panic(err)
	}
	userByID, err := us.ByID(user.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(userByID)
}
