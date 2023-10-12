package psgcapi

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func NewMigration(db *sql.DB) error {
    goose.SetBaseFS(embedMigrations)

    if err := goose.SetDialect("sqlite3"); err != nil {
        return err
    }

    if err := goose.Reset(db, "migrations"); err != nil {
        return err
    }

    if err := goose.Up(db, "migrations"); err != nil {
        return err
    }


	return nil
}
