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
	}
}
