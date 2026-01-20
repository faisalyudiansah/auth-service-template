package apperror

import (
	"errors"

	"github.com/faisalyudiansah/auth-service-template/pkg/constant"
)

func NewLimitError() *AppError {
	msg := constant.TooManyRequestsErrorMessage

	err := errors.New(msg)

	return NewAppError(err, TooManyRequestsErrorCode, msg)
}
