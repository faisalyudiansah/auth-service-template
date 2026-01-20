package apperror

import (
	"errors"

	"github.com/faisalyudiansah/auth-service-template/pkg/constant"
)

func NewInvalidIdError() *AppError {
	msg := constant.InvalidIDErrorMessage

	err := errors.New(msg)

	return NewAppError(err, DefaultClientErrorCode, msg)
}
