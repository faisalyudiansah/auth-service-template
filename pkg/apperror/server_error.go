package apperror

import (
	"github.com/faisalyudiansah/auth-service-template/pkg/constant"
)

func NewServerError(err error) *AppError {
	msg := constant.InternalServerErrorMessage

	return NewAppError(err, DefaultServerErrorCode, msg)
}
