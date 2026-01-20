package middleware

import (
	"net/http"
	"time"

	"github.com/faisalyudiansah/auth-service-template/pkg/apperror"
	"github.com/faisalyudiansah/auth-service-template/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path

		requestID := uuid.NewString()

		reqLogger := logger.Log.WithField("request_id", requestID)

		ctx.Request = ctx.Request.WithContext(
			logger.InjectToContext(ctx.Request.Context(), reqLogger),
		)

		ctx.Next()

		params := map[string]any{
			"path":        path,
			"method":      ctx.Request.Method,
			"code_status": ctx.Writer.Status(),
			"client_ip":   ctx.ClientIP(),
			"latency":     time.Since(start).String(),
		}

		if len(ctx.Errors) == 0 {
			reqLogger.WithFields(params).Info("Incoming Request")
			return
		}

		logErrors(ctx, params, reqLogger)
	}
}

func logErrors(ctx *gin.Context, params map[string]any, log logger.Logger) {
	errors := []error{}
	for _, err := range ctx.Errors {
		switch e := err.Err.(type) {
		case validator.ValidationErrors:
			params["code_status"] = http.StatusBadRequest
			errors = append(errors, err)
		case *apperror.AppError:
			params["code_status"] = codeMap[e.GetCode()]
			errors = append(errors, e.OriginalError())
		default:
			params["code_status"] = http.StatusInternalServerError
			errors = append(errors, err)
		}
	}

	params["errors"] = errors
	log.WithFields(params).Error("Error Request")
}
