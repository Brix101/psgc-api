package api

import (
	"encoding/json"
	"net/http"

	"github.com/Brix101/psgc-api/pkg/generator"
	"github.com/go-chi/chi/v5"
)

type citiesResource struct {
	Cities []generator.GeographicArea
}

// Routes creates a REST router for the cities resource
func (rs citiesResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.With(paginate).Get("/", rs.List) // GET /cities - read a list of cities

	r.Route("/{id}", func(r chi.Router) {
		// r.Use(rs.TodoCtx) // lets have a cities map, and lets actually load/manipulate
		r.Get("/", rs.Get) // GET /cities/{id} - read a single todo by :id
	})

	return r
}

func (rs citiesResource) List(w http.ResponseWriter, r *http.Request) {
	// Get the context from the request
	ctx := r.Context()

	paginationInfo, ok := ctx.Value("pagination").(PaginationInfo)
	if !ok {
		// Handle the case where pagination information is not found in the context
		// You can choose to use default values or return an error response.
		http.Error(w, "Pagination information not found", http.StatusBadRequest)
		return
	}

	// Create the PaginatedResponse using the retrieved data and pagination information
	response := createPaginatedResponse(rs.Cities, paginationInfo)

	// Marshal and send the response
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (rs citiesResource) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("city get"))
}
