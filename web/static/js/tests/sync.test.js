// @jest-environment jsdom
import { SyncManager } from '../sync.js';
import { OfflineDB } from '../db.js';

// Mock fetch globally
global.fetch = jest.fn();

describe('SyncManager', () => {
    let syncManager;
    let mockDB;

    beforeEach(() => {
        // Reset fetch mock
        fetch.mockReset();

        // Create mock DB
        mockDB = {
            getAllHorses: jest.fn(),
            getAllHealthRecords: jest.fn(),
            getAllPregnancyEvents: jest.fn(),
            updateLastSync: jest.fn(),
            clearSyncQueue: jest.fn(),
            clearAll: jest.fn(),
            saveHorse: jest.fn(),
            saveHealthRecord: jest.fn(),
            savePregnancyEvent: jest.fn()
        };

        syncManager = new SyncManager(mockDB);
    });

    describe('initialization', () => {
        test('should initialize with user credentials', async () => {
            await syncManager.init(1, 'test-token');
            expect(syncManager.userId).toBe(1);
            expect(syncManager.authToken).toBe('test-token');
        });
    });

    describe('sync operations', () => {
        beforeEach(() => {
            // Setup mock data
            mockDB.getAllHorses.mockResolvedValue([{ id: 1, name: 'Thunder' }]);
            mockDB.getAllHealthRecords.mockResolvedValue([{ id: 1, horseId: 1, type: 'Checkup' }]);
            mockDB.getAllPregnancyEvents.mockResolvedValue([{ id: 1, horseId: 1, type: 'Conception' }]);
        });

        test('should not sync when offline', async () => {
            // Mock offline status
            Object.defineProperty(navigator, 'onLine', { value: false });

            await syncManager.init(1, 'test-token');
            await syncManager.checkAndSync();

            expect(fetch).not.toHaveBeenCalled();
        });

        test('should sync data when online', async () => {
            // Mock online status
            Object.defineProperty(navigator, 'onLine', { value: true });

            // Mock successful sync
            fetch.mockResolvedValueOnce({
                ok: true,
                json: () => Promise.resolve({ status: 'success' })
            });

            await syncManager.init(1, 'test-token');
            await syncManager.checkAndSync();

            expect(fetch).toHaveBeenCalledWith('/api/sync', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer test-token'
                },
                body: expect.any(String)
            });

            expect(mockDB.updateLastSync).toHaveBeenCalled();
            expect(mockDB.clearSyncQueue).toHaveBeenCalled();
        });

        test('should handle sync failure', async () => {
            // Mock online status
            Object.defineProperty(navigator, 'onLine', { value: true });

            // Mock failed sync
            fetch.mockRejectedValueOnce(new Error('Sync failed'));

            await syncManager.init(1, 'test-token');
            
            // Add event listener for sync failure
            const syncFailedHandler = jest.fn();
            window.addEventListener('sync-failed', syncFailedHandler);

            await syncManager.checkAndSync();

            expect(syncFailedHandler).toHaveBeenCalled();
            expect(mockDB.updateLastSync).not.toHaveBeenCalled();
            expect(mockDB.clearSyncQueue).not.toHaveBeenCalled();
        });
    });

    describe('restore operations', () => {
        test('should restore data from server', async () => {
            const serverData = {
                horses: [{ id: 1, name: 'Thunder' }],
                health: [{ id: 1, horseId: 1, type: 'Checkup' }],
                events: [{ id: 1, horseId: 1, type: 'Conception' }]
            };

            fetch.mockResolvedValueOnce({
                ok: true,
                json: () => Promise.resolve(serverData)
            });

            await syncManager.init(1, 'test-token');
            await syncManager.restoreFromServer();

            expect(mockDB.clearAll).toHaveBeenCalled();
            expect(mockDB.saveHorse).toHaveBeenCalledWith(serverData.horses[0]);
            expect(mockDB.saveHealthRecord).toHaveBeenCalledWith(serverData.health[0]);
            expect(mockDB.savePregnancyEvent).toHaveBeenCalledWith(serverData.events[0]);
        });

        test('should handle restore failure', async () => {
            fetch.mockRejectedValueOnce(new Error('Restore failed'));

            await syncManager.init(1, 'test-token');
            
            await expect(syncManager.restoreFromServer()).rejects.toThrow('Failed to restore data');
            expect(mockDB.clearAll).not.toHaveBeenCalled();
        });
    });

    describe('sync status', () => {
        test('should get sync status', async () => {
            const lastSync = new Date();
            fetch.mockResolvedValueOnce({
                ok: true,
                json: () => Promise.resolve({ last_sync: lastSync, status: 'ok' })
            });

            await syncManager.init(1, 'test-token');
            const status = await syncManager.getLastSyncTime();

            expect(status).toEqual(new Date(lastSync));
        });

        test('should handle status check failure', async () => {
            fetch.mockRejectedValueOnce(new Error('Status check failed'));

            await syncManager.init(1, 'test-token');
            
            await expect(syncManager.getLastSyncTime()).rejects.toThrow('Failed to get sync status');
        });
    });
});
