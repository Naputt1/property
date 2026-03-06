package services

import (
	"backend/internal/config"
	"backend/internal/repository"

	"github.com/hibiken/asynq"
)

type Services struct {
	User      UserService
	Property  PropertyService
	Job       JobService
	Analytics AnalyticsService
}

type Repositories struct {
	User      repository.UserRepository
	Property  repository.PropertyRepository
	Job       repository.JobRepository
	Analytics repository.AnalyticsRepository
}

func NewServices(cfg *config.Config, repos Repositories, asynqClient *asynq.Client) *Services {
	return &Services{
		User:      NewUserService(repos.User),
		Property:  NewPropertyService(repos.Property),
		Job:       NewJobService(repos.Job, asynqClient),
		Analytics: NewAnalyticsService(cfg, repos.Analytics),
	}
}
