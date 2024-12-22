import { openDB } from 'idb';

const DB_NAME = 'horseTrackingDB';
const DB_VERSION = 1;

export class OfflineDB {
    constructor() {
        this.dbPromise = openDB(DB_NAME, DB_VERSION, {
            upgrade(db) {
                // Create stores if they don't exist
                if (!db.objectStoreNames.contains('horses')) {
                    db.createObjectStore('horses', { keyPath: 'id' });
                }
                if (!db.objectStoreNames.contains('healthRecords')) {
                    const healthStore = db.createObjectStore('healthRecords', { keyPath: 'id' });
                    healthStore.createIndex('horseId', 'horseId');
                }
                if (!db.objectStoreNames.contains('pregnancyEvents')) {
                    const pregnancyStore = db.createObjectStore('pregnancyEvents', { keyPath: 'id' });
                    pregnancyStore.createIndex('horseId', 'horseId');
                }
                if (!db.objectStoreNames.contains('syncData')) {
                    db.createObjectStore('syncData', { keyPath: 'key' });
                }
            },
        });
    }

    async saveHorse(horse) {
        const db = await this.dbPromise;
        return db.put('horses', horse);
    }

    async getHorse(id) {
        const db = await this.dbPromise;
        return db.get('horses', id);
    }

    async getAllHorses() {
        const db = await this.dbPromise;
        return db.getAll('horses');
    }

    async saveHealthRecord(record) {
        const db = await this.dbPromise;
        return db.put('healthRecords', record);
    }

    async getHealthRecords(horseId) {
        const db = await this.dbPromise;
        const tx = db.transaction('healthRecords', 'readonly');
        const index = tx.store.index('horseId');
        return index.getAll(horseId);
    }

    async savePregnancyEvent(event) {
        const db = await this.dbPromise;
        return db.put('pregnancyEvents', event);
    }

    async getPregnancyEvents(horseId) {
        const db = await this.dbPromise;
        const tx = db.transaction('pregnancyEvents', 'readonly');
        const index = tx.store.index('horseId');
        return index.getAll(horseId);
    }

    async updateLastSync(timestamp) {
        const db = await this.dbPromise;
        return db.put('syncData', { key: 'lastSync', timestamp });
    }

    async getLastSync() {
        const db = await this.dbPromise;
        const result = await db.get('syncData', 'lastSync');
        return result ? result.timestamp : null;
    }

    async clear() {
        const db = await this.dbPromise;
        const stores = ['horses', 'healthRecords', 'pregnancyEvents', 'syncData'];
        const tx = db.transaction(stores, 'readwrite');
        await Promise.all(stores.map(store => tx.objectStore(store).clear()));
        await tx.done;
    }
}
