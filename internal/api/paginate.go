package api

import (
	"context"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/Brix101/psgc-api/internal/generator"
)

const (
	DefaultPage         = 1
	DefaultPerPage      = 1000
	PaginationParamsKey = "paginationParams"
)

type MetaData struct {
	Page       int `json:"page" example:"1"`
	TotalPages int `json:"totalPages" example:"10"`
	PerPage    int `json:"perPage" example:"1000"`
	TotalItems int `json:"totalItems" example:"10000"`
	ItemCount  int `json:"itemCount" example:"1000"`
} //@name MetaData
//? comment above is for renaming stuct

type PaginatedResponse struct {
	MetaData MetaData                   `json:"metadata"`
	Data     []generator.GeographicArea `json:"data"`
} //@name PaginatedResponse
//? comment above is for renaming stuct

type PaginationParams struct {
	Page    int    `json:"page" example:"1"`
	PerPage int    `json:"perPage" example:"1000"`
	Filter  string `json:"filter" example:"filter"` // This is filter for all the field in the object GeographicArea
} //@name PaginationParams
//? comment above is for renaming stuct

func createPaginatedResponse(data interface{}, PaginationParams PaginationParams) PaginatedResponse {
	page := PaginationParams.Page
	perPage := PaginationParams.PerPage
	filter := PaginationParams.Filter

	// Type assertion to convert the data interface{} to []generator.GeographicArea
	dataList, ok := data.([]generator.GeographicArea)
	if !ok { // Return an empty response if data is not of the expected type
		return PaginatedResponse{}
	}

	// Create a channel for sending filtered items and receiving filtered items
	filterChan := make(chan generator.GeographicArea)

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Define a function for parallel filtering
	filterFunc := func(item generator.GeographicArea) {
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
	slicedData := []generator.GeographicArea{}
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

	return PaginatedResponse{
		MetaData: MetaData{
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

		// Create a context with pagination information and pass it down the chain
		ctx := context.WithValue(r.Context(), PaginationParamsKey, PaginationParams{
			Page:    page,
			PerPage: perPage,
			Filter:  filterParam,
		})

		// Serve the request with the modified context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
