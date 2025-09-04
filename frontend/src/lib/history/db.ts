// Tiny IndexedDB helper without external deps

const DB_NAME = 'grompt@v1';
const DB_VERSION = 1;

export type IDBDatabaseEx = IDBDatabase;

export function openHistoryDB(): Promise<IDBDatabaseEx> {
  return new Promise((resolve, reject) => {
    if (!('indexedDB' in globalThis)) {
      return reject(new Error('IndexedDB not available'));
    }
    const req = indexedDB.open(DB_NAME, DB_VERSION);
    req.onupgradeneeded = () => {
      const db = req.result;
      // sessions
      if (!db.objectStoreNames.contains('sessions')) {
        const s = db.createObjectStore('sessions', { keyPath: 'id' });
        s.createIndex('by_createdAt', 'createdAt');
        s.createIndex('by_name', 'name', { unique: false });
      }
      // entries
      if (!db.objectStoreNames.contains('entries')) {
        const e = db.createObjectStore('entries', { keyPath: 'id' });
        e.createIndex('by_session', 'sessionId', { unique: false });
        e.createIndex('by_createdAt', 'createdAt', { unique: false });
        e.createIndex('by_provider_model', ['provider', 'model'], { unique: false });
      }
      // blobs
      if (!db.objectStoreNames.contains('blobs')) {
        db.createObjectStore('blobs', { keyPath: 'id' });
      }
    };
    req.onsuccess = () => resolve(req.result);
    req.onerror = () => reject(req.error || new Error('Failed to open DB'));
  });
}

export function tx<T extends 'readonly' | 'readwrite'>(db: IDBDatabaseEx, mode: T, stores: string[]) {
  return db.transaction(stores, mode);
}

export function genId(prefix = 'id'): string {
  try {
    // @ts-ignore
    if (globalThis.crypto?.randomUUID) return globalThis.crypto.randomUUID();
  } catch {}
  const rnd = Math.floor(Math.random() * 1e9).toString(36);
  return `${prefix}_${Date.now().toString(36)}_${rnd}`;
}

