# Weather Aggregation Service

A production-ready weather service that combines temperature and humidity data from multiple sources to provide averaged weather information. Built with Go, uses Redis for caching, PostgreSQL for data persistence, and runs on Kubernetes for scalability and reliability.

## System Architecture

```
                                    ┌─────────────────┐
                                    │                 │
                          ┌────────►│  Open-Meteo API │
                          │         │                 │
                          │         └─────────────────┘
                          │
┌──────────────┐    ┌─────────────┐    ┌──────────────┐
│              │    │             │    │              │
│   Client    ─────►│ Weather API ├───►│  Mock API    │
│              │    │             │    │              │
└──────────────┘    └─────────────┘    └──────────────┘
                          │
                          │         ┌─────────────┐
                          │         │             │
                          ├────────►│    Redis    │
                          │         │   (Cache)   │
                          │         │             │
                          │         └─────────────┘
                          │
                          │         ┌─────────────┐
                          │         │             │
                          └────────►│ PostgreSQL  │
                                    │    (DB)     │
                                    │             │
                                    └─────────────┘
```

### Components:

1. **Weather API (Core Service)**
   - Written in Go 1.21
   - Aggregates weather data from multiple providers
   - Implements circuit breakers and timeouts
   - Exposes REST API endpoints
   - Prometheus metrics for monitoring

2. **Redis Cache**
   - Caches weather data
   - Reduces load on external APIs
   - Improves response times
   - Configurable TTL

3. **PostgreSQL Database**
   - Stores historical weather data
   - Maintains user preferences
   - Enables data analysis

4. **Weather Providers**
   - Open-Meteo API (primary provider)
   - Mock Provider (backup/testing)

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

## Deployment Options

### 1. Local Development (Docker)
```bash
# Start everything with Docker Compose
docker compose up --build

# Access the API
curl "http://localhost:8080/weather?city=London"
```

### 2. Production Deployment (Kubernetes)
```bash
# 1. Start Kubernetes cluster
minikube start

# 2. Deploy the application stack
kubectl apply -f k8s/

# 3. Forward the port
kubectl port-forward svc/weather-api 8080:8080 -n weather-app
```

See [QUICKSTART.md](QUICKSTART.md) for detailed deployment instructions.

## Monitoring & Observability

### Health Check
```bash
# Check service health
curl http://localhost:8080/healthz
```

### Metrics
```bash
# View Prometheus metrics
curl http://localhost:8080/metrics
```

Available metrics include:
- Go runtime metrics (memory, goroutines, GC stats)
- Process metrics (CPU, memory, file descriptors)
- HTTP handler metrics
- Custom business metrics

## Configuration

You can customize the service using these settings:

```env
# Redis Configuration
REDIS_ADDR=redis-master:6379

# Database Configuration
DB_URL=postgres://postgres:password@postgres-postgresql:5432/weather

# Web Server Settings
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

## Project Structure

```
weather-agg/
├── cmd/server/           # Application entrypoint
├── internal/            # Core business logic
│   ├── provider/       # Weather data providers
│   ├── agg/           # Data aggregation
│   └── cache/         # Caching logic
├── k8s/               # Kubernetes manifests
│   ├── configmap.yaml    # Configuration
│   ├── secrets.yaml      # Sensitive data
│   ├── deployment.yaml   # Pod configuration
│   ├── service.yaml      # Network routing
│   ├── hpa.yaml          # Auto-scaling
│   └── servicemonitor.yaml # Prometheus integration
├── pkg/                # Shared utilities
└── docker-compose.yml  # Local development
```

Each folder has a specific job:
- `cmd/server`: The main program
- `internal/provider`: Gets weather data
- `internal/agg`: Processes the data
- `internal/cache`: Stores the data
- `pkg/logger`: Helps with logging

## Development

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Kubernetes tools (for production deployment):
  - minikube
  - kubectl
  - helm

### Building
```bash
# Build binary
go build ./cmd/server

# Run tests
go test ./...

# Build Docker image
docker build -t weather-api .
```

### Kubernetes Deployment
See [k8s/BEGINNERS_GUIDE.md](k8s/BEGINNERS_GUIDE.md) for detailed Kubernetes deployment instructions.

## Security Features

- Non-root container execution
- Kubernetes secrets for sensitive data
- Resource limits and quotas
- Network policies
- Proper timeout configurations
- Multi-stage Docker builds
- Regular security updates
- Prometheus metrics for monitoring

## Documentation
- [QUICKSTART.md](QUICKSTART.md): Getting started guide
- [k8s/BEGINNERS_GUIDE.md](k8s/BEGINNERS_GUIDE.md): Kubernetes deployment guide

## License

[Add your chosen license here]