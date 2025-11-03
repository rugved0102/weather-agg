-- Initialize timezone settings
ALTER SYSTEM SET timezone TO 'UTC';
ALTER DATABASE weatherdb SET timezone TO 'UTC';

-- Create the weather_logs table with proper timezone handling
CREATE TABLE IF NOT EXISTS weather_logs (
    id SERIAL PRIMARY KEY,
    city TEXT NOT NULL,
    temp_c REAL NOT NULL,
    humidity REAL NOT NULL,
    retrieved_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on the timestamp for better query performance
CREATE INDEX IF NOT EXISTS idx_weather_logs_retrieved_at ON weather_logs(retrieved_at);

-- Create a view for easier local time querying
CREATE OR REPLACE VIEW weather_logs_local AS
SELECT 
    id,
    city,
    temp_c,
    humidity,
    retrieved_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Kolkata' as local_time,
    retrieved_at as utc_time
FROM weather_logs;