package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "statusCode"},
	)
	SuccessRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_success_rate",
			Help: "Success rate of HTTP requests",
		},
		[]string{"endpoint"},
	)
	TrafficLevels = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_traffic_levels",
			Help: "Number of HTTP requests over time",
		},
		[]string{"endpoint"},
	)
	RequestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_latency_seconds",
			Help:    "Latency of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
	MemoryUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Memory usage in bytes",
		},
	)
)

func init() {
	prometheus.MustRegister(
		RequestCount,
		SuccessRate,
		TrafficLevels,
		RequestLatency,
		MemoryUsage,
	)
}
