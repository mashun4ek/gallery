package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound when resource is not found in database
	ErrNotFound = errors.New("models: resource not found")
	// ErrInvalidID when invalid ID provided (Delete method)
	ErrInvalidID = errors.New("models: ID must be > 0")
	// ErrInvalidPassword is returned when invalid password is used to authenticate a user
	ErrInvalidPassword = errors.New("models: Incorrect password provided")
)

// pepper
const userPwPepper = "unique8!@gallery!"

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
}

type UserService struct {
	db *gorm.DB
}

// DestructiveReset drops user table and rebuilds it
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the user table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// ByID will look up by the id provided
// 1 - user, nil
// 2 - nil, ErrNotFound
// 3 - nil, otherError
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail looks up a user by email
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// Authenticate to Authenticate the user with provided email and password
// if user doesn't exist
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}
	return foundUser, nil
}

// first will query using the provided gorm.DB and get the first item returned and place it into dst
// If nothing is found in the query, it will return ERRNOtFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// Create user
func (us *UserService) Create(user *User) error {
	// if user.ID > 0 {
	// 	return NotNewRecor
	// }
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return us.db.Create(user).Error
}

// Update user
func (us *UserService) Update(user *User) error {
	return us.db.Save(&user).Error
}

// Delete user
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

// Close closes database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"unique;not null"`
	// don't include password to database, just include password hash
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
}
