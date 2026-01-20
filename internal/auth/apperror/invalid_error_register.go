package apperror

import (
	"errors"

	"github.com/faisalyudiansah/auth-service-template/internal/auth/constant"
	"github.com/faisalyudiansah/auth-service-template/pkg/apperror"
)

func NewInvalidEmailAlreadyExists(err error) *apperror.AppError {
	msg := constant.InvalidEmailAlreadyExists

	if err == nil {
		err = errors.New(msg)
	}

	return apperror.NewAppError(err, apperror.DefaultClientErrorCode, msg)
}

func NewInvalidPhoneNumberAlreadyExists() *apperror.AppError {
	msg := constant.InvalidPhoneNumberAlreadyExists
	err := errors.New(msg)
	return apperror.NewAppError(err, apperror.DefaultClientErrorCode, msg)
}
