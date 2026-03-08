package api

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StreamUploadCSV godoc
// @Summary Stream upload CSV file
// @Description Stream a large CSV file directly from the request body to bucket and queue a migration job
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
	cfg := c.MustGet("config").(*config.Config)

	originalFilename := c.DefaultQuery("filename", "upload.csv")

	// Create unique key for the bucket
	bucketKey := fmt.Sprintf("uploads/%s-%s", uuid.New().String(), originalFilename)

	slog.Info("Starting to stream upload to bucket", "filename", originalFilename, "bucketKey", bucketKey)

	// Stream directly from request body to bucket
	contentType := "text/csv"
	if c.Request.Header.Get("Content-Type") != "" {
		contentType = c.Request.Header.Get("Content-Type")
	}

	size := c.Request.ContentLength
	if size <= 0 {
		size = -1
	}

	err := cfg.Bucket.Upload(c.Request.Context(), bucketKey, c.Request.Body, size, contentType)
	if err != nil {
		slog.Error("Failed during streaming to bucket", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed during streaming to bucket"})
		return
	}

	slog.Info("Stream upload to bucket completed, queuing migration job", "bucketKey", bucketKey)

	// Queue job
	req := models.CSVConfigPayload{
		BucketKey: bucketKey,
		HasHeader: true,
	}

	payloadBytes, err := json.Marshal(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to encode payload"})
		return
	}

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

// UploadCSV godoc
// @Summary Upload CSV file directly
// @Description Upload a CSV file via multipart/form-data to bucket and queue a migration job
// @Tags admin
// @Accept multipart/form-data
// @Produce json
// @Security JwtAuth
// @Param file formData file true "CSV File"
// @Success 202 {object} JobResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/upload [post]
func UploadCSV(c *gin.Context) {
	svc := c.MustGet("jobService").(services.JobService)
	cfg := c.MustGet("config").(*config.Config)

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Status: false, Error: "no file uploaded"})
		return
	}

	// Create unique key for the bucket
	bucketKey := fmt.Sprintf("uploads/%s-%s", uuid.New().String(), file.Filename)

	src, err := file.Open()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to open uploaded file"})
		return
	}
	defer src.Close()

	err = cfg.Bucket.Upload(c.Request.Context(), bucketKey, src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		slog.Error("Failed to upload to bucket", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to upload to bucket"})
		return
	}

	// Queue job
	req := models.CSVConfigPayload{
		BucketKey: bucketKey,
		HasHeader: true, // Default to true
	}

	payloadBytes, err := json.Marshal(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to encode payload"})
		return
	}

	job, err := svc.CreateJob(c.Request.Context(), "properties:migrate:csv", payloadBytes)
	if err != nil {
		slog.Error("Failed to create migration job", "error", err)
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

func RegisterAdminRoutes(r *gin.RouterGroup, cfg *config.Config, svc services.JobService) {
	r.Use(func(c *gin.Context) {
		c.Set("jobService", svc)
		c.Next()
	})
	r.POST("/upload", UploadCSV)
	r.POST("/stream-upload", StreamUploadCSV)
}
