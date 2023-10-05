package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Brix101/psgc-api/internal/generator"
	"github.com/go-chi/chi/v5"
)

type provincesResource struct {
	Provinces []generator.GeographicArea
}

// Routes creates a REST router for the provinces resource
func (rs provincesResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.With(paginate).Get("/", rs.List) // GET /provinces - read a list of provinces

	r.Route("/{psgcCode}", func(r chi.Router) {
		r.Use(rs.ProvinceCtx) // lets have a provinces map, and lets actually load/manipulate
		r.Get("/", rs.Get)    // GET /provinces/{psgcCode} - read a single todo by :id
	})

	return r
}

func (rs provincesResource) ProvinceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		psgcCode := chi.URLParam(r, "psgcCode") // Get the {psgcCode} from the route
        
		fmt.Println(psgcCode)
		// Your middleware logic here, for example, loading/manipulating data.
		// ...

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// ShowProvinces godoc
//
//	@Summary		Show list of provinces
//	@Description	get provinces
//	@Tags			provinces
//	@Accept			json
//	@Produce		json
//	@Param			query	query		PaginationParams	false	"Pagination and filter parameters"
//	@Success		200		{object}	PaginatedResponse
//	@Failure		400		{object}	string	"Bad Request"
//	@Failure		500		{object}	string	"Internal Server Error"
//	@Router			/provinces [get]
func (rs provincesResource) List(w http.ResponseWriter, r *http.Request) {
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
	response := createPaginatedResponse(rs.Provinces, pageParams)

	// Marshal and send the response
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (rs provincesResource) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("province get"))
}
