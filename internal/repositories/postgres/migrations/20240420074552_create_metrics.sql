-- +goose Up
-- +goose StatementBegin
CREATE TYPE metric_type AS ENUM ('gauge', 'counter');

CREATE TABLE IF NOT EXISTS metrics (
    name VARCHAR(255) PRIMARY KEY,
    m_type metric_type NOT NULL,
    delta BIGINT,
    value DOUBLE PRECISION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS metrics;
DROP TYPE IF EXISTS metric_type;
-- +goose StatementEnd
