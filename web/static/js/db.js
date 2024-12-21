// IndexedDB wrapper for offline storage
class OfflineDB {
    constructor() {
        this.dbName = 'horseTrackingDB';
        this.version = 1;
        this.db = null;
    }

    async init() {
        return new Promise((resolve, reject) => {
            const request = indexedDB.open(this.dbName, this.version);

            request.onerror = () => reject(request.error);
            request.onsuccess = () => {
                this.db = request.result;
                resolve();
            };

            request.onupgradeneeded = (event) => {
                const db = event.target.result;

                // Create stores for offline data
                if (!db.objectStoreNames.contains('horses')) {
                    db.createObjectStore('horses', { keyPath: 'id', autoIncrement: true });
                }
                if (!db.objectStoreNames.contains('healthRecords')) {
                    db.createObjectStore('healthRecords', { keyPath: 'id', autoIncrement: true });
                }
                if (!db.objectStoreNames.contains('pregnancyEvents')) {
                    db.createObjectStore('pregnancyEvents', { keyPath: 'id', autoIncrement: true });
                }
                if (!db.objectStoreNames.contains('metadata')) {
                    db.createObjectStore('metadata', { keyPath: 'key' });
                }
            };
        });
    }

    async saveHorse(horse) {
        const store = this.db
            .transaction(['horses'], 'readwrite')
            .objectStore('horses');

        await store.put(horse);
    }

    async getHorse(id) {
        const store = this.db
            .transaction(['horses'], 'readonly')
            .objectStore('horses');

        return new Promise((resolve, reject) => {
            const request = store.get(id);
            request.onerror = () => reject(request.error);
            request.onsuccess = () => resolve(request.result);
        });
    }

    async getAllHorses() {
        const store = this.db
            .transaction(['horses'], 'readonly')
            .objectStore('horses');

        return new Promise((resolve, reject) => {
            const request = store.getAll();
            request.onerror = () => reject(request.error);
            request.onsuccess = () => resolve(request.result);
        });
    }

    async saveHealthRecord(record) {
        const store = this.db
            .transaction(['healthRecords'], 'readwrite')
            .objectStore('healthRecords');

        await store.put(record);
    }

    async getAllHealthRecords() {
        const store = this.db
            .transaction(['healthRecords'], 'readonly')
            .objectStore('healthRecords');

        return new Promise((resolve, reject) => {
            const request = store.getAll();
            request.onerror = () => reject(request.error);
            request.onsuccess = () => resolve(request.result);
        });
    }

    async getHealthRecords(horseId) {
        const records = await this.getAllHealthRecords();
        return records.filter(r => r.horseId === horseId);
    }

    async savePregnancyEvent(event) {
        const store = this.db
            .transaction(['pregnancyEvents'], 'readwrite')
            .objectStore('pregnancyEvents');

        await store.put(event);
    }

    async getAllPregnancyEvents() {
        const store = this.db
            .transaction(['pregnancyEvents'], 'readonly')
            .objectStore('pregnancyEvents');

        return new Promise((resolve, reject) => {
            const request = store.getAll();
            request.onerror = () => reject(request.error);
            request.onsuccess = () => resolve(request.result);
        });
    }

    async getPregnancyEvents(horseId) {
        const events = await this.getAllPregnancyEvents();
        return events.filter(e => e.horseId === horseId);
    }

    async updateLastSync(date) {
        const store = this.db
            .transaction(['metadata'], 'readwrite')
            .objectStore('metadata');

        await store.put({
            key: 'lastSync',
            value: date.toISOString()
        });
    }

    async getLastSync() {
        const store = this.db
            .transaction(['metadata'], 'readonly')
            .objectStore('metadata');

        return new Promise((resolve, reject) => {
            const request = store.get('lastSync');
            request.onerror = () => reject(request.error);
            request.onsuccess = () => {
                if (request.result) {
                    resolve(new Date(request.result.value));
                } else {
                    resolve(null);
                }
            };
        });
    }

    async clearAll() {
        const stores = ['horses', 'healthRecords', 'pregnancyEvents'];
        const tx = this.db.transaction(stores, 'readwrite');
        
        await Promise.all(
            stores.map(store => 
                new Promise((resolve, reject) => {
                    const request = tx.objectStore(store).clear();
                    request.onerror = () => reject(request.error);
                    request.onsuccess = () => resolve();
                })
            )
        );
    }
}
