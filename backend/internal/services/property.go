package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"

	"github.com/google/uuid"
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

func (s *propertyService) CreateProperty(ctx context.Context, property *models.Property) error {
	if property.ID == uuid.Nil {
		property.ID = uuid.New()
	}
	return s.repo.Create(ctx, property)
}


func (s *propertyService) UpdateProperty(ctx context.Context, property *models.Property) error {
	return s.repo.Update(ctx, property)
}

func (s *propertyService) DeleteProperty(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
