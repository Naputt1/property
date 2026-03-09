package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type analyticsService struct {
	cfg   *config.Config
	repo  repository.AnalyticsRepository
	redis *redis.Client
}

func NewAnalyticsService(cfg *config.Config, repo repository.AnalyticsRepository) AnalyticsService {
	return &analyticsService{
		cfg:   cfg,
		repo:  repo,
		redis: cfg.Redis,
	}
}

const (
	AnalyticsCachePrefix = "analytics:"
	AnalyticsCacheTTL    = 24 * time.Hour
)

func (s *analyticsService) getCached(ctx context.Context, key string, target interface{}) (bool, error) {
	if s.redis == nil {
		return false, nil
	}

	val, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		slog.Error("redis get error", "key", key, "error", err)
		return false, err
	}

	err = json.Unmarshal([]byte(val), target)
	if err != nil {
		slog.Error("json unmarshal error", "key", key, "error", err)
		return false, err
	}

	return true, nil
}

func (s *analyticsService) setCache(ctx context.Context, key string, value interface{}) {
	if s.redis == nil {
		return
	}

	data, err := json.Marshal(value)
	if err != nil {
		slog.Error("json marshal error", "key", key, "error", err)
		return
	}

	err = s.redis.Set(ctx, key, data, AnalyticsCacheTTL).Err()
	if err != nil {
		slog.Error("redis set error", "key", key, "error", err)
	}
}

func (s *analyticsService) GetMedianPriceByRegion(ctx context.Context, regionType string, year int) ([]models.MedianPriceResult, error) {
	cacheKey := fmt.Sprintf("%smedian_price:%s:%d", AnalyticsCachePrefix, regionType, year)
	var results []models.MedianPriceResult

	found, _ := s.getCached(ctx, cacheKey, &results)
	if found {
		return results, nil
	}

	results, err := s.repo.GetMedianPriceByRegion(ctx, regionType, year)
	if err != nil {
		return nil, err
	}

	s.setCache(ctx, cacheKey, results)
	return results, nil
}

func (s *analyticsService) GetPriceTrend(ctx context.Context, interval string) ([]models.PriceTrendResult, error) {
	cacheKey := fmt.Sprintf("%sprice_trend:%s", AnalyticsCachePrefix, interval)
	var results []models.PriceTrendResult

	found, _ := s.getCached(ctx, cacheKey, &results)
	if found {
		return results, nil
	}

	results, err := s.repo.GetPriceTrend(ctx, interval)
	if err != nil {
		return nil, err
	}

	s.setCache(ctx, cacheKey, results)
	return results, nil
}

func (s *analyticsService) GetAffordability(ctx context.Context) ([]models.AffordabilityResult, error) {
	cacheKey := fmt.Sprintf("%saffordability", AnalyticsCachePrefix)
	var results []models.AffordabilityResult

	found, _ := s.getCached(ctx, cacheKey, &results)
	if found {
		return results, nil
	}

	results, err := s.repo.GetAffordability(ctx)
	if err != nil {
		return nil, err
	}

	s.setCache(ctx, cacheKey, results)
	return results, nil
}

func (s *analyticsService) GetGrowthHotspots(ctx context.Context, regionType string, limit, year int) ([]models.GrowthHotspotResult, error) {
	cacheKey := fmt.Sprintf("%sgrowth_hotspots:%s:%d:%d", AnalyticsCachePrefix, regionType, limit, year)
	var results []models.GrowthHotspotResult

	found, _ := s.getCached(ctx, cacheKey, &results)
	if found {
		return results, nil
	}

	results, err := s.repo.GetGrowthHotspots(ctx, regionType, limit, year)
	if err != nil {
		return nil, err
	}

	s.setCache(ctx, cacheKey, results)
	return results, nil
}

func (s *analyticsService) GetNewBuildPremium(ctx context.Context, regionType string) ([]models.NewBuildPremiumResult, error) {
	cacheKey := fmt.Sprintf("%snew_build_premium:%s", AnalyticsCachePrefix, regionType)
	var results []models.NewBuildPremiumResult

	found, _ := s.getCached(ctx, cacheKey, &results)
	if found {
		return results, nil
	}

	results, err := s.repo.GetNewBuildPremium(ctx, regionType)
	if err != nil {
		return nil, err
	}

	s.setCache(ctx, cacheKey, results)
	return results, nil
}

