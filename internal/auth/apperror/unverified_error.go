package apperror

import (
	"errors"

	"github.com/faisalyudiansah/auth-service-template/internal/auth/constant"
	apperrorPkg "github.com/faisalyudiansah/auth-service-template/pkg/apperror"
)

func NewUnverifiedError() *apperrorPkg.AppError {
	msg := constant.UnverifiedErrorMessage

	err := errors.New(msg)

	return apperrorPkg.NewAppError(err, apperrorPkg.ForbiddenAccessErrorCode, msg)
}
