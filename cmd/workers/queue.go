package workers

import (
	"context"

	"github.com/faisalyudiansah/auth-service-template/internal/gateway/server"
)

func runQueueWorker(ctx context.Context) {
	srv := server.NewQueueServer(cfg)
	go srv.Start()

	<-ctx.Done()
	srv.Shutdown()
}
