package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/timur-danilchenko/metric-bridge/internal/config"
	"github.com/timur-danilchenko/metric-bridge/internal/kafka"
	"go.uber.org/zap"
)

var configPath = "./configs"

func main() {
	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(configPath)
	if err != nil {
		logger.Fatalf("Can't load config: %v", err)
	}

	logger.Infof("Loaded config from %s", configPath)
	logger.Info("MetricBridge worker started. Press Ctrl+C to exit.")

	consumer := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, logger)
	defer consumer.Close()

	go consumer.Start(ctx)

	<-ctx.Done()
	logger.Info("Shutting down gracefully...")
}
