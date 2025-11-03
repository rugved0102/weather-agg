# Database Guide - Weather Analytics

This guide explains how to access and analyze the weather data stored in PostgreSQL.

## Database Details

The service stores all weather queries in PostgreSQL for analytics purposes:

```
Host: localhost
Port: 5432
Database: weatherdb
Username: weather
Password: weatherpass
```

## Connecting with DBeaver

1. Open DBeaver
2. Click "New Database Connection" (plug icon with '+')
3. Choose "PostgreSQL"
4. Enter connection details:
   - Host: localhost
   - Port: 5432
   - Database: weatherdb
   - Username: weather
   - Password: weatherpass
5. Click "Test Connection" to verify
6. Click "Finish"

## Database Schema

The weather data is stored in the `weather_logs` table:

```sql
CREATE TABLE weather_logs (
    id SERIAL PRIMARY KEY,
    city TEXT,
    temp_c REAL,
    humidity REAL,
    retrieved_at TIMESTAMP
);
```

## Useful Queries

### Latest Weather Data
```sql
-- Get the most recent 10 weather records
SELECT city, temp_c, humidity, retrieved_at
FROM weather_logs
ORDER BY retrieved_at DESC
LIMIT 10;
```

### Average Temperature by City
```sql
-- Get average temperature for each city
SELECT 
    city,
    ROUND(AVG(temp_c)::numeric, 2) as avg_temp,
    COUNT(*) as measurements
FROM weather_logs
GROUP BY city
ORDER BY measurements DESC;
```

### Temperature Trends
```sql
-- Get temperature trends over time for a specific city
SELECT 
    date_trunc('hour', retrieved_at) as hour,
    city,
    ROUND(AVG(temp_c)::numeric, 2) as avg_temp,
    ROUND(AVG(humidity)::numeric, 2) as avg_humidity
FROM weather_logs
WHERE city = 'London'  -- Replace with your city
GROUP BY hour, city
ORDER BY hour DESC;
```

### Daily Statistics
```sql
-- Get daily min, max, and average temperatures
SELECT 
    city,
    date_trunc('day', retrieved_at) as date,
    MIN(temp_c) as min_temp,
    MAX(temp_c) as max_temp,
    ROUND(AVG(temp_c)::numeric, 2) as avg_temp,
    ROUND(AVG(humidity)::numeric, 2) as avg_humidity,
    COUNT(*) as measurements
FROM weather_logs
GROUP BY city, date
ORDER BY date DESC, city;
```

### Finding Temperature Extremes
```sql
-- Find highest and lowest recorded temperatures
SELECT 
    city,
    ROUND(MIN(temp_c)::numeric, 2) as lowest_temp,
    ROUND(MAX(temp_c)::numeric, 2) as highest_temp,
    MIN(retrieved_at) as first_record,
    MAX(retrieved_at) as last_record,
    COUNT(*) as total_records
FROM weather_logs
GROUP BY city
ORDER BY total_records DESC;
```

## Maintenance Queries

### Data Cleanup
```sql
-- Remove duplicate entries within the same minute
DELETE FROM weather_logs
WHERE id IN (
    SELECT id
    FROM (
        SELECT id,
            ROW_NUMBER() OVER (
                PARTITION BY city, date_trunc('minute', retrieved_at)
                ORDER BY retrieved_at DESC
            ) as rn
        FROM weather_logs
    ) t
    WHERE t.rn > 1
);
```

### Table Size
```sql
-- Check how much space the table is using
SELECT pg_size_pretty(pg_total_relation_size('weather_logs'));
```

## Using Docker CLI

If you prefer using the command line, you can also access the database through Docker:

```powershell
# Connect to PostgreSQL and run query
docker exec -it postgres psql -U weather -d weatherdb -c "SELECT * FROM weather_logs ORDER BY retrieved_at DESC LIMIT 5;"
```

## Backup and Restore

### Create Backup
```powershell
# Backup the database
docker exec -t postgres pg_dump -U weather weatherdb > weather_backup.sql
```

### Restore from Backup
```powershell
# Restore the database
docker exec -i postgres psql -U weather -d weatherdb < weather_backup.sql
```

## Data Retention

Currently, all weather data is retained indefinitely. If you need to implement data retention policies, you can use queries like:

```sql
-- Delete data older than 30 days
DELETE FROM weather_logs
WHERE retrieved_at < NOW() - INTERVAL '30 days';
```

## Monitoring

To monitor database performance and storage:

```sql
-- Check table statistics
SELECT 
    relname as table_name,
    n_live_tup as row_count,
    pg_size_pretty(pg_total_relation_size(relid)) as total_size
FROM pg_stat_user_tables
WHERE relname = 'weather_logs';
```

## Troubleshooting

1. If connection fails:
   - Ensure Docker containers are running (`docker ps`)
   - Check if PostgreSQL is ready (`docker logs postgres`)
   - Verify port mapping (`docker port postgres`)

2. If data is missing:
   - Check the application logs (`docker logs weather-api`)
   - Verify database connectivity in the application
   - Check for any error messages in the API response