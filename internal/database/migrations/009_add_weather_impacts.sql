-- +goose Up
CREATE TABLE weather_impacts (
    id SERIAL PRIMARY KEY,
    horse_id INTEGER NOT NULL REFERENCES horses(id),
    weather_data_id INTEGER NOT NULL REFERENCES weather_data(id),
    impact_type VARCHAR(50) NOT NULL,
    severity INTEGER CHECK (severity BETWEEN 1 AND 5),
    notes TEXT,
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes
CREATE INDEX idx_weather_impacts_horse ON weather_impacts(horse_id);
CREATE INDEX idx_weather_impacts_weather ON weather_impacts(weather_data_id);
CREATE INDEX idx_weather_impacts_recorded ON weather_impacts(recorded_at);

-- +goose Down
DROP TABLE IF EXISTS weather_impacts;
