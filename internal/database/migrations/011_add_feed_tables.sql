-- +goose Up
-- Create feed_types table
CREATE TABLE feed_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL, -- hay, grain, supplement, mineral, etc.
    unit VARCHAR(20) NOT NULL, -- kg, g, ml, etc.
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create feed_inventory table
CREATE TABLE feed_inventory (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    feed_type_id INTEGER NOT NULL REFERENCES feed_types(id),
    quantity DECIMAL(10,2) NOT NULL,
    unit_price DECIMAL(10,2),
    purchase_date DATE,
    expiry_date DATE,
    supplier VARCHAR(100),
    batch_number VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create feed_requirements table
CREATE TABLE feed_requirements (
    id SERIAL PRIMARY KEY,
    horse_id INTEGER NOT NULL REFERENCES horses(id),
    feed_type_id INTEGER NOT NULL REFERENCES feed_types(id),
    daily_amount DECIMAL(10,2) NOT NULL,
    notes TEXT,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create feed_logs table
CREATE TABLE feed_logs (
    id SERIAL PRIMARY KEY,
    horse_id INTEGER NOT NULL REFERENCES horses(id),
    feed_type_id INTEGER NOT NULL REFERENCES feed_types(id),
    amount DECIMAL(10,2) NOT NULL,
    feeding_time TIMESTAMP WITH TIME ZONE NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add some common feed types
INSERT INTO feed_types (name, category, unit) VALUES
    ('Hay', 'hay', 'kg'),
    ('Grass', 'hay', 'kg'),
    ('Oats', 'grain', 'kg'),
    ('Barley', 'grain', 'kg'),
    ('Salt block', 'mineral', 'kg'),
    ('Mineral supplement', 'mineral', 'g'),
    ('Vitamin supplement', 'supplement', 'g');

-- Add indexes
CREATE INDEX idx_feed_inventory_user ON feed_inventory(user_id);
CREATE INDEX idx_feed_inventory_type ON feed_inventory(feed_type_id);
CREATE INDEX idx_feed_requirements_horse ON feed_requirements(horse_id);
CREATE INDEX idx_feed_logs_horse ON feed_logs(horse_id);
CREATE INDEX idx_feed_logs_time ON feed_logs(feeding_time);

-- +goose Down
DROP TABLE IF EXISTS feed_logs;
DROP TABLE IF EXISTS feed_requirements;
DROP TABLE IF EXISTS feed_inventory;
DROP TABLE IF EXISTS feed_types;
