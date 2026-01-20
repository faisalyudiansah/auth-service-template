package apperror

import (
	"errors"

	"github.com/faisalyudiansah/auth-service-template/pkg/constant"
)

func NewUnauthorizedError() *AppError {
	msg := constant.UnauthorizedErrorMessage
	err := errors.New(msg)
	return NewAppError(err, UnauthorizedErrorCode, msg)
}

func NewDontHavePermissionErrorMessageError() *AppError {
	msg := constant.DontHavePermissionErrorMessage
	err := errors.New(msg)
	return NewAppError(err, DefaultClientErrorCode, msg)
}

func NewExpiredTokenError() *AppError {
	msg := constant.ExpiredTokenErrorMessage
	err := errors.New(msg)
	return NewAppError(err, UnauthorizedErrorCode, msg)
}

func NewSessionExpiredError() *AppError {
	msg := constant.SessionExpiredMessage
	err := errors.New(msg)
	return NewAppError(err, ForbiddenAccessErrorCode, msg)
}
