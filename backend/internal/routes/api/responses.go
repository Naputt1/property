package api

import "backend/internal/models"

// BaseResponse is the standard API response format
type BaseResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message,omitempty"`
}

// PropertyResponse represents a single property response
type PropertyResponse struct {
	BaseResponse
	Data *models.Property `json:"data,omitempty"`
}

// PropertyListResponse represents a list of properties response
type PropertyListResponse struct {
	BaseResponse
	Data  []models.Property `json:"data"`
	Total int64             `json:"total"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	Status bool   `json:"status" example:"true"`
	Token  string `json:"token,omitempty"`
}

// ErrorResponse represents a failure API response
type ErrorResponse struct {
	Status bool   `json:"status" example:"false"`
	Error  string `json:"error"`
}

// JobResponse represents a response after creating a background job
type JobResponse struct {
	BaseResponse
	JobID string `json:"job_id,omitempty"`
}

// JobListResponse represents a list of background jobs response
type JobListResponse struct {
	BaseResponse
	Data  []models.Job `json:"data"`
	Total int64        `json:"total"`
}

// MedianPriceResponse represents a median price analytics response
type MedianPriceResponse struct {
	BaseResponse
	Data []models.MedianPriceResult `json:"data"`
}

// PriceTrendResponse represents a price trend analytics response
type PriceTrendResponse struct {
	BaseResponse
	Data []models.PriceTrendResult `json:"data"`
}

// AffordabilityResponse represents an affordability analytics response
type AffordabilityResponse struct {
	BaseResponse
	Data []models.AffordabilityResult `json:"data"`
}

// GrowthHotspotResponse represents a growth hotspots analytics response
type GrowthHotspotResponse struct {
	BaseResponse
	Data []models.GrowthHotspotResult `json:"data"`
}

type NewBuildPremiumResponse struct {
	BaseResponse
	Data []models.NewBuildPremiumResult `json:"data"`
}

type PropertyTypeDistributionResponse struct {
	BaseResponse
	Data []models.PropertyTypeDistributionResult `json:"data"`
}

type PriceBracketResponse struct {
	BaseResponse
	Data []models.PriceBracketResult `json:"data"`
}

type TopActiveAreaResponse struct {
	BaseResponse
	Data []models.TopActiveAreaResult `json:"data"`
}
