import { DBSchema, openDB } from 'idb';
import { compressionUtils } from '../utils/compression';

interface VersionedData<T> {
	version: number;
	timestamp: number;
	data: T;
	compressed?: boolean;
}

interface StorageDB extends DBSchema {
	data: {
		key: string;
		value: VersionedData<any>;
		indexes: { 'by-version': number };
	};
	migrations: {
		key: number;
		value: {
			applied: boolean;
			timestamp: number;
		};
	};
}

export const storageService = {
	currentVersion: 1,
	compressionThreshold: 1024 * 10, // 10KB

	async initDB() {
		return openDB<StorageDB>('horse-app-storage', this.currentVersion, {
			upgrade: (db, oldVersion, newVersion) => {
				if (!db.objectStoreNames.contains('data')) {
					const store = db.createObjectStore('data', {
						keyPath: 'id',
					});
					store.createIndex('by-version', 'version');
				}
				db.createObjectStore('migrations', { keyPath: 'version' });
			},
		});
	},

	async store<T>(key: string, data: T) {
		const db = await this.initDB();
		const shouldCompress =
			JSON.stringify(data).length > this.compressionThreshold;

		const versionedData: VersionedData<T> = {
			version: this.currentVersion,
			timestamp: Date.now(),
			data: shouldCompress ? compressionUtils.compressData(data) : data,
			compressed: shouldCompress,
		};

		await db.put('data', versionedData, key);
	},

	async retrieve<T>(key: string): Promise<T | null> {
		const db = await this.initDB();
		const versionedData = (await db.get('data', key)) as VersionedData<T>;

		if (!versionedData) return null;

		if (versionedData.compressed) {
			return compressionUtils.decompressData<T>(versionedData.data);
		}

		return versionedData.data;
	},
};
