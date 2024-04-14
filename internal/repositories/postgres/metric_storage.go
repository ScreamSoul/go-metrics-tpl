package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
	"github.com/screamsoul/go-metrics-tpl/pkg/logging"
	"go.uber.org/zap"
)

type PostgresStorage struct {
	db      *sql.DB
	logging *zap.Logger
}

func NewPostgresStorage(dataSourceName string) *PostgresStorage {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		panic(err)
	}

	// Инициализация базы данных
	err = initDB(db)
	if err != nil {
		panic(err)
	}

	return &PostgresStorage{db, logging.GetLogger()}
}

func (storage *PostgresStorage) Add(ctx context.Context, metric metrics.Metrics) error {
	stmt, err := storage.db.PrepareContext(ctx, `
		INSERT INTO metric (name, m_type, delta, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT(name) DO UPDATE SET m_type = excluded.m_type, delta = excluded.delta, value = excluded.value;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, metric.ID, metric.MType, metric.Delta, metric.Value)
	return err
}

func (storage *PostgresStorage) Get(ctx context.Context, metric *metrics.Metrics) error {
	query := `SELECT delta, value FROM metric WHERE name = $1 and m_type = $2`
	row := storage.db.QueryRowContext(ctx, query, metric.ID, metric.MType)

	err := row.Scan(metric.Delta, metric.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("metric with ID %s not found", metric.ID)
		}
		return err
	}
	return nil
}

func (storage *PostgresStorage) List(ctx context.Context) ([]metrics.Metrics, error) {
	query := `SELECT name, m_type, delta, value FROM metric`
	rows, err := storage.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metricsList []metrics.Metrics

	for rows.Next() {
		var metric metrics.Metrics
		err := rows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
		if err != nil {
			return nil, err
		}
		metricsList = append(metricsList, metric)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return metricsList, nil
}

func (storage *PostgresStorage) Ping(ctx context.Context) bool {
	err := storage.db.PingContext(ctx)
	if err != nil {
		storage.logging.Error("db connect error", zap.Error(err))
	}
	return err == nil
}

func (storage *PostgresStorage) Close() {
	err := storage.db.Close()
	if err != nil {
		storage.logging.Error("db close connection error", zap.Error(err))
	}
}
