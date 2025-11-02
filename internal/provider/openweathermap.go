package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type OpenWeatherProvider struct {
	APIKey string
	Client *http.Client
}

func NewOpenWeather(apiKey string) *OpenWeatherProvider {
	return &OpenWeatherProvider{
		APIKey: apiKey,
		Client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (p *OpenWeatherProvider) Name() string { return "openweathermap" }

func (p *OpenWeatherProvider) Fetch(ctx context.Context, city string) (ProviderResponse, error) {
	// If API key missing, return error
	if p.APIKey == "" {
		return ProviderResponse{}, fmt.Errorf("openweathermap: api key not set")
	}

	u := "https://api.openweathermap.org/data/2.5/weather"
	q := url.Values{}
	q.Set("q", city)
	q.Set("appid", p.APIKey)
	q.Set("units", "metric") // we want Celsius

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u+"?"+q.Encode(), nil)
	if err != nil {
		return ProviderResponse{}, err
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return ProviderResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ProviderResponse{}, fmt.Errorf("openweathermap: status %d", resp.StatusCode)
	}

	var body struct {
		Name string `json:"name"`
		Main struct {
			Temp     float64 `json:"temp"`
			Humidity int     `json:"humidity"`
		} `json:"main"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return ProviderResponse{}, err
	}

	return ProviderResponse{
		City:     body.Name,
		TempC:    body.Main.Temp,
		Humidity: body.Main.Humidity,
		Source:   p.Name(),
	}, nil
}
