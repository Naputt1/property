package api

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StreamUploadCSV godoc
// @Summary Stream upload CSV file
// @Description Stream a large CSV file directly from the request body to disk and queue a migration job
// @Tags admin
// @Accept octet-stream
// @Produce json
// @Security JwtAuth
// @Param filename query string true "Original filename"
// @Success 202 {object} JobResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/stream-upload [post]
func StreamUploadCSV(c *gin.Context) {
	svc := c.MustGet("jobService").(services.JobService)
	
	originalFilename := c.DefaultQuery("filename", "upload.csv")

	// Create temp directory if not exists
	tempDir := "tmp/uploads"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		slog.Error("Failed to create upload directory", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to create upload directory"})
		return
	}

	// Save file with unique name
	filename := fmt.Sprintf("%s-%s", uuid.New().String(), originalFilename)
	filePath := filepath.Join(tempDir, filename)
	
	out, err := os.Create(filePath)
	if err != nil {
		slog.Error("Failed to create local file", "error", err, "path", filePath)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to create local file"})
		return
	}
	defer out.Close()

	slog.Info("Starting to stream upload to disk", "filename", originalFilename, "path", filePath)

	// Stream directly from request body to file
	// This is synchronous to ensure file integrity before job creation
	_, err = io.Copy(out, c.Request.Body)
	if err != nil {
		slog.Error("Failed during streaming to disk", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed during streaming to disk"})
		return
	}

	// Get absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to get absolute path"})
		return
	}

	slog.Info("Stream upload completed, queuing migration job", "path", absPath)

	// Queue job
	req := models.CSVConfigPayload{
		FilePath:  absPath,
		HasHeader: true,
	}

	payloadBytes, err := json.Marshal(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to encode payload"})
		return
	}

	// svc.CreateJob now properly enqueues to asynq
	job, err := svc.CreateJob(c.Request.Context(), "properties:migrate:csv", payloadBytes)
	if err != nil {
		slog.Error("Failed to create migration job", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to create migration job"})
		return
	}

	c.JSON(http.StatusAccepted, JobResponse{
		BaseResponse: BaseResponse{
			Status:  true,
			Message: "Stream migration job queued successfully",
		},
		JobID: job.ID,
	})
}

// UploadCSVFile godoc
// @Summary Upload CSV file directly
// @Description Upload a CSV file via multipart/form-data and queue a migration job
// @Tags admin
// @Accept multipart/form-data
// @Produce json
// @Security JwtAuth
// @Param file formData file true "CSV File"
// @Success 202 {object} JobResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/upload-file [post]
func UploadCSVFile(c *gin.Context) {
	svc := c.MustGet("jobService").(services.JobService)

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Status: false, Error: "no file uploaded"})
		return
	}

	// Create temp directory if not exists
	tempDir := "tmp/uploads"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to create upload directory"})
		return
	}

	// Save file with unique name
	filename := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)
	filePath := filepath.Join(tempDir, filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to save uploaded file"})
		return
	}

	// Get absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to get absolute path"})
		return
	}

	// Queue job
	req := models.CSVConfigPayload{
		FilePath:  absPath,
		HasHeader: true, // Default to true
	}

	payloadBytes, err := json.Marshal(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to encode payload"})
		return
	}

	job, err := svc.CreateJob(c.Request.Context(), "properties:migrate:csv", payloadBytes)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to create migration job"})
		return
	}

	c.JSON(http.StatusAccepted, JobResponse{
		BaseResponse: BaseResponse{
			Status:  true,
			Message: "Migration job queued successfully",
		},
		JobID: job.ID,
	})
}

// UploadCSV godoc
// @Summary Upload CSV for migration
// @Description Queue a background job to migrate properties from a CSV file
// @Tags admin
// @Accept json
// @Produce json
// @Security JwtAuth
// @Param body body models.CSVConfigPayload true "CSV Migration Configuration"
// @Success 202 {object} JobResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/upload [post]
func UploadCSV(c *gin.Context) {
	svc := c.MustGet("jobService").(services.JobService)
	var req models.CSVConfigPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	payloadBytes, err := json.Marshal(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to encode payload"})
		return
	}

	job, err := svc.CreateJob(c.Request.Context(), "properties:migrate:csv", payloadBytes)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to create migration job"})
		return
	}

	c.JSON(http.StatusAccepted, JobResponse{
		BaseResponse: BaseResponse{
			Status:  true,
			Message: "Migration job queued successfully",
		},
		JobID: job.ID,
	})
}

// ListJobs godoc
// @Summary List background jobs
// @Description Get a list of background jobs with pagination
// @Tags admin
// @Accept json
// @Produce json
// @Security JwtAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} JobListResponse
// @Router /admin/jobs [get]
func ListJobs(c *gin.Context) {
	svc := c.MustGet("jobService").(services.JobService)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	jobs, total, err := svc.GetJobs(c.Request.Context(), limit, offset)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to list jobs"})
		return
	}

	c.JSON(http.StatusOK, JobListResponse{
		BaseResponse: BaseResponse{Status: true},
		Data:         jobs,
		Total:        total,
	})
}

func RegisterAdminRoutes(r *gin.RouterGroup, cfg *config.Config, svc services.JobService) {
	r.Use(func(c *gin.Context) {
		c.Set("jobService", svc)
		c.Next()
	})
	r.POST("/upload", UploadCSV)
	r.POST("/upload-file", UploadCSVFile)
	r.POST("/stream-upload", StreamUploadCSV)
	r.GET("/jobs", ListJobs)
}
