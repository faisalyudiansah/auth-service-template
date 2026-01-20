package smtputils

import (
	"bytes"
	"context"
	"html/template"
	"sync"

	"github.com/faisalyudiansah/auth-service-template/pkg/config"

	"gopkg.in/gomail.v2"
)

var (
	dialer   *gomail.Dialer
	initOnce sync.Once
)

func initDialer(config *config.SMTPConfig) {
	initOnce.Do(func() {
		dialer = gomail.NewDialer(config.Host, config.Port, config.Email, "")
	})
}

type SMTPUtils interface {
	SendMail(to, subject, body string) error
	SendMailHTML(to, subject string, EmailTemplate EmailTemplate, data map[string]any) error
	SendMailHTMLContext(ctx context.Context, to, subject string, EmailTemplate EmailTemplate, data map[string]any) error
}

type smtpUtils struct {
	config *config.SMTPConfig
}

func NewSMTPUtils(config *config.SMTPConfig) *smtpUtils {
	initDialer(config)

	return &smtpUtils{
		config: config,
	}
}

func (s *smtpUtils) SendMail(to, subject, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", s.config.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	return dialer.DialAndSend(m)
}

func (s *smtpUtils) SendMailHTML(to, subject string, EmailTemplate EmailTemplate, data map[string]any) error {
	tmpl, err := template.ParseFS(EmailHTMLTemplates, string(EmailTemplate))
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	return dialer.DialAndSend(m)
}

func (s *smtpUtils) SendMailHTMLContext(ctx context.Context, to, subject string, EmailTemplate EmailTemplate, data map[string]any) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return s.SendMailHTML(to, subject, EmailTemplate, data)
	}
}
