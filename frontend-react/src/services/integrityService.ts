import { createHash } from 'crypto-browserify';

interface DataChecksum {
	hash: string;
	timestamp: number;
	dataType: string;
}

export const integrityService = {
	async calculateChecksum(data: any): Promise<string> {
		const str = JSON.stringify(data);
		return createHash('sha256').update(str).digest('hex');
	},

	async verifyIntegrity(data: any, storedChecksum: string): Promise<boolean> {
		const currentChecksum = await this.calculateChecksum(data);
		return currentChecksum === storedChecksum;
	},

	async storeChecksum(data: any, type: string) {
		const checksum: DataChecksum = {
			hash: await this.calculateChecksum(data),
			timestamp: Date.now(),
			dataType: type,
		};

		const db = await storageService.initDB();
		await db.put('checksums', checksum);
	},

	validateMeasurements(measurements: any[]) {
		return measurements.every((m) => {
			// Basic validation
			if (!m.value || !m.timestamp) return false;

			// Value range validation
			if (m.type === 'weight' && (m.value < 0 || m.value > 1000))
				return false;
			if (m.type === 'height' && (m.value < 0 || m.value > 200))
				return false;

			// Timestamp validation
			const date = new Date(m.timestamp);
			if (isNaN(date.getTime())) return false;

			return true;
		});
	},
};
