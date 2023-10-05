package cmd

import (
	"context"

	"github.com/Brix101/psgc-api/internal/generator"
	"github.com/Brix101/psgc-api/internal/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func GeneratorCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Args:  cobra.ExactArgs(1),
		Short: "Genrate a new json files.",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := util.NewLogger("generator")
			defer func() { _ = logger.Sync() }()

			jsonGenerator := generator.NewGenerator(args[0])
			if err := jsonGenerator.GenerateJson(ctx, logger); err != nil {
				logger.Error("Generator error:", zap.Error(err))
			}

			return nil
		},
	}
	return cmd
}
