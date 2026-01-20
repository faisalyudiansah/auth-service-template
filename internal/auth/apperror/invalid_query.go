package apperror

import (
	"errors"

	"github.com/faisalyudiansah/auth-service-template/internal/auth/constant"
	"github.com/faisalyudiansah/auth-service-template/pkg/apperror"
)

func NewInvalidQueryLimitError() *apperror.AppError {
	msg := constant.InvalidQueryLimit
	err := errors.New(msg)
	return apperror.NewAppError(err, apperror.DefaultClientErrorCode, msg)
}

func NewInvalidQueryPageError() *apperror.AppError {
	msg := constant.InvalidQueryPage
	err := errors.New(msg)
	return apperror.NewAppError(err, apperror.DefaultClientErrorCode, msg)
}

func NewInvalidQueryRoleError() *apperror.AppError {
	msg := constant.InvalidQueryRole
	err := errors.New(msg)
	return apperror.NewAppError(err, apperror.DefaultClientErrorCode, msg)
}

func NewInvalidQueryisAssignError() *apperror.AppError {
	msg := constant.InvalidQueryIsAssign
	err := errors.New(msg)
	return apperror.NewAppError(err, apperror.DefaultClientErrorCode, msg)
}

func NewInvalidQueryisAssignAndRoleError() *apperror.AppError {
	msg := constant.InvalidQueryisAssignAndRole
	err := errors.New(msg)
	return apperror.NewAppError(err, apperror.DefaultClientErrorCode, msg)
}
