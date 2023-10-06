package domain

import "context"

type Masterlist struct {
	PsgcCode string `csv:"10-digit PSGC"       json:"psgcCode"`
	Name     string `csv:"Name"                json:"name"`
	Code     string `csv:"Correspondence Code" json:"-"`
	Level    string `csv:"Geographic Level"    json:"-"`
} //@name Masterlist
//? comment above is for renaming stuct

// MasterlistRepository represents the masterlist's repository contract
type MasterlistRepository interface {
	GetList(ctx context.Context, params PaginationParams) (PaginatedResponse, error)
	Create(ctx context.Context, mlist *Masterlist) (*string, *error)
	GetBarangayList(ctx context.Context, params PaginationParams) (PaginatedResponse, error)
	GetBarangayById(ctx context.Context, psgcCode string) (Masterlist, error)
	GetCityList(ctx context.Context, params PaginationParams) (PaginatedResponse, error)
	GetCityById(ctx context.Context, psgcCode string) (Masterlist, error)
	GetProvinceList(ctx context.Context, params PaginationParams) (PaginatedResponse, error)
	GetProvinceById(ctx context.Context, psgcCode string) (Masterlist, error)
	GetRegionList(ctx context.Context, params PaginationParams) (PaginatedResponse, error)
	GetRegionById(ctx context.Context, psgcCode string) (Masterlist, error)
}
