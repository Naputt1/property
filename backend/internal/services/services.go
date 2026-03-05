package services

import "backend/internal/repository"

type Services struct {
	User     UserService
	Property PropertyService
	Job      JobService
}

type Repositories struct {
	User     repository.UserRepository
	Property repository.PropertyRepository
	Job      repository.JobRepository
}

func NewServices(repos Repositories) *Services {
	return &Services{
		User:     NewUserService(repos.User),
		Property: NewPropertyService(repos.Property),
		Job:      NewJobService(repos.Job),
	}
}
