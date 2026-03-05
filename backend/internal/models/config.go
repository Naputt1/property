package models

import (
	"time"

	"gorm.io/gorm"
)

type Config struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Key       string         `gorm:"uniqueIndex"`
	Value     string
}

func (Config) TableName() string {
	return "configs"
}
