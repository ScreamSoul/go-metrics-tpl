package postgres

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

func (storage *PostgresStorage) Bootstrap(ctx context.Context) error {
	tx, err := storage.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			storage.logging.Warn("rollback transaction error", zap.Error(err))
		}
	}()

	// SQL-запрос для создания таблицы с метриками
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS metrics (
		name VARCHAR(255) PRIMARY KEY,
		m_type VARCHAR(255) NOT NULL CHECK (m_type IN ('gauge', 'counter')),
		delta BIGINT,
		value DOUBLE PRECISION
	);
	`

	_, err = tx.ExecContext(ctx, createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return tx.Commit()
}
