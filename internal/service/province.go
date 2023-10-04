package service

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Brix101/psgc-api/pkg/generator"
	"go.uber.org/zap"
)

func GetProvinces(logger *zap.Logger) []generator.GeographicArea {
	filePath := fmt.Sprintf("%s/%s.json", generator.JsonFolder, generator.Provinces)
	file, err := os.Open(filePath)
	if err != nil {
		logger.Error("Failed to open", zap.Error(err))
		return []generator.GeographicArea{}
	}
	defer file.Close()
	byteResult, err := io.ReadAll(file)
	if err != nil {
		logger.Error("Failed to read", zap.Error(err))
		return []generator.GeographicArea{}
	}

	var psgcData []generator.GeographicArea // Declare a slice

	err = json.Unmarshal(byteResult, &psgcData) // Unmarshal into psgcData
	if err != nil {
		logger.Error("Failed to read", zap.Error(err))
		return []generator.GeographicArea{}
	}

	return psgcData
}
