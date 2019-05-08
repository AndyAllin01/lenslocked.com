package middleware

import (
	"fmt"
	"net/http"

	"lenslocked.com/models"
)

type RequireUser struct {
	models.UserService
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//if user is logged in call
		fmt.Println("NOW IM HERE")
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			//	http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := mw.UserService.ByRemember(cookie.Value)
		fmt.Println("NOW IM THERE")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			//		http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}
		fmt.Println("User found : ", user)
		next(w, r)
	})
}
