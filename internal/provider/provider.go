package provider

import "context"

// ProviderResponse is a normalized provider result
type ProviderResponse struct {
	City     string
	TempC    float64
	Humidity int
	Source   string // provider name
}

// Provider is the interface each provider must implement
type Provider interface {
	Fetch(ctx context.Context, city string) (ProviderResponse, error)
	Name() string
}
