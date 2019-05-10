package main

import (
	"fmt"
	"net/http"

	"lenslocked.com/middleware"

	"lenslocked.com/models"

	"lenslocked.com/controllers"

	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "andya"
	password = "password"
	dbname   = "lenslocked_test"
)

func main() {
	fmt.Println("RUNNING")

	psqlInfo := "postgres://bond:password@localhost/lenslocked_dev?sslmode=disable"

	services, err := models.NewServices(psqlInfo)
	must(err)
	//TODO - fix this
	defer services.Close()
	services.AutoMigrate()
	//services.DestructiveReset()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery)
	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.FAQ).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	//gallery routes
	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	//r.Handle("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	fmt.Println("STARTING SERVER ######")
	http.ListenAndServe(":8080", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
