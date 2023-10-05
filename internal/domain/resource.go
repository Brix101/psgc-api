package domain

import "github.com/Brix101/psgc-api/internal/generator"

type Resource struct {
	Barangays  []generator.GeographicArea
	Cities     []generator.GeographicArea
	Provinces  []generator.GeographicArea
	Regions    []generator.GeographicArea
	Masterlist []generator.GeographicArea
}
