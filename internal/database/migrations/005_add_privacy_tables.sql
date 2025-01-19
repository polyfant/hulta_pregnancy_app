-- +goose Up
-- Create privacy_preferences table
CREATE TABLE privacy_preferences (
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
    pregnancy_notifications BOOLEAN DEFAULT FALSE,
    system_notifications BOOLEAN DEFAULT FALSE,
    
    -- Data Retention
    data_retention_period INTEGER DEFAULT 365, -- days
    backup_retention_period INTEGER DEFAULT 90  -- days
);

-- Create consent_logs table
CREATE TABLE consent_logs (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    preference_type VARCHAR(50) NOT NULL,
    old_value BOOLEAN,
    new_value BOOLEAN,
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45),
    user_agent TEXT
);

-- Add indexes
CREATE INDEX idx_privacy_preferences_user_id ON privacy_preferences(user_id);
CREATE INDEX idx_consent_logs_user_id ON consent_logs(user_id);
CREATE INDEX idx_consent_logs_changed_at ON consent_logs(changed_at);

-- +goose Down
DROP TABLE IF EXISTS consent_logs;
DROP TABLE IF EXISTS privacy_preferences;
