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

func (r *propertyRepository) GetProperties(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Property, int64, error) {
	var properties []models.Property
	var count int64

	query := r.db.WithContext(ctx).Model(&models.Property{})

	// Apply filters
	if town, ok := filters["town_city"].(string); ok && town != "" {
		query = query.Where("town_city ILIKE ?", town+"%")
	}
	if district, ok := filters["district"].(string); ok && district != "" {
		query = query.Where("district ILIKE ?", district+"%")
	}
	if county, ok := filters["county"].(string); ok && county != "" {
		query = query.Where("county ILIKE ?", county+"%")
	}
	if ptype, ok := filters["property_type"].(string); ok && ptype != "" {
		query = query.Where("property_type = ?", ptype)
	}
	if minPrice, ok := filters["min_price"].(int64); ok && minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice, ok := filters["max_price"].(int64); ok && maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}

	// Optimized count: If no filters are applied, use estimated count from postgres stats
	if len(filters) == 0 {
		var estimatedCount int64
		r.db.Raw("SELECT reltuples::bigint FROM pg_class WHERE relname = 'properties'").Scan(&estimatedCount)
		if estimatedCount > 0 {
			count = estimatedCount
		} else {
			query.Count(&count)
		}
	} else {
		// With filters, we need exact count (or we could use EXPLAIN to estimate, but it's more complex)
		query.Count(&count)
	}

	err := query.Limit(limit).Offset(offset).Order("date_of_transfer DESC, id DESC").Find(&properties).Error
	if err != nil {
		return nil, 0, err
	}

	return properties, count, nil
}
