package provider

import (
	"context"
	"math"
	"time"
)

type MockProvider struct {
	Offset float64
	NameID string
}

func NewMockProvider(name string, offset float64) *MockProvider {
	return &MockProvider{Offset: offset, NameID: name}
}

func (m *MockProvider) Name() string { return m.NameID }

func (m *MockProvider) Fetch(ctx context.Context, city string) (ProviderResponse, error) {
	// produce deterministic fake data derived from city length + offset
	baseTemp := float64(len(city)) + 20.0
	t := math.Round((baseTemp+m.Offset)*10) / 10
	// fake humidity
	h := 50 + int(m.Offset)
	// small artificial delay to simulate network
	time.Sleep(200 * time.Millisecond)
	return ProviderResponse{
		City:     city,
		TempC:    t,
		Humidity: h,
		Source:   m.Name(),
	}, nil
}
