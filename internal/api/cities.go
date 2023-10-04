package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type citiesResource struct{}

// Routes creates a REST router for the cities resource
func (rs citiesResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", rs.List)    // GET /cities - read a list of cities
	r.Post("/", rs.Create) // POST /cities - create a new todo and persist it
	r.Put("/", rs.Delete)

	r.Route("/{id}", func(r chi.Router) {
		// r.Use(rs.TodoCtx) // lets have a cities map, and lets actually load/manipulate
		r.Get("/", rs.Get)       // GET /cities/{id} - read a single todo by :id
		r.Put("/", rs.Update)    // PUT /cities/{id} - update a single todo by :id
		r.Delete("/", rs.Delete) // DELETE /cities/{id} - delete a single todo by :id
		r.Get("/sync", rs.Sync)
	})

	return r
}

func (rs citiesResource) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("cities list of stuff.."))
}

func (rs citiesResource) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("cities create"))
}

func (rs citiesResource) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo get"))
}

func (rs citiesResource) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo update"))
}

func (rs citiesResource) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo delete"))
}

func (rs citiesResource) Sync(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo sync"))
}
