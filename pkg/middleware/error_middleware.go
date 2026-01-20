package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/faisalyudiansah/auth-service-template/pkg/apperror"
	"github.com/faisalyudiansah/auth-service-template/pkg/constant"
	"github.com/faisalyudiansah/auth-service-template/pkg/dto"
	"github.com/faisalyudiansah/auth-service-template/pkg/utils/validationutils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var codeMap = map[int]int{
	apperror.DefaultClientErrorCode:   http.StatusBadRequest,
	apperror.DefaultServerErrorCode:   http.StatusInternalServerError,
	apperror.NotFoundErrorCode:        http.StatusNotFound,
	apperror.RequestTimeoutErrorCode:  http.StatusRequestTimeout,
	apperror.TooManyRequestsErrorCode: http.StatusTooManyRequests,
	apperror.ForbiddenAccessErrorCode: http.StatusForbidden,
	apperror.UnauthorizedErrorCode:    http.StatusUnauthorized,
}

var (
	re = regexp.MustCompile(`"([^"]*)"`)
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		errLen := len(ctx.Errors)
		if errLen > 0 {
			err := ctx.Errors.Last()

			switch e := err.Err.(type) {
			case validator.ValidationErrors:
				handleValidationError(ctx, e)
			case *apperror.FieldBindError:
				ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
					Message: constant.ValidationErrorMessage,
					Errors: []dto.FieldError{
						{
							Field:   e.Field,
							Message: e.Message,
						},
					},
				})
				return
			case *json.SyntaxError:
				handleJsonSyntaxError(ctx)
			case *json.UnmarshalTypeError:
				handleJsonUnmarshalTypeError(ctx, e)
			case *time.ParseError:
				handleParseTimeError(ctx, e)
			case *apperror.AppError:
				var originalErrorMsg *string
				orig := e.OriginalError().Error()
				if e.DisplayMessage() != orig {
					originalErrorMsg = &orig
				}
				ctx.AbortWithStatusJSON(codeMap[e.GetCode()], dto.WebResponse[any]{
					Message: e.DisplayMessage(),
					Error:   originalErrorMsg,
				})
			default:
				if errors.Is(e, io.EOF) {
					ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
						Message: constant.EOFErrorMessage,
					})
					return
				}
				if errors.Is(e, strconv.ErrSyntax) {
					match := re.FindStringSubmatch(e.Error())
					ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
						Message: fmt.Sprintf(constant.StrConvSyntaxErrorMessage, strings.ReplaceAll(match[0], "\"", "")),
					})
					return
				}

				ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.WebResponse[any]{
					Message: constant.InternalServerErrorMessage,
				})
			}
		}
	}
}

func handleJsonSyntaxError(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
		Message: constant.JsonSyntaxErrorMessage,
	})
}

func handleJsonUnmarshalTypeError(ctx *gin.Context, err *json.UnmarshalTypeError) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
		Message: fmt.Sprintf(constant.InvalidJsonValueTypeErrorMessage, err.Field),
	})
}

func handleParseTimeError(ctx *gin.Context, err *time.ParseError) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
		Message: fmt.Sprintf("please send time in format of %s, got: %s", constant.ConvertGoTimeLayoutToReadable(err.Layout), err.Value),
	})
}

func handleValidationError(ctx *gin.Context, err validator.ValidationErrors) {
	ve := []dto.FieldError{}

	for _, fe := range err {
		ve = append(ve, dto.FieldError{
			Field:   fe.Field(),
			Message: validationutils.TagToMsg(fe),
		})
	}

	ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
		Message: constant.ValidationErrorMessage,
		Errors:  ve,
	})
}
