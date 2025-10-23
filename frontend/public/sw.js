/**
 * Grompt Service Worker
 * Provides offline functionality and caching for the PWA
 */

const CACHE_NAME = 'grompt-v1.0.9';
const API_CACHE_NAME = 'grompt-api-v1.0.9';
const DYNAMIC_CACHE_NAME = 'grompt-dynamic-v1.0.9';

// Files to cache immediately
const STATIC_ASSETS = [
  '/',
  '/index.html',
  '/manifest.json',
  '/icons/icon-192x192.png',
  '/icons/icon-512x512.png'
];

// API endpoints to cache
const CACHEABLE_API_ENDPOINTS = [
  '/v1/providers',
  '/v1/health',
  '/api/v1/health'
];

// Maximum cache age (24 hours)
const MAX_CACHE_AGE = 24 * 60 * 60 * 1000;

// Maximum number of dynamic cache entries
const MAX_DYNAMIC_CACHE_SIZE = 50;

/**
 * Install event - cache static assets
 */
self.addEventListener('install', (event) => {
  console.log('🔧 Service Worker: Installing...');

  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => {
        console.log('📦 Service Worker: Caching static assets');
        return cache.addAll(STATIC_ASSETS);
      })
      .then(() => {
        console.log('✅ Service Worker: Installation complete');
        return self.skipWaiting();
      })
      .catch((error) => {
        console.error('❌ Service Worker: Installation failed', error);
      })
  );
});

/**
 * Activate event - clean up old caches
 */
self.addEventListener('activate', (event) => {
  console.log('🚀 Service Worker: Activating...');

  event.waitUntil(
    caches.keys()
      .then((cacheNames) => {
        return Promise.all(
          cacheNames.map((cacheName) => {
            if (
              cacheName !== CACHE_NAME &&
              cacheName !== API_CACHE_NAME &&
              cacheName !== DYNAMIC_CACHE_NAME
            ) {
              console.log('🗑️ Service Worker: Deleting old cache', cacheName);
              return caches.delete(cacheName);
            }
          })
        );
      })
      .then(() => {
        console.log('✅ Service Worker: Activation complete');
        return self.clients.claim();
      })
  );
});

/**
 * Fetch event - handle all network requests
 */
self.addEventListener('fetch', (event) => {
  const { request } = event;
  const url = new URL(request.url);

  // Skip non-GET requests for caching
  if (request.method !== 'GET') {
    return;
  }

  // Skip chrome-extension and other non-http requests
  if (!url.protocol.startsWith('http')) {
    return;
  }

  // Handle different types of requests
  if (isAPIRequest(url)) {
    event.respondWith(handleAPIRequest(request));
  } else if (isStaticAsset(url)) {
    event.respondWith(handleStaticAsset(request));
  } else {
    event.respondWith(handleDynamicContent(request));
  }
});

/**
 * Background sync for offline requests
 */
self.addEventListener('sync', (event) => {
  if (event.tag === 'background-sync') {
    console.log('🔄 Service Worker: Background sync triggered');
    event.waitUntil(syncOfflineRequests());
  }
});

/**
 * Handle push notifications
 */
self.addEventListener('push', (event) => {
  console.log('🔔 Service Worker: Push received');

  const options = {
    body: 'Novo conteúdo disponível no Grompt',
    icon: '/icons/icon-192x192.png',
    badge: '/icons/badge-72x72.png',
    vibrate: [100, 50, 100],
    data: {
      dateOfArrival: Date.now(),
      primaryKey: '1'
    },
    actions: [
      {
        action: 'explore',
        title: 'Abrir Grompt',
        icon: '/icons/action-explore.png'
      },
      {
        action: 'close',
        title: 'Fechar',
        icon: '/icons/action-close.png'
      }
    ]
  };

  event.waitUntil(
    self.registration.showNotification('Grompt', options)
  );
});

/**
 * Handle notification clicks
 */
self.addEventListener('notificationclick', (event) => {
  console.log('🖱️ Service Worker: Notification clicked');

  event.notification.close();

  if (event.action === 'explore') {
    event.waitUntil(
      clients.openWindow('/')
    );
  }
});

