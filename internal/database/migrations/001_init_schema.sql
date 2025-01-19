-- +goose Up
-- Create users table (core table, no dependencies)
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    profile_picture_url TEXT,
    role VARCHAR(50) DEFAULT 'standard',
    last_login TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create horses table (depends on users)
CREATE TABLE IF NOT EXISTS horses (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    name VARCHAR(100) NOT NULL,
    breed VARCHAR(100),
    birth_date DATE,
    conception_date DATE,
    mother_id INTEGER REFERENCES horses(id) ON DELETE SET NULL,
    father_id INTEGER REFERENCES horses(id) ON DELETE SET NULL,
    is_pregnant BOOLEAN DEFAULT FALSE,
    owner_name VARCHAR(100),
    owner_contact VARCHAR(100),
    owner_email VARCHAR(100),
    owner_phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_horses_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create expense type enum
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

-- Create expenses table (depends on users and horses)
CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    horse_id INTEGER,
    amount DECIMAL(10,2) NOT NULL CHECK (amount >= 0),
    description TEXT,
    date DATE NOT NULL,
    expense_type expense_type DEFAULT 'other',
    payment_method VARCHAR(50),
    receipt_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_expenses_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_expenses_horse FOREIGN KEY (horse_id) REFERENCES horses(id) ON DELETE SET NULL
);

-- Create recurring expenses table (depends on users and horses)
CREATE TABLE IF NOT EXISTS recurring_expenses (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    horse_id INTEGER,
    amount DECIMAL(10,2) NOT NULL CHECK (amount >= 0),
    description TEXT,
    frequency VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    next_due_date DATE NOT NULL,
    expense_type expense_type DEFAULT 'other',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_recurring_expenses_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_recurring_expenses_horse FOREIGN KEY (horse_id) REFERENCES horses(id) ON DELETE SET NULL
);

-- Add core indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_horses_user_id ON horses(user_id);
CREATE INDEX idx_horses_name ON horses(name);
CREATE INDEX idx_horses_owner_name ON horses(owner_name);
CREATE INDEX idx_expenses_user_id ON expenses(user_id);
CREATE INDEX idx_expenses_horse_id ON expenses(horse_id);
CREATE INDEX idx_expenses_date ON expenses(date);
CREATE INDEX idx_expenses_expense_type ON expenses(expense_type);
CREATE INDEX idx_recurring_expenses_user_id ON recurring_expenses(user_id);
CREATE INDEX idx_recurring_expenses_horse_id ON recurring_expenses(horse_id);
CREATE INDEX idx_recurring_expenses_next_due_date ON recurring_expenses(next_due_date);
CREATE INDEX idx_recurring_expenses_expense_type ON recurring_expenses(expense_type);

-- +goose Down
DROP TABLE IF EXISTS recurring_expenses;
DROP TABLE IF EXISTS expenses;
DROP TABLE IF EXISTS horses;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS expense_type;
