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
			ID: "2026030501_initial_schema",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&models.User{},
					&models.Config{},
					&models.Job{},
					&models.Property{},
				); err != nil {
					return err
				}

				// insert default admin user if not exists
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
				return tx.Migrator().DropTable(
					&models.User{},
					&models.Config{},
					&models.Job{},
					&models.Property{},
				)
			},
		},
		{
			ID: "2026030602_analytics_indexes",
			Migrate: func(tx *gorm.DB) error {
				// Use raw SQL to create indexes concurrently or just standard indexes if concurrent is not supported
				// Standard indexes for safety across environments
				if err := tx.Exec("CREATE INDEX IF NOT EXISTS idx_properties_property_type ON properties(property_type)").Error; err != nil {
					return err
				}
				if err := tx.Exec("CREATE INDEX IF NOT EXISTS idx_properties_date_of_transfer ON properties(date_of_transfer)").Error; err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Exec("DROP INDEX IF EXISTS idx_properties_property_type").Error; err != nil {
					return err
				}
				if err := tx.Exec("DROP INDEX IF EXISTS idx_properties_date_of_transfer").Error; err != nil {
					return err
				}
				return nil
			},
		},
		{
			ID: "2026030703_optimized_indexes",
			Migrate: func(tx *gorm.DB) error {
				indexes := []string{
					"CREATE INDEX IF NOT EXISTS idx_properties_town_city ON properties(town_city)",
					"CREATE INDEX IF NOT EXISTS idx_properties_district ON properties(district)",
					"CREATE INDEX IF NOT EXISTS idx_properties_county ON properties(county)",
					"CREATE INDEX IF NOT EXISTS idx_properties_postcode ON properties(postcode)",
					"CREATE INDEX IF NOT EXISTS idx_properties_price ON properties(price)",
					"CREATE INDEX IF NOT EXISTS idx_properties_old_new ON properties(old_new)",
					"CREATE INDEX IF NOT EXISTS idx_properties_date_price ON properties(date_of_transfer DESC, price DESC)",
				}
				for _, q := range indexes {
					if err := tx.Exec(q).Error; err != nil {
						return err
					}
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				indexes := []string{
					"DROP INDEX IF EXISTS idx_properties_town_city",
					"DROP INDEX IF EXISTS idx_properties_district",
					"DROP INDEX IF EXISTS idx_properties_county",
					"DROP INDEX IF EXISTS idx_properties_postcode",
					"DROP INDEX IF EXISTS idx_properties_price",
					"DROP INDEX IF EXISTS idx_properties_old_new",
					"DROP INDEX IF EXISTS idx_properties_date_price",
				}
				for _, q := range indexes {
					if err := tx.Exec(q).Error; err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			ID: "2026030704_analytics_mv",
			Migrate: func(tx *gorm.DB) error {
				// Create Materialized View for yearly stats
				// This aggregates data by year and district for fast hotspots analysis
				mvQuery := `
					CREATE MATERIALIZED VIEW IF NOT EXISTS property_yearly_stats AS
					SELECT district,
					       extract(year from date_of_transfer) as year,
					       (percentile_cont(0.5) WITHIN GROUP (ORDER BY price))::bigint as median_price,
					       AVG(price)::bigint as avg_price,
					       COUNT(*) as transaction_count
					FROM properties
					GROUP BY district, year
				`
				if err := tx.Exec(mvQuery).Error; err != nil {
					return err
				}

				// Index for the Materialized View
				if err := tx.Exec("CREATE INDEX IF NOT EXISTS idx_pys_district_year ON property_yearly_stats(district, year)").Error; err != nil {
					return err
				}

				// Add composite indexes to properties table for faster regional aggregates
				extraIndexes := []string{
					"CREATE INDEX IF NOT EXISTS idx_properties_county_price ON properties(county, price)",
					"CREATE INDEX IF NOT EXISTS idx_properties_district_price ON properties(district, price)",
					"CREATE INDEX IF NOT EXISTS idx_properties_town_city_price ON properties(town_city, price)",
				}
				for _, q := range extraIndexes {
					if err := tx.Exec(q).Error; err != nil {
						return err
					}
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Exec("DROP MATERIALIZED VIEW IF EXISTS property_yearly_stats").Error; err != nil {
					return err
				}
				return nil
			},
		},
		{
			ID: "2026030705_comprehensive_analytics_mvs",
			Migrate: func(tx *gorm.DB) error {
				// 1. Monthly Stats per District (Foundation for Trends and Hotspots)
				tx.Exec("DROP MATERIALIZED VIEW IF EXISTS property_yearly_stats")
				
				mvMonthly := `
					CREATE MATERIALIZED VIEW IF NOT EXISTS mv_district_monthly_stats AS
					SELECT district,
					       extract(year from date_of_transfer) as year,
					       extract(month from date_of_transfer) as month,
					       (percentile_cont(0.5) WITHIN GROUP (ORDER BY price))::bigint as median_price,
					       AVG(price)::bigint as avg_price,
					       COUNT(*) as transaction_count
					FROM properties
					GROUP BY district, year, month
				`
				if err := tx.Exec(mvMonthly).Error; err != nil {
					return err
				}
				tx.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_dms_unique ON mv_district_monthly_stats(district, year, month)")

				// 2. Regional Stats (For regional medians and active areas)
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
					return err
				}
				tx.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_rs_unique ON mv_regional_stats(region_type, region_name)")

				// 3. New Build Premium Stats
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
					return err
				}
				tx.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_nbs_unique ON mv_new_build_stats(region_type, region_name, old_new)")

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				tx.Exec("DROP MATERIALIZED VIEW IF EXISTS mv_district_monthly_stats")
				tx.Exec("DROP MATERIALIZED VIEW IF EXISTS mv_regional_stats")
				tx.Exec("DROP MATERIALIZED VIEW IF EXISTS mv_new_build_stats")
				return nil
			},
		},
		{
			ID: "2026030806_address_consolidation",
			Migrate: func(tx *gorm.DB) error {
				// Add address column
				if err := tx.Exec("ALTER TABLE properties ADD COLUMN IF NOT EXISTS address TEXT").Error; err != nil {
					return err
				}

				// Consolidate data
				updateQuery := `
					UPDATE properties 
					SET address = TRIM(CONCAT_WS(', ', 
						NULLIF(TRIM(saon), ''), 
						NULLIF(TRIM(paon), ''), 
						NULLIF(TRIM(street), ''), 
						NULLIF(TRIM(locality), '')
					))
				`
				if err := tx.Exec(updateQuery).Error; err != nil {
					return err
				}

				// Mapping for common Unitary Authorities to Ceremonial Counties (simplified for SQL)
				// This handles the major ones mentioned by the user and others
				countyUpdates := map[string]string{
					"BRIGHTON AND HOVE":                 "EAST SUSSEX",
					"BATH AND NORTH EAST SOMERSET":      "SOMERSET",
					"BOURNEMOUTH, CHRISTCHURCH AND POOLE": "DORSET",
					"BOURNEMOUTH":                       "DORSET",
					"POOLE":                             "DORSET",
					"WEST BERKSHIRE":                    "BERKSHIRE",
					"WINDSOR AND MAIDENHEAD":            "BERKSHIRE",
					"WOKINGHAM":                         "BERKSHIRE",
					"BRACKNELL FOREST":                  "BERKSHIRE",
					"READING":                           "BERKSHIRE",
					"WEST NORTHAMPTONSHIRE":             "NORTHAMPTONSHIRE",
					"NORTH NORTHAMPTONSHIRE":            "NORTHAMPTONSHIRE",
					"CENTRAL BEDFORDSHIRE":              "BEDFORDSHIRE",
					"BEDFORD":                           "BEDFORDSHIRE",
					"CHESHIRE EAST":                     "CHESHIRE",
					"CHESHIRE WEST AND CHESTER":         "CHESHIRE",
					"WESTMORLAND AND FURNESS":           "CUMBRIA",
					"CITY OF BRISTOL":                   "BRISTOL",
					"MILTON KEYNES":                     "BUCKINGHAMSHIRE",
				}

				for unitary, ceremonial := range countyUpdates {
					tx.Exec("UPDATE properties SET county = ? WHERE county = ?", ceremonial, unitary)
				}

				// Drop old columns
				tx.Exec("ALTER TABLE properties DROP COLUMN IF EXISTS saon")
				tx.Exec("ALTER TABLE properties DROP COLUMN IF EXISTS paon")
				tx.Exec("ALTER TABLE properties DROP COLUMN IF EXISTS street")
				tx.Exec("ALTER TABLE properties DROP COLUMN IF EXISTS locality")

				// Refresh views to reflect new county names
				tx.Exec("REFRESH MATERIALIZED VIEW mv_regional_stats")
				tx.Exec("REFRESH MATERIALIZED VIEW mv_new_build_stats")
				tx.Exec("REFRESH MATERIALIZED VIEW mv_district_monthly_stats")

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil // Irreversible column drop
			},
		},
	}
}
