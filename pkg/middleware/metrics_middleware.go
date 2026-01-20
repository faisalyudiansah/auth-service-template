package middleware

import (
	"runtime"
	"strconv"
	"time"

	"github.com/faisalyudiansah/auth-service-template/pkg/metrics"

	"github.com/gin-gonic/gin"
)

func Metrics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()

		latency := time.Since(start)
		fullPath := ctx.FullPath()
		method := ctx.Request.Method
		statusCode := ctx.Writer.Status()

		requestCountMetric(method, fullPath, statusCode)
		successRateMetric(fullPath, statusCode)
		trafficLevelsMetric(fullPath)
		requestLatencyMetric(method, fullPath, latency)
		memoryUsageMetric()
	}
}

func requestCountMetric(method, fullPath string, statusCode int) {
	metrics.RequestCount.WithLabelValues(method, fullPath, strconv.Itoa(statusCode)).Inc()
}

func successRateMetric(fullPath string, statusCode int) {
	if statusCode >= 200 && statusCode < 300 {
		metrics.SuccessRate.WithLabelValues(fullPath).Set(1)
	} else {
		metrics.SuccessRate.WithLabelValues(fullPath).Set(0)
	}
}

func trafficLevelsMetric(fullPath string) {
	metrics.TrafficLevels.WithLabelValues(fullPath).Inc()
}

func requestLatencyMetric(method, fullPath string, latency time.Duration) {
	metrics.RequestLatency.WithLabelValues(method, fullPath).Observe(latency.Seconds())
}

func memoryUsageMetric() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	metrics.MemoryUsage.Set(float64(m.Alloc))
}
