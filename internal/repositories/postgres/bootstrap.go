package postgres

import (
	"context"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (storage *PostgresStorage) Bootstrap(ctx context.Context) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(storage.db.DB, "migrations"); err != nil {
		return err
	}
	return nil
}
