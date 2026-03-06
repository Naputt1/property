package queue

import (
	"testing"
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
