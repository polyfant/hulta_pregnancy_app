// @jest-environment jsdom
import { OfflineDB } from '../db.js';

describe('OfflineDB', () => {
    let db;

    beforeEach(async () => {
        // Setup fake IndexedDB
        const indexedDB = new IDBFactory();
        global.indexedDB = indexedDB;
        
        db = new OfflineDB();
        await db.init();
    });

    afterEach(() => {
        // Clean up
        indexedDB.deleteDatabase('horseTrackingDB');
    });

    describe('Horse Operations', () => {
        test('should save and retrieve a horse', async () => {
            const horse = {
                id: 1,
                name: 'Thunder',
                breed: 'Arabian'
            };

            await db.saveHorse(horse);
            const retrieved = await db.getHorse(1);

            expect(retrieved).toEqual(horse);
        });

        test('should get all horses', async () => {
            const horses = [
                { id: 1, name: 'Thunder' },
                { id: 2, name: 'Storm' }
            ];

            await Promise.all(horses.map(h => db.saveHorse(h)));
            const retrieved = await db.getAllHorses();

            expect(retrieved).toHaveLength(2);
            expect(retrieved).toEqual(expect.arrayContaining(horses));
        });
    });

    describe('Health Record Operations', () => {
        test('should save and retrieve health records', async () => {
            const record = {
                id: 1,
                horseId: 1,
                type: 'Checkup',
                date: new Date().toISOString()
            };

            await db.saveHealthRecord(record);
            const records = await db.getHealthRecords(1);

            expect(records).toHaveLength(1);
            expect(records[0]).toEqual(record);
        });

        test('should filter health records by horse ID', async () => {
            const records = [
                { id: 1, horseId: 1, type: 'Checkup' },
                { id: 2, horseId: 2, type: 'Vaccination' },
                { id: 3, horseId: 1, type: 'Dental' }
            ];

            await Promise.all(records.map(r => db.saveHealthRecord(r)));
            const horse1Records = await db.getHealthRecords(1);

            expect(horse1Records).toHaveLength(2);
            expect(horse1Records.every(r => r.horseId === 1)).toBe(true);
        });
    });

    describe('Pregnancy Event Operations', () => {
        test('should save and retrieve pregnancy events', async () => {
            const event = {
                id: 1,
                horseId: 1,
                type: 'Conception',
                date: new Date().toISOString()
            };

            await db.savePregnancyEvent(event);
            const events = await db.getPregnancyEvents(1);

            expect(events).toHaveLength(1);
            expect(events[0]).toEqual(event);
        });
    });

    describe('Sync Operations', () => {
        test('should update and retrieve last sync time', async () => {
            const syncTime = new Date();
            await db.updateLastSync(syncTime);
            const retrieved = await db.getLastSync();

            expect(retrieved.getTime()).toBe(syncTime.getTime());
        });

        test('should return null if no sync time set', async () => {
            const lastSync = await db.getLastSync();
            expect(lastSync).toBeNull();
        });
    });

    describe('Clear Operations', () => {
        test('should clear all data', async () => {
            // Add some test data
            await db.saveHorse({ id: 1, name: 'Thunder' });
            await db.saveHealthRecord({ id: 1, horseId: 1, type: 'Checkup' });
            await db.savePregnancyEvent({ id: 1, horseId: 1, type: 'Conception' });

            // Clear all data
            await db.clearAll();

            // Verify all stores are empty
            const horses = await db.getAllHorses();
            const health = await db.getAllHealthRecords();
            const events = await db.getAllPregnancyEvents();

            expect(horses).toHaveLength(0);
            expect(health).toHaveLength(0);
            expect(events).toHaveLength(0);
        });
    });
});
