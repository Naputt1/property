package repository

import (
	"backend/internal/models"
	"context"
	"io"

	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type PropertyRepository interface {
	Create(ctx context.Context, property *models.Property) error
	Update(ctx context.Context, property *models.Property) error
	Delete(ctx context.Context, id string) error
	CreateBatch(ctx context.Context, properties []models.Property, batchSize int) error
	GetByID(ctx context.Context, id string) (*models.Property, error)
	GetProperties(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Property, int64, error)
}

type JobRepository interface {
	Create(ctx context.Context, job *models.Job) error
	Update(ctx context.Context, job *models.Job) error
	UpdateProgress(ctx context.Context, id string, progress, total int) error
	GetByID(ctx context.Context, id string) (*models.Job, error)
	GetJobs(ctx context.Context, limit, offset int) ([]models.Job, int64, error)
	GetPendingOrRunningJobsCount(ctx context.Context, taskType string) (int64, error)
}

type AnalyticsRepository interface {
	GetMedianPriceByRegion(ctx context.Context, regionType string) ([]models.MedianPriceResult, error)
	GetPriceTrend(ctx context.Context, interval string) ([]models.PriceTrendResult, error)
	GetAffordability(ctx context.Context) ([]models.AffordabilityResult, error)
	GetGrowthHotspots(ctx context.Context, limit int) ([]models.GrowthHotspotResult, error)
	GetNewBuildPremium(ctx context.Context, regionType string) ([]models.NewBuildPremiumResult, error)
	GetPropertyTypeDistribution(ctx context.Context) ([]models.PropertyTypeDistributionResult, error)
	GetPriceBracketDistribution(ctx context.Context) ([]models.PriceBracketResult, error)
	GetTopActiveAreas(ctx context.Context, regionType string, limit int) ([]models.TopActiveAreaResult, error)
	RefreshMaterializedView(ctx context.Context) error
}

type SocketService interface {
	Broadcast(data any)
	Run()
	ServeWS(c *gin.Context)
}

type BucketService interface {
	Upload(ctx context.Context, key string, body io.Reader, size int64, contentType string) error
	GetObject(ctx context.Context, key string) (io.ReadCloser, int64, error)
	Delete(ctx context.Context, key string) error
	EnsureBucket(ctx context.Context) error
}
