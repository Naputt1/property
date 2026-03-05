package repository

import (
	"backend/internal/models"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type PropertyRepository interface {
	Create(ctx context.Context, property *models.Property) error
	CreateBatch(ctx context.Context, properties []models.Property, batchSize int) error
	GetByID(ctx context.Context, id string) (*models.Property, error)
	GetProperties(ctx context.Context, limit, offset int) ([]models.Property, int64, error)
}

type JobRepository interface {
	Create(ctx context.Context, job *models.Job) error
	Update(ctx context.Context, job *models.Job) error
	GetByID(ctx context.Context, id string) (*models.Job, error)
	GetJobs(ctx context.Context, limit, offset int) ([]models.Job, int64, error)
}
