package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"lenslocked.com/context"

	"lenslocked.com/models"
	"lenslocked.com/views"
)

const (
	ShowGallery = "show_gallery"
	EditGallery = "edit_gallery"
)

//NewUsers creates a newusers controller
//panics if incorrect parse so only use
//during setup
func NewGalleries(gs models.GalleryService, r *mux.Router) *Galleries {
	return &Galleries{
		New:       views.NewView("bootstrap", "galleries/new"),
		ShowView:  views.NewView("bootstrap", "galleries/show"),
		EditView:  views.NewView("bootstrap", "galleries/edit"),
		IndexView: views.NewView("bootstrap", "galleries/index"),
		gs:        gs,
		r:         r,
	}
}

type Galleries struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	gs        models.GalleryService
	r         *mux.Router
}

type GalleryForm struct {
	Title string `schema: "title"`
}

//GET/galleries
func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleries, err := g.gs.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "No gallery for this user", http.StatusInternalServerError)
	}

	var vd views.Data
	vd.Yield = galleries
	//	fmt.Fprintln(w, galleries)
	g.IndexView.Render(w, vd)

}

//GET/gallerries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
	//	fmt.Fprintln(w, gallery)
}

//GET/galleries/id:/edit
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())

	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(w, vd)
}

//POST/galleries/id:/update
func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())

	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var form GalleryForm
	var vd views.Data
	vd.Yield = gallery

	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.EditView.Render(w, vd)
		return
	}
	gallery.Title = form.Title
	//	g.EditView.Render(w, vd)
	err = g.gs.Update(gallery)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, vd)
		return
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Gallery successfully updated",
	}
	g.EditView.Render(w, vd)
	//	fmt.Fprintln(w, gallery)
}

//Create processes the signup form when submitted.
//Create a new user accout
//POST/galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	fmt.Println("Create got the user : ", user)

	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}

	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	url, err := g.r.Get(EditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		//find an appropriate action here
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
	fmt.Fprint(w, gallery)
}

//POST/galleries/id:/delete
func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())

	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	err = g.gs.Delete(gallery.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = gallery
		g.EditView.Render(w, vd)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g *Galleries) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery id", http.StatusNotFound)
		return nil, err
	}
	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Unknown error", http.StatusInternalServerError)
		}
		return nil, err
	}
	return gallery, nil
}
