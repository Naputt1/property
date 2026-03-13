package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             int64          `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string" format:"date-time"`
	Name           string         `json:"name"`
	Username       string     `gorm:"uniqueIndex" json:"username"`
	Password       string     `json:"-"`
	IsAdmin        bool       `json:"is_admin"`
	RefreshVersion int64      `gorm:"default:1" json:"refresh_version"`
}

func (User) TableName() string {
	return "users"
}

type UserJwt struct {
	Id int64 `json:"id"`
}

type UserJwtRefresh struct {
	Id int64 `json:"id"`
}