/**
 * Helper functions
 */

function isAPIRequest(url) {
  return url.pathname.startsWith('/v1/') || url.pathname.startsWith('/api/');
}

function isStaticAsset(url) {
  const staticExtensions = ['.js', '.css', '.png', '.jpg', '.jpeg', '.svg', '.ico', '.woff', '.woff2'];
  return staticExtensions.some(ext => url.pathname.endsWith(ext)) ||
         url.pathname === '/' ||
         url.pathname === '/index.html' ||
         url.pathname === '/manifest.json';
}

async function handleAPIRequest(request) {
  const url = new URL(request.url);

  try {
    // Always try network first for API requests
    const networkResponse = await fetch(request);

    // Cache successful responses for specific endpoints
    if (networkResponse.ok && CACHEABLE_API_ENDPOINTS.some(endpoint => url.pathname.startsWith(endpoint))) {
      const cache = await caches.open(API_CACHE_NAME);
      const responseClone = networkResponse.clone();

      // Add timestamp to cached response
      const modifiedResponse = new Response(responseClone.body, {
        status: responseClone.status,
        statusText: responseClone.statusText,
        headers: {
          ...Object.fromEntries(responseClone.headers.entries()),
          'sw-cached-at': Date.now().toString()
        }
      });

      cache.put(request, modifiedResponse);
    }

    return networkResponse;
  } catch (error) {
    console.log('🔄 Service Worker: Network failed for API, trying cache', url.pathname);

    // Network failed, try cache
    const cache = await caches.open(API_CACHE_NAME);
    const cachedResponse = await cache.match(request);

    if (cachedResponse) {
      const cachedAt = cachedResponse.headers.get('sw-cached-at');
      const isStale = cachedAt && (Date.now() - parseInt(cachedAt)) > MAX_CACHE_AGE;

      if (isStale) {
        console.log('⚠️ Service Worker: Cached API response is stale', url.pathname);
      }

      return cachedResponse;
    }

    // No cache available, return offline response
    return createOfflineAPIResponse(url);
  }
}

async function handleStaticAsset(request) {
  try {
    // Try cache first for static assets
    const cache = await caches.open(CACHE_NAME);
    const cachedResponse = await cache.match(request);

    if (cachedResponse) {
      // Return cached version and update in background
      updateCacheInBackground(request, cache);
      return cachedResponse;
    }

    // Not in cache, fetch from network
    const networkResponse = await fetch(request);

    if (networkResponse.ok) {
      // Cache the response
      cache.put(request, networkResponse.clone());
    }

    return networkResponse;
  } catch (error) {
    console.log('🔄 Service Worker: Network failed for static asset, trying cache');

    // Network failed, try cache again
    const cache = await caches.open(CACHE_NAME);
    const cachedResponse = await cache.match(request);

    if (cachedResponse) {
      return cachedResponse;
    }

    // For HTML requests, return the cached index.html (SPA fallback)
    if (request.headers.get('accept')?.includes('text/html')) {
      const indexResponse = await cache.match('/index.html');
      if (indexResponse) {
        return indexResponse;
      }
    }

    // No cache available
    return new Response('Offline - Recurso não disponível', {
      status: 503,
      statusText: 'Service Unavailable'
    });
  }
}

async function handleDynamicContent(request) {
  const url = new URL(request.url);
  const cache = await caches.open(DYNAMIC_CACHE_NAME);

  try {
    // Try network first
    const networkResponse = await fetch(request);

    if (networkResponse.ok) {
      // Cache successful responses
      await limitCacheSize(cache, MAX_DYNAMIC_CACHE_SIZE);
      cache.put(request, networkResponse.clone());
    }

    return networkResponse;
  } catch (error) {
    console.log('🔄 Service Worker: Network failed for dynamic content, trying cache');

    // Network failed, try cache
    const cachedResponse = await cache.match(request);

    if (cachedResponse) {
      return cachedResponse;
    }

    // No cache, return offline page
    return createOfflinePage();
  }
}

