import { differenceInMinutes } from 'date-fns';
import { DBSchema, openDB } from 'idb';

interface MeasurementDB extends DBSchema {
	measurements: {
		key: string;
		value: {
			id?: string;
			foalId: number;
			weight: number;
			height: number;
			date: string;
			synced: boolean;
		};
		indexes: { 'by-sync': boolean };
	};
}

interface SyncStatus {
	pendingCount: number;
	lastSyncAttempt: Date | null;
	lastSuccessfulSync: Date | null;
	syncErrors: Array<{
		timestamp: Date;
		error: string;
	}>;
}

interface ConflictResolution {
	localData: any;
	serverData: any;
	resolution: 'local' | 'server' | 'merge';
}

interface MergeStrategy {
	type: 'timestamp' | 'value-based' | 'manual';
	resolution: 'local' | 'server' | 'merge';
	mergeFunction?: (local: any, server: any) => any;
}

export const syncService = {
	async initDB() {
		return openDB<MeasurementDB>('foal-growth-db', 1, {
			upgrade(db) {
				const measurementStore = db.createObjectStore('measurements', {
					keyPath: 'id',
					autoIncrement: true,
				});
				measurementStore.createIndex('by-sync', 'synced');
			},
		});
	},

	async saveMeasurement(measurement) {
		const db = await this.initDB();
		await db.add('measurements', {
			...measurement,
			synced: navigator.onLine,
			id: crypto.randomUUID(),
		});

		if (navigator.onLine) {
			await this.syncMeasurement(measurement);
		} else {
			// Register for sync when back online
			await this.registerSync();
		}
	},

	async syncMeasurement(measurement) {
		try {
			const response = await fetch(
				`/api/foals/${measurement.foalId}/measurements`,
				{
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify(measurement),
				}
			);

			if (response.ok) {
				const db = await this.initDB();
				await db.put('measurements', { ...measurement, synced: true });
			}
		} catch (error) {
			console.error('Sync failed:', error);
		}
	},

	async registerSync() {
		if ('serviceWorker' in navigator && 'sync' in registration) {
			const registration = await navigator.serviceWorker.ready;
			await registration.sync.register('sync-measurements');
		}
	},

	async syncPendingMeasurements() {
		const db = await this.initDB();
		const unsyncedMeasurements = await db.getAllFromIndex(
			'measurements',
			'by-sync',
			false
		);

		for (const measurement of unsyncedMeasurements) {
			await this.syncMeasurement(measurement);
		}
	},

	mergeStrategies: {
		measurements: {
			type: 'value-based',
			resolution: 'merge',
			mergeFunction: (local, server) => {
				// Combine measurements and remove duplicates
				const combined = [...local, ...server];
				return Array.from(
					new Map(
						combined.map((item) => [item.timestamp, item])
					).values()
				);
			},
		},
		notes: {
			type: 'timestamp',
			resolution: 'merge',
			mergeFunction: (local, server) => {
				// Merge notes, keeping newest version of conflicts
				return local.map((note) => {
					const serverNote = server.find((s) => s.id === note.id);
					return serverNote &&
						new Date(serverNote.updatedAt) >
							new Date(note.updatedAt)
						? serverNote
						: note;
				});
			},
		},
	},

	async handleConflict(
		local,
		server,
		type = 'measurements'
	): Promise<ConflictResolution> {
		const strategy = this.mergeStrategies[type];

		if (!strategy) {
			return this.defaultConflictResolution(local, server);
		}

		switch (strategy.type) {
			case 'timestamp':
				return this.handleTimestampBasedConflict(local, server);
			case 'value-based':
				return {
					localData: local,
					serverData: server,
					resolution: 'merge',
					mergedData: strategy.mergeFunction(local, server),
				};
			case 'manual':
				// Could show UI for manual resolution
				return this.defaultConflictResolution(local, server);
			default:
				return this.defaultConflictResolution(local, server);
		}
	},

	async handleTimestampBasedConflict(local, server) {
		const timeDiff = differenceInMinutes(
			new Date(server.lastModified),
			new Date(local.lastModified)
		);

		if (Math.abs(timeDiff) < 5) {
			// Changes within 5 minutes, merge them
			return {
				localData: local,
				serverData: server,
				resolution: 'merge',
				mergedData: this.mergeStrategies.notes.mergeFunction(
					local,
					server
				),
			};
		}

		return {
			localData: local,
			serverData: server,
			resolution: timeDiff > 0 ? 'server' : 'local',
		};
	},
};
