package util

import (
	"database/sql"

	psgctool "github.com/Brix101/psgc-tool"
	"github.com/pressly/goose/v3"
)


func NewMigration(db *sql.DB) error {
    goose.SetBaseFS(psgctool.EmbedMigrations)

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
