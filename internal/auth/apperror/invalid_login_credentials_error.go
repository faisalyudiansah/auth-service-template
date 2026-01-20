package apperror

import (
	"errors"

	"github.com/faisalyudiansah/auth-service-template/internal/auth/constant"
	"github.com/faisalyudiansah/auth-service-template/pkg/apperror"
)

func NewInvalidLoginCredentials(err error) *apperror.AppError {
	msg := constant.InvalidLoginCredentialsErrorMessage

	if err == nil {
		err = errors.New(msg)
	}

	return apperror.NewAppError(err, apperror.DefaultClientErrorCode, msg)
}
