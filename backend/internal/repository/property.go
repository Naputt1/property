package repository

import (
	"backend/internal/models"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type propertyRepository struct {
	db *gorm.DB
}

func NewPropertyRepository(db *gorm.DB) PropertyRepository {
	return &propertyRepository{db: db}
}

func (r *propertyRepository) Create(ctx context.Context, property *models.Property) error {
	return r.db.WithContext(ctx).Create(property).Error
}

func (r *propertyRepository) CreateBatch(ctx context.Context, properties []models.Property, batchSize int) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(properties, batchSize).Error
}

func (r *propertyRepository) GetByID(ctx context.Context, id string) (*models.Property, error) {
	var property models.Property
	err := r.db.WithContext(ctx).First(&property, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &property, nil
}

func (r *propertyRepository) GetProperties(ctx context.Context, limit, offset int) ([]models.Property, int64, error) {
	var properties []models.Property
	var count int64

	err := r.db.WithContext(ctx).Model(&models.Property{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("date_of_transfer DESC").Find(&properties).Error
	if err != nil {
		return nil, 0, err
	}

	return properties, count, nil
}
