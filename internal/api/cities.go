package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Brix101/psgc-api/internal/generator"
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

	r.Route("/{psgcCode}", func(r chi.Router) {
		r.Use(rs.CitiesCtx) // lets have a cities map, and lets actually load/manipulate
		r.Get("/", rs.Get)  // GET /cities/{psgcCode} - read a single todo by :id
	})

	return r
}

func (rs citiesResource) CitiesCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		psgcCode := chi.URLParam(r, "psgcCode") // Get the {psgcCode} from the route
        
		fmt.Println(psgcCode)
		// Your middleware logic here, for example, loading/manipulating data.
		// ...

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// ShowCities godoc
//
//	@Summary		Show list of cities
//	@Description	get cities
//	@Tags			cities
//	@Accept			json
//	@Produce		json
//	@Param			query	query		PaginationParams	false	"Pagination and filter parameters"
//	@Success		200		{object}	PaginatedResponse
//	@Failure		400		{object}	string	"Bad Request"
//	@Failure		500		{object}	string	"Internal Server Error"
//	@Router			/cities [get]
func (rs citiesResource) List(w http.ResponseWriter, r *http.Request) {
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
	response := createPaginatedResponse(rs.Cities, pageParams)

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
