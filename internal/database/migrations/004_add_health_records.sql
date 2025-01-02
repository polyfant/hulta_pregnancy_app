-- +goose Up

CREATE TABLE health_records (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    horse_id INTEGER NOT NULL,
    record_type VARCHAR(50) NOT NULL,
    record_date TIMESTAMP NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_horse
        FOREIGN KEY(horse_id) 
        REFERENCES horses(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_health_records_user_id ON health_records(user_id);
CREATE INDEX idx_health_records_horse_id ON health_records(horse_id);

-- +goose Down
DROP TABLE IF EXISTS health_records; 