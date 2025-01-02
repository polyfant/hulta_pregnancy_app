-- Create horses table
CREATE TABLE IF NOT EXISTS horses (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    name VARCHAR(100) NOT NULL,
    breed VARCHAR(100),
    birth_date DATE,
    conception_date DATE,
    mother_id INTEGER REFERENCES horses(id),
    father_id INTEGER REFERENCES horses(id),
    is_pregnant BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create expenses table
CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    horse_id INTEGER REFERENCES horses(id),
    amount DECIMAL(10,2) NOT NULL,
    description TEXT,
    date DATE NOT NULL,
    type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create recurring_expenses table
CREATE TABLE IF NOT EXISTS recurring_expenses (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    horse_id INTEGER REFERENCES horses(id),
    amount DECIMAL(10,2) NOT NULL,
    description TEXT,
    frequency VARCHAR(50),
    next_due_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes
CREATE INDEX idx_horses_user_id ON horses(user_id);
CREATE INDEX idx_expenses_user_id ON expenses(user_id);
CREATE INDEX idx_recurring_expenses_user_id ON recurring_expenses(user_id); 