//go:build integration

package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

// getKafkaBroker возвращает адрес Kafka брокера из env или дефолтное значение
func getKafkaBroker() string {
	broker := os.Getenv("KAFKA_BROKERS")
	if broker == "" {
		return "localhost:9092"
	}
	return broker
}

func TestKafkaToPostgresFlow(t *testing.T) {
	ctx := context.Background()

	brokerAddr := getKafkaBroker()

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddr),
		Topic:    "metrics",
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	payload := `{
		"type": "cpu",
		"value": 99.9,
		"timestamp": "2025-06-25T12:34:56Z"
	}`

	err := writer.WriteMessages(ctx, kafka.Message{
		Value: []byte(payload),
	})
	assert.NoError(t, err, "failed to send metric to Kafka")

	dsn := os.Getenv("DB_INTEGRATION_DSN")
	assert.NotEmpty(t, dsn)

	pool, err := pgxpool.New(ctx, dsn)
	assert.NoError(t, err)
	defer pool.Close()

	// Retry check для появления данных в базе
	found := false
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)

		var count int
		err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM metrics WHERE type='cpu' AND value=99.9`).Scan(&count)
		if err == nil && count > 0 {
			found = true
			break
		}
	}
	assert.True(t, found, "metric not found in Postgres")
}
