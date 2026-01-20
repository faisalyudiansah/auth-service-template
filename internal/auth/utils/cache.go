package utils

import (
	constantAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/constant"

	"fmt"

	"github.com/google/uuid"
)

const (
	resetTokenKey        = "reset"
	verificationTokenKey = "verification"
)

func VerificationTokenCacheKey(email string) string {
	return fmt.Sprintf("%v:%v", email, verificationTokenKey)
}

func ResetTokenCacheKey(email string) string {
	return fmt.Sprintf("%v:%v", email, resetTokenKey)
}

func SessionKey(sessionID uuid.UUID) string {
	return fmt.Sprintf("%v:%v", constantAuth.SESSION_ID, sessionID)
}
