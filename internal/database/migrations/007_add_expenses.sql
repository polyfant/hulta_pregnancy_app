-- +goose Up
-- Create expense tables with constraints and expense types
CREATE TYPE expense_type AS ENUM ('FEED', 'VET', 'FARRIER', 'EQUIPMENT', 'TRAINING', 'COMPETITION', 'OTHER');
CREATE TYPE expense_frequency AS ENUM ('DAILY', 'WEEKLY', 'MONTHLY', 'QUARTERLY', 'YEARLY');

CREATE TABLE expenses (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    horse_id INTEGER REFERENCES horses(id) ON DELETE SET NULL,
    type expense_type NOT NULL,
    amount DECIMAL(10,2) NOT NULL CHECK (amount >= 0),
    date TIMESTAMP NOT NULL,
    description TEXT,
    receipt TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(auth0_id)
);

CREATE TABLE recurring_expenses (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    horse_id INTEGER REFERENCES horses(id) ON DELETE SET NULL,
    type expense_type NOT NULL,
    amount DECIMAL(10,2) NOT NULL CHECK (amount >= 0),
    frequency expense_frequency NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(auth0_id),
    CONSTRAINT check_dates CHECK (end_date IS NULL OR end_date > start_date)
);

-- Create indexes for better query performance
CREATE INDEX idx_expenses_user_id ON expenses(user_id);
CREATE INDEX idx_expenses_horse_id ON expenses(horse_id);
CREATE INDEX idx_expenses_date ON expenses(date);
CREATE INDEX idx_expenses_type ON expenses(type);
CREATE INDEX idx_recurring_expenses_user_id ON recurring_expenses(user_id);
CREATE INDEX idx_recurring_expenses_horse_id ON recurring_expenses(horse_id);
CREATE INDEX idx_recurring_expenses_type ON recurring_expenses(type);

-- +goose Down
DROP TABLE IF EXISTS recurring_expenses;
DROP TABLE IF EXISTS expenses;
DROP TYPE IF EXISTS expense_frequency;
DROP TYPE IF EXISTS expense_type; 