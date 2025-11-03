# How It Works - Technical Guide

This guide explains how the weather service is built and how all the pieces work together.

## The Big Picture

The service is like a weather data collector:
1. It asks different weather providers for data
2. Combines their answers
3. Saves the result for quick access later
4. Serves it through a simple API

## Main Parts and How They Work

### 1. Weather Data Sources (Providers)

Think of providers like weather reporters. Each one:
- Knows how to get weather data
- Speaks the same language (uses the same data format)
- Works independently of others
- Has its own timeout (won't make others wait too long)

We have two providers:
1. Open-Meteo Provider
   - Gets real weather data
   - Free to use, no API key needed
   - Very reliable

2. Mock Provider
   - Creates fake weather data
   - Useful for testing
   - Works offline
   - Backup when needed

How they work:
```go
// Every provider must do these things:
type Provider interface {
    // Get weather for a city
    Fetch(ctx context.Context, city string) (ProviderResponse, error)
    // Tell us its name
    Name() string
}
```

### 2. Data Collector (Aggregator)

The aggregator is like a coordinator that:
1. Asks all providers for weather data at the same time
2. Doesn't wait too long for slow providers
3. Combines all the responses it gets
4. Calculates averages for temperature and humidity

Smart features:
- Asks all providers at once (fast!)
- Still works if some providers fail
- Times out if things take too long
- Gives partial results if needed

For example:
- Provider 1 says: 20°C, 60% humidity
- Provider 2 says: 22°C, 64% humidity
- Aggregator returns: 21°C, 62% humidity (average)

### 3. Quick Memory (Cache)

We use Redis (a fast storage system) to remember recent weather data:

How it works:
1. Someone asks for London's weather
2. We check if we have recent data stored
3. If yes → return stored data (fast!)
4. If no → get new data, store it, then return it

Benefits:
- Super fast for repeated requests
- Reduces load on weather providers
- Data expires after a while (stays fresh)
- Works even if Redis is temporarily down

Storage details:
- Key format: `weather:cityname` (e.g., `weather:london`)
- Data format: JSON (easy to read and write)
- Expiry: Default 10 minutes (changeable)

### 4. Web Server

The web server is the front door of our service. It:
- Listens for requests
- Gives responses in JSON format
- Has safety timeouts (won't hang forever)
- Tells you if something goes wrong
- Has a health check endpoint

## How a Request Works

When you ask for weather data:

1. Your request comes in
   ```http
   GET /weather?city=London
   ```

2. Server checks Redis
   - Found? Return it immediately! ✨
   - Not found? Continue to step 3

3. Get fresh data
   - Ask all providers at once
   - Wait for responses (but not too long)
   - Calculate averages
   - Save in Redis for next time
   - Send back to you

4. You get the response
   ```json
   {
     "city": "London",
     "avg_temp_c": 18.5,
     "avg_humidity": 65
   }
   ```

## Error Handling

### Provider Level
- Network errors
- Invalid responses
- Timeout handling
- City not found

### Aggregator Level
- Partial provider failures
- All providers failed
- Context cancellation
- Timeout handling

### Cache Level
- Connection errors
- Serialization errors
- TTL management

### HTTP Level
- Invalid requests
- Timeout handling
- JSON encoding errors
- 4xx/5xx responses

## Performance Considerations

1. Concurrency
   - Goroutines for provider calls
   - Connection pooling for Redis
   - Context for timeout management

2. Caching
   - Configurable TTL
   - JSON serialization
   - Redis persistence

3. Resource Management
   - HTTP client reuse
   - Timeout configurations
   - Memory efficient data structures

## Deployment

Docker-based deployment with:
- Multi-stage builds
- Non-root user
- Health checks
- Resource constraints
- Network isolation
- Environment configuration

## Monitoring & Logging

Structured JSON logging with:
- Timestamp
- Log level
- Event type
- Module
- Additional context
- Error details

## Configuration

Environment variable based configuration for:
- Server settings
- Timeouts
- Cache TTL
- Redis connection
- Provider settings

## Testing

1. Unit Tests
   - Mock provider testing
   - Aggregator logic
   - Cache operations

2. Integration Tests
   - Provider integration
   - Redis integration
   - Full request flow

## Future Enhancements

1. Technical Improvements
   - Circuit breaker for providers
   - Rate limiting
   - Request tracing
   - Metrics collection
   - API versioning

2. Feature Additions
   - More weather providers
   - Extended weather data
   - Historical data
   - Location validation
   - Weather alerts

3. Operational Improvements
   - Prometheus metrics
   - Grafana dashboards
   - Alert rules
   - Documentation automation