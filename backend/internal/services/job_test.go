package services

import (
	"backend/internal/mocks"
	"backend/internal/models"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobService_UpdateJobStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockJobRepository)
		service := NewJobService(mockRepo, nil, nil)

		job := &models.Job{ID: "job-1", Status: models.JobStatusPending}
		mockRepo.On("GetByID", ctx, "job-1").Return(job, nil)
		mockRepo.On("Update", ctx, job).Return(nil)

		err := service.UpdateJobStatus(ctx, "job-1", models.JobStatusSuccess, "done")
		assert.NoError(t, err)
		assert.Equal(t, models.JobStatusSuccess, job.Status)
		assert.Equal(t, "done", job.Message)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.MockJobRepository)
		service := NewJobService(mockRepo, nil, nil)

		mockRepo.On("GetByID", ctx, "job-2").Return(nil, errors.New("not found"))

		err := service.UpdateJobStatus(ctx, "job-2", models.JobStatusSuccess, "done")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestJobService_GetJobs(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockJobRepository)
	service := NewJobService(mockRepo, nil, nil)

	t.Run("success", func(t *testing.T) {
		jobs := []models.Job{{ID: "1"}, {ID: "2"}}
		mockRepo.On("GetJobs", ctx, 10, 0).Return(jobs, int64(2), nil)

		res, total, err := service.GetJobs(ctx, 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, res, 2)
		mockRepo.AssertExpectations(t)
	})
}
