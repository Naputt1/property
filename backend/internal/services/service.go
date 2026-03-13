package services

import (
	"backend/internal/models"
	"context"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, username, password, name string, isAdmin bool) error
	Authenticate(ctx context.Context, username, password string) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	ListUsers(ctx context.Context) ([]*models.User, error)
	UpdateUser(ctx context.Context, id int64, username, name string, isAdmin bool) error
	DeleteUser(ctx context.Context, id int64) error
	UpdatePassword(ctx context.Context, id int64, newPassword string) error
}

type PropertyService interface {
	GetProperties(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Property, int64, error)
	GetPropertyByID(ctx context.Context, id string) (*models.Property, error)
	CreateProperty(ctx context.Context, property *models.Property) error
	UpdateProperty(ctx context.Context, property *models.Property) error
	DeleteProperty(ctx context.Context, id string) error
	CreateBatch(ctx context.Context, properties []models.Property, batchSize int) error
	Truncate(ctx context.Context) error
}

type JobService interface {
	CreateJob(ctx context.Context, taskType string, payload []byte) (*models.Job, error)
	EnqueueAnalyticsRefresh(ctx context.Context, delay time.Duration) error
	UpdateJobStatus(ctx context.Context, id string, status models.JobStatus, message string) error
	UpdateJobProgress(ctx context.Context, id string, progress, total int) error
	GetJobs(ctx context.Context, limit, offset int) ([]models.Job, int64, error)
	GetJobByID(ctx context.Context, id string) (*models.Job, error)
	Truncate(ctx context.Context) error
	DeleteAllTasks(ctx context.Context, queues []string) error
}

type AnalyticsService interface {
	GetMedianPriceByRegion(ctx context.Context, regionType string, year int) ([]models.MedianPriceResult, error)
	GetPriceTrend(ctx context.Context, interval string) ([]models.PriceTrendResult, error)
	GetAffordability(ctx context.Context) ([]models.AffordabilityResult, error)
	GetGrowthHotspots(ctx context.Context, regionType string, limit, year int) ([]models.GrowthHotspotResult, error)
	GetNewBuildPremium(ctx context.Context, regionType string) ([]models.NewBuildPremiumResult, error)
	GetPropertyTypeDistribution(ctx context.Context) ([]models.PropertyTypeDistributionResult, error)
	GetPriceBracketDistribution(ctx context.Context) ([]models.PriceBracketResult, error)
	GetTopActiveAreas(ctx context.Context, regionType string, limit, year int) ([]models.TopActiveAreaResult, error)
	GetTimeRange(ctx context.Context) (*models.TimeRangeResult, error)
	RefreshMaterializedView(ctx context.Context) error
	PrecomputeCache(ctx context.Context) error
	ClearCache(ctx context.Context) error
}
