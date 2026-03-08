package api

import "backend/internal/models"

// BaseResponse is the standard API response format
type BaseResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message,omitempty"`
}

// --- Raw Responses (used by backend implementation) ---

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

type LoginPayload struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// GraphQLRequest defines the structure for a GraphQL operation
type GraphQLRequest struct {
	// The GraphQL query or mutation string
	Query string `json:"query" example:"query { properties(limit: 5) { items { id price street } } }"`
	// Optional variables for the operation
	Variables map[string]interface{} `json:"variables,omitempty"`
	// Optional operation name if multiple are defined in the query
	OperationName string `json:"operationName,omitempty"`
}

// GraphQLResponse represents a standard GraphQL response
type GraphQLResponse struct {
	// The requested data
	Data interface{} `json:"data"`
	// Any errors encountered during execution
	Errors []GraphQLError `json:"errors,omitempty"`
}

// GraphQLError represents a GraphQL error entry
type GraphQLError struct {
	Message   string                 `json:"message"`
	Path      []interface{}          `json:"path,omitempty"`
	Locations []GraphQLLocation      `json:"locations,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// GraphQLLocation represents a line/column position in a GraphQL document
type GraphQLLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}
