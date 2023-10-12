package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Brix101/psgc-tool/internal/domain"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	ProvCtx = "ProvCtx"
)

type provResource struct {
	logger   *zap.Logger
	provRepo domain.ProvinceRepository
}

// Routes creates a REST router for the provinces resource
func (rs provResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.With(paginate).Get("/", rs.List) // GET /provinces - read a list of provinces

	r.Route("/{psgc_code}", func(r chi.Router) {
		r.Use(rs.ProvinceCtx) // lets have a provinces map, and lets actually load/manipulate
		r.Get("/", rs.Get)    // GET /provinces/{psgc_code} - read a single todo by :id
	})

	return r
}

func (rs provResource) ProvinceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		psgcCode := chi.URLParam(r, "psgc_code") // Get the {psgc_code} from the route

		item, err := rs.provRepo.GetById(ctx, psgcCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		ctx = context.WithValue(ctx, ProvCtx, item)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ShowProvinces godoc
//
//	@Summary		Show list of Provinces
//	@Description	get Provinces
//	@Tags			Provinces
//	@Accept			json
//	@Produce		json
//	@Param			query	query		PaginationParams	false	"Pagination and filter parameters"
//	@Success		200		{object}	PaginatedProvince
//	@Failure		400		{object}	string	"Bad Request"
//	@Failure		500		{object}	string	"Internal Server Error"
//	@Router			/provinces [get]
func (rs provResource) List(w http.ResponseWriter, r *http.Request) {
	// Get the context from the request
	ctx := r.Context()

	pageParams, ok := ctx.Value(PaginationParamsKey).(domain.PaginationParams)
	if !ok {
		// Handle the case where pagination information is not found in the context
		// You can choose to use default values or return an error response.
		http.Error(w, "Pagination information not found", http.StatusBadRequest)
		return
	}

	data, err := rs.provRepo.GetList(ctx, pageParams)
	if err != nil {
		rs.logger.Error("failed to fetch provinces from database", zap.Error(err))
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

// ShowProvinces godoc
//
//	@Summary		Show a Province
//	@Description	get string by PsgcCode
//	@Tags			Provinces
//	@Accept			json
//	@Produce		json
//	@Param			psgc_code	path		string true	"Province PsgcCode"
//	@Success		200			{object}	domain.Province
//	@Failure		400			{object}	string	"Bad Request"
//	@Failure		400			{object}	string	"Item Not Found"
//	@Failure		500			{object}	string	"Internal Server Error"
//	@Router			/provinces/{psgc_code} [get]
func (rs provResource) Get(w http.ResponseWriter, r *http.Request) {
	// Get the context from the request
	ctx := r.Context()

	item, ok := ctx.Value(ProvCtx).(domain.Province)
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
