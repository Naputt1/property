package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"backend/internal/config"
	"backend/internal/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	Config          *config.Config
	PropertyService services.PropertyService
	JobService      services.JobService
}
