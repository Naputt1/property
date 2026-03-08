package queue

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCleanUUID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"{F87E6C09-CA25-4981-8053-90378C6A8D00}", "f87e6c09-ca25-4981-8053-90378c6a8d00"},
		{"F87E6C09-CA25-4981-8053-90378C6A8D00", "f87e6c09-ca25-4981-8053-90378c6a8d00"},
		{"{abc-123}", "abc-123"},
		{"ABC-123", "abc-123"},
	}

	for _, tt := range tests {
		result := cleanUUID(tt.input)
		if result != tt.expected {
			t.Errorf("cleanUUID(%s) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}

func TestMapRecordToProperty(t *testing.T) {
	t.Run("16-column standard PPD", func(t *testing.T) {
		record := []string{
			"{UUID-1}", "250000", "2024-03-07 00:00", "SW1A 1AA", "D", "N", "F", "10", "", "Downing Street", "Westminster", "London", "Greater London", "London", "A", "A",
		}
		prop, err := mapRecordToProperty(record)
		assert.NoError(t, err)
		assert.Equal(t, int64(250000), prop.Price)
		assert.Equal(t, "SW1A", prop.PostcodeOutward)
		assert.Equal(t, "1AA", prop.PostcodeInward)
		assert.Equal(t, "D", prop.PropertyType)
		assert.Equal(t, "London", prop.TownCity)
		assert.Equal(t, "Greater London", prop.District) // Matches input
		assert.Equal(t, "GREATER LONDON", prop.County) // Derived from SW postcode (uppercase from map)
	})

	t.Run("17-column split postcode", func(t *testing.T) {
		record := []string{
			"{UUID-1}", "250000", "2024-03-07 00:00", "SW1A", "1AA", "D", "N", "F", "10", "", "Downing Street", "Westminster", "London", "Greater London", "London", "A", "A",
		}
		prop, err := mapRecordToProperty(record)
		assert.NoError(t, err)
		assert.Equal(t, int64(250000), prop.Price)
		assert.Equal(t, "SW1A", prop.PostcodeOutward)
		assert.Equal(t, "1AA", prop.PostcodeInward)
		assert.Equal(t, "D", prop.PropertyType)
		assert.Equal(t, "London", prop.TownCity)
		assert.Equal(t, "Greater London", prop.District) // Matches input
		assert.Equal(t, "GREATER LONDON", prop.County) // Derived from SW postcode (uppercase from map)
	})

	t.Run("17-column shift verification", func(t *testing.T) {
		// If it's 17 columns, but we used the old mapping (16-col), indices would be wrong.
		// Let's verify our new mapping correctly gets the fields.
		record := []string{
			"ID", "100", "2024-01-01 12:00", "SW1A", "1AA", "T", "Y", "L", "1", "A", "STREET", "LOCALITY", "TOWN", "DISTRICT", "COUNTY", "A", "A",
		}
		prop, err := mapRecordToProperty(record)
		assert.NoError(t, err)
		assert.Equal(t, "SW1A", prop.PostcodeOutward)
		assert.Equal(t, "1AA", prop.PostcodeInward)
		assert.Equal(t, "T", prop.PropertyType)
		assert.Equal(t, "Y", prop.OldNew)
		assert.Equal(t, "L", prop.Duration)
		assert.Equal(t, "TOWN", prop.TownCity)
		assert.Equal(t, "DISTRICT", prop.District)
		assert.Equal(t, "GREATER LONDON", prop.County) // Derived from SW1A
	})
}
