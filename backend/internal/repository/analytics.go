package repository

import (
	"backend/internal/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) GetMedianPriceByRegion(ctx context.Context, regionType string) ([]models.MedianPriceResult, error) {
	var results []models.MedianPriceResult

	validRegions := map[string]bool{"county": true, "district": true, "town_city": true}
	if !validRegions[regionType] {
		return nil, fmt.Errorf("invalid region type: %s", regionType)
	}

	query := `
		SELECT region_name as region, median_price
		FROM mv_regional_stats
		WHERE region_type = ?
		ORDER BY median_price DESC
	`

	err := r.db.WithContext(ctx).Raw(query, regionType).Scan(&results).Error
	return results, err
}

func (r *analyticsRepository) GetPriceTrend(ctx context.Context, interval string) ([]models.PriceTrendResult, error) {
	var results []models.PriceTrendResult

	if interval == "year" {
		query := `
			SELECT year::text as period,
			       AVG(avg_price)::bigint as avg_price,
			       AVG(median_price)::bigint as median_price,
			       SUM(transaction_count) as transaction_count
			FROM mv_district_monthly_stats
			GROUP BY year
			ORDER BY year ASC
		`
		err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
		return results, err
	}

	// Monthly trend
	query := `
		SELECT format('%s-%s-01', year, lpad(month::text, 2, '0')) as period,
		       AVG(avg_price)::bigint as avg_price,
		       AVG(median_price)::bigint as median_price,
		       SUM(transaction_count) as transaction_count
		FROM mv_district_monthly_stats
		GROUP BY year, month
		ORDER BY year ASC, month ASC
	`

	err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
	return results, err
}

func (r *analyticsRepository) GetAffordability(ctx context.Context) ([]models.AffordabilityResult, error) {
	var results []models.AffordabilityResult

	// Relative affordability is calculated as (AvgPrice of type / Overall AvgPrice)
	// Lower value means more affordable compared to average.
	query := `
		WITH overall_avg AS (
			SELECT AVG(price) as total_avg FROM properties
		)
		SELECT property_type,
		       AVG(price)::bigint as avg_price,
		       (AVG(price) / (SELECT total_avg FROM overall_avg)) as relative_affordability
		FROM properties
		GROUP BY property_type
		ORDER BY avg_price ASC
	`

	err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
	return results, err
}

func (r *analyticsRepository) GetGrowthHotspots(ctx context.Context, limit int) ([]models.GrowthHotspotResult, error) {
	var results []models.GrowthHotspotResult

	query := `
		WITH yearly_stats AS (
			SELECT district, year, AVG(median_price)::bigint as median_price
			FROM mv_district_monthly_stats
			GROUP BY district, year
		),
		latest_years AS (
			SELECT DISTINCT year FROM yearly_stats ORDER BY year DESC LIMIT 2
		),
		growth AS (
			SELECT curr.district as region,
			       curr.median_price as current_median,
			       prev.median_price as prev_median,
			       ((curr.median_price - prev.median_price)::float / prev.median_price::float) * 100 as growth_rate
			FROM yearly_stats curr
			JOIN yearly_stats prev ON curr.district = prev.district 
			     AND curr.year = (SELECT year FROM latest_years LIMIT 1)
			     AND prev.year = (SELECT year FROM latest_years OFFSET 1 LIMIT 1)
		)
		SELECT region, growth_rate, prev_median, current_median FROM growth WHERE prev_median > 0 ORDER BY growth_rate DESC LIMIT ?
	`

	err := r.db.WithContext(ctx).Raw(query, limit).Scan(&results).Error
	return results, err
}

func (r *analyticsRepository) RefreshMaterializedView(ctx context.Context) error {
	views := []string{"mv_district_monthly_stats", "mv_regional_stats", "mv_new_build_stats"}
	for _, v := range views {
		if err := r.db.WithContext(ctx).Exec(fmt.Sprintf("REFRESH MATERIALIZED VIEW CONCURRENTLY %s", v)).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *analyticsRepository) GetNewBuildPremium(ctx context.Context, regionType string) ([]models.NewBuildPremiumResult, error) {
	var results []models.NewBuildPremiumResult

	validRegions := map[string]bool{"county": true, "district": true}
	if !validRegions[regionType] {
		return nil, fmt.Errorf("invalid region type: %s", regionType)
	}

	query := `
		SELECT n.region_name as region,
		       n.avg_price as new_avg,
		       o.avg_price as old_avg,
		       ((n.avg_price - o.avg_price)::float / o.avg_price::float) * 100 as premium_percent
		FROM mv_new_build_stats n
		JOIN mv_new_build_stats o ON n.region_name = o.region_name 
		     AND n.region_type = o.region_type
		WHERE n.region_type = ? AND n.old_new = 'Y' AND o.old_new = 'N'
		ORDER BY premium_percent DESC
	`

	err := r.db.WithContext(ctx).Raw(query, regionType).Scan(&results).Error
	return results, err
}

func (r *analyticsRepository) GetPropertyTypeDistribution(ctx context.Context) ([]models.PropertyTypeDistributionResult, error) {
	var results []models.PropertyTypeDistributionResult

	query := `
		SELECT property_type,
		       COUNT(*) as count,
		       (COUNT(*)::float / (SELECT COUNT(*) FROM properties)::float) * 100 as percentage
		FROM properties
		GROUP BY property_type
		ORDER BY count DESC
	`

	err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
	return results, err
}

func (r *analyticsRepository) GetPriceBracketDistribution(ctx context.Context) ([]models.PriceBracketResult, error) {
	var results []models.PriceBracketResult

	query := `
		SELECT 
			CASE 
				WHEN price < 150000 THEN '< £150k'
				WHEN price >= 150000 AND price < 250000 THEN '£150k - £250k'
				WHEN price >= 250000 AND price < 500000 THEN '£250k - £500k'
				WHEN price >= 500000 AND price < 1000000 THEN '£500k - £1M'
				ELSE '> £1M'
			END as bracket,
			COUNT(*) as count,
			(COUNT(*)::float / (SELECT COUNT(*) FROM properties)::float) * 100 as percentage
		FROM properties
		GROUP BY bracket
		ORDER BY MIN(price) ASC
	`

	err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
	return results, err
}

func (r *analyticsRepository) GetTopActiveAreas(ctx context.Context, regionType string, limit int) ([]models.TopActiveAreaResult, error) {
	var results []models.TopActiveAreaResult

	query := `
		SELECT region_name as region, transaction_count, (transaction_count * avg_price)::bigint as total_value
		FROM mv_regional_stats
		WHERE region_type = ?
		ORDER BY transaction_count DESC
		LIMIT ?
	`

	err := r.db.WithContext(ctx).Raw(query, regionType, limit).Scan(&results).Error
	return results, err
}
