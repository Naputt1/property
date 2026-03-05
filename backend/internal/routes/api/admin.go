package api

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func RegisterAdminRoutes(r *gin.RouterGroup, cfg *config.Config, svc services.JobService) {
	r.Use(func(c *gin.Context) {
		c.Set("jobService", svc)
		c.Next()
	})
	r.POST("/upload", UploadCSV)
}
