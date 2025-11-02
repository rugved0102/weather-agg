// Package main implements a weather aggregation service that provides
// consolidated weather information from various sources (currently mock data).
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
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

	// TODO: Replace mock data with real weather service integration
	resp := AggregatedWeather{
		City:        city,
		TempC:       27.3,
		HumidityPct: 62,
		RetrievedAt: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func main() {
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
	log.Println("weather-agg starting on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
