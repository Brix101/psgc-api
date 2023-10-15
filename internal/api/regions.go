package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Brix101/psgc-tool/internal/domain"
	"github.com/Brix101/psgc-tool/internal/util"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type (
	RegCtx      struct{}
	regResource struct {
		logger  *zap.Logger
		regRepo domain.RegionRepository
	}
)

// Routes creates a REST router for the regions resource
func (rs regResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.With(util.Paginate).Get("/", rs.List) // GET /regions - read a list of regions

	r.Route("/{psgc_code}", func(r chi.Router) {
		r.Use(rs.RegionCtx) // lets have a regions map, and lets actually load/manipulate
		r.Get("/", rs.Get)  // GET /regions/{psgc_code} - read a single todo by :id
	})
	return r
}

func (rs regResource) RegionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		psgcCode := chi.URLParam(r, "psgc_code") // Get the {psgc_code} from the route

		item, err := rs.regRepo.GetById(ctx, psgcCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		ctx = context.WithValue(ctx, RegCtx{}, item)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ShowRegions godoc
//
//	@Summary		Show list of Regions
//	@Description	get Regions
//	@Tags			Regions
//	@Accept			json
//	@Produce		json
//	@Param			query	query		PaginationParams	false	"Pagination and filter parameters"
//	@Success		200		{object}	PaginatedRegion
//	@Failure		400		{object}	string	"Bad Request"
//	@Failure		500		{object}	string	"Internal Server Error"
//	@Router			/regions [get]
func (rs regResource) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pageParams, ok := ctx.Value(util.PaginateCtx{}).(domain.PaginationParams)
	if !ok {
		http.Error(w, "Pagination information not found", http.StatusBadRequest)
		return
	}

	data, err := rs.regRepo.GetAll(ctx, pageParams)
	if err != nil {
		rs.logger.Error("failed to fetch regions from database", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

// ShowRegions godoc
//
//	@Summary		Show a Region
//	@Description	get string by PsgcCode
//	@Tags			Regions
//	@Accept			json
//	@Produce		json
//	@Param			psgc_code	path		string true	"Region PsgcCode"
//	@Success		200			{object}	domain.Region
//	@Failure		400			{object}	string	"Bad Request"
//	@Failure		404			{object}	string	"Item Not Found"
//	@Failure		500			{object}	string	"Internal Server Error"
//	@Router			/regions/{psgc_code} [get]
func (rs regResource) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	item, ok := ctx.Value(RegCtx{}).(domain.Region)
	if !ok {
		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	res, err := json.Marshal(item)
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
