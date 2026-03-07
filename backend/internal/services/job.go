package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/datatypes"
)

type jobService struct {
	repo        repository.JobRepository
	asynqClient *asynq.Client
}

func NewJobService(repo repository.JobRepository, asynqClient *asynq.Client) JobService {
	return &jobService{
		repo:        repo,
		asynqClient: asynqClient,
	}
}

func (s *jobService) CreateJob(ctx context.Context, taskType string, payload []byte) (*models.Job, error) {
	jobID := uuid.New().String()

	// If it's a CSV migration, we need to inject the JobID into the payload
	// so the handler can update the job status later.
	if taskType == "properties:migrate:csv" {
		var csvPayload models.CSVConfigPayload
		if err := json.Unmarshal(payload, &csvPayload); err == nil {
			csvPayload.JobID = jobID
			newPayload, _ := json.Marshal(csvPayload)
			payload = newPayload
		}
	}

	job := &models.Job{
		ID:       jobID,
		TaskType: taskType,
		Status:   models.JobStatusPending,
		Payload:  datatypes.JSON(payload),
	}

	if err := s.repo.Create(ctx, job); err != nil {
		return nil, err
	}

	// Enqueue to asynq
	task := asynq.NewTask(taskType, payload)

	var opts []asynq.Option
	if taskType == "properties:migrate:csv" {
		opts = append(opts, asynq.Timeout(2*time.Hour))
	}

	info, err := s.asynqClient.Enqueue(task, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue task: %w", err)
	}

	fmt.Printf("Enqueued task: id=%s queue=%s\n", info.ID, info.Queue)

	return job, nil
}

func (s *jobService) EnqueueAnalyticsRefresh(ctx context.Context, delay time.Duration) error {
	task := asynq.NewTask("analytics:refresh_mvs", nil)

	opts := []asynq.Option{
		asynq.ProcessIn(delay),
		asynq.TaskID("analytics_refresh_mv_singleton"),
		asynq.MaxRetry(3),
		asynq.Unique(delay),
	}

	_, err := s.asynqClient.Enqueue(task, opts...)
	if err != nil {
		if fmt.Sprintf("%v", err) == "task already exists" {
			return nil
		}
		return fmt.Errorf("failed to enqueue analytics refresh: %w", err)
	}

	return nil
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

func (s *jobService) UpdateJobProgress(ctx context.Context, id string, progress, total int) error {
	return s.repo.UpdateProgress(ctx, id, progress, total)
}

func (s *jobService) GetJobs(ctx context.Context, limit, offset int) ([]models.Job, int64, error) {
	return s.repo.GetJobs(ctx, limit, offset)
}

func (s *jobService) GetJobByID(ctx context.Context, id string) (*models.Job, error) {
	return s.repo.GetByID(ctx, id)
}
