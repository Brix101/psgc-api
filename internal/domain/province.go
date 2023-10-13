package domain

import "context"

type Province struct {
	PsgcCode string `json:"psgc_code"`
	RegCode  string `json:"regCode"`
	Name     string `json:"name"`
} //@name Province
//? comment above is for renaming stuct

type PaginatedProvince struct {
	MetaData MetaData   `json:"metadata"`
	Data     []Province `json:"data"`
} //@name PaginatedProvince
//? comment above is for renaming stuct

// ProvinceRepository represents the province's repository contract
type ProvinceRepository interface {
	GetAll(ctx context.Context, params PaginationParams) (PaginatedProvince, error)
	GetById(ctx context.Context, psgcCode string) (Province, error)
	
	Create(ctx context.Context, reg *Masterlist) error
}
