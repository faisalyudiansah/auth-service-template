package apperror

import (
	"errors"
)

const (
	DefaultServerErrorCode = iota
	DefaultClientErrorCode
	NotFoundErrorCode
	RequestTimeoutErrorCode
	TooManyRequestsErrorCode
	ForbiddenAccessErrorCode
	UnauthorizedErrorCode
)

type AppError struct {
	err  error
	code int
	msg  string
}

func NewAppError(err error, code int, msg string) *AppError {
	return &AppError{
		err:  err,
		code: code,
		msg:  msg,
	}
}

func (e AppError) Error() string {
	if e.msg == "" {
		return e.OriginalMessage()
	}
	return e.msg
}

func (e AppError) GetCode() int {
	return e.code
}

func (e AppError) OriginalError() error {
	var currErr AppError

	currErr = e

	for {
		nextErr := currErr.err
		if nextErr == nil {
			break
		}

		var appErr AppError
		if !errors.As(nextErr, &appErr) {
			return nextErr
		}
		currErr = appErr
	}

	return e
}

func (e AppError) OriginalMessage() string {
	return e.OriginalError().Error()
}

func (e AppError) DisplayMessage() string {
	return e.msg
}

type FieldBindError struct {
	Field   string
	Message string
}

func (e *FieldBindError) Error() string {
	return e.Message
}
