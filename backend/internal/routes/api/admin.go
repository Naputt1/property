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
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden (Admin only)"
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
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden (Admin only)"
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

// ResetBackend godoc
// @Summary Reset backend state
// @Description Truncate properties and jobs tables, clear queues, and clear analytics cache
// @Tags admin
// @Produce json
// @Security JwtAuth
// @Success 200 {object} BaseResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden (Admin only)"
// @Failure 500 {object} ErrorResponse
// @Router /admin/reset [post]
func ResetBackend(c *gin.Context) {
	svcs := c.MustGet("services").(*services.Services)

	// 1. Truncate properties
	if err := svcs.Property.Truncate(c.Request.Context()); err != nil {
		slog.Error("Failed to truncate properties", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to truncate properties"})
		return
	}

	// 2. Truncate jobs in DB
	if err := svcs.Job.Truncate(c.Request.Context()); err != nil {
		slog.Error("Failed to truncate jobs", "error", err)
	}

	// 3. Clear Asynq queues
	if err := svcs.Job.DeleteAllTasks(c.Request.Context(), []string{"default", "migration"}); err != nil {
		slog.Error("Failed to clear queues", "error", err)
	}

	// 4. Clear analytics cache
	if err := svcs.Analytics.ClearCache(c.Request.Context()); err != nil {
		slog.Error("Failed to clear analytics cache", "error", err)
	}

	c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Backend state reset successfully",
	})
}

// MigrateExisting godoc
// @Summary Trigger migration for existing bucket file
// @Description Queue a migration job for a file already in the bucket
// @Tags admin
// @Produce json
// @Security JwtAuth
// @Param bucketKey query string true "Key in bucket"
// @Param hasHeader query boolean false "CSV has header"
// @Success 202 {object} JobResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden (Admin only)"
// @Failure 500 {object} ErrorResponse
// @Router /admin/migrate-existing [post]
func MigrateExisting(c *gin.Context) {
	svcs := c.MustGet("services").(*services.Services)

	bucketKey := c.Query("bucketKey")
	if bucketKey == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Status: false, Error: "bucketKey is required"})
		return
	}

	hasHeader := c.DefaultQuery("hasHeader", "true") == "true"

	// Queue job
	req := models.CSVConfigPayload{
		BucketKey: bucketKey,
		HasHeader: hasHeader,
	}

	payloadBytes, err := json.Marshal(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to encode payload"})
		return
	}

	job, err := svcs.Job.CreateJob(c.Request.Context(), "properties:migrate:csv", payloadBytes)
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

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"is_admin"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"is_admin"`
}

// ListUsers godoc
// @Summary List all users
// @Description Get a list of all users. Admin only.
// @Tags admin
// @Produce json
// @Security JwtAuth
// @Success 200 {array} models.User
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse
// @Router /admin/users [get]
func ListUsers(c *gin.Context) {
	svcs := c.MustGet("services").(*services.Services)
	users, err := svcs.User.ListUsers(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to list users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with specified role. Admin only.
// @Tags admin
// @Accept json
// @Produce json
// @Security JwtAuth
// @Param body body CreateUserRequest true "User details"
// @Success 201 {object} BaseResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse
// @Router /admin/users [post]
func CreateUser(c *gin.Context) {
	svcs := c.MustGet("services").(*services.Services)
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	err := svcs.User.CreateUser(c.Request.Context(), req.Username, req.Password, req.Name, req.IsAdmin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, BaseResponse{Status: true, Message: "user created successfully"})
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update user details. Admin only.
// @Tags admin
// @Accept json
// @Produce json
// @Security JwtAuth
// @Param id path int true "User ID"
// @Param body body UpdateUserRequest true "User details"
// @Success 200 {object} BaseResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse
// @Router /admin/users/{id} [put]
func UpdateUser(c *gin.Context) {
	svcs := c.MustGet("services").(*services.Services)
	idStr := c.Param("id")
	var id int64
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Status: false, Error: "invalid user id"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	err := svcs.User.UpdateUser(c.Request.Context(), id, req.Username, req.Name, req.IsAdmin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, BaseResponse{Status: true, Message: "user updated successfully"})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID. Admin only.
// @Tags admin
// @Produce json
// @Security JwtAuth
// @Param id path int true "User ID"
// @Success 200 {object} BaseResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse
// @Router /admin/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	svcs := c.MustGet("services").(*services.Services)
	idStr := c.Param("id")
	var id int64
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Status: false, Error: "invalid user id"})
		return
	}

	err := svcs.User.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, BaseResponse{Status: true, Message: "user deleted successfully"})
}

func RegisterAdminRoutes(r *gin.RouterGroup, cfg *config.Config, svcs *services.Services) {
	r.Use(func(c *gin.Context) {
		c.Set("jobService", svcs.Job)
		c.Set("services", svcs)
		c.Next()
	})
	r.POST("/upload", UploadCSV)
	r.POST("/stream-upload", StreamUploadCSV)
	r.POST("/reset", ResetBackend)
	r.POST("/migrate-existing", MigrateExisting)

	r.GET("/users", ListUsers)
	r.POST("/users", CreateUser)
	r.PUT("/users/:id", UpdateUser)
	r.DELETE("/users/:id", DeleteUser)
}
