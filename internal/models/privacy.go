package models

import "time"

// PrivacyPreferences stores all opt-in preferences for privacy-sensitive features
type PrivacyPreferences struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"uniqueIndex"`
	UpdatedAt time.Time `json:"updated_at"`

	// Location and Environmental Data
	WeatherTrackingEnabled bool `json:"weather_tracking_enabled" gorm:"default:false"`
	LocationSharingEnabled bool `json:"location_sharing_enabled" gorm:"default:false"`
	
	// Health Data Sharing
	ShareHealthData      bool `json:"share_health_data" gorm:"default:false"`
	SharePregnancyData  bool `json:"share_pregnancy_data" gorm:"default:false"`
	ShareGeneticData    bool `json:"share_genetic_data" gorm:"default:false"`
	
	// Analytics and Tracking
	AllowAnonymousAnalytics bool `json:"allow_anonymous_analytics" gorm:"default:false"`
	AllowUsageTracking      bool `json:"allow_usage_tracking" gorm:"default:false"`
	
	// Notifications
	WeatherNotifications bool `json:"weather_notifications" gorm:"default:false"`
	HealthNotifications  bool `json:"health_notifications" gorm:"default:false"`
	EventNotifications   bool `json:"event_notifications" gorm:"default:false"`
	
	// Data Retention
	DataRetentionDays   int  `json:"data_retention_days" gorm:"default:365"`
	AutoDeleteOldData   bool `json:"auto_delete_old_data" gorm:"default:false"`
	
	// External Services
	AllowThirdPartySharing bool `json:"allow_third_party_sharing" gorm:"default:false"`
	AllowDataExport        bool `json:"allow_data_export" gorm:"default:false"`

	// Preferences for specific features
	WeatherPrefs WeatherPrivacyPreferences `json:"weather_preferences" gorm:"embedded"`
	HealthPrefs  HealthPreferences  `json:"health_preferences" gorm:"embedded"`
}

// WeatherPrivacyPreferences stores detailed weather tracking preferences
type WeatherPrivacyPreferences struct {
	UpdateFrequency    string  `json:"update_frequency" gorm:"default:'daily'"`
	DefaultLatitude    float64 `json:"default_latitude,omitempty"`
	DefaultLongitude   float64 `json:"default_longitude,omitempty"`
	StoreHistoricalData bool   `json:"store_historical_data" gorm:"default:false"`
}

// HealthPreferences stores detailed health tracking preferences
type HealthPreferences struct {
	ShareWithVets       bool `json:"share_with_vets" gorm:"default:false"`
	ShareWithBreeders   bool `json:"share_with_breeders" gorm:"default:false"`
	ShareWithResearchers bool `json:"share_with_researchers" gorm:"default:false"`
	StoreGeneticHistory bool `json:"store_genetic_history" gorm:"default:false"`
}

// PrivacyChangeLog records changes to privacy preferences
type PrivacyChangeLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"index"`
	Field     string    `json:"field"`
	OldValue  bool      `json:"old_value"`
	NewValue  bool      `json:"new_value"`
	Timestamp time.Time `json:"timestamp"`
}

// GetDefaultPrivacyPreferences returns default privacy settings
func GetDefaultPrivacyPreferences(userID string) PrivacyPreferences {
	return PrivacyPreferences{
		UserID:    userID,
		UpdatedAt: time.Now(),
		// All opt-in features default to false
		DataRetentionDays: 365, // Default to 1 year retention
	}
}

// IsFeatureEnabled checks if a specific privacy-sensitive feature is enabled
func (p *PrivacyPreferences) IsFeatureEnabled(feature string) bool {
	switch feature {
	case "weather_tracking":
		return p.WeatherTrackingEnabled
	case "location_sharing":
		return p.LocationSharingEnabled
	case "health_data":
		return p.ShareHealthData
	case "pregnancy_data":
		return p.SharePregnancyData
	case "genetic_data":
		return p.ShareGeneticData
	case "analytics":
		return p.AllowAnonymousAnalytics
	case "usage_tracking":
		return p.AllowUsageTracking
	case "third_party_sharing":
		return p.AllowThirdPartySharing
	default:
		return false
	}
}