function createOfflineAPIResponse(url) {
  const offlineData = {
    error: 'offline',
    message: 'Aplicação em modo offline. Algumas funcionalidades podem estar limitadas.',
    offline: true,
    timestamp: Date.now()
  };

  // Customize response based on endpoint
  if (url.pathname.includes('/providers')) {
    offlineData.data = [
      {
        name: 'offline-template',
        available: true,
        type: 'template',
        defaultModel: 'template-based'
      }
    ];
  } else if (url.pathname.includes('/health')) {
    offlineData.status = 'offline';
    offlineData.service = 'grompt-v1';
    offlineData.version = '1.0.9';
  }

  return new Response(JSON.stringify(offlineData), {
    status: 200,
    statusText: 'OK (Offline)',
    headers: {
      'Content-Type': 'application/json',
      'SW-Offline': 'true'
    }
  });
}

function createOfflinePage() {
  const offlineHTML = `
    <!DOCTYPE html>
    <html lang="pt-BR">
    <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <title>Grompt - Offline</title>
      <style>
        body {
          font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
          background: #030712;
          color: white;
          margin: 0;
          padding: 20px;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          min-height: 100vh;
          text-align: center;
        }
        .offline-icon {
          font-size: 4rem;
          margin-bottom: 1rem;
        }
        h1 {
          color: #9333ea;
          margin-bottom: 1rem;
        }
        p {
          color: #d1d5db;
          line-height: 1.6;
          max-width: 400px;
        }
        .retry-btn {
          background: #9333ea;
          color: white;
          border: none;
          padding: 12px 24px;
          border-radius: 8px;
          font-size: 1rem;
          cursor: pointer;
          margin-top: 1rem;
          transition: background 0.3s;
        }
        .retry-btn:hover {
          background: #7c3aed;
        }
      </style>
    </head>
    <body>
      <div class="offline-icon">🔌</div>
      <h1>Modo Offline</h1>
      <p>
        Você está navegando offline. O Grompt continua funcionando com
        funcionalidades limitadas usando templates locais e dados em cache.
      </p>
      <button class="retry-btn" onclick="window.location.reload()">
        Tentar Novamente
      </button>
      <script>
        // Auto-reload when back online
        window.addEventListener('online', () => {
          window.location.reload();
        });
      </script>
    </body>
    </html>
  `;

  return new Response(offlineHTML, {
    status: 200,
    statusText: 'OK (Offline)',
    headers: {
      'Content-Type': 'text/html',
      'SW-Offline': 'true'
    }
  });
}

async function updateCacheInBackground(request, cache) {
  try {
    const networkResponse = await fetch(request);
    if (networkResponse.ok) {
      cache.put(request, networkResponse);
    }
  } catch (error) {
    // Silent background update failure
    console.log('🔄 Service Worker: Background cache update failed');
  }
}

async function limitCacheSize(cache, maxSize) {
  const keys = await cache.keys();
  if (keys.length >= maxSize) {
    // Remove oldest entries
    const entriesToDelete = keys.slice(0, keys.length - maxSize + 1);
    await Promise.all(entriesToDelete.map(key => cache.delete(key)));
  }
}

async function syncOfflineRequests() {
  try {
    // This would integrate with the IndexedDB offline queue
    // from the enhancedAPI service to sync pending requests
    console.log('🔄 Service Worker: Syncing offline requests...');

    // Send message to main thread to trigger sync
    const clients = await self.clients.matchAll();
    clients.forEach(client => {
      client.postMessage({
        type: 'SYNC_OFFLINE_QUEUE'
      });
    });
  } catch (error) {
    console.error('❌ Service Worker: Sync failed', error);
  }
}

/**
 * Handle messages from main thread
 */
self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING') {
    self.skipWaiting();
  }

  if (event.data && event.data.type === 'CACHE_NEW_ROUTE') {
    const cache = caches.open(DYNAMIC_CACHE_NAME);
    cache.then(c => c.add(event.data.url));
  }
});

console.log('🚀 Service Worker: Script loaded');