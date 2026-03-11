package config

import (
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"time"

	"backend/internal/repository"

	"github.com/caarlos0/env/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type OptionDatabase struct {
	Host     string `env:"DATABASE_HOST" envDefault:"localhost"`
	Port     int    `env:"DATABASE_PORT" envDefault:"5432"`
	Username string `env:"DATABASE_USERNAME" envDefault:"postgres"`
	Password string `env:"DATABASE_PASSWORD" envDefault:"password"`
	Name     string `env:"DATABASE_NAME" envDefault:"property"`
}

type OptionRedis struct {
	URL string `env:"REDIS_URL" envDefault:"localhost:6379"`
	DB  int    `env:"REDIS_DB" envDefault:"0"`
}

type OptionBucket struct {
	Endpoint   string `env:"BUCKET_ENDPOINT" envDefault:"http://localhost:9000"`
	AccessKey  string `env:"BUCKET_ACCESS_KEY" envDefault:"rustfsadmin"`
	SecretKey  string `env:"BUCKET_SECRET_KEY" envDefault:"rustfsadmin"`
	UseSSL     bool   `env:"BUCKET_USE_SSL" envDefault:"false"`
	BucketName string `env:"BUCKET_NAME" envDefault:"property-data"`
}

type Option struct {
	DB           OptionDatabase
	Redis        OptionRedis
	Bucket       OptionBucket
	SecretKey    string `env:"SECRET_KEY" envDefault:"secret"`
	TokenVersion int    `env:"TOKEN_VERSION" envDefault:"1"`
	Port         string `env:"PORT" envDefault:"8080"`
	IsProd       bool   `env:"IS_PROD" envDefault:"false"`
}

type Config struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Bucket repository.BucketService
	Opt    Option
}

func (c *Config) Close() {
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	if c.Redis != nil {
		c.Redis.Close()
	}
}

func (c *Config) RegisterBucket(bucket repository.BucketService) {
	c.Bucket = bucket
}

func CreateConfig() (*Config, error) {
	opt := &Option{}
	err := env.Parse(opt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		_ = os.Mkdir(logDir, 0755)
	}
	logFile, err := os.OpenFile(fmt.Sprintf("%s/backend.log", logDir), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Setup slog as default logger with JSON output to both stdout and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	loggerJson := slog.New(slog.NewJSONHandler(multiWriter, nil))
	slog.SetDefault(loggerJson)

	// Database Connection
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", opt.DB.Username, url.QueryEscape(opt.DB.Password), opt.DB.Host, opt.DB.Port, opt.DB.Name)
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql db: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	// Redis client initialization
	rdb := redis.NewClient(&redis.Options{
		Addr: opt.Redis.URL,
		DB:   opt.Redis.DB,
	})

	return &Config{
		DB:    db,
		Redis: rdb,
		Opt:   *opt,
	}, nil
}

const (
	CONTEXT_USER            = "user"
	CONTEXT_VERSION         = "version"
	CONTEXT_REFRESH_VERSION = "version"
)

type WsMessageType string

const (
	WsMessageTypeJobUpdate WsMessageType = "job_update"
)

const DEFAULT_USER = "admin"
const DEFAULT_PASSWORD = "admin"
