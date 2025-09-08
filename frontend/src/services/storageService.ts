import { HistoryItem, Draft } from '../types';

const DB_NAME = 'GromptDB';
const HISTORY_STORE = 'history';
const DRAFT_STORE = 'drafts';
const DB_VERSION = 1;

// LocalStorage fallback keys
const HISTORY_FALLBACK_KEY = 'history_fallback';
const DRAFT_FALLBACK_KEY = 'draft_fallback';

let dbPromise: Promise<IDBDatabase> | null = null;

const openDB = (): Promise<IDBDatabase> => {
  if (dbPromise) {
    return dbPromise;
  }
  dbPromise = new Promise((resolve, reject) => {
    try {
      const request = indexedDB.open(DB_NAME, DB_VERSION);
      request.onerror = () => {
        console.warn("IndexedDB could not be opened. It might be blocked or unavailable.");
        reject(new Error("IndexedDB unavailable"));
      };
      request.onsuccess = () => resolve(request.result);
      request.onupgradeneeded = (event) => {
        const db = (event.target as IDBOpenDBRequest).result;
        if (!db.objectStoreNames.contains(HISTORY_STORE)) {
          db.createObjectStore(HISTORY_STORE, { keyPath: 'id' });
        }
        if (!db.objectStoreNames.contains(DRAFT_STORE)) {
          db.createObjectStore(DRAFT_STORE);
        }
      };
    } catch (error) {
        console.error("Error initiating IndexedDB open request.", error);
        reject(error);
    }
  });
  return dbPromise;
};

const getFromStore = <T>(storeName: string, key: IDBValidKey): Promise<T | undefined> => {
    return openDB().then(db => {
        return new Promise<T | undefined>((resolve, reject) => {
            const transaction = db.transaction(storeName, 'readonly');
            const store = transaction.objectStore(storeName);
            const request = store.get(key);
            request.onerror = () => reject(new Error(`Error getting item ${key} from ${storeName}`));
            request.onsuccess = () => resolve(request.result as T | undefined);
        });
    });
};

const getAllFromStore = <T>(storeName: string): Promise<T[]> => {
    return openDB().then(db => {
        return new Promise<T[]>((resolve, reject) => {
            const transaction = db.transaction(storeName, 'readonly');
            const store = transaction.objectStore(storeName);
            const request = store.getAll();
            request.onerror = () => reject(new Error(`Error getting all items from ${storeName}`));
            request.onsuccess = () => resolve(request.result as T[]);
        });
    });
};

const putInStore = <T>(storeName: string, value: T, key?: IDBValidKey): Promise<void> => {
    return openDB().then(db => {
        return new Promise<void>((resolve, reject) => {
            const transaction = db.transaction(storeName, 'readwrite');
            const store = transaction.objectStore(storeName);
            const request = key ? store.put(value, key) : store.put(value);
            request.onerror = () => reject(new Error(`Error putting item into ${storeName}`));
            request.onsuccess = () => resolve();
        });
    });
};

// --- History Persistence ---

export const saveHistory = async (history: HistoryItem[]): Promise<void> => {
    try {
        const db = await openDB();
        const transaction = db.transaction(HISTORY_STORE, 'readwrite');
        const store = transaction.objectStore(HISTORY_STORE);
        store.clear(); // Clear old history
        for (const item of history) {
            store.put(item); // Add new items
        }
    } catch (error) {
        console.warn('Failed to save history to IndexedDB, using localStorage fallback.', error);
        try {
            localStorage.setItem(HISTORY_FALLBACK_KEY, JSON.stringify(history));
        } catch (lsError) {
            console.error('Failed to save history to localStorage fallback.', lsError);
        }
    }
};

export const getHistory = async (): Promise<HistoryItem[]> => {
    try {
        const history = await getAllFromStore<HistoryItem>(HISTORY_STORE);
        return history.sort((a, b) => b.timestamp - a.timestamp);
    } catch (error) {
        console.warn('Failed to get history from IndexedDB, using localStorage fallback.', error);
        try {
            const storedHistory = localStorage.getItem(HISTORY_FALLBACK_KEY);
            return storedHistory ? JSON.parse(storedHistory) : [];
        } catch (lsError) {
            console.error('Failed to get history from localStorage fallback.', lsError);
            return [];
        }
    }
};

// --- Draft Persistence ---

export const saveDraft = async (draft: Draft): Promise<void> => {
    try {
        await putInStore(DRAFT_STORE, draft, 'currentDraft');
    } catch (error) {
        console.warn('Failed to save draft to IndexedDB, using localStorage fallback.', error);
        try {
            localStorage.setItem(DRAFT_FALLBACK_KEY, JSON.stringify(draft));
        } catch (lsError) {
            console.error('Failed to save draft to localStorage fallback.', lsError);
        }
    }
};

export const getDraft = async (): Promise<Draft | undefined> => {
    try {
        return await getFromStore<Draft>(DRAFT_STORE, 'currentDraft');
    } catch (error) {
        console.warn('Failed to get draft from IndexedDB, using localStorage fallback.', error);
        try {
            const storedDraft = localStorage.getItem(DRAFT_FALLBACK_KEY);
            return storedDraft ? JSON.parse(storedDraft) : undefined;
        } catch (lsError) {
            console.error('Failed to get draft from localStorage fallback.', lsError);
            return undefined;
        }
    }
};

// --- Initialization and Migration ---

const migrateFromLocalStorage = async () => {
    try {
        // Check for history fallback
        const historyFallback = localStorage.getItem(HISTORY_FALLBACK_KEY);
        if (historyFallback) {
            console.log('Found history in localStorage, migrating to IndexedDB...');
            const history: HistoryItem[] = JSON.parse(historyFallback);
            await saveHistory(history); // This will save to IndexedDB if available
            localStorage.removeItem(HISTORY_FALLBACK_KEY);
            console.log('History migration successful.');
        }

        // Check for draft fallback
        const draftFallback = localStorage.getItem(DRAFT_FALLBACK_KEY);
        if (draftFallback) {
            console.log('Found draft in localStorage, migrating to IndexedDB...');
            const draft: Draft = JSON.parse(draftFallback);
            await saveDraft(draft); // This will save to IndexedDB if available
            localStorage.removeItem(DRAFT_FALLBACK_KEY);
            console.log('Draft migration successful.');
        }
    } catch (error) {
        console.error('Failed to migrate data from localStorage to IndexedDB.', error);
    }
};


let isInitialized = false;
export const initStorage = () => {
    if (isInitialized) return;

    // Attempt to open DB and migrate on successful connection
    openDB()
        .then(() => {
            console.log("IndexedDB connection successful. Checking for migrations...");
            migrateFromLocalStorage();
        })
        .catch(err => {
            console.warn("IndexedDB is not available on initial load.", err.message);
        });

    isInitialized = true;
};