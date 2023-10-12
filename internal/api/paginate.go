package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Brix101/psgc-tool/internal/domain"
	"github.com/go-playground/validator/v10"
)

const (
	DefaultPage         = 1
	DefaultPerPage      = 1000
	PaginationParamsKey = "paginationParams"
)

func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the "page", "perPage", and "filter" query parameters from the URL
		pageParam := r.URL.Query().Get("page")
		perPageParam := r.URL.Query().Get("perPage")
		filterParam := r.URL.Query().Get("filter")

		// Parse the "page", "perPage", and "filter" query parameters
		page, err := strconv.Atoi(pageParam)
		if err != nil || page <= 0 {
			page = DefaultPage
		}

		perPage, err := strconv.Atoi(perPageParam)
		if err != nil || perPage <= 0 {
			perPage = DefaultPerPage
		}

		prms := domain.PaginationParams{
			Page:    page,
			PerPage: perPage,
			Filter:  filterParam,
		}

		validate := validator.New()
		if err := validate.Struct(prms); err != nil {
			validationErr, isValidationErr := err.(validator.ValidationErrors)
			if isValidationErr {
				fieldName := validationErr[0].Namespace()
				fieldName = strings.ToLower(fieldName[strings.LastIndex(fieldName, ".")+1:])
				message := fmt.Sprintf(
					"%s should be less than %s.",
					fieldName,
					validationErr[0].Param(),
				)

				http.Error(w, message, http.StatusBadRequest)
				return
			}

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create a context with pagination information and pass it down the chain
		ctx := context.WithValue(r.Context(),
			PaginationParamsKey,
			prms)

		// Serve the request with the modified context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
