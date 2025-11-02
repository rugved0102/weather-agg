# Build stage
FROM golang:1.21-alpine3.18 AS builder

# Add necessary security updates and CA certificates
RUN apk update && \
    apk add --no-cache ca-certificates && \
    update-ca-certificates

# Create non-root user
RUN adduser -D appuser

# Set working directory
WORKDIR /app

# Copy go.mod first
COPY go.mod ./
RUN go mod download && go mod verify

# Copy the rest of the code
COPY . .

# Build the binary with security flags
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags='-w -s -extldflags=-static' -o weather-api ./cmd/server

# Final stage
FROM alpine:3.18

# Copy only the necessary files from builder
COPY --from=builder /app/weather-api /app/weather-api
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Use non-root user
USER appuser

# Set the binary as entrypoint
ENTRYPOINT ["/app/weather-api"]