package mocks

import (
	"backend/internal/models"
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, username, password, name string, isAdmin bool) error {
	args := m.Called(ctx, username, password, name, isAdmin)
	return args.Error(0)
}

func (m *MockUserService) Authenticate(ctx context.Context, username, password string) (*models.User, error) {
	args := m.Called(ctx, username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id int64, username, name string, isAdmin bool) error {
	args := m.Called(ctx, id, username, name, isAdmin)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) UpdatePassword(ctx context.Context, id int64, newPassword string) error {
	args := m.Called(ctx, id, newPassword)
	return args.Error(0)
}

type MockPropertyService struct {
	mock.Mock
}

func (m *MockPropertyService) GetProperties(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Property, int64, error) {
	args := m.Called(ctx, filters, limit, offset)
	return args.Get(0).([]models.Property), args.Get(1).(int64), args.Error(2)
}

func (m *MockPropertyService) GetPropertyByID(ctx context.Context, id string) (*models.Property, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Property), args.Error(1)
}

func (m *MockPropertyService) CreateProperty(ctx context.Context, property *models.Property) error {
	args := m.Called(ctx, property)
	return args.Error(0)
}

func (m *MockPropertyService) UpdateProperty(ctx context.Context, property *models.Property) error {
	args := m.Called(ctx, property)
	return args.Error(0)
}

func (m *MockPropertyService) DeleteProperty(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPropertyService) CreateBatch(ctx context.Context, properties []models.Property, batchSize int) error {
	args := m.Called(ctx, properties, batchSize)
	return args.Error(0)
}

func (m *MockPropertyService) Truncate(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockAnalyticsService struct {
	mock.Mock
}

func (m *MockAnalyticsService) GetMedianPriceByRegion(ctx context.Context, regionType string, year int) ([]models.MedianPriceResult, error) {
	args := m.Called(ctx, regionType, year)
	return args.Get(0).([]models.MedianPriceResult), args.Error(1)
}

func (m *MockAnalyticsService) GetPriceTrend(ctx context.Context, interval string) ([]models.PriceTrendResult, error) {
	args := m.Called(ctx, interval)
	return args.Get(0).([]models.PriceTrendResult), args.Error(1)
}

func (m *MockAnalyticsService) GetAffordability(ctx context.Context) ([]models.AffordabilityResult, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.AffordabilityResult), args.Error(1)
}

func (m *MockAnalyticsService) GetGrowthHotspots(ctx context.Context, regionType string, limit, year int) ([]models.GrowthHotspotResult, error) {
	args := m.Called(ctx, regionType, limit, year)
	return args.Get(0).([]models.GrowthHotspotResult), args.Error(1)
}

func (m *MockAnalyticsService) GetNewBuildPremium(ctx context.Context, regionType string) ([]models.NewBuildPremiumResult, error) {
	args := m.Called(ctx, regionType)
	return args.Get(0).([]models.NewBuildPremiumResult), args.Error(1)
}

func (m *MockAnalyticsService) GetPropertyTypeDistribution(ctx context.Context) ([]models.PropertyTypeDistributionResult, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.PropertyTypeDistributionResult), args.Error(1)
}

func (m *MockAnalyticsService) GetPriceBracketDistribution(ctx context.Context) ([]models.PriceBracketResult, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.PriceBracketResult), args.Error(1)
}

func (m *MockAnalyticsService) GetTopActiveAreas(ctx context.Context, regionType string, limit, year int) ([]models.TopActiveAreaResult, error) {
	args := m.Called(ctx, regionType, limit, year)
	return args.Get(0).([]models.TopActiveAreaResult), args.Error(1)
}

func (m *MockAnalyticsService) GetTimeRange(ctx context.Context) (*models.TimeRangeResult, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TimeRangeResult), args.Error(1)
}

func (m *MockAnalyticsService) RefreshMaterializedView(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAnalyticsService) PrecomputeCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAnalyticsService) ClearCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockJobService struct {
	mock.Mock
}

func (m *MockJobService) CreateJob(ctx context.Context, taskType string, payload []byte) (*models.Job, error) {
	args := m.Called(ctx, taskType, payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Job), args.Error(1)
}

func (m *MockJobService) EnqueueAnalyticsRefresh(ctx context.Context, delay time.Duration) error {
	args := m.Called(ctx, delay)
	return args.Error(0)
}

func (m *MockJobService) UpdateJobStatus(ctx context.Context, id string, status models.JobStatus, message string) error {
	args := m.Called(ctx, id, status, message)
	return args.Error(0)
}

func (m *MockJobService) UpdateJobProgress(ctx context.Context, id string, progress, total int) error {
	args := m.Called(ctx, id, progress, total)
	return args.Error(0)
}

func (m *MockJobService) GetJobs(ctx context.Context, limit, offset int) ([]models.Job, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]models.Job), args.Get(1).(int64), args.Error(2)
}

func (m *MockJobService) GetJobByID(ctx context.Context, id string) (*models.Job, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Job), args.Error(1)
}

func (m *MockJobService) Truncate(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockJobService) DeleteAllTasks(ctx context.Context, queues []string) error {
	args := m.Called(ctx, queues)
	return args.Error(0)
}
