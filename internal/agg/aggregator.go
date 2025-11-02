package agg

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/rugved0102/weather-agg/internal/provider"
)

// AggregatedWeather is the API-level output
type AggregatedWeather struct {
	City          string    `json:"city"`
	AvgTempC      float64   `json:"avg_temp_c"`
	AvgHumidity   int       `json:"avg_humidity"`
	Providers     []string  `json:"providers"`
	RetrievedAt   time.Time `json:"retrieved_at"`
	ProviderCount int       `json:"provider_count"`
}

type Aggregator struct {
	providers []provider.Provider
	timeout   time.Duration
}

func NewAggregator(providers []provider.Provider, timeout time.Duration) *Aggregator {
	return &Aggregator{providers: providers, timeout: timeout}
}

// Aggregate concurrently calls all providers and returns averaged results.
func (a *Aggregator) Aggregate(ctx context.Context, city string) (AggregatedWeather, error) {
	if len(a.providers) == 0 {
		return AggregatedWeather{}, errors.New("no providers configured")
	}

	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	type res struct {
		r   provider.ProviderResponse
		err error
	}

	ch := make(chan res, len(a.providers))
	var wg sync.WaitGroup

	for _, p := range a.providers {
		wg.Add(1)
		go func(pr provider.Provider) {
			defer wg.Done()
			resp, err := pr.Fetch(ctx, city)
			ch <- res{r: resp, err: err}
		}(p)
	}

	// wait in a goroutine then close channel
	go func() {
		wg.Wait()
		close(ch)
	}()

	var sumTemp float64
	var sumHum int
	var count int
	var providerNames []string
	var errs []error

	for rr := range ch {
		if rr.err != nil {
			errs = append(errs, rr.err)
			continue
		}
		sumTemp += rr.r.TempC
		sumHum += rr.r.Humidity
		count++
		providerNames = append(providerNames, rr.r.Source)
	}

	if count == 0 {
		return AggregatedWeather{}, fmt.Errorf("all providers failed: %v", errs)
	}

	avgTemp := sumTemp / float64(count)
	avgHum := sumHum / count

	agg := AggregatedWeather{
		City:          city,
		AvgTempC:      avgTemp,
		AvgHumidity:   avgHum,
		Providers:     providerNames,
		RetrievedAt:   time.Now().UTC(),
		ProviderCount: count,
	}
	return agg, nil
}
