import { compressionUtils } from '../utils/compression';
import { integrityService } from './integrityService';
import { storageService } from './storageService';

export const backupService = {
	async createBackup() {
		const db = await storageService.initDB();
		const allData = await db.getAll('data');

		const backup = {
			timestamp: new Date().toISOString(),
			version: storageService.currentVersion,
			data: allData,
			checksums: await Promise.all(
				allData.map((d) => integrityService.calculateChecksum(d))
			),
		};

		const compressed = compressionUtils.compressData(backup);

		// Create downloadable file
		const blob = new Blob([compressed], { type: 'application/json' });
		const url = URL.createObjectURL(blob);

		// Trigger download
		const a = document.createElement('a');
		a.href = url;
		a.download = `horse-app-backup-${new Date().toISOString()}.json`;
		a.click();

		URL.revokeObjectURL(url);

		// Also store in cloud if available
		if (navigator.onLine) {
			await this.uploadToCloud(compressed);
		}
	},

	async restore(backupFile: File) {
		const reader = new FileReader();

		return new Promise((resolve, reject) => {
			reader.onload = async (e) => {
				try {
					const compressed = e.target.result as string;
					const backup = compressionUtils.decompressData(compressed);

					// Verify integrity
					const valid = await Promise.all(
						backup.data.map((d, i) =>
							integrityService.verifyIntegrity(
								d,
								backup.checksums[i]
							)
						)
					);

					if (!valid.every((v) => v)) {
						throw new Error('Backup integrity check failed');
					}

					// Restore data
					const db = await storageService.initDB();
					await Promise.all(
						backup.data.map((d) => db.put('data', d))
					);

					resolve(backup);
				} catch (error) {
					reject(error);
				}
			};
			reader.readAsText(backupFile);
		});
	},

	async scheduleAutoBackup(intervalHours = 24) {
		setInterval(async () => {
			if (navigator.onLine) {
				await this.createBackup();
			}
		}, intervalHours * 60 * 60 * 1000);
	},
};
