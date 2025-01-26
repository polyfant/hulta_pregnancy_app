-- Add weather-related fields to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS weather_enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS location_sharing BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS default_latitude DOUBLE PRECISION;
ALTER TABLE users ADD COLUMN IF NOT EXISTS default_longitude DOUBLE PRECISION;
ALTER TABLE users ADD COLUMN IF NOT EXISTS weather_notifications_on BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS weather_update_frequency VARCHAR(10) DEFAULT 'hourly';

-- Create locations table
CREATE TABLE IF NOT EXISTS locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    altitude DOUBLE PRECISION,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create weather_data table
CREATE TABLE IF NOT EXISTS weather_data (
    id SERIAL PRIMARY KEY,
    location_id INTEGER REFERENCES locations(id),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    temperature DOUBLE PRECISION NOT NULL,
    humidity DOUBLE PRECISION NOT NULL,
    wind_speed DOUBLE PRECISION NOT NULL,
    pressure DOUBLE PRECISION NOT NULL,
    rain_amount DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create weather_impacts table
CREATE TABLE IF NOT EXISTS weather_impacts (
    id SERIAL PRIMARY KEY,
    horse_id INTEGER REFERENCES horses(id),
    weather_data_id INTEGER REFERENCES weather_data(id),
    stress_level DOUBLE PRECISION NOT NULL,
    comfort_index DOUBLE PRECISION NOT NULL,
    exercise_safe BOOLEAN NOT NULL,
    recommendations JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_weather_data_location_timestamp ON weather_data(location_id, timestamp);
CREATE INDEX IF NOT EXISTS idx_weather_impacts_horse ON weather_impacts(horse_id);
CREATE INDEX IF NOT EXISTS idx_weather_impacts_created_at ON weather_impacts(created_at);
