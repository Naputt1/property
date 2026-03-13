package services

import (
	"backend/internal/mocks"
	"backend/internal/models"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPropertyService_CreateProperty(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockPropertyRepository)
		mockAnalytics := new(mocks.MockAnalyticsService)
		mockJob := new(mocks.MockJobService)
		service := NewPropertyService(mockRepo, mockAnalytics, mockJob)

		property := &models.Property{
			Address: "123 Test St",
		}

		mockRepo.On("Create", ctx, mock.MatchedBy(func(p *models.Property) bool {
			return p.Address == "123 Test St" && p.ID != uuid.Nil
		})).Return(nil)

		// These are called in a goroutine, so we might need a small sleep to verify
		mockAnalytics.On("ClearCache", mock.Anything).Return(nil)
		mockJob.On("EnqueueAnalyticsRefresh", mock.Anything, 30*time.Second).Return(nil)

		err := service.CreateProperty(ctx, property)
		assert.NoError(t, err)

		// Wait a bit for goroutine
		time.Sleep(100 * time.Millisecond)

		mockRepo.AssertExpectations(t)
		mockAnalytics.AssertExpectations(t)
		mockJob.AssertExpectations(t)
	})
}

func TestPropertyService_GetProperties(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockPropertyRepository)
	mockAnalytics := new(mocks.MockAnalyticsService)
	mockJob := new(mocks.MockJobService)
	service := NewPropertyService(mockRepo, mockAnalytics, mockJob)

	t.Run("success", func(t *testing.T) {
		properties := []models.Property{{Address: "A"}, {Address: "B"}}
		filters := map[string]interface{}{"city": "London"}
		mockRepo.On("GetProperties", ctx, filters, 10, 0).Return(properties, int64(2), nil)

		res, total, err := service.GetProperties(ctx, filters, 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, res, 2)
		mockRepo.AssertExpectations(t)
	})
}
