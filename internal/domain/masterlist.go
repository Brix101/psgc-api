package domain

import "context"

type Masterlist struct {
	PsgcCode string `csv:"10-digit PSGC"       json:"psgc_code"`
	Name     string `csv:"Name"                json:"name"`
	Code     string `csv:"Correspondence Code" json:"-"`
	Level    string `csv:"Geographic Level"    json:"-"`
} //@name Masterlist
//? comment above is for renaming stuct

type PaginatedMasterlist struct {
	MetaData MetaData     `json:"metadata"`
	Data     []Masterlist `json:"data"`
} //@name PaginatedMasterlist
//? comment above is for renaming stuct

// MasterlistRepository represents the masterlist's repository contract
type MasterlistRepository interface {
	GetList(ctx context.Context, params PaginationParams) (PaginatedMasterlist, error)
	Create(ctx context.Context, data *Masterlist) error
	CreateBatch(ctx context.Context, datas []*Masterlist) error
	GetBarangayList(ctx context.Context, params PaginationParams) (PaginatedMasterlist, error)
	GetBarangayById(ctx context.Context, psgcCode string) (Masterlist, error)
	GetCityList(ctx context.Context, params PaginationParams) (PaginatedMasterlist, error)
	GetCityById(ctx context.Context, psgcCode string) (Masterlist, error)
	GetProvinceList(ctx context.Context, params PaginationParams) (PaginatedMasterlist, error)
	GetProvinceById(ctx context.Context, psgcCode string) (Masterlist, error)
	GetRegionList(ctx context.Context, params PaginationParams) (PaginatedMasterlist, error)
	GetRegionById(ctx context.Context, psgcCode string) (Masterlist, error)
}
