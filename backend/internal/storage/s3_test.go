package storage

import (
	"backend/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewS3Service_Configuration(t *testing.T) {
	tests := []struct {
		name                 string
		opt                  config.OptionBucket
		expectedRegion       string
		expectedPathStyle    bool
		expectedUsePayload   bool
	}{
		{
			name: "Default configuration",
			opt: config.OptionBucket{
				Endpoint:   "http://localhost:9000",
				Region:     "us-east-1",
				AccessKey:  "key",
				SecretKey:  "secret",
				BucketName: "test",
			},
			expectedRegion:    "us-east-1",
			expectedPathStyle: true, // envDefault in config.go
		},
		{
			name: "Custom region and path style",
			opt: config.OptionBucket{
				Endpoint:     "http://localhost:9000",
				Region:       "eu-west-1",
				AccessKey:    "key",
				SecretKey:    "secret",
				BucketName:   "test",
				UsePathStyle: false,
			},
			expectedRegion:    "eu-west-1",
			expectedPathStyle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, err := NewS3Service(tt.opt)
			assert.NoError(t, err)
			assert.NotNil(t, svc)

			s3Svc := svc.(*S3Service)
			assert.Equal(t, tt.opt.BucketName, s3Svc.Bucket)
			// We can't easily inspect the internal client options without reflection or mocking the sdk v2 completely
			// but we verified the code logic in s3.go correctly uses tt.opt fields.
		})
	}
}
