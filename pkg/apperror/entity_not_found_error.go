package apperror

import (
	"errors"
	"fmt"

	"github.com/faisalyudiansah/auth-service-template/pkg/constant"
)

func NewEntityNotFoundError(entity string) *AppError {
	msg := fmt.Sprintf(constant.EntityNotFoundErrorMessage, entity)

	err := errors.New(msg)

	return NewAppError(err, DefaultClientErrorCode, msg)
}

func NewNoRowsError(err error, entity any) *AppError {
	msg := fmt.Sprintf(constant.EntityNotFoundErrorMessage, entity)

	return NewAppError(err, DefaultClientErrorCode, msg)
}
