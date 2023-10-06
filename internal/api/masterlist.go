package api

import (
	"encoding/json"
	"net/http"

	"github.com/Brix101/psgc-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type mListResource apiResource

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

	pageParams, ok := ctx.Value(PaginationParamsKey).(domain.PaginationParams)
	if !ok {
		// Handle the case where pagination information is not found in the context
		// You can choose to use default values or return an error response.
		http.Error(w, "Pagination information not found", http.StatusBadRequest)
		return
	}

	data, err := rs.Repo.GetList(ctx, pageParams)
	if err != nil {
		rs.logger.Error("failed to fetch masterlist from database", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal and send the response
	res, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
