-- +goose Up
CREATE TABLE pregnancy_events (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    pregnancy_id INTEGER NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    event_date TIMESTAMP NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_pregnancy
        FOREIGN KEY(pregnancy_id) 
        REFERENCES pregnancies(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_pregnancy_events_user_id ON pregnancy_events(user_id);
CREATE INDEX idx_pregnancy_events_pregnancy_id ON pregnancy_events(pregnancy_id);

-- +goose Down
DROP TABLE IF EXISTS pregnancy_events; 