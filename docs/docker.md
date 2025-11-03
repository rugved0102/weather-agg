# Docker Guide for Weather Aggregation Service

This guide explains Docker concepts and commands used in this project, specifically designed for beginners.

## Basic Concepts

### What is Docker?
- A platform for building, running, and shipping applications consistently
- Uses containers to package applications and their dependencies
- Ensures the app works the same everywhere

### Key Terms
- **Container**: A running instance of an application and its environment
- **Image**: A template for creating containers (like a class in programming)
- **Dockerfile**: Instructions for building an image
- **docker-compose.yml**: Configuration for running multiple containers together

## Common Commands

### Starting and Stopping Services

```bash
# Start all services in the background
docker compose up -d

# Start with rebuilding images (after code changes)
docker compose up --build

# Stop all services
docker compose down

# Stop and remove volumes (clean slate)
docker compose down -v
```

When to use:
- `docker compose up -d`: Daily development, when code hasn't changed
- `docker compose up --build`: After changing code or Dockerfile
- `docker compose down`: End of day cleanup
- `docker compose down -v`: When you need to reset all data

### Viewing Logs and Status

```bash
# View logs of all services
docker compose logs

# Follow logs in real-time
docker compose logs -f

# View logs of a specific service
docker compose logs weather-api
docker compose logs postgres
docker compose logs redis

# List running containers
docker ps

# List all containers (including stopped)
docker ps -a
```

When to use:
- `docker compose logs`: Debugging issues
- `docker compose logs -f`: Monitoring in real-time
- `docker ps`: Checking what's running

### Working with Containers

```bash
# Execute command in a running container
docker exec -it postgres psql -U weather -d weatherdb

# View container logs
docker logs postgres

# Restart a single service
docker compose restart weather-api

# Stop a single service
docker compose stop redis
```

When to use:
- `docker exec`: Accessing database, debugging inside containers
- `docker logs`: Troubleshooting specific containers
- `docker compose restart`: When a service is misbehaving

### Managing Data and Volumes

```bash
# List volumes
docker volume ls

# Remove a specific volume
docker volume rm weather-agg_pgdata

# Clean up unused volumes
docker volume prune
```

When to use:
- Before major database changes
- When you need to start fresh
- When troubleshooting data issues

## Our Project Structure

### Services Overview
1. **weather-api** (Main Service)
   - Built from our Go code
   - Handles HTTP requests
   - Aggregates weather data

2. **postgres** (Database)
   - Stores historical weather data
   - Persists data in a volume
   - Used for analytics

3. **redis** (Cache)
   - Caches weather data
   - Improves response time
   - In-memory storage

### Common Workflows

#### 1. Regular Development
```bash
# Start everything
docker compose up -d

# Make code changes...

# Rebuild and restart
docker compose up --build

# Check logs for errors
docker compose logs -f
```

#### 2. Database Operations
```bash
# Connect to PostgreSQL
docker exec -it postgres psql -U weather -d weatherdb

# View Redis data
docker exec -it redis redis-cli

# Backup PostgreSQL data
docker exec -t postgres pg_dump -U weather weatherdb > backup.sql
```

#### 3. Troubleshooting
```bash
# Check service status
docker ps

# View specific service logs
docker compose logs weather-api

# Restart problematic service
docker compose restart weather-api

# Full reset if needed
docker compose down -v
docker compose up --build
```

## Common Issues and Solutions

### 1. Container Won't Start
Check logs:
```bash
docker compose logs [service-name]
```

Common causes:
- Port conflicts
- Missing environment variables
- Database initialization errors

### 2. Database Connection Issues
```bash
# Check if container is running
docker ps | grep postgres

# Check postgres logs
docker compose logs postgres

# Try connecting manually
docker exec -it postgres psql -U weather -d weatherdb
```

### 3. Changes Not Reflecting
```bash
# Rebuild images
docker compose up --build

# If still not working
docker compose down
docker compose up --build
```

### 4. Volume Issues
```bash
# List volumes
docker volume ls

# Remove volume and start fresh
docker compose down -v
docker compose up --build
```

## Best Practices

1. **Always Use docker-compose.yml**
   - Keeps service configuration in one place
   - Makes it easy to manage multiple services
   - Version control friendly

2. **Use Volumes for Persistence**
   - Store important data in volumes
   - Our project uses `pgdata` for PostgreSQL

3. **Environment Variables**
   - Keep sensitive data in .env file
   - Use environment variables for configuration
   - Don't commit .env files to git

4. **Regular Cleanup**
   ```bash
   # Remove unused containers
   docker container prune

   # Remove unused images
   docker image prune

   # Remove unused volumes
   docker volume prune
   ```

5. **Monitor Resource Usage**
   ```bash
   # View container resource usage
   docker stats
   ```

## Security Considerations

1. **Never expose database ports** in production
2. **Use non-root users** in containers
3. **Keep images updated** for security patches
4. **Don't store secrets** in images
5. **Use specific versions** instead of 'latest' tag

## Dockerfile Tips

Our project's Dockerfile uses multi-stage builds:
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
# ... building code ...

# Final stage
FROM alpine:3.18
# ... copying only necessary files ...
```

Benefits:
- Smaller final image
- Doesn't include build tools
- More secure

## Getting Help

1. Check container logs first
2. Look for error messages
3. Check port conflicts
4. Verify environment variables
5. Try a clean rebuild

Common commands for help:
```bash
docker --help
docker compose --help
docker [command] --help
```