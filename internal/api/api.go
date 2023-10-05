package api

import (
	"context"
	"fmt"
	"net/http"

	_ "github.com/Brix101/psgc-api/docs"
	"github.com/Brix101/psgc-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type api struct {
	logger *zap.Logger

	barangayApi   barangaysResource
	cityApi       citiesResource
	provinceApi   provincesResource
	regionApi     regionsResource
	masterlistApi masterlistResource
}

func NewAPI(_ context.Context, logger *zap.Logger) *api {
	barangays := service.GetBarangays(logger)
	cities := service.GetCities(logger)
	provinces := service.GetProvinces(logger)
	regions := service.GetRegions(logger)
	masterlist := service.GetMasterlist(logger)

	return &api{
		logger: logger,

		barangayApi: barangaysResource{
			Barangays: barangays,
		},
		cityApi: citiesResource{
			Cities: cities,
		},
		provinceApi: provincesResource{
			Provinces: provinces,
		},
		regionApi: regionsResource{
			Regions: regions,
		},
		masterlistApi: masterlistResource{
			Masterlist: masterlist,
		},
	}
}

func (a *api) Server(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.Routes(),
	}
}

func (a *api) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Heartbeat("/ping"))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"), // The url pointing to API definition
	))

	r.Route("/api", func(r chi.Router) {
		r.Mount("/barangays", a.barangayApi.Routes())
		r.Mount("/cities", a.cityApi.Routes())
		r.Mount("/provinces", a.provinceApi.Routes())
		r.Mount("/regions", a.regionApi.Routes())
		r.Mount("/masterlist", a.masterlistApi.Routes())
	})

	// Catch-all route for 404 errors, redirect to Swagger
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/index.html", http.StatusFound)
	})

	return r
}
