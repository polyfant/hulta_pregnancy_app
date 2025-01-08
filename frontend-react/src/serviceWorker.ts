/// <reference lib="webworker" />

import { syncService } from './services/syncService';

declare const self: ServiceWorkerGlobalScope;

self.addEventListener('sync', (event) => {
	if (event.tag === 'sync-measurements') {
		event.waitUntil(syncService.syncPendingMeasurements());
	}
});

self.addEventListener('fetch', (event) => {
	// Add offline-first caching strategy
	event.respondWith(
		caches.match(event.request).then((response) => {
			return response || fetch(event.request);
		})
	);
});
