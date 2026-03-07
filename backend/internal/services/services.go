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
	Socket    repository.SocketService
}

type Repositories struct {
	User      repository.UserRepository
	Property  repository.PropertyRepository
	Job       repository.JobRepository
	Analytics repository.AnalyticsRepository
}

func NewServices(cfg *config.Config, repos Repositories, asynqClient *asynq.Client, socket repository.SocketService) *Services {
	analyticsSvc := NewAnalyticsService(cfg, repos.Analytics)
	jobSvc := NewJobService(repos.Job, asynqClient)
	return &Services{
		User:      NewUserService(repos.User),
		Property:  NewPropertyService(repos.Property, analyticsSvc, jobSvc),
		Job:       jobSvc,
		Analytics: analyticsSvc,
		Socket:    socket,
	}
}
