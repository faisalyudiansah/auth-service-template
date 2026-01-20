package route

import (
	"github.com/faisalyudiansah/auth-service-template/internal/queue/processor"
	"github.com/faisalyudiansah/auth-service-template/internal/queue/tasks"

	"github.com/hibiken/asynq"
)

func EmailTaskRoute(mux *asynq.ServeMux, processor *processor.EmailTaskProcessor) {
	mux.HandleFunc(tasks.TypeEmailVerification, processor.HandleVerificationEmail)
	mux.HandleFunc(tasks.TypeEmailForgotPassword, processor.HandleForgotPasswordEmail)
}
