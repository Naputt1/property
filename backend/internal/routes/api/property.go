package api

import (
	"backend/internal/config"
	"backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListProperties godoc
// @Summary List properties
// @Description Get a list of UK housing properties with pagination and filtering
// @Tags property
// @Accept json
// @Produce json
// @Security JwtAuth
// @Success 200 {object} PropertyListResponse
// @Router /property [get]
func ListProperties(c *gin.Context) {
	// Example usage of svc
	c.JSON(http.StatusOK, PropertyListResponse{
		BaseResponse: BaseResponse{Status: true},
		Data:         []interface{}{},
	})
}

// GetProperty godoc
// @Summary Get property by ID
// @Description Get detailed information for a specific property
// @Tags property
// @Accept json
// @Produce json
// @Security JwtAuth
// @Param id path int true "Property ID"
// @Success 200 {object} PropertyResponse
// @Failure 400 {object} ErrorResponse
// @Router /property/{id} [get]
func GetProperty(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Status: false, Error: "invalid id"})
		return
	}
	c.JSON(http.StatusOK, PropertyResponse{
		BaseResponse: BaseResponse{Status: true},
		Data:         map[string]interface{}{"id": id},
	})
}

func RegisterPropertyRoutes(r *gin.RouterGroup, cfg *config.Config, svc services.PropertyService) {
	r.GET("/", ListProperties)
	r.GET("/:id", GetProperty)
}
