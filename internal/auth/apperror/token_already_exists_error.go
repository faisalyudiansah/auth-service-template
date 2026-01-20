package apperror

import (
	"errors"

	"github.com/faisalyudiansah/auth-service-template/internal/auth/constant"
	"github.com/faisalyudiansah/auth-service-template/pkg/apperror"
)

func NewTokenAlreadyExistsError() *apperror.AppError {
	msg := constant.TokenAlreadyExistsErrorMessage

	err := errors.New(msg)

	return apperror.NewAppError(err, apperror.TooManyRequestsErrorCode, msg)
}
