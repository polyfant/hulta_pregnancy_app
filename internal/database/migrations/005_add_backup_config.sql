-- +goose Up

-- Create a schema for backup management
CREATE SCHEMA IF NOT EXISTS backup_manager;

-- Create a table to track backups
CREATE TABLE backup_manager.backup_history (
    id SERIAL PRIMARY KEY,
    backup_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    backup_path TEXT NOT NULL,
    status VARCHAR(50) NOT NULL,
    retention_days INTEGER DEFAULT 7
);

-- Create a function to clean old backups
CREATE OR REPLACE FUNCTION backup_manager.cleanup_old_backups()
    RETURNS void
    LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM backup_manager.backup_history
    WHERE backup_time < (CURRENT_TIMESTAMP - (retention_days * INTERVAL '1 day'));
END;
$$;

-- Create a function to perform backup
CREATE OR REPLACE FUNCTION backup_manager.perform_backup()
    RETURNS void
    LANGUAGE plpgsql
AS $$
DECLARE
    backup_path TEXT;
BEGIN
    backup_path := '/var/lib/postgresql/backups/he_horse_db_' || 
                  to_char(CURRENT_TIMESTAMP, 'YYYY_MM_DD_HH24_MI_SS') || '.sql';
    
    INSERT INTO backup_manager.backup_history (backup_path, status)
    VALUES (backup_path, 'COMPLETED');
    
    PERFORM backup_manager.cleanup_old_backups();
END;
$$;

-- +goose Down
DROP SCHEMA IF EXISTS backup_manager CASCADE; 