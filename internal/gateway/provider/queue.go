package provider

import (
	"github.com/faisalyudiansah/auth-service-template/internal/queue/processor"
	"github.com/faisalyudiansah/auth-service-template/internal/queue/route"
	"github.com/faisalyudiansah/auth-service-template/internal/queue/tasks"
	"github.com/faisalyudiansah/auth-service-template/pkg/config"

	"github.com/hibiken/asynq"
)

var (
	emailTask tasks.EmailTask
)

var (
	emailTaskProcessor *processor.EmailTaskProcessor
)

func ProvideQueueModule(client *asynq.Client, mux *asynq.ServeMux, cfg *config.Config) {
	injectQueueModuleTask(client)
	injectQueueModuleProcessor(cfg)

	route.EmailTaskRoute(mux, emailTaskProcessor)
}

func injectQueueModuleTask(client *asynq.Client) {
	if emailTask == nil {
		emailTask = tasks.NewEmailTask(client)
	}
}

func injectQueueModuleProcessor(cfg *config.Config) {
	emailTaskProcessor = processor.NewEmailTaskProcessor(base64Encryptor, smtpUtil, cfg)
}
