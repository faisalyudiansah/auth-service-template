package apperror

import (
	"errors"

	"github.com/faisalyudiansah/auth-service-template/internal/auth/constant"
	"github.com/faisalyudiansah/auth-service-template/pkg/apperror"
)

func NewExpiredTokenError() *apperror.AppError {
	msg := constant.ExpiredTokenErrorMessage

	err := errors.New(msg)

	return apperror.NewAppError(err, apperror.DefaultClientErrorCode, msg)
}
