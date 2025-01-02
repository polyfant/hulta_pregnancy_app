-- +goose Up
CREATE TABLE breeding_records (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    mare_id INTEGER NOT NULL,
    stallion_id INTEGER,
    external_stallion_name VARCHAR(255),
    breeding_date TIMESTAMP NOT NULL,
    success_status VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_mare
        FOREIGN KEY(mare_id) 
        REFERENCES horses(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_stallion
        FOREIGN KEY(stallion_id) 
        REFERENCES horses(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_breeding_records_user_id ON breeding_records(user_id);
CREATE INDEX idx_breeding_records_mare_id ON breeding_records(mare_id);
CREATE INDEX idx_breeding_records_stallion_id ON breeding_records(stallion_id);

-- +goose Down
DROP TABLE IF EXISTS breeding_records; 