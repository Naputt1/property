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

	// Validate regionType to prevent SQL injection (though GORM handles simple params, column names need care)
	validRegions := map[string]bool{"county": true, "district": true, "town_city": true}
	if !validRegions[regionType] {
		return nil, fmt.Errorf("invalid region type: %s", regionType)
	}

	query := fmt.Sprintf(`
		SELECT %s as region, 
		       (percentile_cont(0.5) WITHIN GROUP (ORDER BY price))::bigint as median_price
		FROM properties
		WHERE %s IS NOT NULL AND %s != ''
		GROUP BY region
		ORDER BY median_price DESC
	`, regionType, regionType, regionType)

	err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
	return results, err
}

func (r *analyticsRepository) GetPriceTrend(ctx context.Context, interval string) ([]models.PriceTrendResult, error) {
	var results []models.PriceTrendResult

	validIntervals := map[string]bool{"month": true, "year": true}
	if !validIntervals[interval] {
		return nil, fmt.Errorf("invalid interval: %s", interval)
	}

	query := fmt.Sprintf(`
		SELECT to_char(date_trunc('%s', date_of_transfer), 'YYYY-MM-DD') as period,
		       AVG(price)::bigint as avg_price,
		       (percentile_cont(0.5) WITHIN GROUP (ORDER BY price))::bigint as median_price,
		       COUNT(*) as transaction_count
		FROM properties
		GROUP BY period
		ORDER BY period ASC
	`, interval)

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

	// Compare median prices of the latest year vs the previous year per district
	// This query assumes we have data for at least two distinct years.
	query := `
		WITH yearly_median AS (
			SELECT district,
			       extract(year from date_of_transfer) as year,
			       (percentile_cont(0.5) WITHIN GROUP (ORDER BY price))::bigint as median_price
			FROM properties
			GROUP BY district, year
		),
		latest_years AS (
			SELECT DISTINCT year FROM yearly_median ORDER BY year DESC LIMIT 2
		),
		growth AS (
			SELECT curr.district as region,
			       curr.median_price as current_median,
			       prev.median_price as prev_median,
			       ((curr.median_price - prev.median_price)::float / prev.median_price::float) * 100 as growth_rate
			FROM yearly_median curr
			JOIN yearly_median prev ON curr.district = prev.district 
			     AND curr.year = (SELECT year FROM latest_years LIMIT 1)
			     AND prev.year = (SELECT year FROM latest_years OFFSET 1 LIMIT 1)
		)
		SELECT region, growth_rate, prev_median, current_median FROM growth WHERE prev_median > 0 ORDER BY growth_rate DESC LIMIT ?
	`

	err := r.db.WithContext(ctx).Raw(query, limit).Scan(&results).Error
	return results, err
}

func (r *analyticsRepository) GetNewBuildPremium(ctx context.Context, regionType string) ([]models.NewBuildPremiumResult, error) {
	var results []models.NewBuildPremiumResult

	validRegions := map[string]bool{"county": true, "district": true, "town_city": true}
	if !validRegions[regionType] {
		return nil, fmt.Errorf("invalid region type: %s", regionType)
	}

	query := fmt.Sprintf(`
		SELECT %s as region,
		       AVG(CASE WHEN old_new = 'Y' THEN price END)::bigint as new_avg,
		       AVG(CASE WHEN old_new = 'N' THEN price END)::bigint as old_avg,
		       CASE 
		           WHEN AVG(CASE WHEN old_new = 'N' THEN price END) > 0 
		           THEN ((AVG(CASE WHEN old_new = 'Y' THEN price END) - AVG(CASE WHEN old_new = 'N' THEN price END)) / AVG(CASE WHEN old_new = 'N' THEN price END)) * 100
		           ELSE 0 
		       END as premium_percent
		FROM properties
		WHERE %s IS NOT NULL AND %s != ''
		GROUP BY region
		HAVING AVG(CASE WHEN old_new = 'Y' THEN price END) IS NOT NULL
		   AND AVG(CASE WHEN old_new = 'N' THEN price END) IS NOT NULL
		ORDER BY premium_percent DESC
	`, regionType, regionType, regionType)

	err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
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

	validRegions := map[string]bool{"county": true, "district": true, "town_city": true}
	if !validRegions[regionType] {
		return nil, fmt.Errorf("invalid region type: %s", regionType)
	}

	query := fmt.Sprintf(`
		SELECT %s as region,
		       COUNT(*) as transaction_count,
		       SUM(price)::bigint as total_value
		FROM properties
		WHERE %s IS NOT NULL AND %s != ''
		GROUP BY region
		ORDER BY transaction_count DESC
		LIMIT ?
	`, regionType, regionType, regionType)

	err := r.db.WithContext(ctx).Raw(query, limit).Scan(&results).Error
	return results, err
}
