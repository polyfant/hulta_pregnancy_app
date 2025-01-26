import { useQuery, QueryFunction } from '@tanstack/react-query';
import { openDB } from 'idb';
import { 
  FoalGrowthData, 
  getChartDataKey,
  ChartDataKey,
} from '../types/indexeddb';
import { syncService } from '../services/syncService';

interface DBSchema {
  chartData: {
    key: ChartDataKey;
    value: FoalGrowthData;
    indexes: {};
  };
}

const STORE_NAME = 'chartData';

const fetchFoalGrowthData: QueryFunction<FoalGrowthData, ['foal-growth', number]> = async ({ queryKey }) => {
  const [, foalId] = queryKey;
  const dataKey = getChartDataKey(foalId);
  
  try {
    // Try online first
    const response = await fetch(`/api/foals/${foalId}/growth-data`);
    const sqlData: Array<{
      weight: number;
      height: number;
      measurementDate: string;
    }> = await response.json();

    // Transform SQL data to IndexedDB schema
    const data: FoalGrowthData = {
      foalId,
      weightData: sqlData.map(item => item.weight),
      heightData: sqlData.map(item => item.height),
      timestamp: Date.now()
    };

    // Cache the transformed data locally
    const db = await openDB<DBSchema>('foal-growth-cache', 1, {
      upgrade(database) {
        if (!database.objectStoreNames.contains(STORE_NAME)) {
          database.createObjectStore(STORE_NAME);
        }
      },
    });
    
    const tx = db.transaction(STORE_NAME, 'readwrite');
    const store = tx.objectStore(STORE_NAME);
    
    // Store the data with proper structure
    await store.put({ key: dataKey, value: data, indexes: {} }, dataKey);
    await tx.done;

    // If offline, register for sync when online
    if (!navigator.onLine) {
      await syncService.registerSync();
    }

    return data;
  } catch (error) {
    // If online fetch fails, try to get cached data
    const db = await openDB<DBSchema>('foal-growth-cache', 1);
    const tx = db.transaction(STORE_NAME, 'readonly');
    const store = tx.objectStore(STORE_NAME);
    const cachedData = await store.get(dataKey);
    await tx.done;

    // Ensure we return only FoalGrowthData
    if (cachedData) {
      const actualData: FoalGrowthData = 'value' in cachedData ? cachedData.value : cachedData;
      
      // Transform cached data to match the expected format
      const transformedData: FoalGrowthData = {
        foalId,
        weightData: actualData.weightData,
        heightData: actualData.heightData,
        timestamp: actualData.timestamp,
      };

      // If cached data exists but we're online, try to sync
      if (navigator.onLine) {
        try {
          await syncService.syncMeasurement(transformedData);
        } catch (syncError) {
          console.warn('Background sync failed', syncError);
        }
      }
      return transformedData;
    }

    throw error;
  }
};

export function useOfflineCharts(foalId: number) {
  return useQuery({
    queryKey: ['foal-growth', foalId] as const,
    queryFn: fetchFoalGrowthData,
    staleTime: 5 * 60 * 1000, // Consider data fresh for 5 minutes
    gcTime: Infinity, // Replace deprecated cacheTime with gcTime
  });
}
