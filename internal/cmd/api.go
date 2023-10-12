package cmd

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/Brix101/psgc-tool/internal/api"
	"github.com/Brix101/psgc-tool/internal/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func APICmd(ctx context.Context) *cobra.Command {
	var port int
	cmd := &cobra.Command{
		Use:   "api",
		Args:  cobra.ExactArgs(0),
		Short: "Runs the RESTful API.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if os.Getenv("PORT") != "" {
				port, _ = strconv.Atoi(os.Getenv("PORT"))
			}

			logger := util.NewLogger("api")
			defer func() { _ = logger.Sync() }()

			db, err := util.NewSQLitePool(ctx)
			if err != nil {
				return err
			}
			defer db.Close()

			api := api.NewAPI(ctx, logger, db)
			server := api.Server(port)

			// Graceful shutdown with a 30-second timeout
			shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
			defer shutdownCancel()

			// Run the server
			go func() { _ = server.ListenAndServe() }()
			logger.Info("🚀🚀🚀 Server at port:", zap.Int("port", port))

			// Wait for the context to be canceled (due to signal or other reasons)
			<-ctx.Done()

			// Trigger graceful shutdown
			err = server.Shutdown(shutdownCtx)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().IntVarP(&port, "port", "P", 5000, "Port number")

	return cmd
}
