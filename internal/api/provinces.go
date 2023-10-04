package api

import (
	"encoding/json"
	"net/http"

	"github.com/Brix101/psgc-api/pkg/generator"
	"github.com/go-chi/chi/v5"
)

type provincesResource struct {
	Provinces []generator.GeographicArea
}

// Routes creates a REST router for the provinces resource
func (rs provincesResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", rs.List)    // GET /provinces - read a list of provinces
	r.Post("/", rs.Create) // POST /provinces - create a new todo and persist it
	r.Put("/", rs.Delete)

	r.Route("/{id}", func(r chi.Router) {
		// r.Use(rs.TodoCtx) // lets have a provinces map, and lets actually load/manipulate
		r.Get("/", rs.Get)       // GET /provinces/{id} - read a single todo by :id
		r.Put("/", rs.Update)    // PUT /provinces/{id} - update a single todo by :id
		r.Delete("/", rs.Delete) // DELETE /provinces/{id} - delete a single todo by :id
		r.Get("/sync", rs.Sync)
	})

	return r
}

func (rs provincesResource) List(w http.ResponseWriter, r *http.Request) {
	d := rs.Provinces

	res, _ := json.Marshal(d)

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (rs provincesResource) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("provinces create"))
}

func (rs provincesResource) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo get"))
}

func (rs provincesResource) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo update"))
}

func (rs provincesResource) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo delete"))
}

func (rs provincesResource) Sync(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo sync"))
}
