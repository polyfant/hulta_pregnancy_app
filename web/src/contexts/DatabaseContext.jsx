import { createContext, useContext, useEffect, useState } from 'react';
import { OfflineDB } from '../lib/db';
import { SyncManager } from '../lib/sync';

const DatabaseContext = createContext(null);

export function DatabaseProvider({ children }) {
    const [db, setDb] = useState(null);
    const [syncManager, setSyncManager] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const initDb = async () => {
            try {
                const offlineDb = new OfflineDB();
                const sync = new SyncManager(offlineDb);
                setDb(offlineDb);
                setSyncManager(sync);
            } catch (err) {
                setError(err);
                console.error('Failed to initialize database:', err);
            } finally {
                setIsLoading(false);
            }
        };

        initDb();
    }, []);

    if (isLoading) {
        return <div>Loading database...</div>;
    }

    if (error) {
        return <div>Error initializing database: {error.message}</div>;
    }

    return (
        <DatabaseContext.Provider value={{ db, syncManager }}>
            {children}
        </DatabaseContext.Provider>
    );
}

export function useDatabase() {
    const context = useContext(DatabaseContext);
    if (!context) {
        throw new Error('useDatabase must be used within a DatabaseProvider');
    }
    return context;
}
