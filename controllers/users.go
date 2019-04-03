package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/views"
)

//NewUsers creates a newusers controller
//panics if incorrect parse so only use
//during setup
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

type Users struct {
	NewView *views.View
}

//New renders form for new user account
//GET/signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

//Create processes the signup form when submitted.
//Create a new user accout
//POST/signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, r.PostForm["email"])
	fmt.Fprintln(w, r.PostFormValue("email"))
	fmt.Fprintln(w, r.PostForm["password"])
	fmt.Fprintln(w, r.PostFormValue("password"))
	fmt.Fprintln(w, "temp response")
}
