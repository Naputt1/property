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
	UpdateJobProgress(ctx context.Context, id string, progress, total int) error
	GetJobs(ctx context.Context, limit, offset int) ([]models.Job, int64, error)
	GetJobByID(ctx context.Context, id string) (*models.Job, error)
}

type AnalyticsService interface {
	GetMedianPriceByRegion(ctx context.Context, regionType string) ([]models.MedianPriceResult, error)
	GetPriceTrend(ctx context.Context, interval string) ([]models.PriceTrendResult, error)
	GetAffordability(ctx context.Context) ([]models.AffordabilityResult, error)
	GetGrowthHotspots(ctx context.Context, limit int) ([]models.GrowthHotspotResult, error)
	PrecomputeCache(ctx context.Context) error
	ClearCache(ctx context.Context) error
}
