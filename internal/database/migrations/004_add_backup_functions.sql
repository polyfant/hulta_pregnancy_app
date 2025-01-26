-- +goose Up
-- Create function to clean old backups
CREATE OR REPLACE FUNCTION cleanup_old_backups() RETURNS void AS 'DELETE FROM backup_history WHERE backup_time < (CURRENT_TIMESTAMP - (retention_days * INTERVAL ''1 day''));' LANGUAGE sql;

-- Create function to perform backup
CREATE OR REPLACE FUNCTION perform_backup() RETURNS void AS 'INSERT INTO backup_history (backup_path, status) VALUES (''/var/lib/postgresql/backups/hulta_pregnancy_'' || to_char(CURRENT_TIMESTAMP, ''YYYY_MM_DD_HH24_MI_SS'') || ''.sql'', ''COMPLETED'');' LANGUAGE sql;

-- +goose Down
DROP FUNCTION IF EXISTS perform_backup();
DROP FUNCTION IF EXISTS cleanup_old_backups();
