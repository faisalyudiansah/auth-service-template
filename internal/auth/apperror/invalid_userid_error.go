package apperror

import (
	"errors"

	"github.com/faisalyudiansah/auth-service-template/internal/auth/constant"
	"github.com/faisalyudiansah/auth-service-template/pkg/apperror"
)

func NewInvalidUserIdError() *apperror.AppError {
	msg := constant.InvalidUserId

	err := errors.New(msg)

	return apperror.NewAppError(err, apperror.DefaultClientErrorCode, msg)
}

func NewUserDetailNotExistsError() *apperror.AppError {
	msg := constant.UserDetailNotFoundErrorMessage

	err := errors.New(msg)

	return apperror.NewAppError(err, apperror.NotFoundErrorCode, msg)
}

func NewAccoundIsNotValidError() *apperror.AppError {
	msg := constant.AccountIsNotValidErrorMessage

	err := errors.New(msg)

	return apperror.NewAppError(err, apperror.DefaultClientErrorCode, msg)
}
