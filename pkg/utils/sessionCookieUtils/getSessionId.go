package sessioncookieutils

import (
	constantAuth "github.com/faisalyudiansah/auth-service-template/internal/auth/constant"
	"github.com/faisalyudiansah/auth-service-template/pkg/apperror"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetSessionIDFromCookie(ctx *gin.Context) (uuid.UUID, error) {
	sessionIDstr, err := ctx.Cookie(constantAuth.SESSION_NAME)
	if err != nil {
		return uuid.Nil, apperror.NewForbiddenAccessError()
	}

	sessionID, err := uuid.Parse(sessionIDstr)
	if err != nil {
		ctx.Error(apperror.NewForbiddenAccessError())
		return uuid.Nil, apperror.NewForbiddenAccessError()
	}
	return sessionID, nil
}
