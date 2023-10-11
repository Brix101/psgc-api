package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/Brix101/psgc-api/docs"
	"github.com/Brix101/psgc-api/internal/domain"
	"github.com/Brix101/psgc-api/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type apiResource struct {
	logger *zap.Logger
	Repo   domain.MasterlistRepository
}

type api struct {
	logger *zap.Logger

	bgryApi  brgyResource
	cityApi  cityResource
	provApi  provResource
	regApi   regResource
	mListApi mListResource
}

func NewAPI(_ context.Context, logger *zap.Logger, db *sql.DB) *api {
	dbMl := repository.NewDBMasterlist(db)

	aRs := apiResource{
		logger: logger,
		Repo:   dbMl,
	}

	return &api{
		logger: logger,

		bgryApi:  brgyResource(aRs),
		cityApi:  cityResource(aRs),
		provApi:  provResource(aRs),
		regApi:   regResource(aRs),
		mListApi: mListResource(aRs),
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

	fmt.Println(1 * time.Minute)
	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"), // The url pointing to API definition
	))

	r.Route("/api", func(r chi.Router) {
		r.Mount("/barangays", a.bgryApi.Routes())
		r.Mount("/cities", a.cityApi.Routes())
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
