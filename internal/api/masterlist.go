package api

import (
	"encoding/json"
	"net/http"

	"github.com/Brix101/psgc-api/internal/domain"
	"github.com/go-chi/chi/v5"
)

type mListResource domain.Resource

// Routes creates a REST router for the masterlist resource
func (rs mListResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.With(paginate).Get("/", rs.List) // GET /masterlist - read a list of masterlist

	return r
}

// ShowMasterlist godoc
//
//	@Summary		Show list of masterlist
//	@Description	get masterlist
//	@Tags			masterlist
//	@Accept			json
//	@Produce		json
//	@Param			query	query		PaginationParams	false	"Pagination and filter parameters"
//	@Success		200		{object}	PaginatedResponse
//	@Failure		400		{object}	string	"Bad Request"
//	@Failure		500		{object}	string	"Internal Server Error"
//	@Router			/masterlist [get]
func (rs mListResource) List(w http.ResponseWriter, r *http.Request) {
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
	response := createPaginatedResponse(rs.Masterlist, pageParams)

	// Marshal and send the response
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
