package main

import (
	"github.com/mashun4ek/webdevcalhoun/gallery/email"
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/mashun4ek/webdevcalhoun/gallery/controllers"
	"github.com/mashun4ek/webdevcalhoun/gallery/middleware"
	"github.com/mashun4ek/webdevcalhoun/gallery/models"
	"github.com/mashun4ek/webdevcalhoun/gallery/rand"
)

func main() {
	boolPtr := flag.Bool("prod", false, "Provide this flag in production. This ensures that a .config file is provided before the app starts.")
	flag.Parse()

	cfg := LoadConfig(*boolPtr)
	dbCfg := cfg.Database
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMAC),
		models.WithGallery(),
		models.WithImage(),
	)
	must(err)
	defer services.Close()
	// services.DestructiveReset()
	services.AutoMigrate()

	mgCfg := cfg.Mailgun
	emailer := email.NewClient(
		email.WithSender("MyGravitation Support", "support@sandbox0a473d6bb8c94265aa28ced67515768c.mailgun.org"),
		email.WithMailgun(mgCfg.domain. mgCfg.APIKey)
	)

	r := mux.NewRouter()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User, emailer)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	b, err := rand.Bytes(32)
	must(err)
	// 32 bytes generated auth key
	csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProd()))

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
	r.HandleFunc("/logout", requireUserMw.RUMFn(usersC.Logout)).Methods("POST")

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

	fmt.Printf("Starting the server on :%d...\n", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), csrfMw(userMw.RequireUserMiddleware(r)))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
