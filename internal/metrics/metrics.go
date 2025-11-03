package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestCounter tracks total number of requests
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "weather_requests_total",
			Help: "Total number of weather API requests",
		},
		[]string{"city"},
	)

	// CacheHits tracks cache hits/misses
	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "weather_cache_hits_total",
			Help: "Number of cache hits/misses",
		},
		[]string{"type"}, // hit or miss
	)

	// RequestDuration tracks request latency
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "weather_request_duration_seconds",
			Help:    "Time taken to process weather requests",
			Buckets: []float64{0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"city"},
	)

	// ProviderSuccessCounter tracks successful provider responses
	ProviderSuccessCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "weather_provider_success_total",
			Help: "Number of successful provider responses",
		},
		[]string{"provider"},
	)

	// ProviderErrorCounter tracks provider errors
	ProviderErrorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "weather_provider_errors_total",
			Help: "Number of provider errors",
		},
		[]string{"provider"},
	)
)
