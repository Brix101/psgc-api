package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type regionsResource struct{}

// Routes creates a REST router for the regions resource
func (rs regionsResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", rs.List)    // GET /regions - read a list of regions
	r.Post("/", rs.Create) // POST /regions - create a new todo and persist it
	r.Put("/", rs.Delete)

	r.Route("/{id}", func(r chi.Router) {
		// r.Use(rs.TodoCtx) // lets have a regions map, and lets actually load/manipulate
		r.Get("/", rs.Get)       // GET /regions/{id} - read a single todo by :id
		r.Put("/", rs.Update)    // PUT /regions/{id} - update a single todo by :id
		r.Delete("/", rs.Delete) // DELETE /regions/{id} - delete a single todo by :id
		r.Get("/sync", rs.Sync)
	})

	return r
}

func (rs regionsResource) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("regions list of stuff.."))
}

func (rs regionsResource) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("regions create"))
}

func (rs regionsResource) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo get"))
}

func (rs regionsResource) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo update"))
}

func (rs regionsResource) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo delete"))
}

func (rs regionsResource) Sync(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo sync"))
}
