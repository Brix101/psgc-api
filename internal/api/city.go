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
	CityCtx      struct{}
	cityResource struct {
		logger       *zap.Logger
		cityMuniRepo domain.CityMuniRepository
	}
)

// Routes creates a REST router for the cities resource
func (rs cityResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.With(util.Paginate).Get("/", rs.List) // GET /city - read a list of cities

	r.Route("/{psgc_code}", func(r chi.Router) {
		r.Use(rs.CitiesCtx) // lets have a cities map, and lets actually load/manipulate
		r.Get("/", rs.Get)  // GET /city/{psgc_code} - read a single todo by :id
	})

	return r
}

func (rs cityResource) CitiesCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		psgcCode := chi.URLParam(r, "psgc_code") // Get the {psgc_code} from the route

		item, err := rs.cityMuniRepo.GetCityById(ctx, psgcCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		ctx = context.WithValue(ctx, CityCtx{}, item)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ShowCities godoc
//
//	@Summary		Show list of Cities
//	@Description	get Cities
//	@Tags			Cities
//	@Accept			json
//	@Produce		json
//	@Param			query	query		PaginationParams	false	"Pagination and filter parameters"
//	@Success		200		{object}	PaginatedCityMuni
//	@Failure		400		{object}	string	"Bad Request"
//	@Failure		500		{object}	string	"Internal Server Error"
//	@Router			/cities [get]
func (rs cityResource) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pageParams, ok := ctx.Value(util.PaginateCtx{}).(domain.PaginationParams)
	if !ok {
		http.Error(w, "Pagination information not found", http.StatusBadRequest)
		return
	}

	data, err := rs.cityMuniRepo.GetAllCity(ctx, pageParams)
	if err != nil {
		rs.logger.Error("failed to fetch cities from database", zap.Error(err))
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

// ShowCities godoc
//
//	@Summary		Show a City
//	@Description	get string by PsgcCode
//	@Tags			Cities
//	@Accept			json
//	@Produce		json
//	@Param			psgc_code	path		string true	"City PsgcCode"
//	@Success		200			{object}	domain.CityMuni
//	@Failure		400			{object}	string	"Bad Request"
//	@Failure		404			{object}	string	"Item Not Found"
//	@Failure		500			{object}	string	"Internal Server Error"
//	@Router			/cities/{psgc_code} [get]
func (rs cityResource) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	item, ok := ctx.Value(CityCtx{}).(domain.CityMuni)
	if !ok {

		http.Error(w, "Item not found", http.StatusNotFound)
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
