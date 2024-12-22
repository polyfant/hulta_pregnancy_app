export class SyncManager {
    constructor(db) {
        this.db = db;
        this.apiBaseUrl = '/api';
    }

    async sync() {
        if (!navigator.onLine) {
            throw new Error('No internet connection');
        }

        try {
            // Get all local data
            const horses = await this.db.getAllHorses() || [];
            const healthRecords = await Promise.all(
                horses.map(horse => this.db.getHealthRecords(horse.id))
            ).then(records => records.flat());
            const pregnancyEvents = await Promise.all(
                horses.map(horse => this.db.getPregnancyEvents(horse.id))
            ).then(events => events.flat());

            // Send data to server
            const response = await fetch(`${this.apiBaseUrl}/sync`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    horses,
                    healthRecords,
                    pregnancyEvents
                })
            });

            if (!response.ok) {
                throw new Error('API Error');
            }

            // Update last sync time
            await this.db.updateLastSync(new Date().toISOString());

            return response.json();
        } catch (error) {
            console.error('Sync error:', error);
            throw error;
        }
    }

    async restore() {
        try {
            const response = await fetch(`${this.apiBaseUrl}/restore`);
            if (!response.ok) {
                throw new Error('API Error');
            }

            const data = await response.json();

            // Clear existing data
            await this.db.clear();

            // Save new data
            await Promise.all([
                ...(data.horses || []).map(horse => this.db.saveHorse(horse)),
                ...(data.healthRecords || []).map(record => this.db.saveHealthRecord(record)),
                ...(data.pregnancyEvents || []).map(event => this.db.savePregnancyEvent(event))
            ]);

            // Update last sync time
            await this.db.updateLastSync(new Date().toISOString());

            return data;
        } catch (error) {
            console.error('Restore error:', error);
            throw error;
        }
    }

    async getSyncStatus() {
        const lastSync = await this.db.getLastSync();
        return {
            lastSync,
            isOnline: navigator.onLine
        };
    }
}
