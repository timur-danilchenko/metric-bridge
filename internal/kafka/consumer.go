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
	reader    *kafka.Reader
	logger    *zap.SugaredLogger
	processor *processor.Processor
}

func NewConsumer(brokers []string, topic string, logger *zap.SugaredLogger, processor *processor.Processor) *Consumer {
	readerCfg := kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "metric-bridge-group",
		MinBytes: 1e3, // 1KB
		MaxBytes: 1e6, // 1MB
	}

	reader := kafka.NewReader(readerCfg)

	return &Consumer{
		reader:    reader,
		logger:    logger,
		processor: processor,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				c.logger.Info("Kafka consumer context canceled, shutting down.")
				return
			}
			c.logger.Errorf("Failed to read message: %v", err)
			continue
		}

		var metric model.Metric
		if err := json.Unmarshal(msg.Value, &metric); err != nil {
			c.logger.Errorf("Invalid JSON metric: %v", err)
			continue
		}

		if err := c.processor.Handle(ctx, metric); err != nil {
			c.logger.Errorf("Failed to process metric: %v", err)
			continue
		}

		c.logger.Infof("Metric processed: %+v", metric)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
