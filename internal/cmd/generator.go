package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/Brix101/psgc-api/internal/generator"
	"github.com/Brix101/psgc-api/internal/util"
	"github.com/spf13/cobra"
)

func GeneratorCmd(ctx context.Context) *cobra.Command {
	year := time.Now().Year()
	file := fmt.Sprintf("%s/psgc_%d.csv", generator.CsvFolder, year)

	cmd := &cobra.Command{
		Use:   "generate",
		Args:  cobra.ExactArgs(0),
		Short: "Generate a new JSON file.",
		Long:  "Generate a new JSON file from a CSV input file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := util.NewLogger("generator")
			defer func() { _ = logger.Sync() }()

			db, err := util.NewSQLitePool(ctx)
			if err != nil {
				return err
			}
			defer db.Close()

			jsonGenerator := generator.NewGenerator(file, db)
			if err := jsonGenerator.GenerateData(ctx, logger); err != nil {
				return err
			}

			<-ctx.Done()

			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "CSV file location")

	return cmd
}
