package types

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string `json:"message"`
}

// WeatherResponse represents a weather data response
type WeatherResponse struct {
	Temperature    float64 `json:"temperature"`
	Humidity      float64 `json:"humidity"`
	WindSpeed     float64 `json:"wind_speed"`
	Precipitation float64 `json:"precipitation"`
	UpdatedAt     string  `json:"updated_at"`
}

// PrivacySettingsResponse represents a response containing privacy settings
type PrivacySettingsResponse struct {
	WeatherTrackingEnabled bool `json:"weather_tracking_enabled"`
	LocationSharingEnabled bool `json:"location_sharing_enabled"`
	ShareHealthData       bool `json:"share_health_data"`
	SharePregnancyData   bool `json:"share_pregnancy_data"`
	ShareGeneticData     bool `json:"share_genetic_data"`
	AllowAnonymousAnalytics bool `json:"allow_anonymous_analytics"`
	AllowUsageTracking    bool `json:"allow_usage_tracking"`
	WeatherNotifications  bool `json:"weather_notifications"`
	HealthNotifications   bool `json:"health_notifications"`
	EventNotifications    bool `json:"event_notifications"`
	DataRetentionDays    int  `json:"data_retention_days"`
	AutoDeleteOldData    bool `json:"auto_delete_old_data"`
	AllowThirdPartySharing bool `json:"allow_third_party_sharing"`
	AllowDataExport      bool `json:"allow_data_export"`
}
