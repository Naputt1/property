package services

import (
	"backend/internal/models"
	"context"
)

type UserService interface {
	CreateUser(ctx context.Context, username, password string, isAdmin bool) error
	Authenticate(ctx context.Context, username, password string) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
}

type PropertyService interface {
	GetProperties(ctx context.Context, limit, offset int) ([]models.Property, int64, error)
	GetPropertyByID(ctx context.Context, id string) (*models.Property, error)
	CreateBatch(ctx context.Context, properties []models.Property, batchSize int) error
}

type JobService interface {
	CreateJob(ctx context.Context, taskType string, payload []byte) (*models.Job, error)
	UpdateJobStatus(ctx context.Context, id string, status models.JobStatus, message string) error
	GetJobs(ctx context.Context, limit, offset int) ([]models.Job, int64, error)
}
