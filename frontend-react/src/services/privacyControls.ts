interface PrivacySettings {
	dataRetention: {
		measurementHistory: number; // days to keep
		environmentalData: number;
		locationHistory: number;
		autoDeleteEnabled: boolean;
	};
	dataMasking: {
		maskHorseNames: boolean;
		maskLocation: boolean;
		maskBreedInfo: boolean;
		maskMeasurements: boolean;
	};
	dataSharing: {
		shareAnonymized: boolean;
		shareAggregated: boolean;
		contributeToBenchmarks: boolean;
		allowResearch: boolean;
	};
}

export const privacyControls = {
	async getPrivacySettings(): Promise<PrivacySettings> {
		const stored = localStorage.getItem('privacy-settings');
		return stored ? JSON.parse(stored) : this.getDefaultPrivacySettings();
	},

	getDefaultPrivacySettings(): PrivacySettings {
		return {
			dataRetention: {
				measurementHistory: 365, // 1 year
				environmentalData: 30, // 30 days
				locationHistory: 7, // 1 week
				autoDeleteEnabled: true,
			},
			dataMasking: {
				maskHorseNames: true,
				maskLocation: true,
				maskBreedInfo: false,
				maskMeasurements: false,
			},
			dataSharing: {
				shareAnonymized: false,
				shareAggregated: false,
				contributeToBenchmarks: false,
				allowResearch: false,
			},
		};
	},

	maskSensitiveData(data: any, settings: PrivacySettings): any {
		const masked = { ...data };
		if (settings.dataMasking.maskHorseNames) {
			masked.horseName = `Horse-${data.id.slice(-4)}`;
		}
		if (settings.dataMasking.maskLocation) {
			masked.location = this.generalizeLocation(data.location);
		}
		return masked;
	},
};
