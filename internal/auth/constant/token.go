package constant

import "time"

var (
	ResetTokenExpireDuration        = 5 * time.Minute
	VerificationTokenExpireDuration = 5 * time.Minute
)

var (
	ResetTokenCooldownDuration        = 1 * time.Minute
	VerificationTokenCooldownDuration = 1 * time.Minute
)
