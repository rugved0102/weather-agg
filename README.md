# Weather Aggregation Service

A simple weather service that combines temperature and humidity data from multiple sources to provide averaged weather information. Built with Go, it uses Redis for caching and Docker for easy deployment.

## What Does It Do?

1. Gets weather data from:
   - Open-Meteo API (real weather data, no API key needed)
   - Backup mock provider (for testing/fallback)
2. Combines and averages the data
3. Caches results in Redis for faster responses
4. Provides simple REST API endpoints

## Getting Started

### What You Need
- Docker and Docker Compose
  - That's it! Docker handles everything else
- (Optional) Go 1.21+ if you want to run without Docker

### Quick Start (Using Docker)
```bash
# 1. Get the code
git clone https://github.com/rugved0102/weather-agg.git
cd weather-agg

# 2. Start everything
docker compose up --build

# 3. Try it out
curl "http://localhost:8080/weather?city=London"
```

### Developer Setup (Without Docker)
```bash
# 1. Get the code
git clone https://github.com/rugved0102/weather-agg.git
cd weather-agg

# 2. Start Redis
docker run -d --name redis -p 6379:6379 redis:alpine

# 3. Run the service
go run ./cmd/server
```

## API Guide

### Getting Weather Data
```http
GET /weather?city=CityName
```

This gets weather data for any city. For example:
```bash
# Using curl
curl "http://localhost:8080/weather?city=London"

# Using PowerShell
Invoke-WebRequest "http://localhost:8080/weather?city=London"
```

You'll get back something like this:
```json
{
    "city": "London",
    "avg_temp_c": 18.5,           # Temperature in Celsius
    "avg_humidity": 65,           # Humidity percentage
    "providers": [                # Where the data came from
        "openmeteo",
        "mock-A"
    ],
    "provider_count": 2,          # How many providers responded
    "retrieved_at": "2025-11-02T12:34:56Z"  # When the data was fetched
}
```

### Checking If Service Is Working
```http
GET /healthz
```

Returns "ok" if everything is working. Use this to monitor service health.

## Configuration

You can customize the service using these settings (all are optional):

```env
# Where to find Redis
REDIS_ADDR=redis:6379   # Format: hostname:port

# Web server settings
PORT=8080              # What port to run on
CACHE_TTL=10m         # How long to cache results (e.g., 10m = 10 minutes)

# Timeouts (using Go duration format: 5s = 5 seconds, 1m = 1 minute)
SERVER_READ_TIMEOUT=5s    # How long to wait for client requests
SERVER_WRITE_TIMEOUT=10s  # How long to wait when sending responses
SERVER_IDLE_TIMEOUT=60s   # How long to keep connections alive
PROVIDER_TIMEOUT=6s       # How long to wait for weather providers
REQUEST_TIMEOUT=8s        # Total time limit for each request
```

To use these:
1. Create a `.env` file with your settings, or
2. Set them in your terminal before running the service

## Project Layout

Here's how the code is organized:

```
weather-agg/
├── cmd/server/              # Main application
│   └── main.go             # Entry point, sets everything up
│
├── internal/               # Core functionality
│   ├── provider/          # Weather data providers
│   │   ├── provider.go    # Interface for all providers
│   │   ├── openmeteo.go   # Real weather data provider
│   │   └── mockprovider.go# Test/backup provider
│   │
│   ├── agg/              # Data processing
│   │   └── aggregator.go # Combines weather data
│   │
│   └── cache/            # Data storage
│       ├── convert.go    # Data conversion helpers
│       └── redis.go      # Redis storage implementation
│
├── pkg/                   # Shared utilities
│   └── logger/           # Logging tools
│
└── docker-compose.yml    # Docker configuration
```

Each folder has a specific job:
- `cmd/server`: The main program
- `internal/provider`: Gets weather data
- `internal/agg`: Processes the data
- `internal/cache`: Stores the data
- `pkg/logger`: Helps with logging

## Development

### Building

```bash
# Build the binary
go build ./cmd/server

# Run tests
go test ./...
```

### Docker Build

```bash
docker build -t weather-api .
```

## Security

- Runs as non-root user in Docker
- Proper timeout configurations
- Multi-stage Docker builds
- Redis only accessible within Docker network

## License

[Add your chosen license here]