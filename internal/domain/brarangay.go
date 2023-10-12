package domain

import "context"

type Barangay struct {
	PsgcCode     string `json:"psgc_code"`
	CityMuniCode string `json:"city_muni_code"`
	Name         string `json:"name"`
} //@name Barangay
//? comment above is for renaming stuct

type PaginatedBarangay struct {
	MetaData MetaData   `json:"metadata"`
	Data     []Barangay `json:"data"`
} //@name PaginatedBarangay
//? comment above is for renaming stuct

// BarangayRepository represents the barangay's repository contract
type BarangayRepository interface {
	GetList(ctx context.Context, params PaginationParams) (PaginatedBarangay, error)
	Create(ctx context.Context, reg *Masterlist) error
	GetById(ctx context.Context, psgcCode string) (Barangay, error)
}
