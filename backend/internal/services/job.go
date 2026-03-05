package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type jobService struct {
	repo repository.JobRepository
}

func NewJobService(repo repository.JobRepository) JobService {
	return &jobService{repo: repo}
}

func (s *jobService) CreateJob(ctx context.Context, taskType string, payload []byte) (*models.Job, error) {
	job := &models.Job{
		ID:       uuid.New().String(),
		TaskType: taskType,
		Status:   models.JobStatusPending,
		Payload:  datatypes.JSON(payload),
	}

	if err := s.repo.Create(ctx, job); err != nil {
		return nil, err
	}

	return job, nil
}

func (s *jobService) UpdateJobStatus(ctx context.Context, id string, status models.JobStatus, message string) error {
	job, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	job.Status = status
	job.Message = message

	return s.repo.Update(ctx, job)
}

func (s *jobService) GetJobs(ctx context.Context, limit, offset int) ([]models.Job, int64, error) {
	return s.repo.GetJobs(ctx, limit, offset)
}
