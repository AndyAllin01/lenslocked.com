package main

import (
	"fmt"
	"net/http"

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

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.FAQ).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	fmt.Println("STARTING SERVER ######")
	http.ListenAndServe(":8080", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
