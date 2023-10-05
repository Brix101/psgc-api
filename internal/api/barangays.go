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
	FilteredBrgy = "filteredBrgy"
)

type brgyResource domain.Resource

// Routes creates a REST router for the barangays resource
func (rs brgyResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.With(paginate).Get("/", rs.List) // GET /barangays - read a list of barangays

	r.Route("/{psgcCode}", func(r chi.Router) {
		r.Use(rs.BarangayCtx) // lets have a barangays map, and lets actually load/manipulate
		r.Get("/", rs.Get)    // GET /barangays/{psgcCode} - read a single todo by :id
	})

	return r
}

func (rs brgyResource) BarangayCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		psgcCode := chi.URLParam(r, "psgcCode")                // Get the {psgcCode} from the route
		filteredItem := make(chan generator.GeographicArea, 1) // Create a channel with buffer size 1

		// Create a goroutine to filter the data
		go func() {
			defer close(filteredItem)
			for _, item := range rs.Barangays {
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

		ctx := context.WithValue(r.Context(), FilteredBrgy, item)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ShowBarangays godoc
//
//	@Summary		Show list of Barangays
//	@Description	get barangays
//	@Tags			barangays
//	@Accept			json
//	@Produce		json
//	@Param			query	query		PaginationParams	false	"Pagination and filter parameters"
//	@Success		200		{object}	PaginatedResponse
//	@Failure		400		{object}	string	"Bad Request"
//	@Failure		500		{object}	string	"Internal Server Error"
//	@Router			/barangays [get]
func (rs brgyResource) List(w http.ResponseWriter, r *http.Request) {
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
	response := createPaginatedResponse(rs.Barangays, pageParams)

	// Marshal and send the response
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (rs brgyResource) Get(w http.ResponseWriter, r *http.Request) {
	// Get the context from the request
	ctx := r.Context()

	item, ok := ctx.Value(FilteredBrgy).(generator.GeographicArea)
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
