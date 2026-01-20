package tasks

import (
	"context"
	"encoding/json"
	"time"

	"github.com/faisalyudiansah/auth-service-template/internal/queue/payload"

	"github.com/hibiken/asynq"
)

const (
	TypeEmailVerification   = "email:verification"
	TypeEmailForgotPassword = "email:forgot-password"
)

type EmailTask interface {
	QueueVerificationEmail(ctx context.Context, payload *payload.VerificationEmailPayload) error
	QueueForgotPasswordEmail(ctx context.Context, payload *payload.ForgotPasswordEmailPayload) error
}

type emailTaskImpl struct {
	client *asynq.Client
}

func NewEmailTask(client *asynq.Client) *emailTaskImpl {
	return &emailTaskImpl{
		client: client,
	}
}

func (t *emailTaskImpl) QueueVerificationEmail(ctx context.Context, payload *payload.VerificationEmailPayload) error {
	enqueueCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeEmailVerification, data, asynq.Timeout(5*time.Second), asynq.MaxRetry(10))
	_, err = t.client.EnqueueContext(enqueueCtx, task)

	return err
}

func (t *emailTaskImpl) QueueForgotPasswordEmail(ctx context.Context, payload *payload.ForgotPasswordEmailPayload) error {
	enqueueCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeEmailForgotPassword, data, asynq.Timeout(5*time.Second), asynq.MaxRetry(10))
	_, err = t.client.EnqueueContext(enqueueCtx, task)

	return err
}
