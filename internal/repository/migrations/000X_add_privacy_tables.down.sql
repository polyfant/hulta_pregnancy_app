-- Drop trigger
DROP TRIGGER IF EXISTS privacy_audit_trigger ON privacy_preferences;

-- Drop function
DROP FUNCTION IF EXISTS log_privacy_change();

-- Drop indexes
DROP INDEX IF EXISTS idx_privacy_log_changed_at;
DROP INDEX IF EXISTS idx_privacy_log_user;
DROP INDEX IF EXISTS idx_privacy_prefs_user;

-- Drop tables (in correct order due to foreign key constraints)
DROP TABLE IF EXISTS privacy_change_log;
DROP TABLE IF EXISTS privacy_preferences;
