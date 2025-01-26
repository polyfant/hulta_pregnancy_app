import { encryptionService } from './encryptionService';

export const dataManager = {
	async exportData(format: 'json' | 'csv' = 'json', encrypt = true) {
		const data = {
			measurements: await this.getAllMeasurements(),
			settings: await this.getAllSettings(),
			logs: await this.getLocalLogs(),
		};

		if (encrypt) {
			const key = await encryptionService.generateKey(
				Date.now().toString()
			);
			return {
				data: encryptionService.encrypt(data, key),
				key,
			};
		}

		return format === 'json'
			? JSON.stringify(data, null, 2)
			: this.convertToCSV(data);
	},

	async purgeData(options: {
		measurements?: boolean;
		settings?: boolean;
		logs?: boolean;
		environmentalData?: boolean;
	}) {
		const promises = [];
		if (options.measurements) {
			promises.push(this.clearMeasurements());
		}
		if (options.settings) {
			promises.push(this.clearSettings());
		}
		if (options.logs) {
			promises.push(this.clearLogs());
		}
		if (options.environmentalData) {
			promises.push(this.clearEnvironmentalData());
		}

		await Promise.all(promises);
		await auditLogger.logEvent({
			type: 'DATA_PURGE',
			component: 'DataManager',
			action: 'purge',
			success: true,
			anonymousMetrics: {
				purgeOptions: Object.keys(options).length,
			},
		});
	},
};
