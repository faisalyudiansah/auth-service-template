package constant

import constantAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/constant"

type contextKey string

const (
	ContextSessionID contextKey = constantAuth.SESSION_ID
	ContextUserID    contextKey = "user_id"
	ContextRole      contextKey = "role"
	ContextJTI       contextKey = "jti"
)
