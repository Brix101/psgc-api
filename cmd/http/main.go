package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Brix101/psgc-api/internal/api"
	"github.com/Brix101/psgc-api/internal/util"
	// "github.com/Brix101/psgc-api/pkg/generator"
	"go.uber.org/zap"
)

func main() {
	port := 5000
	if os.Getenv("PORT") != "" {
		port, _ = strconv.Atoi(os.Getenv("PORT"))
	}

	// jsonGenerator := generator.InitGenerator("psgc_2023.csv")
	// if err := jsonGenerator.GenerateJson(); err != nil {
	// 	log.Fatal(err)
	// }

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	logger := util.NewLogger("api")
	defer func() { _ = logger.Sync() }()

	api := api.NewAPI(serverCtx, logger)
	server := api.Server(port)

	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownStopCtx := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
			shutdownStopCtx()
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	logger.Info("🚀🚀🚀 Server at port: ", zap.Int("port", port))
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
