package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
	"github.com/screamsoul/go-metrics-tpl/pkg/backoff"
	"github.com/screamsoul/go-metrics-tpl/pkg/logging"
	"go.uber.org/zap"
)

type PostgresStorage struct {
	db               *sqlx.DB
	logging          *zap.Logger
	backoffInteraval []time.Duration
}

func NewPostgresStorage(dataSourceName string, backoffInteraval []time.Duration) *PostgresStorage {
	db := sqlx.MustOpen("pgx", dataSourceName)

	return &PostgresStorage{db, logging.GetLogger(), backoffInteraval}
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

	exec := func() error {
		_, err = stmt.ExecContext(ctx, metric.ID, metric.MType, metric.Delta, metric.Value)
		return err
	}

	if storage.backoffInteraval != nil {
		err = backoff.RetryWithBackoff(storage.backoffInteraval, IsTemporaryConnectionError, exec)
		if err != nil {
			err = fmt.Errorf("failed retries db request, %w", err)
		}
		return err
	}
	return exec()

}

func (storage *PostgresStorage) Get(ctx context.Context, metric *metrics.Metrics) error {
	query := `SELECT value, delta FROM metrics WHERE name = $1 AND m_type = $2`
	var value sql.NullFloat64
	var delta sql.NullInt64

	row := storage.db.QueryRowContext(ctx, query, metric.ID, metric.MType)
	var err error

	exec := func() error {
		err := row.Scan(&value, &delta)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("metric with Name %s not found", metric.ID)
			}
			return err
		}
		return nil
	}

	if storage.backoffInteraval != nil {
		err = backoff.RetryWithBackoff(storage.backoffInteraval, IsTemporaryConnectionError, exec)
		if err != nil {
			err = fmt.Errorf("failed retries db request, %w", err)
		}
	} else {
		err = exec()
	}

	if err != nil {
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
	exec := func() error {
		return storage.db.SelectContext(ctx, &metricsList, query)
	}

	if storage.backoffInteraval != nil {
		err = backoff.RetryWithBackoff(storage.backoffInteraval, IsTemporaryConnectionError, exec)
		if err != nil {
			err = fmt.Errorf("failed retries db request, %w", err)
		}
	} else {
		err = exec()
	}

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
	defer func() {
		if err := tx.Rollback(); err != nil {
			storage.logging.Warn("rollback transaction error", zap.Error(err))
		}
	}()

	stmt, err := tx.PreparexContext(ctx, `
		INSERT INTO metrics (name, m_type, delta, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name) DO UPDATE SET
			delta = CASE WHEN metrics.m_type = 'counter' THEN metrics.delta + excluded.delta ELSE metrics.delta END,
			value = CASE WHEN metrics.m_type = 'gauge' THEN excluded.value ELSE metrics.value END;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, metric := range metricList {
		var delta sql.NullInt64
		var value sql.NullFloat64

		if metric.Delta != nil {
			delta.Int64 = *metric.Delta
			delta.Valid = true
		}

		if metric.Value != nil {
			value.Float64 = *metric.Value
			value.Valid = true
		}

		_, err = stmt.ExecContext(ctx, metric.ID, metric.MType, delta, value)
		if err != nil {
			return err
		}
	}

	exec := func() error {
		return tx.Commit()
	}

	if storage.backoffInteraval != nil {
		err = backoff.RetryWithBackoff(storage.backoffInteraval, IsTemporaryConnectionError, exec)
		if err != nil {
			err = fmt.Errorf("failed retries db request, %w", err)
		}
		return err
	}
	return tx.Commit()
}
