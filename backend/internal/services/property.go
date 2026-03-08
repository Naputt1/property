package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type propertyService struct {
	repo      repository.PropertyRepository
	analytics AnalyticsService
	job       JobService
}

func NewPropertyService(repo repository.PropertyRepository, analytics AnalyticsService, job JobService) PropertyService {
	return &propertyService{repo: repo, analytics: analytics, job: job}
}

func (s *propertyService) GetProperties(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Property, int64, error) {
	return s.repo.GetProperties(ctx, filters, limit, offset)
}

func (s *propertyService) GetPropertyByID(ctx context.Context, id string) (*models.Property, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *propertyService) CreateBatch(ctx context.Context, properties []models.Property, batchSize int) error {
	return s.repo.CreateBatch(ctx, properties, batchSize)
}

func (s *propertyService) CreateProperty(ctx context.Context, property *models.Property) error {
	if property.ID == uuid.Nil {
		property.ID = uuid.New()
	}
	err := s.repo.Create(ctx, property)
	if err == nil {
		go func() {
			_ = s.analytics.ClearCache(context.Background())
			_ = s.job.EnqueueAnalyticsRefresh(context.Background(), 30*time.Second)
		}()
	}
	return err
}

func (s *propertyService) UpdateProperty(ctx context.Context, property *models.Property) error {
	err := s.repo.Update(ctx, property)
	if err == nil {
		go func() {
			_ = s.analytics.ClearCache(context.Background())
			_ = s.job.EnqueueAnalyticsRefresh(context.Background(), 30*time.Second)
		}()
	}
	return err
}

func (s *propertyService) DeleteProperty(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err == nil {
		go func() {
			_ = s.analytics.ClearCache(context.Background())
			_ = s.job.EnqueueAnalyticsRefresh(context.Background(), 30*time.Second)
		}()
	}
	return err
}
