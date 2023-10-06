package domain

type MetaData struct {
	Page       int `json:"page"       example:"1"`
	TotalPages int `json:"totalPages" example:"10"`
	PerPage    int `json:"perPage"    example:"1000"`
	TotalItems int `json:"totalItems" example:"10000"`
	ItemCount  int `json:"itemCount"  example:"1000"`
} //@name MetaData
//? comment above is for renaming stuct

type PaginatedResponse struct {
	MetaData MetaData     `json:"metadata"`
	Data     []Masterlist `json:"data"`
} //@name PaginatedResponse
//? comment above is for renaming stuct

type PaginationParams struct {
	Page    int    `json:"page"    example:"1"`
	PerPage int    `json:"perPage" example:"1000"`
	Filter  string `json:"filter"  example:"filter"` // This is filter for all the field in the Masterlist object
} //@name PaginationParams
// INFO? comment above is for renaming stuct
