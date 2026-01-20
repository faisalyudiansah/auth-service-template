package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/faisalyudiansah/auth-service-template/internal/queue/payload"
	"github.com/faisalyudiansah/auth-service-template/pkg/config"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/encryptutils"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/smtputils"

	"github.com/hibiken/asynq"
)

type EmailTaskProcessor struct {
	base64Encryptor encryptutils.Base64Encryptor
	smtpUtil        smtputils.SMTPUtils
	cfg             *config.Config
}

func NewEmailTaskProcessor(
	base64Encryptor encryptutils.Base64Encryptor,
	smtpUtil smtputils.SMTPUtils,
	cfg *config.Config,
) *EmailTaskProcessor {
	return &EmailTaskProcessor{
		base64Encryptor: base64Encryptor,
		smtpUtil:        smtpUtil,
		cfg:             cfg,
	}
}

func (p *EmailTaskProcessor) HandleVerificationEmail(ctx context.Context, t *asynq.Task) error {
	payload := new(payload.VerificationEmailPayload)
	if err := json.Unmarshal(t.Payload(), payload); err != nil {
		return err
	}

	encodedEmail := p.base64Encryptor.EncodeURL(payload.Email)
	encodedToken := p.base64Encryptor.EncodeURL(payload.Token)
	err := p.smtpUtil.SendMailHTMLContext(
		ctx,
		payload.Email,
		smtputils.VerificationSubject,
		smtputils.VerificationTemplate,
		map[string]any{
			"Link": fmt.Sprintf("%s/verify-account?token=%v&email=%v", p.cfg.URLClientConfig.URLCientVerificationEmail, encodedToken, encodedEmail),
		},
	)

	return err
}

func (p *EmailTaskProcessor) HandleForgotPasswordEmail(ctx context.Context, t *asynq.Task) error {
	payload := new(payload.ForgotPasswordEmailPayload)
	if err := json.Unmarshal(t.Payload(), payload); err != nil {
		return err
	}

	encodedEmail := p.base64Encryptor.EncodeURL(payload.Email)
	encodedToken := p.base64Encryptor.EncodeURL(payload.Token)
	err := p.smtpUtil.SendMailHTMLContext(
		ctx,
		payload.Email,
		smtputils.ResetPasswordSubject,
		smtputils.ResetPasswordTemplate,
		map[string]any{
			"Link": fmt.Sprintf("%s/reset-password?token=%v&email=%v", p.cfg.URLClientConfig.URLClientForgotPassword, encodedToken, encodedEmail),
		},
	)

	return err
}
