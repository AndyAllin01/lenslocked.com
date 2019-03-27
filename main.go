package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	fmt.Fprint(w, "<h1>First Calhoun Web Page</h1>")
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	fmt.Fprintf(w, "to get in touc please email : "+
		"to <a href=\"mailto:support@lenslocked.com\">"+
		"support@lenslocked.com</a>.")
}
func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	fmt.Fprintf(w, "you dare to question the almighty holio?")
}

func notFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "the almighty holio says \n\n\n        N O T    F O U N D")
}

func main() {
	fmt.Println("RUNNING")
	var h http.Handler = http.HandlerFunc(notFound)
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/faq", faq)
	r.NotFoundHandler = h
	http.ListenAndServe(":8080", r)
}
