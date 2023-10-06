package api

import (
	"context"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/Brix101/psgc-api/internal/domain"
)

const (
	DefaultPage         = 1
	DefaultPerPage      = 1000
	PaginationParamsKey = "paginationParams"
)

func createPaginatedResponse(
	data interface{},
	PaginationParams domain.PaginationParams,
) domain.PaginatedResponse {
	page := PaginationParams.Page
	perPage := PaginationParams.PerPage
	filter := PaginationParams.Filter

	// Type assertion to convert the data interface{} to []domain.Masterlist
	dataList, ok := data.([]domain.Masterlist)
	if !ok { // Return an empty response if data is not of the expected type
		return domain.PaginatedResponse{}
	}

	// Create a channel for sending filtered items and receiving filtered items
	filterChan := make(chan domain.Masterlist)

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Define a function for parallel filtering
	filterFunc := func(item domain.Masterlist) {
		defer wg.Done()
		v := reflect.ValueOf(item)
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			value, ok := field.Interface().(string)
			if !ok || len(value) <= 0 {
				continue
			}

			if strings.Contains(strings.ToLower(value), strings.ToLower(filter)) {
				filterChan <- item // Pass teh filtered item to the chanel
				break
			}
		}
	}

	// Start goroutines for filtering
	for _, item := range dataList {
		wg.Add(1)
		go filterFunc(item)
	}

	// Close the output channel when all goroutines are done
	go func() {
		wg.Wait()
		close(filterChan)
	}()

	totalItems := len(dataList)
	slicedData := []domain.Masterlist{}
	itemCount := 0

	// Iterate through filtered results and perform pagination
	for filteredItem := range filterChan {
		itemCount++
		if itemCount > (page-1)*perPage && itemCount <= page*perPage {
			slicedData = append(slicedData, filteredItem)
		}
	}

	totalPages := (itemCount + perPage - 1) / perPage

	sort.Slice(slicedData, func(i, j int) bool {
		return slicedData[i].Name < slicedData[j].Name
	})

	return domain.PaginatedResponse{
		MetaData: domain.MetaData{
			Page:       page,
			TotalPages: totalPages,
			PerPage:    perPage,
			TotalItems: totalItems,
			ItemCount:  itemCount,
		},
		Data: slicedData,
	}
}

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

		if filterParam == "" {
			filterParam = ""
		}

		// Create a context with pagination information and pass it down the chain
		ctx := context.WithValue(r.Context(), PaginationParamsKey, domain.PaginationParams{
			Page:    page,
			PerPage: perPage,
			Filter:  filterParam,
		})

		// Serve the request with the modified context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
