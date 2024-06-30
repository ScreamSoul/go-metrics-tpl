package postgres

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/screamsoul/go-metrics-tpl/internal/models/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

func getRandomMetric() metrics.Metrics {
	// Create a new random generator with a time-based seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	id := fmt.Sprintf("metric_%d", r.Intn(1000))                             // Generate a random ID
	mType := []metrics.MetricType{metrics.Gauge, metrics.Counter}[r.Intn(2)] // Randomly select Gauge or Counter

	var delta int64
	var value float64

	switch mType {
	case metrics.Gauge:
		value = r.Float64() // Generate a random float64 value for Gauge
	case metrics.Counter:
		delta = r.Int63n(100) // Generate a random int64 value for Counter
	}

	return metrics.Metrics{
		ID:    id,
		MType: mType,
		Delta: &delta,
		Value: &value,
	}
}

type PostgresStorageTestSuite struct {
	suite.Suite
	mockDB  *sqlx.DB
	mock    sqlmock.Sqlmock
	storage *PostgresStorage
}

func TestPostgresStorageTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresStorageTestSuite))
}

func (suite *PostgresStorageTestSuite) SetupSuite() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.Fail("Failed to create mock DB")
	}
	suite.mockDB = sqlx.NewDb(db, "sqlmock")
	suite.mock = mock

	suite.storage = &PostgresStorage{
		suite.mockDB, zap.NewNop(), []time.Duration{},
	}
}

func (suite *PostgresStorageTestSuite) TearDownSuite() {
	suite.mockDB.Close()
}

func (suite *PostgresStorageTestSuite) TestAdd() {
	ctx := context.Background()
	metric := getRandomMetric()

	// Mock expected insert statement
	// Expecting INSERT statement
	suite.mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO metrics`))
	suite.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO metrics`)).
		WithArgs(metric.ID, metric.MType, metric.Delta, metric.Value).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute the Add method
	err := suite.storage.Add(ctx, metric)
	assert.NoError(suite.T(), err)

	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresStorageTestSuite) TestGet() {
	rows := sqlmock.NewRows([]string{"value", "delta"}).AddRow(123.45, 6789)

	metric := &metrics.Metrics{ID: "test_id", MType: metrics.Gauge}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT value, delta`)).
		WithArgs(metric.ID, metric.MType).
		WillReturnRows(rows)

	err := suite.storage.Get(context.Background(), metric)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), float64(123.45), *metric.Value)
	assert.Equal(suite.T(), int64(6789), *metric.Delta)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresStorageTestSuite) TestList() {
	count := rand.Intn(10) + 1
	metricsExpected := make([]metrics.Metrics, count)
	valuesExpected := make([][]driver.Value, count)

	var m metrics.Metrics

	for i := 0; i < count; i++ {
		m = getRandomMetric()
		metricsExpected[i] = m
		valuesExpected[i] = []driver.Value{
			m.ID, string(m.MType), m.Delta, m.Value,
		}
	}

	rows := sqlmock.NewRows([]string{"name", "m_type", "delta", "value"}).
		AddRows(valuesExpected...)

	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT name, m_type, delta, value`)).
		WillReturnRows(rows)

	metricsActual, err := suite.storage.List(context.Background())

	require.NoError(suite.T(), err)

	assert.EqualValues(suite.T(), metricsExpected, metricsActual)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresStorageTestSuite) TestPing() {
	suite.mock.ExpectPing()

	isConnect := suite.storage.Ping(context.Background())
	require.Equal(suite.T(), true, isConnect)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresStorageTestSuite) TestBulkAdd() {
	count := rand.Intn(10) + 1
	metricsExpected := make([]metrics.Metrics, count)

	var m metrics.Metrics

	for i := 0; i < count; i++ {
		m = getRandomMetric()
		metricsExpected[i] = m
	}

	suite.mock.ExpectBegin()
	suite.mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO metrics (name, m_type, delta, value)`))
	for i := 0; i < count; i++ {
		suite.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO metrics (name, m_type, delta, value)`)).
			WithArgs(metricsExpected[i].ID, metricsExpected[i].MType, metricsExpected[i].Delta, metricsExpected[i].Value).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
	suite.mock.ExpectCommit()

	err := suite.storage.BulkAdd(context.Background(), metricsExpected)
	require.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}
