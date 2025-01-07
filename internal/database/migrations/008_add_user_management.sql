-- Create users table with Auth0 integration
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,  -- Auth0 sub claim
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

-- Update existing tables to use users table as foreign key
ALTER TABLE horses 
    DROP CONSTRAINT IF EXISTS horses_user_id_fkey,
    ADD CONSTRAINT horses_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE expenses 
    DROP CONSTRAINT IF EXISTS expenses_user_id_fkey,
    ADD CONSTRAINT expenses_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE recurring_expenses 
    DROP CONSTRAINT IF EXISTS recurring_expenses_user_id_fkey,
    ADD CONSTRAINT recurring_expenses_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id);

-- Add indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
