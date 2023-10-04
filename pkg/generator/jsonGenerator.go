package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gocarina/gocsv"
)

type Psgc struct {
	PsgcCode     string `csv:"10-digit PSGC" json:"psgcCode"`
	RegionCode   string `json:"regionCode,omitempty"`
	ProvinceCode string `json:"provinceCode,omitempty"`
	CityMunCode  string `json:"cityMunCode,omitempty"`
	Name         string `csv:"Name" json:"name"`
	Code         string `csv:"Correspondence Code" json:"-"`
	Level        string `csv:"Geographic Level" json:"-"`
}

type Generator struct {
	Path         string
	OutputFolder string
}

func InitGenerator(Path string, OutputFolders ...string) *Generator {
	var outputFolder string

	if len(OutputFolders) > 0 {
		outputFolder = OutputFolders[0]
	} else {
		outputFolder = "files" // Default value
	}

	return &Generator{
		Path:         Path,
		OutputFolder: outputFolder,
	}
}

func (g *Generator) GenerateJson() error {
	path := g.Path
	outputFolder := g.OutputFolder

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	psgcData := []*Psgc{}

	if err := gocsv.Unmarshal(file, &psgcData); err != nil {
		return err
	}

	year := time.Now().Year()
	publication := "aaa-Publication-Datafile"

	// Create the output folder if it doesn't exist
	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		return err
	}

	var wg sync.WaitGroup
	doneCh := make(chan struct{})

	// Define a function to create and write a JSON file
	createJSONFile := func(level string, data []*Psgc, doneCh chan<- struct{}) {
		defer wg.Done()
		var formatLevel string

		switch level {
		case "Reg":
			formatLevel = "regions"
		case "Prov":
			for i, item := range data {
				psgcCode := item.PsgcCode
				data[i].RegionCode = psgcCode[:2] + strings.Repeat("0", len(psgcCode)-2)
			}
			formatLevel = "provinces"
		case "City":
			for i, item := range data {
				psgcCode := item.PsgcCode
				data[i].RegionCode = psgcCode[:2] + strings.Repeat("0", len(psgcCode)-2)
				data[i].ProvinceCode = psgcCode[:5] + strings.Repeat("0", len(psgcCode)-5)
			}
			formatLevel = "cities"
		case "Bgy":
			for i, item := range data {
				psgcCode := item.PsgcCode
				data[i].RegionCode = psgcCode[:2] + strings.Repeat("0", len(psgcCode)-2)
				data[i].ProvinceCode = psgcCode[:5] + strings.Repeat("0", len(psgcCode)-5)
				data[i].CityMunCode = psgcCode[:7] + strings.Repeat("0", len(psgcCode)-7)
			}
			formatLevel = "barangays"
		default:
			formatLevel = level
		}

		if formatLevel != level || level == publication {
			filename := fmt.Sprintf("%s/%d-%s.json", outputFolder, year, formatLevel)

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

			fmt.Printf("%d Data for level '%s' written to %s\n", len(data), level, filename)
		}
		// Notify that this goroutine is done
		doneCh <- struct{}{}
	}

	// Group data by level and start creating JSON files concurrently

	groupedData := make(map[string][]*Psgc)
	for _, item := range psgcData {
		level := item.Level

		if len(level) >= 3 {
			if level == "Mun" {
				level = "City"
			}
			groupedData[level] = append(groupedData[level], item)
		}
		groupedData[publication] = append(groupedData[publication], item)
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

	return nil
}
