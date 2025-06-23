-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS metrics (
    id SERIAL PRIMARY KEY,
    type TEXT NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    timestamp TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS metrics;
-- +goose StatementEnd
