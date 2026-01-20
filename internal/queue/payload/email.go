package payload

type VerificationEmailPayload struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type ForgotPasswordEmailPayload struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
