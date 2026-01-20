package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/faisalyudiansah/auth-service-template/internal/gateway/provider"
	"github.com/faisalyudiansah/auth-service-template/pkg/config"
	"github.com/faisalyudiansah/auth-service-template/pkg/logger"

	"github.com/hibiken/asynq"
)

type QueueServer struct {
	cfg    *config.Config
	mux    *asynq.ServeMux
	server *asynq.Server
}

func NewQueueServer(cfg *config.Config) *QueueServer {
	redisOpt := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%v:%v", cfg.Redis.Host, cfg.Redis.Port),
	}

	client := asynq.NewClient(redisOpt)
	mux := asynq.NewServeMux()
	provider.ProvideQueueDependency(client, mux)

	return &QueueServer{
		cfg: cfg,
		mux: mux,
		server: asynq.NewServer(
			redisOpt,
			asynq.Config{
				Concurrency: 10,
				Queues: map[string]int{
					"critical": 6,
					"default":  3,
					"low":      1,
				},
				ShutdownTimeout: time.Duration(cfg.HttpServer.GracePeriod) * time.Second,
			},
		),
	}
}

func (s *QueueServer) Start() {
	logger.Log.Info("Running queue server...")
	if err := s.server.Run(s.mux); err != nil && !errors.Is(err, asynq.ErrServerClosed) {
		logger.Log.Fatal("Error while queue server listening:", err)
	}
	logger.Log.Info("Queue server is not receiving new requests...")
}

func (s *QueueServer) Shutdown() {
	logger.Log.Info("Attempting to shut down the queue server...")
	s.server.Shutdown()
	logger.Log.Info("Queue server shut down gracefully")
}
