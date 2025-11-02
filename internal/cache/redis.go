package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(addr, password string) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &Redis{client: client}
}

func (r *Redis) GetAggregated(ctx context.Context, city string) (AggregatedWeather, error) {
	data, err := r.client.Get(ctx, "weather:"+city).Result()
	if err != nil {
		return AggregatedWeather{}, err
	}

	var weather AggregatedWeather
	if err := json.Unmarshal([]byte(data), &weather); err != nil {
		return AggregatedWeather{}, err
	}

	return weather, nil
}

func (r *Redis) SetAggregated(ctx context.Context, weather AggregatedWeather, ttl time.Duration) error {
	data, err := json.Marshal(weather)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, "weather:"+weather.City, data, ttl).Err()
}

func (r *Redis) StoreAggregated(ctx context.Context, weather AggregatedWeather) error {
	// Store in Redis with 1 hour TTL
	return r.SetAggregated(ctx, weather, time.Hour)
}
