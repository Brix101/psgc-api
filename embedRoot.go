package psgctool

import (
	"embed"
)

//go:embed migrations/*.sql
var EmbedMigrations embed.FS

//go:embed db/*.db
var EmbedDB embed.FS
