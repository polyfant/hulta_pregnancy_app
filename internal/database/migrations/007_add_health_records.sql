-- +goose Up
CREATE TABLE health_records (
    id SERIAL PRIMARY KEY,
    horse_id INTEGER REFERENCES horses(id),
    type VARCHAR(50) NOT NULL,
    description TEXT,
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add index
CREATE INDEX idx_health_records_horse_id ON health_records(horse_id);

-- +goose Down
DROP TABLE IF EXISTS health_records;
