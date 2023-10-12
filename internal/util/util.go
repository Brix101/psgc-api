package util

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	psgctool "github.com/Brix101/psgc-tool"
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

	entries, err := fs.ReadDir(psgctool.EmbedDB, "db")
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no .db files found in embedded data")
	}

	latestEntry := entries[len(entries)-1]
	dbFile := "db/" + latestEntry.Name()

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
