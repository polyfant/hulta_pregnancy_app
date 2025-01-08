import { useQuery } from '@tanstack/react-query';
import { openDB } from 'idb';

export function useOfflineCharts(foalId: number) {
	return useQuery({
		queryKey: ['foal-growth', foalId],
		queryFn: async () => {
			try {
				// Try online first
				const response = await fetch(
					`/api/foals/${foalId}/growth-data`
				);
				const data = await response.json();

				// Cache the data
				const db = await openDB('foal-growth-cache', 1, {
					upgrade(db) {
						db.createObjectStore('chartData');
					},
				});
				await db.put('chartData', data, `foal-${foalId}`);

				return data;
			} catch (error) {
				// If offline, try to get cached data
				const db = await openDB('foal-growth-cache', 1);
				const cachedData = await db.get('chartData', `foal-${foalId}`);

				if (cachedData) {
					return cachedData;
				}
				throw new Error('No cached data available');
			}
		},
		staleTime: 5 * 60 * 1000, // Consider data fresh for 5 minutes
		cacheTime: Infinity, // Keep cache forever
	});
}
