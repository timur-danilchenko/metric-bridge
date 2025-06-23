package processor

import (
	"context"
	"errors"

	"github.com/timur-danilchenko/metric-bridge/internal/model"
	"github.com/timur-danilchenko/metric-bridge/internal/storage"
	"go.uber.org/zap"
)

type Processor struct {
	logger *zap.SugaredLogger
	repo   *storage.Repository
}

func NewProcessor(logger *zap.SugaredLogger, repo *storage.Repository) *Processor {
	return &Processor{
		logger: logger,
		repo:   repo,
	}
}

func (p *Processor) Handle(ctx context.Context, metric model.Metric) error {
	if metric.Type == "" {
		return errors.New("Received metric with empty type")
	}

	p.repo.SaveMetric(ctx, metric)

	p.logger.Infof("Processed metric: type=%s value=%.2f timestamp=%s",
		metric.Type,
		metric.Value,
		metric.Timestamp.Format("2006-01-02 15:04:05"),
	)

	return nil
}
