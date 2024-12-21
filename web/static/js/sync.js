// Handles synchronization of offline data
class SyncManager {
    constructor(db) {
        this.db = db;
        this.isSyncing = false;
        this.userId = null;
        this.authToken = null;
    }

    async init(userId, authToken) {
        this.userId = userId;
        this.authToken = authToken;
        
        // Start periodic sync check
        setInterval(() => this.checkAndSync(), 60000); // Check every minute
        window.addEventListener('online', () => this.checkAndSync());
    }

    async checkAndSync() {
        if (!navigator.onLine || this.isSyncing || !this.authToken) {
            return;
        }

        this.isSyncing = true;
        try {
            // Get all local data
            const horses = await this.db.getAllHorses();
            const healthRecords = await this.db.getAllHealthRecords();
            const pregnancyEvents = await this.db.getAllPregnancyEvents();

            // Prepare sync data
            const syncData = {
                userId: this.userId,
                timestamp: new Date().toISOString(),
                horses: horses,
                health: healthRecords,
                events: pregnancyEvents
            };

            // Send to server
            const response = await fetch('/api/sync', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${this.authToken}`
                },
                body: JSON.stringify(syncData)
            });

            if (!response.ok) {
                throw new Error(`Sync failed: ${response.statusText}`);
            }

            // Update last sync time
            await this.db.updateLastSync(new Date());
            
            // Clear sync queue
            await this.db.clearSyncQueue();
            
            // Trigger sync success event
            window.dispatchEvent(new CustomEvent('sync-completed'));
        } catch (error) {
            console.error('Sync failed:', error);
            window.dispatchEvent(new CustomEvent('sync-failed', { detail: error }));
        } finally {
            this.isSyncing = false;
        }
    }

    async restoreFromServer() {
        if (!this.authToken) {
            throw new Error('Not authenticated');
        }

        const response = await fetch('/api/sync/restore', {
            headers: {
                'Authorization': `Bearer ${this.authToken}`
            }
        });

        if (!response.ok) {
            throw new Error('Failed to restore data');
        }

        const serverData = await response.json();

        // Clear local database
        await this.db.clearAll();

        // Restore all data
        for (const horse of serverData.horses) {
            await this.db.saveHorse(horse);
        }
        for (const record of serverData.health) {
            await this.db.saveHealthRecord(record);
        }
        for (const event of serverData.events) {
            await this.db.savePregnancyEvent(event);
        }

        window.dispatchEvent(new CustomEvent('restore-completed'));
    }

    async getLastSyncTime() {
        if (!this.authToken) {
            return null;
        }

        const response = await fetch('/api/sync/status', {
            headers: {
                'Authorization': `Bearer ${this.authToken}`
            }
        });

        if (!response.ok) {
            throw new Error('Failed to get sync status');
        }

        const status = await response.json();
        return new Date(status.last_sync);
    }
}
