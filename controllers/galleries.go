package controllers

import (
	"lenslocked.com/models"
	"lenslocked.com/views"
)

//NewUsers creates a newusers controller
//panics if incorrect parse so only use
//during setup
func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

type Galleries struct {
	New *views.View
	gs  models.GalleryService
}
