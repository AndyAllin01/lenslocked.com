package controllers

import (
	"lenslocked.com/views"
)

func NewStatic() *Static {
	return &Static{
		/*		Home:    views.NewView("bootstrap", "views/static/home.gohtml"),
				Contact: views.NewView("bootstrap", "views/static/contact.gohtml"),*/
		Home:    views.NewView("bootstrap", "static/home"),
		Contact: views.NewView("bootstrap", "static/contact"),
		FAQ:     views.NewView("bootstrap", "static/faq"),
	}
}

type Static struct {
	Home    *views.View
	Contact *views.View
	FAQ     *views.View
}
