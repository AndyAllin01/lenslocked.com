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
)

//NewUsers creates a newusers controller
//panics if incorrect parse so only use
//during setup
func NewGalleries(gs models.GalleryService, r *mux.Router) *Galleries {
	return &Galleries{
		New:      views.NewView("bootstrap", "galleries/new"),
		ShowView: views.NewView("bootstrap", "galleries/show"),
		gs:       gs,
		r:        r,
	}
}

type Galleries struct {
	New      *views.View
	ShowView *views.View
	gs       models.GalleryService
	r        *mux.Router
}

type GalleryForm struct {
	Title string `schema: "title"`
}

//GET/gallerries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery id", http.StatusNotFound)
	}
	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Unknown error", http.StatusInternalServerError)
		}
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
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
	url, err := g.r.Get(ShowGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		//find an appropriate action here
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
	fmt.Fprint(w, gallery)
}
