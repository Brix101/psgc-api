package service

import (
	"context"

	"github.com/Brix101/psgc-api/internal/domain"
	"go.uber.org/zap"
)

type Services struct {
	logger *zap.Logger
}

func NewServices(ctx context.Context, logger *zap.Logger) *Services {
	return &Services{
		logger: logger,
	}
}

func (s *Services) GetResources() *domain.Resource {
	barangays := s.getBarangays()
	cities := s.getCities()
	provinces := s.getProvinces()
	regions := s.getRegions()
	masterlist := s.getMasterList()

	return &domain.Resource{
		Barangays:  barangays,
		Cities:     cities,
		Provinces:  provinces,
		Regions:    regions,
		Masterlist: masterlist,
	}
}
