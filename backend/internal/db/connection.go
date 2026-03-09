package db

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) {
	m := gormigrate.New(cfg.DB, gormigrate.DefaultOptions, migrations())

	if err := m.Migrate(); err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}
}

func migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "2026030801_initial_consolidated_schema",
			Migrate: func(tx *gorm.DB) error {
				// 1. Initial Table Creation
				if err := tx.AutoMigrate(
					&models.User{},
					&models.Config{},
					&models.Job{},
					&models.Property{},
				); err != nil {
					return err
				}

				// 2. Optimized Indexes
				indexes := []string{
					"CREATE INDEX IF NOT EXISTS idx_properties_town_city ON properties(town_city)",
					"CREATE INDEX IF NOT EXISTS idx_properties_district ON properties(district)",
					"CREATE INDEX IF NOT EXISTS idx_properties_county ON properties(county)",
					"CREATE INDEX IF NOT EXISTS idx_properties_postcode_outward ON properties(postcode_outward)",
					"CREATE INDEX IF NOT EXISTS idx_properties_postcode_inward ON properties(postcode_inward)",
					"CREATE INDEX IF NOT EXISTS idx_properties_price ON properties(price)",
					"CREATE INDEX IF NOT EXISTS idx_properties_old_new ON properties(old_new)",
					"CREATE INDEX IF NOT EXISTS idx_properties_property_type ON properties(property_type)",
					"CREATE INDEX IF NOT EXISTS idx_properties_date_price ON properties(date_of_transfer DESC, price DESC)",
					"CREATE INDEX IF NOT EXISTS idx_properties_county_price ON properties(county, price)",
					"CREATE INDEX IF NOT EXISTS idx_properties_district_price ON properties(district, price)",
					"CREATE INDEX IF NOT EXISTS idx_properties_town_city_price ON properties(town_city, price)",
				}
				for _, q := range indexes {
					if err := tx.Exec(q).Error; err != nil {
						return err
					}
				}

				// 3. Materialized Views for Analytics

				// Drop old views if they have the old schema (optional, but good for this migration)
				// tx.Exec("DROP MATERIALIZED VIEW IF EXISTS mv_monthly_regional_stats")
				
				// Monthly Stats per Region (Consolidated with UNION for better granularity)
				mvMonthly := `
					CREATE MATERIALIZED VIEW IF NOT EXISTS mv_monthly_regional_stats AS
					SELECT 'county' as region_type, county as region_name,
					       extract(year from date_of_transfer) as year,
					       extract(month from date_of_transfer) as month,
					       (percentile_cont(0.5) WITHIN GROUP (ORDER BY price))::bigint as median_price,
					       AVG(price)::bigint as avg_price,
					       COUNT(*) as transaction_count
					FROM properties WHERE county IS NOT NULL AND county != ''
					GROUP BY county, year, month
					UNION ALL
					SELECT 'district' as region_type, district as region_name,
					       extract(year from date_of_transfer) as year,
					       extract(month from date_of_transfer) as month,
					       (percentile_cont(0.5) WITHIN GROUP (ORDER BY price))::bigint as median_price,
					       AVG(price)::bigint as avg_price,
					       COUNT(*) as transaction_count
					FROM properties WHERE district IS NOT NULL AND district != ''
					GROUP BY district, year, month
					UNION ALL
					SELECT 'town_city' as region_type, town_city as region_name,
					       extract(year from date_of_transfer) as year,
					       extract(month from date_of_transfer) as month,
					       (percentile_cont(0.5) WITHIN GROUP (ORDER BY price))::bigint as median_price,
					       AVG(price)::bigint as avg_price,
					       COUNT(*) as transaction_count
					FROM properties WHERE town_city IS NOT NULL AND town_city != ''
					GROUP BY town_city, year, month
				`
				if err := tx.Exec(mvMonthly).Error; err != nil {
					// If it fails because it already exists with a different schema, we might need to drop it.
					// For a safer automated approach in this dev environment:
					tx.Exec("DROP MATERIALIZED VIEW IF EXISTS mv_monthly_regional_stats CASCADE")
					if err := tx.Exec(mvMonthly).Error; err != nil {
						return err
					}
				}
				tx.Exec("CREATE INDEX IF NOT EXISTS idx_mv_mrs_region ON mv_monthly_regional_stats(region_type, region_name)")
				tx.Exec("CREATE INDEX IF NOT EXISTS idx_mv_mrs_date ON mv_monthly_regional_stats(year, month)")

				// Regional Stats (For regional medians and active areas)
				mvRegional := `
					CREATE MATERIALIZED VIEW IF NOT EXISTS mv_regional_stats AS
					SELECT 'county' as region_type, county as region_name,
					       (percentile_cont(0.5) WITHIN GROUP (ORDER BY price))::bigint as median_price,
					       AVG(price)::bigint as avg_price,
					       COUNT(*) as transaction_count
					FROM properties WHERE county IS NOT NULL AND county != '' GROUP BY county
					UNION ALL
					SELECT 'district' as region_type, district as region_name,
					       (percentile_cont(0.5) WITHIN GROUP (ORDER BY price))::bigint as median_price,
					       AVG(price)::bigint as avg_price,
					       COUNT(*) as transaction_count
					FROM properties WHERE district IS NOT NULL AND district != '' GROUP BY district
					UNION ALL
					SELECT 'town_city' as region_type, town_city as region_name,
					       (percentile_cont(0.5) WITHIN GROUP (ORDER BY price))::bigint as median_price,
					       AVG(price)::bigint as avg_price,
					       COUNT(*) as transaction_count
					FROM properties WHERE town_city IS NOT NULL AND town_city != '' GROUP BY town_city
				`
				if err := tx.Exec(mvRegional).Error; err != nil {
					tx.Exec("DROP MATERIALIZED VIEW IF EXISTS mv_regional_stats CASCADE")
					if err := tx.Exec(mvRegional).Error; err != nil {
						return err
					}
				}
				tx.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_rs_unique ON mv_regional_stats(region_type, region_name)")

				// New Build Premium Stats
				mvNewBuild := `
					CREATE MATERIALIZED VIEW IF NOT EXISTS mv_new_build_stats AS
					SELECT 'county' as region_type, county as region_name, old_new,
					       AVG(price)::bigint as avg_price
					FROM properties WHERE county IS NOT NULL AND county != '' GROUP BY county, old_new
					UNION ALL
					SELECT 'district' as region_type, district as region_name, old_new,
					       AVG(price)::bigint as avg_price
					FROM properties WHERE district IS NOT NULL AND district != '' GROUP BY district, old_new
				`
				if err := tx.Exec(mvNewBuild).Error; err != nil {
					tx.Exec("DROP MATERIALIZED VIEW IF EXISTS mv_new_build_stats CASCADE")
					if err := tx.Exec(mvNewBuild).Error; err != nil {
						return err
					}
				}
				tx.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_nbs_unique ON mv_new_build_stats(region_type, region_name, old_new)")

				// 4. Default Admin User
				var user models.User
				if err := tx.First(&user).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						hashedPassword, err := utils.HashPassword(config.DEFAULT_PASSWORD)
						if err != nil {
							return err
						}
						user := models.User{
							Name:     config.DEFAULT_USER,
							Username: config.DEFAULT_USER,
							Password: *hashedPassword,
							IsAdmin:  true,
						}
						if err := tx.Create(&user).Error; err != nil {
							return err
						}
					} else {
						return err
					}
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				tx.Exec("DROP MATERIALIZED VIEW IF EXISTS mv_monthly_regional_stats")
				tx.Exec("DROP MATERIALIZED VIEW IF EXISTS mv_regional_stats")
				tx.Exec("DROP MATERIALIZED VIEW IF EXISTS mv_new_build_stats")
				return tx.Migrator().DropTable(
					&models.User{},
					&models.Config{},
					&models.Job{},
					&models.Property{},
				)
			},
		},
	}
}
