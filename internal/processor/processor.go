package processor

import (
	"github.com/timur-danilchenko/metric-bridge/internal/model"
	"go.uber.org/zap"
)

func Handle(metric model.Metric, logger *zap.SugaredLogger) {
	if metric.Type == "" {
		logger.Warn("Received metric with empty type")
		return
	}

	logger.Infof("Processed metric: type=%s value=%.2f timestamp=%s",
		metric.Type,
		metric.Value,
		metric.Timestamp.Format("2006-01-02 15:04:05"),
	)
}
