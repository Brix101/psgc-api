package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Brix101/psgc-tool/internal/domain"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type (
	MunicipalityCtx struct{}
	munResource     struct {
		logger       *zap.Logger
		cityMuniRepo domain.CityMuniRepository
	}
)

// Routes creates a REST router for the municipalities resource
func (rs munResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.With(paginate).Get("/", rs.List) // GET /municipality - read a list of municipalities

	r.Route("/{psgc_code}", func(r chi.Router) {
		r.Use(rs.MunicipalitiesCtx) // lets have a municipalities map, and lets actually load/manipulate
		r.Get("/", rs.Get)          // GET /municipality/{psgc_code} - read a single todo by :id
	})

	return r
}

func (rs munResource) MunicipalitiesCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		psgcCode := chi.URLParam(r, "psgc_code") // Get the {psgc_code} from the route

		item, err := rs.cityMuniRepo.GetMunicipalityById(ctx, psgcCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		ctx = context.WithValue(ctx, MunicipalityCtx{}, item)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ShowMunicipalities godoc
//
//	@Summary		Show list of Municipalities
//	@Description	get Municipalities
//	@Tags			Municipalities
//	@Accept			json
//	@Produce		json
//	@Param			query	query		PaginationParams	false	"Pagination and filter parameters"
//	@Success		200		{object}	PaginatedCityMuni
//	@Failure		400		{object}	string	"Bad Request"
//	@Failure		500		{object}	string	"Internal Server Error"
//	@Router			/municipalities [get]
func (rs munResource) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pageParams, ok := ctx.Value(PaginationParamsKey{}).(domain.PaginationParams)
	if !ok {
		http.Error(w, "Pagination information not found", http.StatusBadRequest)
		return
	}

	data, err := rs.cityMuniRepo.GetAllMunicipality(ctx, pageParams)
	if err != nil {
		rs.logger.Error("failed to fetch municipalities from database", zap.Error(err))
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

// ShowMunicipalities godoc
//
//	@Summary		Show a Municipality
//	@Description	get string by PsgcCode
//	@Tags			Municipalities
//	@Accept			json
//	@Produce		json
//	@Param			psgc_code	path		string true	"Municipality PsgcCode"
//	@Success		200			{object}	domain.CityMuni
//	@Failure		400			{object}	string	"Bad Request"
//	@Failure		400			{object}	string	"Item Not Found"
//	@Failure		500			{object}	string	"Internal Server Error"
//	@Router			/municipalities/{psgc_code} [get]
func (rs munResource) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	item, ok := ctx.Value(MunicipalityCtx{}).(domain.CityMuni)
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
