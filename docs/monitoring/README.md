# Monitoring Guide - Prometheus & Grafana

This guide explains how we implemented monitoring in our weather service using Prometheus and Grafana.

## Implementation Overview

### 1. Added Metrics Collection

We added five key metrics to track service performance:

```go
// RequestCounter tracks total number of requests by city
weather_requests_total{city="London"} 42

// CacheHits tracks cache effectiveness
weather_cache_hits_total{type="hit"} 35
weather_cache_hits_total{type="miss"} 7

// RequestDuration tracks request latency
weather_request_duration_seconds{city="London"} 0.123

// ProviderSuccessCounter tracks successful provider responses
weather_provider_success_total{provider="openmeteo"} 50

// ProviderErrorCounter tracks provider errors
weather_provider_errors_total{provider="openmeteo"} 2
```

### 2. Code Changes Made

1. **Added Metrics Package** (`internal/metrics/metrics.go`):
   - Defined Prometheus metrics using `prometheus.Counter` and `prometheus.Histogram`
   - Added labels for better data segmentation

2. **Modified Main Server** (`cmd/server/main.go`):
   ```go
   // Added Prometheus endpoint
   mux.Handle("/metrics", promhttp.Handler())

   // Added metric collection in request handler
   metrics.RequestCounter.WithLabelValues(city).Inc()
   metrics.RequestDuration.WithLabelValues(city).Observe(time.Since(start).Seconds())
   ```

3. **Added Docker Services**:
   - Prometheus for metrics collection
   - Grafana for visualization

## How to Use the Metrics

### 1. Accessing Raw Metrics

Visit http://localhost:8080/metrics to see raw metrics data. You'll see something like:
```
# HELP weather_requests_total Total number of weather API requests
weather_requests_total{city="London"} 42
```

### 2. Using Prometheus (http://localhost:9090)

#### Common Queries:
1. Request Rate (most used):
   ```
   rate(weather_requests_total[5m])
   ```
   This shows the rate of requests over the last 5 minutes.

2. Cache Hit Ratio:
   ```
   sum(weather_cache_hits_total{type="hit"}) / sum(weather_cache_hits_total)
   ```

3. Average Response Time:
   ```
   rate(weather_request_duration_seconds_sum[5m]) / rate(weather_request_duration_seconds_count[5m])
   ```

### 3. Using Grafana (http://localhost:3000)

Access credentials:
- Username: admin
- Password: weatherpass

Available dashboards show:
- Request rates by city
- Cache hit/miss ratio
- Average response times
- Provider success rates

## Screenshot Gallery

Screenshots are stored in the `docs/monitoring/screenshots` folder:

1. Raw Metrics: Shows the raw metrics from /metrics endpoint
2. Prometheus Queries: Shows example query results
3. Grafana Dashboard: Shows the complete monitoring dashboard
4. Database Results: Shows PostgreSQL query results from weather_logs table

### Screenshot Naming Convention

```
raw_metrics_YYYYMMDD.png        # Raw Prometheus metrics
prometheus_query_YYYYMMDD.png   # Prometheus query results
grafana_dashboard_YYYYMMDD.png  # Grafana dashboard view
postgres_query_YYYYMMDD.png     # PostgreSQL weather_logs data

Example: postgres_query_20251103.png
```

### Database Query Results

The PostgreSQL results show:
- Actual weather data stored in the database
- Timestamp of data collection
- Temperature and humidity values
- City names

Command used to get database results:
```bash
docker exec -it postgres psql -U weather -d weatherdb
```

Example queries shown in screenshots:
```sql
-- Show timezone setting
SHOW timezone;

-- Show table structure
\d weather_logs

-- Show recent weather data with local time
SELECT 
    city,
    temp_c,
    humidity,
    retrieved_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Kolkata' as local_time
FROM weather_logs
ORDER BY retrieved_at DESC
LIMIT 5;

-- Show statistics
SELECT 
    count(*) as total_records,
    count(distinct city) as unique_cities
FROM weather_logs;
```

This data helps correlate the metrics we see in Prometheus/Grafana with the actual weather data being collected and stored.

## Understanding the Metrics

### Request Rate Analysis
The query `rate(weather_requests_total[5m])` shows:
- Request frequency per city
- Traffic patterns
- Peak usage times

Example interpretation:
```
weather_requests_total{city="London"} 0.2
```
This means we're getting about 0.2 requests per second (or 12 requests per minute) for London.

### Cache Performance
Watch the cache hit ratio to:
- Evaluate cache effectiveness
- Adjust cache TTL if needed
- Identify frequently requested cities

### Response Times
Monitor request duration to:
- Identify slow responses
- Track performance degradation
- Set alerting thresholds

### Provider Health
Track provider metrics to:
- Monitor provider reliability
- Detect failing providers
- Balance provider usage

## Best Practices

1. **Regular Monitoring**
   - Check dashboards daily
   - Review error rates
   - Monitor cache effectiveness

2. **Alert Setup**
   - Set up alerts for high error rates
   - Monitor slow response times
   - Watch for cache misses

3. **Dashboard Usage**
   - Use time range selectors
   - Compare metrics across cities
   - Export data for reporting

## Future Enhancements

1. **Additional Metrics**
   - Database query timings
   - Provider response times
   - Memory usage stats

2. **Alert Rules**
   - High error rate alerts
   - Slow response alerts
   - Provider failure alerts

3. **Custom Dashboards**
   - City-specific views
   - Provider performance
   - Cache analysis

## Troubleshooting

Common issues and solutions:
1. Missing metrics
   - Check /metrics endpoint
   - Verify metric registration
   - Check label values

2. No data in Grafana
   - Check Prometheus connection
   - Verify metrics collection
   - Check time range selection

## Screenshots Guide

Store your screenshots in the `screenshots` folder with these naming conventions:
1. `raw_metrics_YYYYMMDD.png`
2. `prometheus_query_YYYYMMDD.png`
3. `grafana_dashboard_YYYYMMDD.png`

This helps track changes over time and maintain a historical record of your monitoring setup.