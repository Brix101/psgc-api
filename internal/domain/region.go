package domain

import "context"

type Region struct {
	PsgcCode string `json:"psgc_code"`
	Name     string `json:"name"`
} //@name Region
//? comment above is for renaming stuct

type PaginatedRegion struct {
	MetaData MetaData `json:"metadata"`
	Data     []Region `json:"data"`
} //@name PaginatedRegion
//? comment above is for renaming stuct

// RegionRepository represents the region's repository contract
type RegionRepository interface {
	GetAll(ctx context.Context, params PaginationParams) (PaginatedRegion, error)
	GetById(ctx context.Context, psgcCode string) (Region, error)
	
	Create(ctx context.Context, reg *Masterlist) error
}
