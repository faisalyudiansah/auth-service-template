package utils

import (
	constantAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/constant"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetSessionCookie(ctx *gin.Context, sessionID uuid.UUID) {
	ctx.SetCookie(
		constantAuth.SESSION_NAME,
		sessionID.String(),
		86400, // 1 hari
		"/",
		"",
		true, // secure
		true, // httpOnly
	)
}

func ClearSessionCookie(ctx *gin.Context) {
	ctx.SetCookie(
		constantAuth.SESSION_NAME,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
}
