-- +goose Up
-- Create expense type enum if not exists
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'expense_type') THEN
        CREATE TYPE expense_type AS ENUM (
            'feed',
            'veterinary',
            'farrier',
            'equipment',
            'training',
            'competition',
            'transport',
            'insurance',
            'boarding',
            'other'
        );
    END IF;
END$$;

-- Create tables if they don't exist
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS horses (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    breed VARCHAR(255),
    date_of_birth DATE,
    gender VARCHAR(50),
    color VARCHAR(100),
    height DECIMAL(5,2),
    weight DECIMAL(6,2),
    registration_number VARCHAR(100),
    microchip_number VARCHAR(100),
    passport_number VARCHAR(100),
    insurance_number VARCHAR(100),
    insurance_company VARCHAR(255),
    insurance_expiry DATE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    horse_id INTEGER REFERENCES horses(id),
    type expense_type NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    date DATE NOT NULL,
    description TEXT,
    receipt_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS recurring_expenses (
    id SERIAL PRIMARY KEY,
    horse_id INTEGER REFERENCES horses(id),
    type expense_type NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    frequency VARCHAR(50) NOT NULL, -- monthly, yearly, etc.
    start_date DATE NOT NULL,
    end_date DATE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes if they don't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE c.relname = 'idx_horses_user_id') THEN
        CREATE INDEX idx_horses_user_id ON horses(user_id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE c.relname = 'idx_expenses_horse_id') THEN
        CREATE INDEX idx_expenses_horse_id ON expenses(horse_id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE c.relname = 'idx_recurring_expenses_horse_id') THEN
        CREATE INDEX idx_recurring_expenses_horse_id ON recurring_expenses(horse_id);
    END IF;
END$$;

-- +goose Down
DROP TABLE IF EXISTS recurring_expenses;
DROP TABLE IF EXISTS expenses;
DROP TABLE IF EXISTS horses;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS expense_type;
