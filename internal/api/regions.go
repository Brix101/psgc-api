package api

import (
	"encoding/json"
	"net/http"

	"github.com/Brix101/psgc-api/internal/generator"
	"github.com/go-chi/chi/v5"
)

type regionsResource struct {
	Regions []generator.GeographicArea
}

// Routes creates a REST router for the regions resource
func (rs regionsResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.With(paginate).Get("/", rs.List) // GET /regions - read a list of regions

	r.Route("/{id}", func(r chi.Router) {
		r.Use(rs.RegionCtx) // lets have a regions map, and lets actually load/manipulate
		r.Get("/", rs.Get)  // GET /regions/{id} - read a single todo by :id
	})

	return r
}

func (rs regionsResource) RegionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Your middleware logic here, for example, loading/manipulating data.
		// ...

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// ShowRegions godoc
//
//	@Summary		Show list of regions
//	@Description	get regions
//	@Tags			regions
//	@Accept			json
//	@Produce		json
//	@Param			query	query		PaginationParams	false	"Pagination and filter parameters"
//	@Success		200		{object}	PaginatedResponse
//	@Failure		400		{object}	string	"Bad Request"
//	@Failure		500		{object}	string	"Internal Server Error"
//	@Router			/regions [get]
func (rs regionsResource) List(w http.ResponseWriter, r *http.Request) {
	// Get the context from the request
	ctx := r.Context()

	pageParams, ok := ctx.Value(PaginationParamsKey).(PaginationParams)
	if !ok {
		// Handle the case where pagination information is not found in the context
		// You can choose to use default values or return an error response.
		http.Error(w, "Pagination information not found", http.StatusBadRequest)
		return
	}

	// Create the PaginatedResponse using the retrieved data and pagination information
	response := createPaginatedResponse(rs.Regions, pageParams)

	// Marshal and send the response
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (rs regionsResource) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo get"))
}
