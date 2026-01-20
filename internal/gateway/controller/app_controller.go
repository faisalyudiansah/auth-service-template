package controller

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/faisalyudiansah/auth-service-template/pkg/dto"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

type AppController struct {
	db    *sql.DB
	redis *redis.Client
}

func NewAppController(db *sql.DB, redis *redis.Client) *AppController {
	return &AppController{
		db:    db,
		redis: redis,
	}
}

func (c *AppController) Route(r *gin.Engine) {
	r.NoRoute(c.RouteNotFound)
	r.NoMethod(c.MethodNotAllowed)
	r.GET("", c.Root)
	r.GET("/health", c.Health)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	pprof.Register(r)
}

func (c *AppController) Root(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dto.WebResponse[any]{
		Message: "server is running...",
	})
}

func (c *AppController) RouteNotFound(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, dto.WebResponse[any]{
		Message: "route not found",
	})
}

func (c *AppController) MethodNotAllowed(ctx *gin.Context) {
	ctx.JSON(http.StatusMethodNotAllowed, dto.WebResponse[any]{
		Message: "method not allowed",
	})
}

type HealthStatus struct {
	Name   string  `json:"name"`
	Status string  `json:"status"`
	Error  *string `json:"error"`
}

func (c *AppController) Health(ctx *gin.Context) {
	ctxTimeout, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	checks := []struct {
		name string
		fn   func() error
	}{
		{
			name: "database",
			fn: func() error {
				return c.db.PingContext(ctxTimeout)
			},
		},
		{
			name: "redis",
			fn: func() error {
				return c.redis.Ping(ctxTimeout).Err()
			},
		},
	}

	results := make([]HealthStatus, 0, len(checks))
	isReady := true

	for _, check := range checks {
		status := buildHealthStatus(check.name, check.fn())
		if status.Status == "down" {
			isReady = false
		}
		results = append(results, status)
	}

	httpStatus := http.StatusOK
	message := "service is ready"

	if !isReady {
		httpStatus = http.StatusServiceUnavailable
		message = "service not ready"
	}

	ctx.JSON(httpStatus, dto.WebResponse[[]HealthStatus]{
		Message: message,
		Data:    results,
	})
}

func buildHealthStatus(name string, err error) HealthStatus {
	if err != nil {
		msg := err.Error()
		return HealthStatus{
			Name:   name,
			Status: "down",
			Error:  &msg,
		}
	}

	return HealthStatus{
		Name:   name,
		Status: "up",
		Error:  nil,
	}
}
