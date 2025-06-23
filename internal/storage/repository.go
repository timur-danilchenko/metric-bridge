package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/timur-danilchenko/metric-bridge/internal/config"
	"github.com/timur-danilchenko/metric-bridge/internal/model"
)

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(ctx context.Context, cfg config.PostgresConfig) (*Repository, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,
		cfg.Sslmode,
	)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &Repository{conn: conn}, nil
}

func (r *Repository) Close(ctx context.Context) error {
	return r.conn.Close(ctx)
}

func (r *Repository) SaveMetric(ctx context.Context, metric model.Metric) error {
	_, err := r.conn.Exec(ctx,
		`INSERT INTO metrics (type, value, timestamp) VALUES ($1, $2, $3)`,
		metric.Type, metric.Value, metric.Timestamp)
	return err
}
