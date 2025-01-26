-- +goose Up
-- Create breeding_costs table
CREATE TABLE breeding_costs (
    id SERIAL PRIMARY KEY,
    horse_id INTEGER NOT NULL REFERENCES horses(id),
    amount DECIMAL(10,2) NOT NULL,
    description TEXT,
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes
CREATE INDEX idx_breeding_costs_horse_id ON breeding_costs(horse_id);

-- +goose Down
DROP TABLE IF EXISTS breeding_costs;
