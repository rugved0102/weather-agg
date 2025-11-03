# Quick Start Guide

## Option 1: Running with Docker (Simple Way)

```powershell
# 1. Build the application
docker build -t weather-app .

# 2. Start Redis
docker run --name redis -d redis

# 3. Start PostgreSQL
docker run --name postgres -e POSTGRES_PASSWORD=yourpassword -d postgres

# 4. Start your application
docker run -p 8080:8080 \
  -e REDIS_ADDR=redis:6379 \
  -e DB_URL=postgres://postgres:yourpassword@postgres:5432/weather \
  --link redis --link postgres \
  weather-app
```

## Option 2: Running with Kubernetes (What we set up)

```powershell
# 1. Start minikube (your local kubernetes cluster)
minikube start

# 2. Deploy everything at once
kubectl apply -f k8s/

# 3. Forward the port to access the application
kubectl port-forward svc/weather-api 8080:8080 -n weather-app
```

## Main Differences

### Docker (Option 1)
- Simpler to understand
- Good for development
- You manage each piece manually
- Have to connect services manually
- No automatic scaling or recovery

### Kubernetes (Option 2)
- More complex but powerful
- Good for production
- Everything managed together
- Automatic service connection
- Automatic scaling and recovery
- Built-in monitoring
- Can handle multiple copies of your app

## Which One Should You Use?

- **For Development**: Start with Docker (Option 1)
- **For Production**: Use Kubernetes (Option 2)

## Testing Your Application

Both options will let you access your application the same way:

```powershell
# Check if application is healthy
curl http://localhost:8080/healthz

# Get weather for a city
curl http://localhost:8080/weather?city=London

# View monitoring metrics
curl http://localhost:8080/metrics
```

## Need Help?

- For Docker issues: Check `docker ps` and `docker logs`
- For Kubernetes issues: Check `kubectl get pods` and `kubectl logs`
- See detailed guide in `k8s/BEGINNERS_GUIDE.md`