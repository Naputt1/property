package graph

import (
	"backend/internal/config"
	"backend/internal/models"
	"context"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPropertyService struct {
	mock.Mock
}

func (m *MockPropertyService) GetProperties(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Property, int64, error) {
	args := m.Called(ctx, filters, limit, offset)
	return args.Get(0).([]models.Property), args.Get(1).(int64), args.Error(2)
}

func (m *MockPropertyService) GetPropertyByID(ctx context.Context, id string) (*models.Property, error) {
	args := m.Called(ctx, id)
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

type MockJobService struct {
	mock.Mock
}

func (m *MockJobService) CreateJob(ctx context.Context, taskType string, payload []byte) (*models.Job, error) {
	args := m.Called(ctx, taskType, payload)
	return args.Get(0).(*models.Job), args.Error(1)
}

func (m *MockJobService) GetJobs(ctx context.Context, limit, offset int) ([]models.Job, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]models.Job), args.Get(1).(int64), args.Error(2)
}

func (m *MockJobService) GetJobByID(ctx context.Context, id string) (*models.Job, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Job), args.Error(1)
}

func (m *MockJobService) UpdateJobStatus(ctx context.Context, id string, status models.JobStatus, message string) error {
	args := m.Called(ctx, id, status, message)
	return args.Error(0)
}

func (m *MockJobService) UpdateJobProgress(ctx context.Context, id string, progress, total int) error {
	args := m.Called(ctx, id, progress, total)
	return args.Error(0)
}

func (m *MockJobService) EnqueueAnalyticsRefresh(ctx context.Context, delay time.Duration) error {
	args := m.Called(ctx, delay)
	return args.Error(0)
}

func TestPropertiesQuery(t *testing.T) {
	mockPropSvc := new(MockPropertyService)
	mockJobSvc := new(MockJobService)

	resolver := &Resolver{
		PropertyService: mockPropSvc,
		JobService:      mockJobSvc,
	}

	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))
	c := client.New(srv)

	t.Run("fetch properties", func(t *testing.T) {
		expectedProps := []models.Property{
			{ID: uuid.New(), Price: 100000, TownCity: "London"},
		}
		mockPropSvc.On("GetProperties", mock.Anything, mock.Anything, 10, 0).Return(expectedProps, int64(1), nil).Once()

		var resp struct {
			Properties struct {
				Items []struct {
					Id    string
					Price int
				}
				Total int
			}
		}
		c.MustPost(`query { properties(limit: 10) { items { id price } total } }`, &resp)

		assert.Equal(t, 1, len(resp.Properties.Items))
		assert.Equal(t, 100000, resp.Properties.Items[0].Price)
		assert.Equal(t, 1, resp.Properties.Total)
		mockPropSvc.AssertExpectations(t)
	})

	t.Run("fetch jobs", func(t *testing.T) {
		expectedJobs := []models.Job{
			{ID: "job-1", TaskType: "test-task", Status: models.JobStatusSuccess},
		}
		mockJobSvc.On("GetJobs", mock.Anything, 10, 0).Return(expectedJobs, int64(1), nil).Once()

		var resp struct {
			Jobs struct {
				Items []struct {
					Id       string
					TaskType string
					Status   string
				}
				Total int
			}
		}
		c.MustPost(`query { jobs(limit: 10) { items { id taskType status } total } }`, &resp)

		assert.Equal(t, 1, len(resp.Jobs.Items))
		assert.Equal(t, "SUCCESS", resp.Jobs.Items[0].Status)
		mockJobSvc.AssertExpectations(t)
	})
}

func TestMutations(t *testing.T) {
	mockPropSvc := new(MockPropertyService)
	mockJobSvc := new(MockJobService)
	cfg := &config.Config{} // Empty config for now

	resolver := &Resolver{
		Config:          cfg,
		PropertyService: mockPropSvc,
		JobService:      mockJobSvc,
	}

	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))
	c := client.New(srv)

	t.Run("delete property unauthorized", func(t *testing.T) {
		// No admin context set
		var resp struct {
			DeleteProperty bool
		}
		err := c.Post(`mutation { deleteProperty(id: "some-uuid") }`, &resp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "forbidden: admin access required")
	})
}
