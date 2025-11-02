// Package main implements a weather aggregation service that provides
// consolidated weather information from various sources (currently mock data).
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rugved0102/weather-agg/pkg/logger"
)

// AggregatedWeather represents the consolidated weather information
// returned by the API. All fields are exported for JSON serialization.
type AggregatedWeather struct {
	// City name for which weather is reported
	City string `json:"city"`

	// Temperature in Celsius
	TempC float64 `json:"temp_c"`

	// Humidity percentage (0-100)
	HumidityPct int `json:"humidity"`

	// ISO 8601 timestamp when the data was retrieved
	RetrievedAt string `json:"retrieved_at"`
}

var (
	// Redis client instance
	rdb *redis.Client
)

// initRedis initializes the Redis client
func initRedis() error {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "redis" // default to service name in docker-compose
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: redisHost + ":6379",
		DB:   0,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return err
	}

	logger.Info("redis_connect", "cache", "Successfully connected to Redis", map[string]interface{}{
		"host": redisHost,
	})
	return nil
}

// getWeatherFromCache attempts to retrieve weather data from Redis cache
func getWeatherFromCache(ctx context.Context, city string) (*AggregatedWeather, error) {
	data, err := rdb.Get(ctx, "weather:"+city).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, err
	}

	var weather AggregatedWeather
	if err := json.Unmarshal([]byte(data), &weather); err != nil {
		return nil, err
	}

	return &weather, nil
}

// setWeatherInCache stores weather data in Redis with a 1-hour expiration
func setWeatherInCache(ctx context.Context, weather *AggregatedWeather) error {
	data, err := json.Marshal(weather)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, "weather:"+weather.City, data, time.Hour).Err()
}

// getMockWeather generates mock weather data (to be replaced with real API calls)
func getMockWeather(city string) *AggregatedWeather {
	return &AggregatedWeather{
		City:        city,
		TempC:       27.3,
		HumidityPct: 62,
		RetrievedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// weatherHandler processes GET requests to the /weather endpoint.
// It accepts an optional 'city' query parameter and returns weather
// information in JSON format.
//
// Currently returns mock data. Future versions will aggregate real
// weather data from multiple providers.
func weatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		city = "Unknown"
	}

	// Try to get from cache first
	weather, err := getWeatherFromCache(r.Context(), city)
	if err != nil {
		logger.Error("cache_error", "cache", "Failed to get weather from cache", map[string]interface{}{
			"city":  city,
			"error": err.Error(),
		})
		// Continue with mock data if cache fails
	}

	// If not in cache, get mock data and cache it
	if weather == nil {
		weather = getMockWeather(city)
		if err := setWeatherInCache(r.Context(), weather); err != nil {
			logger.Error("cache_set", "cache", "Failed to cache weather data", map[string]interface{}{
				"city":  city,
				"error": err.Error(),
			})
			// Continue anyway, just won't be cached
		} else {
			logger.Info("cache_set", "cache", "Successfully cached weather data", map[string]interface{}{
				"city": city,
				"ttl":  "1h",
			})
		}
	} else {
		logger.Info("cache_hit", "cache", "Retrieved weather from cache", map[string]interface{}{
			"city": city,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weather); err != nil {
		logger.Error("response_encode", "http", "Failed to encode JSON response", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Initialize Redis connection
	if err := initRedis(); err != nil {
		logger.Error("startup", "main", "Failed to connect to Redis", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}

	// Create a new server mux for routing
	mux := http.NewServeMux()

	// Register the weather endpoint with timeout handling
	mux.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		// Set a timeout for the entire request processing
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()

		// Update request context with timeout
		r = r.WithContext(ctx)
		weatherHandler(w, r)
	})

	// Configure the HTTP server with timeouts
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,  // Time to read request headers and body
		WriteTimeout: 10 * time.Second, // Time to write response
		IdleTimeout:  60 * time.Second, // Time to keep idle connections
	}

	// Start the server
	logger.Info("startup", "main", "Starting weather-agg server", map[string]interface{}{
		"address": ":8080",
	})

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("startup", "main", "Server failed", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}
}
