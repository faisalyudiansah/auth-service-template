package ginutils

import (
	"net/http"
	"reflect"

	"github.com/faisalyudiansah/auth-service-template/pkg/constant"
	"github.com/faisalyudiansah/auth-service-template/pkg/dto"

	"github.com/gin-gonic/gin"
)

func ResponseOK[T any](ctx *gin.Context, data T) {
	ResponseJSON(ctx, http.StatusOK, constant.ResponseSuccessMessage, data, nil)
}

func ResponseOKPlain(ctx *gin.Context) {
	ResponseOK[any](ctx, nil)
}

func ResponseOKPagination[T any](ctx *gin.Context, data T, paging *dto.PageMetaData) {
	ResponseJSON(ctx, http.StatusOK, constant.ResponseSuccessMessage, data, paging)
}

func ResponseOKSeekPagination[T any](ctx *gin.Context, data T, paging *dto.SeekPageMetaData) {
	ResponseSeekJSON(ctx, http.StatusOK, constant.ResponseSuccessMessage, data, paging)
}

func ResponseCreated[T any](ctx *gin.Context, data T) {
	ResponseJSON(ctx, http.StatusCreated, constant.ResponseSuccessMessage, data, nil)
}

func ResponseCreatedPlain(ctx *gin.Context) {
	ResponseCreated[any](ctx, nil)
}

func ResponseNoContent(ctx *gin.Context) {
	ResponseJSON[any](ctx, http.StatusNoContent, constant.ResponseSuccessMessage, nil, nil)
}

func ResponseSeekJSON[T any](ctx *gin.Context, statusCode int, message string, data T, paging *dto.SeekPageMetaData) {
	ctx.JSON(statusCode, dto.WebResponseSeek[T]{
		Message: message,
		Data:    data,
		Paging:  paging,
	})
}

func ResponseJSON[T any](ctx *gin.Context, statusCode int, message string, data T, paging *dto.PageMetaData) {
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Slice && v.IsNil() {
		ctx.JSON(statusCode, dto.WebResponseDataEmptyArray[[]T]{
			Message: message,
			Data:    []T{},
			Paging:  paging,
		})
		return
	}

	if v.Kind() == reflect.Slice && v.Len() == 0 {
		ctx.JSON(statusCode, dto.WebResponseDataEmptyArray[T]{
			Message: message,
			Data:    data,
			Paging:  paging,
		})
		return
	}

	ctx.JSON(statusCode, dto.WebResponse[T]{
		Message: message,
		Data:    data,
		Paging:  paging,
	})
}
