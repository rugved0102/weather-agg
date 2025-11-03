// Package main implements a weather aggregation service that provides
// consolidated weather information from various sources (currently mock data).
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rugved0102/weather-agg/internal/agg"
	"github.com/rugved0102/weather-agg/internal/cache"
	"github.com/rugved0102/weather-agg/internal/metrics"
	"github.com/rugved0102/weather-agg/internal/provider"
	"github.com/rugved0102/weather-agg/internal/store"
)

func main() {
	// Initialize Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisClient := cache.NewRedis(redisAddr, "") // assumes NewRedis(addr, pwd)

	// Initialize PostgreSQL
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://weather:weatherpass@localhost:5432/weatherdb?sslmode=disable"
	}

	db, err := store.NewStore(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// build providers
	providers := []provider.Provider{
		provider.NewOpenMeteo(),                 // Real provider - no API key needed
		provider.NewMockProvider("mock-A", 0.3), // Backup mock provider
	}

	aggregator := agg.NewAggregator(providers, 6*time.Second)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		city := r.URL.Query().Get("city")
		if city == "" {
			http.Error(w, "missing city", http.StatusBadRequest)
			return
		}

		// Track request count
		metrics.RequestCounter.WithLabelValues(city).Inc()

		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()

		// check cache first
		if aggCached, err := redisClient.GetAggregated(ctx, city); err == nil {
			metrics.CacheHits.WithLabelValues("hit").Inc()
			log.Printf("Cache hit for %s", city)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(aggCached)
			metrics.RequestDuration.WithLabelValues(city).Observe(time.Since(start).Seconds())
			return
		}
		metrics.CacheHits.WithLabelValues("miss").Inc()

		// not cached -> aggregate
		aggRes, err := aggregator.Aggregate(ctx, city)
		if err != nil {
			http.Error(w, "aggregation error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Track successful provider responses
		for _, prov := range aggRes.Providers {
			metrics.ProviderSuccessCounter.WithLabelValues(prov).Inc()
		}

		// transform to cache shape used earlier
		cacheAgg := cache.ModelToAggregatedWeather(aggRes)
		_ = redisClient.SetAggregated(ctx, cacheAgg, 10*time.Minute)
		_ = redisClient.StoreAggregated(ctx, cacheAgg)

		// Save to PostgreSQL for analytics
		if err := db.SaveWeather(city, aggRes.AvgTempC, float64(aggRes.AvgHumidity), len(aggRes.Providers)); err != nil {
			log.Printf("Failed to save weather data to database: %v", err)
			// Don't fail the request if database save fails
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(aggRes)
		metrics.RequestDuration.WithLabelValues(city).Observe(time.Since(start).Seconds())
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:         ":" + getenv("PORT", "8080"),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("weather-agg starting on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

// simple getenv helper
func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
