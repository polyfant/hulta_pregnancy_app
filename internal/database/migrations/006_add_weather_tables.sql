-- +goose Up
-- Add weather-related fields to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS weather_enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS location_sharing BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS default_latitude DOUBLE PRECISION;
ALTER TABLE users ADD COLUMN IF NOT EXISTS default_longitude DOUBLE PRECISION;
ALTER TABLE users ADD COLUMN IF NOT EXISTS weather_notifications_on BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS weather_update_frequency VARCHAR(10) DEFAULT 'hourly';

-- Create locations table
CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    altitude DOUBLE PRECISION,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create weather_data table
CREATE TABLE weather_data (
    id SERIAL PRIMARY KEY,
    location_id INTEGER NOT NULL REFERENCES locations(id),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    temperature DOUBLE PRECISION,
    humidity DOUBLE PRECISION,
    wind_speed DOUBLE PRECISION,
    wind_direction INTEGER,
    precipitation DOUBLE PRECISION,
    pressure DOUBLE PRECISION,
    conditions VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes
CREATE INDEX idx_weather_data_location_id ON weather_data(location_id);
CREATE INDEX idx_weather_data_timestamp ON weather_data(timestamp);
CREATE INDEX idx_locations_coordinates ON locations(latitude, longitude);

-- +goose Down
DROP TABLE IF EXISTS weather_data;
DROP TABLE IF EXISTS locations;
ALTER TABLE users 
    DROP COLUMN IF EXISTS weather_enabled,
    DROP COLUMN IF EXISTS location_sharing,
    DROP COLUMN IF EXISTS default_latitude,
    DROP COLUMN IF EXISTS default_longitude,
    DROP COLUMN IF EXISTS weather_notifications_on,
    DROP COLUMN IF EXISTS weather_update_frequency;
