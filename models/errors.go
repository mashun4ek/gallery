package models

import "strings"

const (
	// ErrNotFound when resource is not found in database
	ErrNotFound modelError = "models: Resource not found"
	// ErrPasswordIncorrect is returned when invalid password is used to authenticate a user
	ErrPasswordIncorrect modelError = "models: Incorrect password provided"
	// ErrEmailTaken if email is already in database
	ErrEmailTaken modelError = "models: Email address is already taken"
	// ErrPasswordTooShort on update and create
	ErrPasswordTooShort modelError = "models: Password should be at least 8 characters long"
	// ErrEmailRequired when creating user - email is required
	ErrEmailRequired modelError = "models: Email is required"
	// ErrEmailInvalid if email is not valid
	ErrEmailInvalid modelError = "models: Email address is not valid"
	// ErrPasswordRequired for create and update
	ErrPasswordRequired modelError = "models: Password is required"
	// ErrTitleRequired is returned when title is not provided when gallery is creating
	ErrTitleRequired modelError = "models: title of the gallery is required"

	// ErrIDInvalid when invalid ID provided (Delete method)
	ErrIDInvalid privateError = "models: ID must be > 0"
	// ErrUserIDRequired is returned when userID is not provided when gallery is creating
	ErrUserIDRequired privateError = "models: UserID is required"
	// ErrRememberRequired is returned when a create or update is attempted without a user remember token hash
	ErrRememberRequired privateError = "models: remember token is required"
	// ErrRememberTooShort is returned when the remember token is less than 32 bytes
	ErrRememberTooShort privateError = "models: Remember token must be at least 32 bytes"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	return s
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
