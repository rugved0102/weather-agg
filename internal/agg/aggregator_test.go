package agg

import (
	"context"
	"testing"
	"time"

	"github.com/rugved0102/weather-agg/internal/provider"
)

func TestAggregator_Aggregate(t *testing.T) {
	provs := []provider.Provider{
		provider.NewMockProvider("m1", 0.1),
		provider.NewMockProvider("m2", -0.2),
	}
	a := NewAggregator(provs, 3*time.Second)
	ctx := context.Background()
	res, err := a.Aggregate(ctx, "Pune")
	if err != nil {
		t.Fatalf("aggregate failed: %v", err)
	}
	if res.ProviderCount != 2 {
		t.Fatalf("expected 2 providers, got %d", res.ProviderCount)
	}
	if res.AvgTempC == 0 {
		t.Fatalf("expected non-zero temp")
	}
}
