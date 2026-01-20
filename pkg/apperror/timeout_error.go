package apperror

import (
	"errors"

	"github.com/faisalyudiansah/auth-service-template/pkg/constant"
)

func NewTimeoutError() *AppError {
	msg := constant.RequestTimeoutErrorMessage

	err := errors.New(msg)

	return NewAppError(err, RequestTimeoutErrorCode, msg)
}
