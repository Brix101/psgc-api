package util

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func NewLogger(service string) *zap.Logger {
	env := os.Getenv("ENV")

	logger, _ := zap.NewProduction(zap.Fields(
		zap.String("env", env),
		zap.String("service", service),
	))

	if env == "" || env == "development" {
		logger, _ = zap.NewDevelopment()
	}

	return logger
}

func NewSQLitePool(ctx context.Context) (*sql.DB, error) {
	year := time.Now().Year()
	dbFile := fmt.Sprintf("db/psgc_%d.db", year)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	// Set the journal mode to "WAL" (Write-Ahead Logging)
	_, err = db.ExecContext(ctx, "PRAGMA journal_mode=WAL;")
	if err != nil {
		db.Close() // Close the database if there's an error
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		db.Close() // Close the database if there's an error
		return nil, err
	}
	// Set the maximum number of open connections in the pool
	// db.SetMaxOpenConns(maxConns)
	return db, nil
}
