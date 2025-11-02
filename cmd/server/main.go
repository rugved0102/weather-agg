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

	"github.com/rugved0102/weather-agg/internal/agg"
	"github.com/rugved0102/weather-agg/internal/cache"
	"github.com/rugved0102/weather-agg/internal/provider"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisClient := cache.NewRedis(redisAddr, "") // assumes NewRedis(addr, pwd)

	// build providers
	providers := []provider.Provider{
		provider.NewOpenMeteo(),                 // Real provider - no API key needed
		provider.NewMockProvider("mock-A", 0.3), // Backup mock provider
	}

	aggregator := agg.NewAggregator(providers, 6*time.Second)

	mux := http.NewServeMux()
	mux.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		city := r.URL.Query().Get("city")
		if city == "" {
			http.Error(w, "missing city", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()

		// check cache first
		if aggCached, err := redisClient.GetAggregated(ctx, city); err == nil {
			log.Printf("Cache hit for %s", city)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(aggCached)
			return
		}

		// not cached -> aggregate
		aggRes, err := aggregator.Aggregate(ctx, city)
		if err != nil {
			http.Error(w, "aggregation error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// transform to cache shape used earlier
		cacheAgg := cache.ModelToAggregatedWeather(aggRes)
		_ = redisClient.SetAggregated(ctx, cacheAgg, 10*time.Minute)
		_ = redisClient.StoreAggregated(ctx, cacheAgg)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(aggRes)
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
