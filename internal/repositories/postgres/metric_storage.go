package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
	"github.com/screamsoul/go-metrics-tpl/pkg/logging"
	"go.uber.org/zap"
)

type PostgresStorage struct {
	db      *sqlx.DB
	logging *zap.Logger
}

func NewPostgresStorage(dataSourceName string) *PostgresStorage {
	db := sqlx.MustOpen("pgx", dataSourceName)

	return &PostgresStorage{db, logging.GetLogger()}
}

func (storage *PostgresStorage) Add(ctx context.Context, metric metrics.Metrics) error {
	stmt, err := storage.db.PrepareContext(ctx, `
		INSERT INTO metrics (name, m_type, delta, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name) DO UPDATE SET
			delta = CASE WHEN metrics.m_type = 'counter' THEN metrics.delta + excluded.delta ELSE excluded.delta END,
			value = excluded.value;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, metric.ID, metric.MType, metric.Delta, metric.Value)
	return err
}

func (storage *PostgresStorage) Get(ctx context.Context, metric *metrics.Metrics) error {
	query := `SELECT value, delta FROM metrics WHERE name = $1 AND m_type = $2`
	var value sql.NullFloat64
	var delta sql.NullInt64

	row := storage.db.QueryRowContext(ctx, query, metric.ID, metric.MType)

	err := row.Scan(&value, &delta)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("metric with Name %s not found", metric.ID)
		}
		return err
	}

	if value.Valid {
		metric.Value = &value.Float64
	}
	if delta.Valid {
		metric.Delta = &delta.Int64
	}

	return nil
}

func (storage *PostgresStorage) List(ctx context.Context) (metricsList []metrics.Metrics, err error) {
	query := `SELECT name, m_type, delta, value FROM metrics`

	err = storage.db.SelectContext(ctx, &metricsList, query)
	return
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

func (storage *PostgresStorage) BulkAdd(ctx context.Context, metricList []metrics.Metrics) error {
	tx, err := storage.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, `
		INSERT INTO metrics (name, m_type, delta, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name) DO UPDATE SET
			delta = CASE WHEN metrics.m_type = 'counter' THEN metrics.delta + excluded.delta ELSE excluded.delta END,
			value = excluded.value;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, metric := range metricList {
		_, err = stmt.ExecContext(ctx, metric.ID, metric.MType, metric.Delta, metric.Value)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
