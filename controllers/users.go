package controllers

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"

	"lenslocked.com/views"
)

//NewUsers creates a newusers controller
//panics if incorrect parse so only use
//during setup
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
	Name     string `schema: "name"`
	Email    string `schema: "email"`
	Password string `schema: "password"`
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
	signIn(w, &user)
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

//Login used to verify user credentials and if correct log in user
//POST / login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid email address")
		case models.ErrInvalidPassword:
			fmt.Fprintln(w, "Invalid password provided")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	signIn(w, user)
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

//sign in the supplied user via cookie
func signIn(w http.ResponseWriter, user *models.User) {
	cookie := http.Cookie{
		Name:  "email",
		Value: user.Email,
	}
	http.SetCookie(w, &cookie)
}

func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("email")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Email is ", cookie.Value)
	fmt.Fprintln(w, cookie)

}
