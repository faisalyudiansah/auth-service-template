package apperror

import (
	"errors"

	"github.com/faisalyudiansah/auth-service-template/internal/auth/constant"
	"github.com/faisalyudiansah/auth-service-template/pkg/apperror"
)

func NewEmailNotExistsError() *apperror.AppError {
	msg := constant.EmailNotExistsErrorMessage

	err := errors.New(msg)

	return apperror.NewAppError(err, apperror.DefaultClientErrorCode, msg)
}
