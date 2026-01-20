package workers

import (
	"context"

	"github.com/faisalyudiansah/auth-service-template/internal/gateway/server"
)

func runHttpWorker(ctx context.Context) {
	srv := server.NewHttpServer(cfg)
	go srv.Start()

	<-ctx.Done()
	srv.Shutdown()
}
