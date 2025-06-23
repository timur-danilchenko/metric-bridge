package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"github.com/timur-danilchenko/metric-bridge/internal/model"
	"github.com/timur-danilchenko/metric-bridge/internal/processor"
	"go.uber.org/zap"
)

type Consumer struct {
	Reader *kafka.Reader
	Logger *zap.SugaredLogger
}

func NewConsumer(brokers []string, topic string, logger *zap.SugaredLogger) *Consumer {
	readerCfg := kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "metric-bridge-group",
		MinBytes: 1e3, // 1KB
		MaxBytes: 1e6, // 1MB
	}

	r := kafka.NewReader(readerCfg)

	return &Consumer{
		Reader: r,
		Logger: logger,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		m, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				c.Logger.Info("Kafka consumer context canceled, shutting down.")
				return
			}
			c.Logger.Errorf("Failed to read message: %v", err)
			continue
		}

		var metric model.Metric
		if err := json.Unmarshal(m.Value, &metric); err != nil {
			c.Logger.Errorf("Invalid JSON metric: %v", err)
			continue
		}

		processor.Handle(metric, c.Logger)
		// c.Logger.Infof("Received message: %s", string(m.Value))
	}
}

func (c *Consumer) Close() error {
	return c.Reader.Close()
}
