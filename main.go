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
	if err != nil {
		panic(err)
	}
	//	must(err)
	//TODO - fix this
	defer services.Close()
	services.AutoMigrate()
	//services.DestructiveReset()

	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, r)

	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{}

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.FAQ).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	//r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	//gallery routes
	r.Handle("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET")
	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")

	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)
	fmt.Println("STARTING SERVER ######")
	http.ListenAndServe(":8080", userMw.Apply(r))
}

/*func must(err error) {
	if err != nil {
		panic(err)
	}
}*/
