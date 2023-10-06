package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Brix101/psgc-api/internal/domain"
	"github.com/Brix101/psgc-api/internal/generator"
	"github.com/go-chi/chi/v5"
)
const (
	FilteredProv = "filteredProv"
)

type provResource domain.Resource

// Routes creates a REST router for the provinces resource
func (rs provResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.With(paginate).Get("/", rs.List) // GET /provinces - read a list of provinces

	r.Route("/{psgcCode}", func(r chi.Router) {
		r.Use(rs.ProvinceCtx) // lets have a provinces map, and lets actually load/manipulate
		r.Get("/", rs.Get)    // GET /provinces/{psgcCode} - read a single todo by :id
	})

	return r
}

func (rs provResource) ProvinceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		psgcCode := chi.URLParam(r, "psgcCode")                // Get the {psgcCode} from the route
		filteredItem := make(chan generator.GeographicArea, 1) // Create a channel with buffer size 1

		// Create a goroutine to filter the data
		go func() {
			defer close(filteredItem)
			for _, item := range rs.Provinces {
				if item.PsgcCode == psgcCode {
					filteredItem <- item
					return // Exit the goroutine once a matching item is found
				}
			}
		}()

		// Receive the filtered item from the channel
		item, found := <-filteredItem
		if !found {
			// No matching item found, return a custom "not found" message
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), FilteredProv, item)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ShowProvinces godoc
//
//	@Summary		Show list of provinces
//	@Description	get provinces
//	@Tags			provinces
//	@Accept			json
//	@Produce		json
//	@Param			query	query		PaginationParams	false	"Pagination and filter parameters"
//	@Success		200		{object}	PaginatedResponse
//	@Failure		400		{object}	string	"Bad Request"
//	@Failure		500		{object}	string	"Internal Server Error"
//	@Router			/provinces [get]
func (rs provResource) List(w http.ResponseWriter, r *http.Request) {
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
	response := createPaginatedResponse(rs.Provinces, pageParams)

	// Marshal and send the response
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

// ShowProvinces godoc
//	@Summary		Show a province
//	@Description	get string by PsgcCode
//	@Tags			provinces
//	@Accept			json
//	@Produce		json
//	@Param			psgcCode	path		string true	"Province PsgcCode"
//	@Success		200			{object}	generator.GeographicArea
//	@Failure		400			{object}	string	"Bad Request"
//	@Failure		400			{object}	string	"Item Not Found"
//	@Failure		500			{object}	string	"Internal Server Error"
//	@Router			/provinces/{psgcCode} [get]
func (rs provResource) Get(w http.ResponseWriter, r *http.Request) {
	// Get the context from the request
	ctx := r.Context()

	item, ok := ctx.Value(FilteredProv).(generator.GeographicArea)
	if !ok {
		// Handle the case where item is not found in the context
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	// Marshal and send the response
	res, err := json.Marshal(item)
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
