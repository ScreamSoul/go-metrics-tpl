package postgres

import (
	"database/sql"
	"fmt"
)

func initDB(db *sql.DB) error {
	// SQL-запрос для создания таблицы с метриками
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS metric (
		name VARCHAR(255) PRIMARY KEY,
		m_type VARCHAR(255) NOT NULL CHECK (m_type IN ('gauge', 'counter')),
		delta BIGINT,
		value DOUBLE PRECISION
	);
	`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}
