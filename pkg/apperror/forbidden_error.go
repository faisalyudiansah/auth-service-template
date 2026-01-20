package apperror

import (
	"errors"

	"github.com/faisalyudiansah/auth-service-template/pkg/constant"
)

func NewForbiddenAccessError() *AppError {
	msg := constant.ForbiddenAccessErrorMessage

	err := errors.New(msg)

	return NewAppError(err, ForbiddenAccessErrorCode, msg)
}
