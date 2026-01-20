package smtputils

import (
	"embed"
)

//go:embed templates/*.html
var EmailHTMLTemplates embed.FS

const (
	ResetPasswordSubject = "authservice - Please reset your password"
	VerificationSubject  = "authservice - Verify your account"
)

type EmailTemplate string

const (
	ResetPasswordTemplate EmailTemplate = "templates/forgot-password.html"
	VerificationTemplate  EmailTemplate = "templates/verification.html"
)
