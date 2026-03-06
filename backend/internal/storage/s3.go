package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"

	"backend/internal/config"
	"backend/internal/repository"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type S3Service struct {
	Client   *s3.Client
	Uploader *manager.Uploader
	Bucket   string
}

func NewS3Service(opt config.OptionBucket) (repository.BucketService, error) {
	u, err := url.Parse(opt.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid RustFS URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("RustFS URL must start with http:// or https://, got: %s", opt.Endpoint)
	}

	// AWS SDK v2 BaseEndpoint cannot have a path component.
	endpoint := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	creds := credentials.NewStaticCredentialsProvider(opt.AccessKey, opt.SecretKey, "")

	// Some S3-compatible backends prefer an empty region
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithCredentialsProvider(creds),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
	// 	o.BaseEndpoint = aws.String(fmt.Sprintf("%s://%s", u.Scheme, u.Host))
	// 	o.UsePathStyle = true
	// 	// S3 compatible backends often need unsigned payload for streaming/multipart
	// 	o.APIOptions = append(o.APIOptions, v4.SwapComputePayloadSHA256ForUnsignedPayloadMiddleware)
	// })

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	uploader := manager.NewUploader(s3Client)

	return &S3Service{
		Client:   s3Client,
		Uploader: uploader,
		Bucket:   opt.BucketName,
	}, nil
}

func containsProtocol(endpoint string) bool {
	return len(endpoint) > 7 && (endpoint[:7] == "http://" || (len(endpoint) > 8 && endpoint[:8] == "https://"))
}

func (s *S3Service) EnsureBucket(ctx context.Context) error {
	_, err := s.Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.Bucket),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.(type) {
			case *types.NotFound:
				_, err = s.Client.CreateBucket(ctx, &s3.CreateBucketInput{
					Bucket: aws.String(s.Bucket),
				})
				return err
			}
		}
		// Fallback for S3-compatible that might not return types.NotFound correctly
		_, err = s.Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(s.Bucket),
		})
		return err
	}
	return nil
}

func (s *S3Service) Upload(ctx context.Context, key string, body io.Reader, size int64, contentType string) error {
	_, err := s.Uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	return err
}

func (s *S3Service) GetObject(ctx context.Context, key string) (io.ReadCloser, int64, error) {
	output, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, 0, err
	}

	size := int64(0)
	if output.ContentLength != nil {
		size = *output.ContentLength
	}

	return output.Body, size, nil
}

func (s *S3Service) Delete(ctx context.Context, key string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	return err
}
