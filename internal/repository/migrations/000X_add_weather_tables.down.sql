-- Drop indexes
DROP INDEX IF EXISTS idx_weather_impacts_created_at;
DROP INDEX IF EXISTS idx_weather_impacts_horse;
DROP INDEX IF EXISTS idx_weather_data_location_timestamp;

-- Drop tables
DROP TABLE IF EXISTS weather_impacts;
DROP TABLE IF EXISTS weather_data;
DROP TABLE IF EXISTS locations;

-- Remove weather-related columns from users
ALTER TABLE users DROP COLUMN IF EXISTS weather_update_frequency;
ALTER TABLE users DROP COLUMN IF EXISTS weather_notifications_on;
ALTER TABLE users DROP COLUMN IF EXISTS default_longitude;
ALTER TABLE users DROP COLUMN IF EXISTS default_latitude;
ALTER TABLE users DROP COLUMN IF EXISTS location_sharing;
ALTER TABLE users DROP COLUMN IF EXISTS weather_enabled;
