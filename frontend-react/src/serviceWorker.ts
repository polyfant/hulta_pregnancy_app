/// <reference lib="webworker" />

import { syncService } from './services/syncService';
import { precacheAndRoute } from 'workbox-precaching';
import { registerRoute } from 'workbox-routing';
import { NetworkFirst, CacheFirst } from 'workbox-strategies';
import { ExpirationPlugin } from 'workbox-expiration';
import { CacheableResponsePlugin } from 'workbox-cacheable-response';

declare const self: ServiceWorkerGlobalScope;

// Configuration for offline-first strategy
const CACHE_NAMES = {
  STATIC_ASSETS: 'static-assets-v1',
  API_CACHE: 'api-cache-v1',
  IMAGES_CACHE: 'images-cache-v1',
};

// Precache critical static assets
precacheAndRoute(self.__WB_MANIFEST || []);

// Cache static assets with cache-first strategy
registerRoute(
  ({ request }) => 
    request.destination === 'script' ||
    request.destination === 'style' ||
    request.destination === 'font',
  new CacheFirst({
    cacheName: CACHE_NAMES.STATIC_ASSETS,
    plugins: [
      new ExpirationPlugin({
        maxEntries: 60,
        maxAgeSeconds: 30 * 24 * 60 * 60, // 30 Days
      }),
      new CacheableResponsePlugin({
        statuses: [0, 200],
      }),
    ],
  })
);

// Network-first strategy for API calls with fallback cache
registerRoute(
  ({ request }) => request.destination === 'fetch' && request.url.includes('/api/'),
  new NetworkFirst({
    cacheName: CACHE_NAMES.API_CACHE,
    plugins: [
      new ExpirationPlugin({
        maxEntries: 10,
        maxAgeSeconds: 24 * 60 * 60, // 24 Hours
      }),
      new CacheableResponsePlugin({
        statuses: [0, 200],
      }),
    ],
    networkTimeoutSeconds: 3, // Fallback to cache after 3 seconds
  })
);

// Cache images with stale-while-revalidate
registerRoute(
  ({ request }) => request.destination === 'image',
  new CacheFirst({
    cacheName: CACHE_NAMES.IMAGES_CACHE,
    plugins: [
      new ExpirationPlugin({
        maxEntries: 100,
        maxAgeSeconds: 7 * 24 * 60 * 60, // 7 Days
      }),
      new CacheableResponsePlugin({
        statuses: [0, 200],
      }),
    ],
  })
);

// Background sync for offline data
self.addEventListener('sync', (event: ExtendedSyncEvent) => {
  if (event.tag === 'sync-measurements') {
    event.waitUntil(syncMeasurements());
  }
});

// Offline data synchronization
async function syncMeasurements() {
  try {
    // Retrieve pending measurements from IndexedDB
    const db = await openDB('foal-growth-db', 1);
    const pendingMeasurements = await db.getAllFromIndex('measurements', 'by-sync', false);
    
    for (const measurement of pendingMeasurements) {
      try {
        // Attempt to sync each measurement
        await fetch('/api/measurements', {
          method: 'POST',
          body: JSON.stringify(measurement),
        });
        
        // Mark as synced if successful
        await db.put('measurements', { ...measurement, synced: true });
      } catch (error) {
        console.error('Sync failed for measurement', measurement);
      }
    }
  } catch (error) {
    console.error('Background sync failed', error);
  }
}

// Install and activate handlers
self.addEventListener('install', (event) => {
  self.skipWaiting(); // Activate service worker immediately
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames
          .filter((name) => 
            name !== CACHE_NAMES.STATIC_ASSETS && 
            name !== CACHE_NAMES.API_CACHE && 
            name !== CACHE_NAMES.IMAGES_CACHE
          )
          .map((name) => caches.delete(name))
      );
    })
  );
});

// Typescript type extension for sync event
interface ExtendedSyncEvent extends Event {
  tag: string;
  waitUntil(promise: Promise<any>): void;
}

// Offline-first caching strategy
self.addEventListener('fetch', (event) => {
  event.respondWith(
    caches.match(event.request).then((response) => {
      return response || fetch(event.request);
    })
  );
});
