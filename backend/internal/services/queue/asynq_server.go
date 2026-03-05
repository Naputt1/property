package queue

import (
	"backend/internal/config"
	"backend/internal/services"
	"context"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
)

type AsynqServer struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewAsynqServer(cfg *config.Config, svcs *services.Services) (*AsynqServer, error) {
	redisConnOpt := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", cfg.Opt.Redis.Host, cfg.Opt.Redis.Port),
		Password: cfg.Opt.Redis.Password,
		DB:       cfg.Opt.Redis.DB,
	}

	srv := asynq.NewServer(
		redisConnOpt,
		asynq.Config{
			Concurrency: 5,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.Use(loggingMiddleware)

	// Register handlers
	migrationHandler := NewMigrationHandler(svcs)
	mux.HandleFunc("properties:migrate:csv", migrationHandler.HandleCSVMigrateTask)

	server := &AsynqServer{
		server: srv,
		mux:    mux,
	}

	return server, nil
}

func (s *AsynqServer) Start() error {
	return s.server.Start(s.mux)
}

func (s *AsynqServer) Stop() {
	s.server.Stop()
}

func loggingMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		slog.Info("Start processing task", "type", t.Type())
		err := h.ProcessTask(ctx, t)
		if err != nil {
			slog.Error("Failed to process task", "type", t.Type(), "error", err)
		} else {
			slog.Info("Successfully processed task", "type", t.Type())
		}
		return err
	})
}
