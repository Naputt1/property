package api

// BaseResponse is the standard API response format
type BaseResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message,omitempty"`
}

// PropertyResponse represents a single property response
type PropertyResponse struct {
	BaseResponse
	Data interface{} `json:"data,omitempty"`
}

// PropertyListResponse represents a list of properties response
type PropertyListResponse struct {
	BaseResponse
	Data []interface{} `json:"data,omitempty"`
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
