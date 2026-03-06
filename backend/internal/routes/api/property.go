package api

import (
	"backend/internal/config"
	"backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PropertyHandler struct {
	svc services.PropertyService
}

func RegisterPropertyRoutes(r *gin.RouterGroup, cfg *config.Config, svc services.PropertyService) {
	h := &PropertyHandler{svc: svc}

	r.GET("", h.ListProperties)
	r.GET("/:id", h.GetProperty)
}

// ListProperties godoc
// @Summary List properties
// @Description Get a list of UK housing properties with pagination and filtering
// @Tags property
// @Accept json
// @Produce json
// @Security JwtAuth
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} PropertyListResponse
// @Failure 500 {object} ErrorResponse
// @Router /property [get]
func (h *PropertyHandler) ListProperties(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	limit := pageSize

	properties, total, err := h.svc.GetProperties(c.Request.Context(), limit, offset)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: "failed to list properties"})
		return
	}

	c.JSON(http.StatusOK, PropertyListResponse{
		BaseResponse: BaseResponse{Status: true},
		Data:         properties,
		Total:        total,
	})
}

// GetProperty godoc
// @Summary Get property by ID
// @Description Get detailed information for a specific property
// @Tags property
// @Accept json
// @Produce json
// @Security JwtAuth
// @Param id path string true "Property ID"
// @Success 200 {object} PropertyResponse
// @Failure 404 {object} ErrorResponse
// @Router /property/{id} [get]
func (h *PropertyHandler) GetProperty(c *gin.Context) {
	id := c.Param("id")
	
	property, err := h.svc.GetPropertyByID(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Status: false, Error: "property not found"})
		return
	}

	c.JSON(http.StatusOK, PropertyResponse{
		BaseResponse: BaseResponse{Status: true},
		Data:         property,
	})
}
