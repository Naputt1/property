package api

import "backend/internal/models"

// BaseResponse is the standard API response format
type BaseResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message,omitempty"`
}

// --- Raw Responses (used by backend implementation) ---

type PropertyResponse struct {
	BaseResponse
	Data *models.Property `json:"data,omitempty"`
}

type PropertyListResponse struct {
	BaseResponse
	Data  []models.Property `json:"data"`
	Total int64             `json:"total"`
}

type LoginResponse struct {
	Status bool         `json:"status" example:"true"`
	Token  string       `json:"token,omitempty"`
	User   *models.User `json:"user,omitempty"`
}

type ErrorResponse struct {
	Status bool   `json:"status" example:"false"`
	Error  string `json:"error"`
}

type JobResponse struct {
	BaseResponse
	JobID string `json:"job_id,omitempty"`
}

type JobListResponse struct {
	BaseResponse
	Data  []models.Job `json:"data"`
	Total int64        `json:"total"`
}

type MedianPriceResponse struct {
	BaseResponse
	Data []models.MedianPriceResult `json:"data"`
}

type PriceTrendResponse struct {
	BaseResponse
	Data []models.PriceTrendResult `json:"data"`
}

type AffordabilityResponse struct {
	BaseResponse
	Data []models.AffordabilityResult `json:"data"`
}

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

// --- Payload Types (used for Swagger documentation to match unwrapped frontend types) ---

type PropertyListPayload struct {
	Data  []models.Property `json:"data"`
	Total int64             `json:"total"`
}

type JobListPayload struct {
	Data  []models.Job `json:"data"`
	Total int64        `json:"total"`
}

type LoginPayload struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}
