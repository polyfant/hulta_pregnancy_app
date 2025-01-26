-- Create privacy_preferences table
CREATE TABLE IF NOT EXISTS privacy_preferences (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL UNIQUE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Location and Environmental Data
    weather_tracking_enabled BOOLEAN DEFAULT FALSE,
    location_sharing_enabled BOOLEAN DEFAULT FALSE,
    
    -- Health Data Sharing
    share_health_data BOOLEAN DEFAULT FALSE,
    share_pregnancy_data BOOLEAN DEFAULT FALSE,
    share_genetic_data BOOLEAN DEFAULT FALSE,
    
    -- Analytics and Tracking
    allow_anonymous_analytics BOOLEAN DEFAULT FALSE,
    allow_usage_tracking BOOLEAN DEFAULT FALSE,
    
    -- Notifications
    weather_notifications BOOLEAN DEFAULT FALSE,
    health_notifications BOOLEAN DEFAULT FALSE,
    event_notifications BOOLEAN DEFAULT FALSE,
    
    -- Data Retention
    data_retention_days INTEGER DEFAULT 365,
    auto_delete_old_data BOOLEAN DEFAULT FALSE,
    
    -- External Services
    allow_third_party_sharing BOOLEAN DEFAULT FALSE,
    allow_data_export BOOLEAN DEFAULT FALSE,

    -- Weather Preferences
    weather_update_frequency VARCHAR(20) DEFAULT 'daily',
    default_latitude DOUBLE PRECISION,
    default_longitude DOUBLE PRECISION,
    store_historical_data BOOLEAN DEFAULT FALSE,

    -- Health Preferences
    share_with_vets BOOLEAN DEFAULT FALSE,
    share_with_breeders BOOLEAN DEFAULT FALSE,
    share_with_researchers BOOLEAN DEFAULT FALSE,
    store_genetic_history BOOLEAN DEFAULT FALSE,

    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- Create audit log for privacy changes
CREATE TABLE IF NOT EXISTS privacy_change_log (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    setting_name VARCHAR(100) NOT NULL,
    old_value TEXT,
    new_value TEXT,
    change_reason TEXT,
    
    CONSTRAINT fk_user_audit
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_privacy_prefs_user ON privacy_preferences(user_id);
CREATE INDEX IF NOT EXISTS idx_privacy_log_user ON privacy_change_log(user_id);
CREATE INDEX IF NOT EXISTS idx_privacy_log_changed_at ON privacy_change_log(changed_at);

-- Create function to log privacy changes
CREATE OR REPLACE FUNCTION log_privacy_change()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO privacy_change_log (
        user_id,
        setting_name,
        old_value,
        new_value
    )
    SELECT
        OLD.user_id,
        key,
        CASE 
            WHEN OLD.* IS NULL THEN NULL 
            ELSE OLD.*::text
        END,
        CASE 
            WHEN NEW.* IS NULL THEN NULL 
            ELSE NEW.*::text
        END
    FROM (
        SELECT unnest(ARRAY[
            'weather_tracking_enabled',
            'location_sharing_enabled',
            'share_health_data',
            'share_pregnancy_data',
            'share_genetic_data',
            'allow_anonymous_analytics',
            'allow_usage_tracking',
            'allow_third_party_sharing'
        ]) AS key
    ) k
    WHERE OLD.*::text IS DISTINCT FROM NEW.*::text;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for logging privacy changes
CREATE TRIGGER privacy_audit_trigger
AFTER UPDATE ON privacy_preferences
FOR EACH ROW
EXECUTE FUNCTION log_privacy_change();

-- Add comment for documentation
COMMENT ON TABLE privacy_preferences IS 'Stores user privacy preferences with all sensitive features defaulting to opt-out';
COMMENT ON TABLE privacy_change_log IS 'Audit log for tracking changes to privacy settings';
