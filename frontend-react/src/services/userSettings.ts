interface FeatureSettings {
	environmentalMonitoring: {
		enabled: boolean;
		airQuality: boolean;
		weatherAlerts: boolean;
		dataCollection: boolean;
	};
	locationTracking: {
		enabled: boolean;
		precision: 'off' | 'city' | 'precise';
	};
	dataSharing: {
		anonymous: boolean;
		breedComparisons: boolean;
	};
}

export const userSettings = {
	async getSettings(): Promise<FeatureSettings> {
		const stored = localStorage.getItem('user-settings');
		return stored ? JSON.parse(stored) : this.getDefaultSettings();
	},

	getDefaultSettings(): FeatureSettings {
		return {
			environmentalMonitoring: {
				enabled: false,
				airQuality: false,
				weatherAlerts: false,
				dataCollection: false,
			},
			locationTracking: {
				enabled: false,
				precision: 'off',
			},
			dataSharing: {
				anonymous: false,
				breedComparisons: false,
			},
		};
	},

	async updateSettings(settings: Partial<FeatureSettings>) {
		const current = await this.getSettings();
		const updated = { ...current, ...settings };
		localStorage.setItem('user-settings', JSON.stringify(updated));

		// If features are disabled, clean up any stored data
		if (!updated.environmentalMonitoring.enabled) {
			this.cleanupEnvironmentalData();
		}
		if (!updated.locationTracking.enabled) {
			this.cleanupLocationData();
		}
	},

	async cleanupEnvironmentalData() {
		localStorage.removeItem('environmental-cache');
		localStorage.removeItem('weather-data');
		// Clear any stored API keys or tokens
		localStorage.removeItem('weather-api-token');
	},

	async cleanupLocationData() {
		localStorage.removeItem('location-cache');
		localStorage.removeItem('location-history');
		// Ensure geolocation is stopped
		if (navigator.geolocation) {
			// Clear any active watchers
			const watchId = localStorage.getItem('location-watch-id');
			if (watchId) {
				navigator.geolocation.clearWatch(Number(watchId));
				localStorage.removeItem('location-watch-id');
			}
		}
	},
};
