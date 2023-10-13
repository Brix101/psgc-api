package domain

import "context"

type CityMuni struct {
	PsgcCode string `json:"psgc_code"`
	ProvCode string `json:"prov_code"`
	Name     string `json:"name"`
	Level    string `json:"level"`
} //@name CityMuni
//? comment above is for renaming stuct

type PaginatedCityMuni struct {
	MetaData MetaData   `json:"metadata"`
	Data     []CityMuni `json:"data"`
} //@name PaginatedCityMuni
//? comment above is for renaming stuct

// CityMuniRepository represents the cityMuni's repository contract
type CityMuniRepository interface {
	GetAll(ctx context.Context, params PaginationParams) (PaginatedCityMuni, error)
	GetById(ctx context.Context, psgcCode string) (CityMuni, error)
	GetAllCity(ctx context.Context, params PaginationParams) (PaginatedCityMuni, error)
	GetCityById(ctx context.Context, psgcCode string) (CityMuni, error)
	GetAllMunicipality(ctx context.Context, params PaginationParams) (PaginatedCityMuni, error)
	GetMunicipalityById(ctx context.Context, psgcCode string) (CityMuni, error)

	Create(ctx context.Context, reg *Masterlist) error
}
