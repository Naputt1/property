package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type propertyService struct {
	repo repository.PropertyRepository
}

func NewPropertyService(repo repository.PropertyRepository) PropertyService {
	return &propertyService{repo: repo}
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
