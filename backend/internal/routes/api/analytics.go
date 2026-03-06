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
}

// GetMedianPriceByRegion godoc
// @Summary Get median price by region
// @Description Get median price grouped by county, district, or town_city
// @Tags analytics
// @Accept json
// @Produce json
// @Param by query string false "Region type (county, district, town_city)" default(county)
// @Success 200 {object} MedianPriceResponse
// @Failure 500 {object} ErrorResponse
// @Router /analytics/median-price [get]
// @Security JwtAuth
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
// @Success 200 {object} PriceTrendResponse
// @Failure 500 {object} ErrorResponse
// @Router /analytics/price-trend [get]
// @Security JwtAuth
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
// @Success 200 {object} AffordabilityResponse
// @Failure 500 {object} ErrorResponse
// @Router /analytics/affordability [get]
// @Security JwtAuth
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
// @Success 200 {object} GrowthHotspotResponse
// @Failure 500 {object} ErrorResponse
// @Router /analytics/growth-hotspots [get]
// @Security JwtAuth
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
