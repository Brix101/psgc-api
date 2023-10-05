package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Brix101/psgc-api/internal/cmd"
)

// @title         Philippine Standard Geographic Code (PSGC) API
// @version       1.0
// @description   This API is based on the Philippine Standard Geographic Code (PSGC), which is a systematic classification and coding of geographic areas in the Philippines. Its units of classification are based on the four well-established levels of geographical-political subdivisions of the country, including the region, the province, the municipality/city, and the barangay.
// @BasePath      /api
// @externalDocs.description  Data used in this API is sourced from PSGC main page
// @externalDocs.url  https://psa.gov.ph/classification/psgc
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	ret := cmd.Execute(ctx)
	os.Exit(ret)
}
