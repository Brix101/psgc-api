package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/gocarina/gocsv"
	"go.uber.org/zap"
)

const (
	Regions     = "regions"
	Provinces   = "provinces"
	Cities      = "cities"
	Barangays   = "barangays"
	Publication = "aaa-Publication-Datafile"
	JsonFolder  = "files/json"
	CsvFolder   = "files/csv"
)

type GeographicArea struct {
	PsgcCode     string `csv:"10-digit PSGC" json:"psgcCode"`
	RegionCode   string `json:"regionCode,omitempty"`
	ProvinceCode string `json:"provinceCode,omitempty"`
	CityCode     string `json:"cityCode,omitempty"`
	Name         string `csv:"Name" json:"name"`
	Code         string `csv:"Correspondence Code" json:"-"`
	Level        string `csv:"Geographic Level" json:"-"`
}

type Generator struct {
	Filename string
}

func NewGenerator(Filename string) *Generator {
	return &Generator{
		Filename: Filename,
	}
}

func (g *Generator) GenerateJson(ctx context.Context, logger *zap.Logger) error {
	file, err := os.Open(g.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	psgcData := []*GeographicArea{}

	if err := gocsv.Unmarshal(file, &psgcData); err != nil {
		return err
	}

	// Create the output folder if it doesn't exist
	if err := os.MkdirAll(JsonFolder, os.ModePerm); err != nil {
		return err
	}

	var wg sync.WaitGroup
	doneCh := make(chan struct{})

	// Define a function to create and write a JSON file
	createJSONFile := func(level string, data []*GeographicArea, doneCh chan<- struct{}) {
		defer wg.Done()
		var formatLevel string

		switch level {
		case "Reg":
			formatLevel = Regions
		case "Prov":
			for i, item := range data {
				psgcCode := item.PsgcCode
				data[i].RegionCode = psgcCode[:2] + strings.Repeat("0", len(psgcCode)-2)
			}
			formatLevel = Provinces
		case "City":
			for i, item := range data {
				psgcCode := item.PsgcCode
				data[i].RegionCode = psgcCode[:2] + strings.Repeat("0", len(psgcCode)-2)
				data[i].ProvinceCode = psgcCode[:5] + strings.Repeat("0", len(psgcCode)-5)
			}
			formatLevel = Cities
		case "Bgy":
			for i, item := range data {
				psgcCode := item.PsgcCode
				data[i].RegionCode = psgcCode[:2] + strings.Repeat("0", len(psgcCode)-2)
				data[i].ProvinceCode = psgcCode[:5] + strings.Repeat("0", len(psgcCode)-5)
				data[i].CityCode = psgcCode[:7] + strings.Repeat("0", len(psgcCode)-7)
			}
			formatLevel = Barangays
		default:
			formatLevel = level
		}

		if formatLevel != level || level == Publication {
			filename := fmt.Sprintf("%s/%s.json", JsonFolder, formatLevel)

			// Remove the existing JSON file if it exists
			if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
				panic(err)
			}

			// Create a new JSON file for writing
			createdFile, err := os.Create(filename)
			if err != nil {
				panic(err)
			}
			defer createdFile.Close()

			// Create a JSON encoder
			encoder := json.NewEncoder(createdFile)

			// Sort the data by PsgcCode before encoding
			sort.Slice(data, func(i, j int) bool {
				return data[i].PsgcCode < data[j].PsgcCode
			})

			// Encode and write the sorted data to the JSON file
			if err := encoder.Encode(data); err != nil {
				panic(err)
			}

			message := fmt.Sprintf("%d Data for level '%s' written to %s\n", len(data), level, filename)
			logger.Info(message)
		}
		// Notify that this goroutine is done
		doneCh <- struct{}{}
	}

	// Group data by level and start creating JSON files concurrently

	groupedData := make(map[string][]*GeographicArea)
	for _, item := range psgcData {
		level := item.Level

		if len(level) >= 3 {
			if level == "Mun" {
				level = "City"
			}
			groupedData[level] = append(groupedData[level], item)
		}
		groupedData[Publication] = append(groupedData[Publication], item)
	}

	for level, data := range groupedData {
		wg.Add(1)
		go createJSONFile(level, data, doneCh)
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	// Wait for all goroutines to complete and receive their notifications
	for range groupedData {
		<-doneCh
	}

	os.Exit(0)
	return nil
}
