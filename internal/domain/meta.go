package domain

type MetaData struct {
	Page       int `json:"page"       example:"1"`
	TotalPages int `json:"total_pages" example:"10"`
	PerPage    int `json:"per_page"    example:"1000"`
	TotalItems int `json:"total_items" example:"10000"`
	ItemCount  int `json:"item_count"  example:"1000"`
} //@name MetaData
//? comment above is for renaming stuct

type PaginationParams struct {
	Page    int    `json:"page"    example:"1"      validate:"gte=0"`
	PerPage int    `json:"per_page" example:"1000"   validate:"lte=1000"`
	Filter  string `json:"filter"  example:"filter"` // This is filter for all the field in the  object
} //@name PaginationParams
// INFO? comment above is for renaming stuct
