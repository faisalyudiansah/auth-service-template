package apperror

import (
	"github.com/faisalyudiansah/auth-service-template/pkg/constant"
)

func NewClientError(err error, msg *string) *AppError {
	errMsg := constant.ClientDefaultErrorMessage

	if msg != nil {
		errMsg = *msg
	}

	return NewAppError(err, DefaultClientErrorCode, errMsg)
}
