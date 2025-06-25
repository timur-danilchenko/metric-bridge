package processor

import (
	"context"
	"errors"
	"time"

	"github.com/timur-danilchenko/metric-bridge/internal/metrics"
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
	start := time.Now()

	if metric.Type == "" {
		p.logger.Warn("Received metric with empty type")
		metrics.ProcessingErrors.Inc()
		return errors.New("invalid metric")
	}

	err := p.repo.SaveMetric(ctx, metric)
	if err != nil {
		metrics.ProcessingErrors.Inc()
		return err
	}

	p.logger.Infof("Processed metric: type=%s value=%.2f timestamp=%s",
		metric.Type,
		metric.Value,
		metric.Timestamp.Format("2006-01-02 15:04:05"),
	)

	metrics.ProcessedMessages.Inc()
	metrics.ProcessingDuration.Observe(time.Since(start).Seconds())

	return nil
}
