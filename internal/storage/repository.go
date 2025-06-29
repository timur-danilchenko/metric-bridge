package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/timur-danilchenko/metric-bridge/internal/config"
	"github.com/timur-danilchenko/metric-bridge/internal/model"
)

type Repo struct {
	conn *pgx.Conn
}

func NewRepository(ctx context.Context, cfg config.PostgresConfig) (*Repo, error) {
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
	return &Repo{conn: conn}, nil
}

func (r *Repo) Close(ctx context.Context) error {
	return r.conn.Close(ctx)
}

func (r *Repo) SaveMetric(ctx context.Context, metric model.Metric) error {
	_, err := r.conn.Exec(ctx,
		`INSERT INTO metrics (type, value, timestamp) VALUES ($1, $2, $3)`,
		metric.Type,
		metric.Value,
		metric.Timestamp,
	)

	return err
}

// Compile-time проверка, что *Repo реализует интерфейс Repository
var _ Repository = (*Repo)(nil)
