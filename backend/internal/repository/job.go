package repository

import (
	"backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Create(ctx context.Context, job *models.Job) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *jobRepository) Update(ctx context.Context, job *models.Job) error {
	return r.db.WithContext(ctx).Save(job).Error
}

func (r *jobRepository) UpdateProgress(ctx context.Context, id string, progress, total int) error {
	return r.db.WithContext(ctx).Model(&models.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"progress": progress,
		"total":    total,
	}).Error
}

func (r *jobRepository) GetByID(ctx context.Context, id string) (*models.Job, error) {
	var job models.Job
	err := r.db.WithContext(ctx).First(&job, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) GetJobs(ctx context.Context, limit, offset int) ([]models.Job, int64, error) {
	var jobs []models.Job
	var count int64

	err := r.db.WithContext(ctx).Model(&models.Job{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&jobs).Error
	if err != nil {
		return nil, 0, err
	}

	return jobs, count, nil
}
