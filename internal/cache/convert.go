package cache

import (
	"time"

	"github.com/rugved0102/weather-agg/internal/agg"
)

// AggregatedWeather represents the cached weather information
type AggregatedWeather struct {
	City          string    `json:"city"`
	TempC         float64   `json:"temp_c"`
	TempF         float64   `json:"temp_f"`
	Humidity      int       `json:"humidity"`
	ProviderCount int       `json:"provider_count"`
	RetrievedAt   time.Time `json:"retrieved_at"`
}

// ModelToAggregatedWeather converts from agg.AggregatedWeather to cache.AggregatedWeather
func ModelToAggregatedWeather(a agg.AggregatedWeather) AggregatedWeather {
	return AggregatedWeather{
		City:          a.City,
		TempC:         a.AvgTempC,
		TempF:         a.AvgTempC*9/5 + 32,
		Humidity:      a.AvgHumidity,
		ProviderCount: a.ProviderCount,
		RetrievedAt:   a.RetrievedAt,
	}
}
