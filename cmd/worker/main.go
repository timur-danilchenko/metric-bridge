package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/timur-danilchenko/metric-bridge/internal/config"
	"github.com/timur-danilchenko/metric-bridge/internal/kafka"
	"github.com/timur-danilchenko/metric-bridge/internal/processor"
	"github.com/timur-danilchenko/metric-bridge/internal/storage"
	"go.uber.org/zap"
)

var configPath = "./configs"

func main() {
	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Inintializing config
	cfg, err := config.Load(configPath)
	if err != nil {
		logger.Fatalf("Can't load config: %v", err)
	}
	logger.Infof("Loaded config from %s", configPath)

	// Initializing repository
	repo, err := storage.NewRepository(ctx, cfg.Postgres)
	if err != nil {
		logger.Fatalf("Failed to init repository", err)
	}
	defer repo.Close(ctx)

	// Initializing processor
	processor := processor.NewProcessor(logger, repo)

	logger.Info("MetricBridge worker started. Press Ctrl+C to exit.")

	// Creating consumer
	consumer := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.Topic,
		logger,
		processor,
	)
	defer consumer.Close()

	go consumer.Start(ctx)

	<-ctx.Done()
	logger.Info("Shutting down gracefully...")
}
