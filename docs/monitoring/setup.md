# Monitoring Setup Guide

This guide explains how to set up and test monitoring for the Weather Aggregator service using Prometheus and Grafana.

## Prerequisites

- Kubernetes cluster running
- `kubectl` configured
- Helm installed
- Weather API service deployed

## Setup Instructions

### 1. Install Prometheus and Grafana Stack

```bash
# Create monitoring namespace
kubectl create namespace monitoring

# Add Prometheus Helm repository
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Install Prometheus stack (includes Grafana)
helm install monitoring prometheus-community/kube-prometheus-stack --namespace monitoring
```

### 2. Configure ServiceMonitor

Apply the ServiceMonitor configuration to enable Prometheus to scrape metrics from the weather-api:

```bash
kubectl apply -f k8s/servicemonitor.yaml
```

### 3. Access Monitoring Dashboards

#### Prometheus:
```bash
kubectl port-forward svc/monitoring-kube-prometheus-prometheus -n monitoring 9090:9090
```
Access at: http://localhost:9090

#### Grafana:
```bash
kubectl port-forward svc/monitoring-grafana -n monitoring 3001:80
```
Access at: http://localhost:3001

Grafana credentials:
- Username: `admin`
- Password: Get using:
  ```bash
  kubectl get secret monitoring-grafana -n monitoring -o jsonpath="{.data.admin-password}" | base64 -d; echo
  ```

## Available Metrics

The Weather API exposes the following custom metrics:

1. `weather_requests_total{city}`
   - Type: Counter
   - Description: Total number of weather API requests by city

2. `weather_cache_hits_total{type}`
   - Type: Counter
   - Description: Cache hits and misses
   - Labels: type="hit" or type="miss"

3. `weather_request_duration_seconds{city}`
   - Type: Histogram
   - Description: Request duration in seconds by city
   - Buckets: 0.1, 0.25, 0.5, 1, 2.5, 5, 10

4. `weather_provider_success_total{provider}`
   - Type: Counter
   - Description: Successful requests by provider

## Testing the Setup

### 1. Verify Metrics Exposition

```bash
# Port-forward the weather-api service
kubectl port-forward svc/weather-api -n weather-app 8080:8080

# In another terminal, check metrics endpoint
curl http://localhost:8080/metrics
```

### 2. Verify Prometheus Target Discovery

1. Access Prometheus UI at http://localhost:9090
2. Go to Status -> Targets
3. Look for the `serviceMonitor/weather-app/weather-api` target
4. Status should be "UP"

### 3. Test Metric Collection

1. Make some API requests to generate metrics:
```bash
curl http://localhost:8080/weather?city=London
curl http://localhost:8080/weather?city=Pune
```

2. Check metrics in Prometheus:
   - Go to http://localhost:9090
   - Enter queries like:
     ```
     weather_requests_total
     weather_cache_hits_total
     rate(weather_request_duration_seconds_count[5m])
     ```

### 4. Grafana Dashboard Setup

1. Log into Grafana
2. Click "+" -> "Import" to create a new dashboard
3. Add the following panels:

   a. Request Rate Panel:
   ```
   rate(weather_requests_total[5m])
   ```

   b. Cache Performance Panel:
   ```
   weather_cache_hits_total
   ```

   c. Response Time Panel:
   ```
   rate(weather_request_duration_seconds_sum[5m]) / rate(weather_request_duration_seconds_count[5m])
   ```

   d. Provider Success Panel:
   ```
   weather_provider_success_total
   ```

## Troubleshooting

1. If metrics are not showing up:
   - Check if ServiceMonitor is properly labeled
   - Verify weather-api pods have Prometheus annotations
   - Check Prometheus target status

2. If Grafana shows no data:
   - Verify Prometheus data source is configured
   - Check if metrics exist in Prometheus first
   - Verify time range selection

3. Common issues:
   - Port conflicts in port-forwarding: Use different local ports
   - Service discovery issues: Check ServiceMonitor namespace and labels
   - Authentication issues: Reset Grafana admin password if needed