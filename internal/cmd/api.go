package cmd

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Brix101/psgc-api/internal/api"
	"github.com/Brix101/psgc-api/internal/util"
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

			api := api.NewAPI(ctx, logger)
			server := api.Server(port)

			// Graceful shutdown with a 30-second timeout
			shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
			defer shutdownCancel()

			// Run the server
			logger.Info("🚀🚀🚀 Server at port:", zap.Int("port", port))
			err := server.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				logger.Error("Server error:", zap.Error(err))
			}

			// Wait for the context to be canceled (due to signal or other reasons)
			<-ctx.Done()

			// Trigger graceful shutdown
			err = server.Shutdown(shutdownCtx)
			if err != nil {
				logger.Error("Server shutdown error:", zap.Error(err))
			}

			return nil
		},
	}

	cmd.Flags().IntVarP(&port, "port", "P", 5000, "Port number")

	return cmd
}
