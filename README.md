# Weather Aggregation Service

A Go-based microservice that aggregates weather information from multiple sources (currently using mock data).

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Docker (optional)

### Running Locally

```bash
# Clone the repository
git clone https://github.com/rugved0102/weather-agg.git
cd weather-agg

# Run the server
go run cmd/server/main.go
```

### Using Docker

```bash
# Build the Docker image
docker build -t weather-api .

# Run the container
docker run -p 8080:8080 weather-api
```

## API Documentation

### GET /weather

Retrieves weather information for a specified city.

**Parameters:**
- `city` (query parameter, optional) - Name of the city to get weather for
  - If not provided, defaults to "Unknown"

**Response:**
```json
{
    "city": "string",
    "temp_c": number,
    "humidity": number,
    "retrieved_at": "string (ISO 8601)"
}
```

**Example Request:**
```bash
curl "http://localhost:8080/weather?city=London"
```

**Example Response:**
```json
{
    "city": "London",
    "temp_c": 27.3,
    "humidity": 62,
    "retrieved_at": "2025-11-02T12:34:56Z"
}
```

## Project Structure

```
weather-agg/
├── cmd/
│   └── server/
│       └── main.go    # Main application entry point
├── Dockerfile         # Container configuration
├── README.md         # Project documentation
└── go.mod            # Go module definition
```

## Configuration

The service runs with the following default configuration:
- Port: 8080
- Read Timeout: 5 seconds
- Write Timeout: 10 seconds
- Idle Timeout: 60 seconds
- Request Timeout: 8 seconds

## Development

### Building

```bash
# Build the binary
go build ./cmd/server

# Run tests (when added)
go test ./...
```

### Docker Build

```bash
docker build -t weather-api .
```

## Future Enhancements

- [ ] Integration with real weather APIs
- [ ] Caching layer for weather data
- [ ] Rate limiting
- [ ] Metrics and monitoring
- [ ] Authentication
- [ ] Environment-based configuration
- [ ] Unit and integration tests

## Security

The service includes several security measures:
- Runs as non-root user in Docker
- Includes proper timeout configurations
- Uses multi-stage Docker builds
- Includes SSL certificates for HTTPS support

## License

[Add your chosen license here]