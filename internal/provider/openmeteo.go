package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type OpenMeteoProvider struct {
	Client *http.Client
}

func NewOpenMeteo() *OpenMeteoProvider {
	return &OpenMeteoProvider{
		Client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (p *OpenMeteoProvider) Name() string { return "openmeteo" }

func (p *OpenMeteoProvider) Fetch(ctx context.Context, city string) (ProviderResponse, error) {
	// First get coordinates for the city using geocoding API
	geoURL := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1", url.QueryEscape(city))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, geoURL, nil)
	if err != nil {
		return ProviderResponse{}, err
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return ProviderResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ProviderResponse{}, fmt.Errorf("openmeteo geocoding: status %d", resp.StatusCode)
	}

	var geoResp struct {
		Results []struct {
			Name      string  `json:"name"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return ProviderResponse{}, err
	}

	if len(geoResp.Results) == 0 {
		return ProviderResponse{}, fmt.Errorf("city not found: %s", city)
	}

	location := geoResp.Results[0]

	// Now get weather data using coordinates
	weatherURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.6f&longitude=%.6f&current=temperature_2m,relative_humidity_2m",
		location.Latitude,
		location.Longitude,
	)

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, weatherURL, nil)
	if err != nil {
		return ProviderResponse{}, err
	}

	resp, err = p.Client.Do(req)
	if err != nil {
		return ProviderResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ProviderResponse{}, fmt.Errorf("openmeteo weather: status %d", resp.StatusCode)
	}

	var weatherResp struct {
		Current struct {
			Temperature float64 `json:"temperature_2m"`
			Humidity    float64 `json:"relative_humidity_2m"`
		} `json:"current"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return ProviderResponse{}, err
	}

	return ProviderResponse{
		City:     location.Name,
		TempC:    weatherResp.Current.Temperature,
		Humidity: int(weatherResp.Current.Humidity),
		Source:   p.Name(),
	}, nil
}
