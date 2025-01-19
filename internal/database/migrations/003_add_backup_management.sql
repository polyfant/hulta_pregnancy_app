-- +goose Up
-- Create backup history table
CREATE TABLE backup_history (
    id SERIAL PRIMARY KEY,
    backup_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    backup_path TEXT NOT NULL,
    status VARCHAR(50) NOT NULL,
    retention_days INTEGER DEFAULT 7,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes for backup history
CREATE INDEX idx_backup_history_backup_time ON backup_history(backup_time);
CREATE INDEX idx_backup_history_status ON backup_history(status);

-- +goose Down
DROP TABLE IF EXISTS backup_history;
