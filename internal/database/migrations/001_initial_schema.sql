-- +goose Up
-- Create users table first
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    auth0_id VARCHAR(255) NOT NULL UNIQUE,  -- This will store the Auth0 user ID
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_auth0_id ON users(auth0_id);

-- Modify horses table to reference users
CREATE TABLE horses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    breed VARCHAR(255),
    gender VARCHAR(50),
    birth_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_horses_user_id ON horses(user_id);
CREATE UNIQUE INDEX idx_horses_name_per_user ON horses(user_id, name);  -- Ensure unique horse names per user

CREATE TABLE pregnancies (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    horse_id INTEGER NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_horse
        FOREIGN KEY(horse_id) 
        REFERENCES horses(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_pregnancies_user_id ON pregnancies(user_id);
CREATE INDEX idx_pregnancies_horse_id ON pregnancies(horse_id);

-- +goose Down
DROP DATABASE IF EXISTS HE_horse_db; 