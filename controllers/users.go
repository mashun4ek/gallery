package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"

	"github.com/mashun4ek/webdevcalhoun/gallery/models"
	"github.com/mashun4ek/webdevcalhoun/gallery/views"
)

// NewUser is used to create a new Users controller
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// renders the form where user can create a new user account
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// create a new user account POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}

func (u *Users) Update(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	dec := schema.NewDecoder()
	var form SignupForm
	if err := dec.Decode(&form, r.PostForm); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, form)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Login is used to verify provided email address and password and then
// log the user in if they are correct
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	switch err {
	case models.ErrNotFound:
		fmt.Fprintln(w, "Invalid email address.")
	case models.ErrInvalidPassword:
		fmt.Fprintln(w, "Invalid password provided.")
	case nil:
		fmt.Fprintln(w, user)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
