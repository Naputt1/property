package mocks

import (
	"backend/internal/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context) ([]*models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockPropertyRepository struct {
	mock.Mock
}

func (m *MockPropertyRepository) Create(ctx context.Context, property *models.Property) error {
	args := m.Called(ctx, property)
	return args.Error(0)
}

func (m *MockPropertyRepository) Update(ctx context.Context, property *models.Property) error {
	args := m.Called(ctx, property)
	return args.Error(0)
}

func (m *MockPropertyRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPropertyRepository) CreateBatch(ctx context.Context, properties []models.Property, batchSize int) error {
	args := m.Called(ctx, properties, batchSize)
	return args.Error(0)
}

func (m *MockPropertyRepository) GetByID(ctx context.Context, id string) (*models.Property, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Property), args.Error(1)
}

func (m *MockPropertyRepository) GetProperties(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Property, int64, error) {
	args := m.Called(ctx, filters, limit, offset)
	return args.Get(0).([]models.Property), args.Get(1).(int64), args.Error(2)
}

func (m *MockPropertyRepository) Truncate(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockJobRepository struct {
	mock.Mock
}

func (m *MockJobRepository) Create(ctx context.Context, job *models.Job) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *MockJobRepository) Update(ctx context.Context, job *models.Job) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *MockJobRepository) UpdateProgress(ctx context.Context, id string, progress, total int) error {
	args := m.Called(ctx, id, progress, total)
	return args.Error(0)
}

func (m *MockJobRepository) GetByID(ctx context.Context, id string) (*models.Job, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Job), args.Error(1)
}

func (m *MockJobRepository) GetJobs(ctx context.Context, limit, offset int) ([]models.Job, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]models.Job), args.Get(1).(int64), args.Error(2)
}

func (m *MockJobRepository) GetPendingOrRunningJobsCount(ctx context.Context, taskType string) (int64, error) {
	args := m.Called(ctx, taskType)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockJobRepository) Truncate(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
