package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
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
		c.Logger.Infof("Received message: %s", string(m.Value))
	}
}

func (c *Consumer) Close() error {
	return c.Reader.Close()
}
