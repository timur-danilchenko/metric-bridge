package processor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/timur-danilchenko/metric-bridge/internal/model"
	"go.uber.org/zap/zaptest"
)

// Мок для репозитория
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) SaveMetric(ctx context.Context, metric model.Metric) error {
	args := m.Called(ctx, metric)
	return args.Error(0)
}

func TestProcessor_Handle_ValidMetric(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	mockRepo := new(MockRepo)
	processor := NewProcessor(logger, mockRepo)

	ctx := context.Background()
	metric := model.Metric{
		Type:      "cpu",
		Value:     42.0,
		Timestamp: time.Now(),
	}

	mockRepo.On("SaveMetric", ctx, metric).Return(nil)

	err := processor.Handle(ctx, metric)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
