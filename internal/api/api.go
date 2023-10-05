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

	bgryApi  brgyResource
	cityApi  cityResource
	provApi  provResource
	regApi   regResource
	mListApi mListResource
}

func NewAPI(ctx context.Context, logger *zap.Logger) *api {
	ns := service.NewServices(ctx, logger)
	re := ns.GetResources()

	return &api{
		logger: logger,

		bgryApi:  brgyResource(*re),
		cityApi:  cityResource(*re),
		provApi:  provResource(*re),
		regApi:   regResource(*re),
		mListApi: mListResource(*re),
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
