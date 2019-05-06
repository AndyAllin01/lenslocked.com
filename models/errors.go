package models

import "strings"

const (
	ErrNotFound          modelError = "models: resource not found"
	ErrIDInvalid         modelError = "models: ID provided invalid"
	ErrInvalidEmail      modelError = "models: incorrect email provided"
	ErrPasswordIncorrect modelError = "models: incorrect password provided"
	ErrEmailRequired     modelError = "Email address is required"
	ErrEmailInvalid      modelError = "Email address is not valid"
	ErrEmailTaken        modelError = "models: email address is already taken"
	ErrPasswordTooShort  modelError = "models: password must be at least 8 characters"
	ErrPasswordRequired  modelError = "models: password is required"
	ErrRememberTooShort  modelError = "models: remember token must be at least 32 bytes"
	ErrRememberRequired  modelError = "models: invlid remember token hassh"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	//capitalize only first char of sentence
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}
