package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type JobStatus string

const (
	JobStatusPending JobStatus = "PENDING"
	JobStatusRunning JobStatus = "RUNNING"
	JobStatusSuccess JobStatus = "SUCCESS"
	JobStatusFailed  JobStatus = "FAILED"
)

type Job struct {
	ID        string         `gorm:"primarykey;type:uuid" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	TaskType string         `json:"task_type"`
	Status   JobStatus      `json:"status"`
	Message  string         `json:"message"`          // For error messages or success info
	Payload  datatypes.JSON `json:"payload"`          // Custom data the job was initiated with
	Result   datatypes.JSON `json:"result,omitempty"` // Data resulting from completion/failure
}

func (Job) TableName() string {
	return "jobs"
}

type CSVConfigPayload struct {
	JobID         string            `json:"job_id"`
	FilePath      string            `json:"file_path" binding:"required"`
	ColumnMapping map[string]string `json:"column_mapping"`
	HasHeader     bool              `json:"has_header"`
}
