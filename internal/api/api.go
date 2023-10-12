package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/Brix101/psgc-tool/docs"
	"github.com/Brix101/psgc-tool/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type api struct {
	logger *zap.Logger

	bgyApi  bryResource
	cityApi  citiMuniResource
	provApi  provResource
	regApi   regResource
	mListApi mListResource
}

func NewAPI(_ context.Context, logger *zap.Logger, db *sql.DB) *api {
	mListRepo := repository.NewDBMasterlist(db)
	regRepo := repository.NewDBRegion(db)
	provRepo := repository.NewDBProvince(db)
	brgyRepo := repository.NewDBBarangay(db)
	cityMuniRepo := repository.NewDBCityMuni(db)

	return &api{
		logger: logger,

		bgyApi: bryResource{
			logger:   logger,
			bgyRepo: brgyRepo,
		},
		cityApi: citiMuniResource{
			logger:       logger,
			cityMuniRepo: cityMuniRepo,
		},
		provApi: provResource{
			logger:   logger,
			provRepo: provRepo,
		},
		regApi: regResource{
			logger:  logger,
			regRepo: regRepo,
		},
		mListApi: mListResource{
			logger:    logger,
			mListRepo: mListRepo,
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

	// Rate limit by IP and URL path (aka endpoint)
	r.Use(httprate.Limit(
		500,           // requests
		1*time.Minute, // per duration
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))

	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"), // The url pointing to API definition
	))

	r.Route("/api", func(r chi.Router) {
		r.Mount("/barangays", a.bgyApi.Routes())
		r.Mount("/citi_muni", a.cityApi.Routes())
		r.Mount("/provinces", a.provApi.Routes())
		r.Mount("/regions", a.regApi.Routes())
		r.Mount("/masterlist", a.mListApi.Routes())
	})

	// Catch-all route for 404 errors, redirect to Swagger
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/index.html", http.StatusFound)
	})

	return r
}
