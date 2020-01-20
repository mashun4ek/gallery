package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/mashun4ek/webdevcalhoun/gallery/controllers"
	"github.com/mashun4ek/webdevcalhoun/gallery/middleware"
	"github.com/mashun4ek/webdevcalhoun/gallery/models"
	"github.com/mashun4ek/webdevcalhoun/gallery/rand"
)

// temporary here
const (
	host     = "localhost"
	port     = 5432
	user     = "mariaker"
	password = "nopassword"
	dbname   = "gallery"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	services, err := models.NewServices(psqlInfo)
	must(err)
	defer services.Close()
	// services.DestructiveReset()
	services.AutoMigrate()

	r := mux.NewRouter()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	// TODO update it to be config variable
	isProd := false
	b, err := rand.Bytes(32)
	must(err)
	// 32 bytes generated auth key
	csrfMw := csrf.Protect(b, csrf.Secure(isProd))

	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{User: userMw}

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/signup", usersC.NewView).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	// Image routes
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images", imageHandler))

	// Gallery routes
	r.Handle("/galleries", requireUserMw.RUMFn(galleriesC.Index)).Methods("GET")
	r.Handle("/galleries/new", requireUserMw.RequireUserMiddleware(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.RUMFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.RUMFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.RUMFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.RUMFn(galleriesC.Delete)).Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.RUMFn(galleriesC.ImageUpload)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMw.RUMFn(galleriesC.ImageDelete)).Methods("POST")

	// Assets
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	fmt.Println("Starting the server on :8080...")
	http.ListenAndServe(":8080", csrfMw(userMw.RequireUserMiddleware(r)))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
