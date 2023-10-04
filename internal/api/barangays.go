package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type barangaysResource struct{}

// Routes creates a REST router for the barangays resource
func (rs barangaysResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", rs.List)    // GET /barangays - read a list of barangays
	r.Post("/", rs.Create) // POST /barangays - create a new todo and persist it
	r.Put("/", rs.Delete)

	r.Route("/{id}", func(r chi.Router) {
		// r.Use(rs.TodoCtx) // lets have a barangays map, and lets actually load/manipulate
		r.Get("/", rs.Get)       // GET /barangays/{id} - read a single todo by :id
		r.Put("/", rs.Update)    // PUT /barangays/{id} - update a single todo by :id
		r.Delete("/", rs.Delete) // DELETE /barangays/{id} - delete a single todo by :id
		r.Get("/sync", rs.Sync)
	})

	return r
}

func (rs barangaysResource) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("barangays list of stuff.."))
}

func (rs barangaysResource) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("barangays create"))
}

func (rs barangaysResource) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo get"))
}

func (rs barangaysResource) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo update"))
}

func (rs barangaysResource) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo delete"))
}

func (rs barangaysResource) Sync(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo sync"))
}
