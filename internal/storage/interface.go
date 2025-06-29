package storage

import (
	"context"

	"github.com/timur-danilchenko/metric-bridge/internal/model"
)

type Repository interface {
	SaveMetric(ctx context.Context, metric model.Metric) error
}
