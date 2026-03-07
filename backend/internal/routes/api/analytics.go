package api

import (
	"backend/internal/config"
	"backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	svc services.AnalyticsService
}

func RegisterAnalyticsRoutes(rg *gin.RouterGroup, cfg *config.Config, svc services.AnalyticsService) {
	h := &AnalyticsHandler{svc: svc}

	rg.GET("/median-price", h.GetMedianPriceByRegion)
	rg.GET("/price-trend", h.GetPriceTrend)
	rg.GET("/affordability", h.GetAffordability)
	rg.GET("/growth-hotspots", h.GetGrowthHotspots)
	rg.GET("/new-build-premium", h.GetNewBuildPremium)
	rg.GET("/property-type-distribution", h.GetPropertyTypeDistribution)
	rg.GET("/price-bracket-distribution", h.GetPriceBracketDistribution)
	rg.GET("/top-active-areas", h.GetTopActiveAreas)
}

// GetMedianPriceByRegion godoc
// @Summary Get median price by region
// @Description Get median price grouped by county, district, or town_city
// @Tags analytics
// @Accept json
// @Produce json
// @Param by query string false "Region type (county, district, town_city)" default(county)
// @Success 200 {array} backend_internal_models.MedianPriceResult
// @Failure 500 {object} ErrorResponse
// @Router /analytics/median-price [get]
func (h *AnalyticsHandler) GetMedianPriceByRegion(c *gin.Context) {
	regionType := c.DefaultQuery("by", "county")
	results, err := h.svc.GetMedianPriceByRegion(c.Request.Context(), regionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, MedianPriceResponse{
		BaseResponse: BaseResponse{Status: true},
		Data:         results,
	})
}

// GetPriceTrend godoc
// @Summary Get price trend analysis
// @Description Get average and median price trends over time
// @Tags analytics
// @Accept json
// @Produce json
// @Param interval query string false "Time interval (month, year)" default(month)
// @Success 200 {array} backend_internal_models.PriceTrendResult
// @Failure 500 {object} ErrorResponse
// @Router /analytics/price-trend [get]
func (h *AnalyticsHandler) GetPriceTrend(c *gin.Context) {
	interval := c.DefaultQuery("interval", "month")
	results, err := h.svc.GetPriceTrend(c.Request.Context(), interval)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, PriceTrendResponse{
		BaseResponse: BaseResponse{Status: true},
		Data:         results,
	})
}

// GetAffordability godoc
// @Summary Get affordability index
// @Description Get relative affordability by property type
// @Tags analytics
// @Accept json
// @Produce json
// @Success 200 {array} backend_internal_models.AffordabilityResult
// @Failure 500 {object} ErrorResponse
// @Router /analytics/affordability [get]
func (h *AnalyticsHandler) GetAffordability(c *gin.Context) {
	results, err := h.svc.GetAffordability(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, AffordabilityResponse{
		BaseResponse: BaseResponse{Status: true},
		Data:         results,
	})
}

// GetGrowthHotspots godoc
// @Summary Get growth hotspots
// @Description Get top districts with highest price growth rate
// @Tags analytics
// @Accept json
// @Produce json
// @Param limit query int false "Number of results" default(10)
// @Success 200 {array} backend_internal_models.GrowthHotspotResult
// @Failure 500 {object} ErrorResponse
// @Router /analytics/growth-hotspots [get]
func (h *AnalyticsHandler) GetGrowthHotspots(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	results, err := h.svc.GetGrowthHotspots(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, GrowthHotspotResponse{
		BaseResponse: BaseResponse{Status: true},
		Data:         results,
	})
}

// GetNewBuildPremium godoc
// @Summary Get new build premium
// @Description Get average prices of new builds vs established properties by region
// @Tags analytics
// @Accept json
// @Produce json
// @Param by query string false "Region type (county, district, town_city)" default(county)
// @Success 200 {array} backend_internal_models.NewBuildPremiumResult
// @Failure 500 {object} ErrorResponse
// @Router /analytics/new-build-premium [get]
func (h *AnalyticsHandler) GetNewBuildPremium(c *gin.Context) {
	regionType := c.DefaultQuery("by", "county")
	results, err := h.svc.GetNewBuildPremium(c.Request.Context(), regionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, NewBuildPremiumResponse{
		BaseResponse: BaseResponse{Status: true},
		Data:         results,
	})
}

// GetPropertyTypeDistribution godoc
// @Summary Get property type distribution
// @Description Get distribution of properties by type (detached, semi, flat, etc.)
// @Tags analytics
// @Accept json
// @Produce json
// @Success 200 {array} backend_internal_models.PropertyTypeDistributionResult
// @Failure 500 {object} ErrorResponse
// @Router /analytics/property-type-distribution [get]
func (h *AnalyticsHandler) GetPropertyTypeDistribution(c *gin.Context) {
	results, err := h.svc.GetPropertyTypeDistribution(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, PropertyTypeDistributionResponse{
		BaseResponse: BaseResponse{Status: true},
		Data:         results,
	})
}

// GetPriceBracketDistribution godoc
// @Summary Get price bracket distribution
// @Description Get distribution of properties by price ranges
// @Tags analytics
// @Accept json
// @Produce json
// @Success 200 {array} backend_internal_models.PriceBracketResult
// @Failure 500 {object} ErrorResponse
// @Router /analytics/price-bracket-distribution [get]
func (h *AnalyticsHandler) GetPriceBracketDistribution(c *gin.Context) {
	results, err := h.svc.GetPriceBracketDistribution(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, PriceBracketResponse{
		BaseResponse: BaseResponse{Status: true},
		Data:         results,
	})
}

// GetTopActiveAreas godoc
// @Summary Get top active areas
// @Description Get regions with highest transaction volume
// @Tags analytics
// @Accept json
// @Produce json
// @Param by query string false "Region type (county, district, town_city)" default(district)
// @Param limit query int false "Number of results" default(10)
// @Success 200 {array} backend_internal_models.TopActiveAreaResult
// @Failure 500 {object} ErrorResponse
// @Router /analytics/top-active-areas [get]
func (h *AnalyticsHandler) GetTopActiveAreas(c *gin.Context) {
	regionType := c.DefaultQuery("by", "district")
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	results, err := h.svc.GetTopActiveAreas(c.Request.Context(), regionType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, TopActiveAreaResponse{
		BaseResponse: BaseResponse{Status: true},
		Data:         results,
	})
}
