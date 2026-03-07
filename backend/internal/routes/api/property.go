package api

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/routes/middlewares"
	"backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PropertyHandler struct {
	svc services.PropertyService
}

func RegisterPropertyRoutes(r *gin.RouterGroup, cfg *config.Config, svc services.PropertyService) {
	h := &PropertyHandler{svc: svc}

	r.GET("", h.ListProperties)
	r.GET("/:id", h.GetProperty)

	// Admin only routes
	admin := r.Group("")
	admin.Use(middlewares.AdminMiddleware())
	{
		admin.POST("", h.CreateProperty)
		admin.PUT("/:id", h.UpdateProperty)
		admin.DELETE("/:id", h.DeleteProperty)
	}
}

// ListProperties godoc
// @Summary List properties
// @Description Get a list of UK housing properties with pagination and filtering
// @Tags property
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param town_city query string false "Town/City"
// @Param district query string false "District"
// @Param county query string false "County"
// @Param property_type query string false "Property Type"
// @Param min_price query int false "Minimum Price"
// @Param max_price query int false "Maximum Price"
// @Success 200 {object} PropertyListPayload
// @Failure 500 {object} ErrorResponse
// @Router /property [get]
func (h *PropertyHandler) ListProperties(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	filters := make(map[string]interface{})
	if town := c.Query("town_city"); town != "" {
		filters["town_city"] = town
	}
	if district := c.Query("district"); district != "" {
		filters["district"] = district
	}
	if county := c.Query("county"); county != "" {
		filters["county"] = county
	}
	if ptype := c.Query("property_type"); ptype != "" {
		filters["property_type"] = ptype
	}
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseInt(minPriceStr, 10, 64); err == nil {
			filters["min_price"] = minPrice
		}
	}
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseInt(maxPriceStr, 10, 64); err == nil {
			filters["max_price"] = maxPrice
		}
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	limit := pageSize

	properties, total, err := h.svc.GetProperties(c.Request.Context(), filters, limit, offset)
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
// @Param id path string true "Property ID"
// @Success 200 {object} models.Property
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

// CreateProperty godoc
// @Summary Create property
// @Description Create a new property record
// @Tags property
// @Accept json
// @Produce json
// @Security JwtAuth
// @Param property body models.Property true "Property object"
// @Success 201 {object} models.Property
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /property [post]
func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	var property models.Property
	if err := c.ShouldBindJSON(&property); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	if err := h.svc.CreateProperty(c.Request.Context(), &property); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, PropertyResponse{
		BaseResponse: BaseResponse{Status: true},
		Data:         &property,
	})
}

// UpdateProperty godoc
// @Summary Update property
// @Description Update an existing property record
// @Tags property
// @Accept json
// @Produce json
// @Security JwtAuth
// @Param id path string true "Property ID"
// @Param property body models.Property true "Property object"
// @Success 200 {object} models.Property
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /property/{id} [put]
func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Status: false, Error: "invalid property id"})
		return
	}

	var property models.Property
	if err := c.ShouldBindJSON(&property); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	property.ID = id
	if err := h.svc.UpdateProperty(c.Request.Context(), &property); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PropertyResponse{
		BaseResponse: BaseResponse{Status: true},
		Data:         &property,
	})
}

// DeleteProperty godoc
// @Summary Delete property
// @Description Delete a property record
// @Tags property
// @Accept json
// @Produce json
// @Security JwtAuth
// @Param id path string true "Property ID"
// @Success 200 {object} BaseResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /property/{id} [delete]
func (h *PropertyHandler) DeleteProperty(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.DeleteProperty(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "property deleted successfully",
	})
}
