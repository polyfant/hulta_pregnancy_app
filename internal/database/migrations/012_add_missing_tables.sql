-- +goose Up
-- Drop old columns from weather_impacts
ALTER TABLE weather_impacts 
    DROP COLUMN IF EXISTS impact_type,
    DROP COLUMN IF EXISTS severity,
    DROP COLUMN IF EXISTS notes,
    DROP COLUMN IF EXISTS recorded_at;

-- Add new columns to weather_impacts
ALTER TABLE weather_impacts 
    ADD COLUMN IF NOT EXISTS stress_level DECIMAL(5,2) CHECK (stress_level >= 0 AND stress_level <= 100),
    ADD COLUMN IF NOT EXISTS comfort_index DECIMAL(5,2) CHECK (comfort_index >= 0 AND comfort_index <= 100),
    ADD COLUMN IF NOT EXISTS exercise_safe BOOLEAN,
    ADD COLUMN IF NOT EXISTS recommendations TEXT[],
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_weather_impacts_horse_id ON weather_impacts(horse_id);
CREATE INDEX IF NOT EXISTS idx_weather_impacts_weather_data_id ON weather_impacts(weather_data_id);
CREATE INDEX IF NOT EXISTS idx_weather_impacts_stress_level ON weather_impacts(stress_level);

-- +goose Down
-- Restore old columns
ALTER TABLE weather_impacts 
    ADD COLUMN IF NOT EXISTS impact_type VARCHAR(50),
    ADD COLUMN IF NOT EXISTS severity INTEGER CHECK (severity BETWEEN 1 AND 5),
    ADD COLUMN IF NOT EXISTS notes TEXT,
    ADD COLUMN IF NOT EXISTS recorded_at TIMESTAMP WITH TIME ZONE;

-- Drop new columns
ALTER TABLE weather_impacts 
    DROP COLUMN IF EXISTS stress_level,
    DROP COLUMN IF EXISTS comfort_index,
    DROP COLUMN IF EXISTS exercise_safe,
    DROP COLUMN IF EXISTS recommendations,
    DROP COLUMN IF EXISTS updated_at;

DROP INDEX IF EXISTS idx_weather_impacts_horse_id;
DROP INDEX IF EXISTS idx_weather_impacts_weather_data_id;
DROP INDEX IF EXISTS idx_weather_impacts_stress_level;
