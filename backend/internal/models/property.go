package models

import (
	"time"

	"gorm.io/gorm"
)

type Property struct {
	// A unique identifier for each property sale
	ID        string         `gorm:"primarykey;type:uuid" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Price          int64     `json:"price"` // In GBP
	DateOfTransfer time.Time `json:"date_of_transfer"`
	Postcode       string    `gorm:"index" json:"postcode"`

	// D = Detached, S = Semi-Detached, T = Terraced, F = Flats/Maisonettes, O = Other
	PropertyType string `json:"property_type"`

	// Y = a newly built property, N = an established residential building
	OldNew string `json:"old_new"`

	// F = Freehold, L = Leasehold
	Duration string `json:"duration"`

	// Primary Addressable Object Name. Typically the house number or name.
	PAON string `json:"paon"`

	// Secondary Addressable Object Name. Where a property has been divided into separate units (for example, flats).
	SAON string `json:"saon"`

	Street   string `json:"street"`
	Locality string `json:"locality"`
	TownCity string `gorm:"index" json:"town_city"`
	District string `gorm:"index" json:"district"`
	County   string `gorm:"index" json:"county"`

	// A = Standard Price Paid entry, B = Additional Price Paid entry
	PPDCategoryType string `json:"ppd_category_type"`

	// A = Addition, C = Change, D = Delete
	RecordStatus string `json:"record_status"`
}

func (Property) TableName() string {
	return "properties"
}