func (s *analyticsService) GetPropertyTypeDistribution(ctx context.Context) ([]models.PropertyTypeDistributionResult, error) {
	cacheKey := fmt.Sprintf("%sproperty_type_distribution", AnalyticsCachePrefix)
	var results []models.PropertyTypeDistributionResult

	found, _ := s.getCached(ctx, cacheKey, &results)
	if found {
		return results, nil
	}

	results, err := s.repo.GetPropertyTypeDistribution(ctx)
	if err != nil {
		return nil, err
	}

	s.setCache(ctx, cacheKey, results)
	return results, nil
}

func (s *analyticsService) GetPriceBracketDistribution(ctx context.Context) ([]models.PriceBracketResult, error) {
	cacheKey := fmt.Sprintf("%sprice_bracket_distribution", AnalyticsCachePrefix)
	var results []models.PriceBracketResult

	found, _ := s.getCached(ctx, cacheKey, &results)
	if found {
		return results, nil
	}

	results, err := s.repo.GetPriceBracketDistribution(ctx)
	if err != nil {
		return nil, err
	}

	s.setCache(ctx, cacheKey, results)
	return results, nil
}

func (s *analyticsService) GetTopActiveAreas(ctx context.Context, regionType string, limit, year int) ([]models.TopActiveAreaResult, error) {
	cacheKey := fmt.Sprintf("%stop_active_areas:%s:%d:%d", AnalyticsCachePrefix, regionType, limit, year)
	var results []models.TopActiveAreaResult

	found, _ := s.getCached(ctx, cacheKey, &results)
	if found {
		return results, nil
	}

	results, err := s.repo.GetTopActiveAreas(ctx, regionType, limit, year)
	if err != nil {
		return nil, err
	}

	s.setCache(ctx, cacheKey, results)
	return results, nil
}

func (s *analyticsService) GetTimeRange(ctx context.Context) (*models.TimeRangeResult, error) {
	cacheKey := fmt.Sprintf("%stime_range", AnalyticsCachePrefix)
	var result models.TimeRangeResult

	found, _ := s.getCached(ctx, cacheKey, &result)
	if found {
		return &result, nil
	}

	res, err := s.repo.GetTimeRange(ctx)
	if err != nil {
		return nil, err
	}

	s.setCache(ctx, cacheKey, res)
	return res, nil
}

func (s *analyticsService) RefreshMaterializedView(ctx context.Context) error {
	return s.repo.RefreshMaterializedView(ctx)
}

func (s *analyticsService) PrecomputeCache(ctx context.Context) error {
	slog.Info("Starting analytics cache pre-computation")
	start := time.Now()

	regions := []string{"county", "district", "town_city"}

	// Precompute median price by region
	for _, region := range regions {
		_, err := s.GetMedianPriceByRegion(ctx, region, 0)
		if err != nil {
			slog.Error("failed to precompute median price", "region", region, "error", err)
		}

		_, err = s.GetNewBuildPremium(ctx, region)
		if err != nil {
			slog.Error("failed to precompute new build premium", "region", region, "error", err)
		}

		_, err = s.GetTopActiveAreas(ctx, region, 10, 0)
		if err != nil {
			slog.Error("failed to precompute top active areas", "region", region, "error", err)
		}
	}

	// Precompute price trends
	intervals := []string{"month", "year"}
	for _, interval := range intervals {
		_, err := s.GetPriceTrend(ctx, interval)
		if err != nil {
			slog.Error("failed to precompute price trend", "interval", interval, "error", err)
		}
	}

	// Precompute affordability
	_, err := s.GetAffordability(ctx)
	if err != nil {
		slog.Error("failed to precompute affordability", "error", err)
	}

	// Precompute growth hotspots
	for _, region := range []string{"county", "district"} {
		_, err = s.GetGrowthHotspots(ctx, region, 100, 0)
		if err != nil {
			slog.Error("failed to precompute growth hotspots", "region", region, "error", err)
		}
	}

	// Precompute distributions
	_, err = s.GetPropertyTypeDistribution(ctx)
	if err != nil {
		slog.Error("failed to precompute property type distribution", "error", err)
	}

	_, err = s.GetPriceBracketDistribution(ctx)
	if err != nil {
		slog.Error("failed to precompute price bracket distribution", "error", err)
	}

	slog.Info("Analytics cache pre-computation completed", "duration", time.Since(start))
	return nil
}

func (s *analyticsService) ClearCache(ctx context.Context) error {
	if s.redis == nil {
		return nil
	}

	iter := s.redis.Scan(ctx, 0, AnalyticsCachePrefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		if err := s.redis.Del(ctx, iter.Val()).Err(); err != nil {
			slog.Error("failed to delete cache key", "key", iter.Val(), "error", err)
		}
	}
	if err := iter.Err(); err != nil {
		slog.Error("redis scan error during clear cache", "error", err)
		return err
	}

	slog.Info("Analytics cache cleared")
	return nil
}
