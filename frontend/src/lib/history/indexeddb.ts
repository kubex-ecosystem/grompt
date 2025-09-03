import { EntryFull, EntryMeta, IHistoryAdapter, SaveEntryInput, Session } from './types';
import { genId, openHistoryDB, tx } from './db';

export class IndexedDBAdapter implements IHistoryAdapter {
  private db?: IDBDatabase;
  private defaultSessionId?: string;

  async init(): Promise<void> {
    if (this.db) return;
    this.db = await openHistoryDB();
  }

  async ensureDefaultSession(name = 'Padrão'): Promise<Session> {
    await this.init();
    const db = this.db!;
    if (this.defaultSessionId) {
      const existing = await this.getSession(this.defaultSessionId);
      if (existing) return existing;
    }
    // Try to find any session named as default
    const s = await this.findSessionByName(name);
    if (s) {
      this.defaultSessionId = s.id;
      return s;
    }
    const created = await this.createSession(name);
    this.defaultSessionId = created.id;
    return created;
  }

  private getSession(id: string): Promise<Session | undefined> {
    const db = this.db!;
    return new Promise((resolve, reject) => {
      const t = tx(db, 'readonly', ['sessions']);
      const req = t.objectStore('sessions').get(id);
      req.onsuccess = () => resolve(req.result as Session | undefined);
      req.onerror = () => reject(req.error);
    });
  }

  private findSessionByName(name: string): Promise<Session | undefined> {
    const db = this.db!;
    return new Promise((resolve, reject) => {
      const t = tx(db, 'readonly', ['sessions']);
      const idx = t.objectStore('sessions').index('by_name');
      const req = idx.openCursor();
      let found: Session | undefined;
      req.onsuccess = () => {
        const cursor = req.result;
        if (!cursor) return resolve(found);
        const value = cursor.value as Session;
        if (value.name === name) {
          found = value;
          return resolve(found);
        }
        cursor.continue();
      };
      req.onerror = () => reject(req.error);
    });
  }

  async createSession(name: string, agentId?: string): Promise<Session> {
    await this.init();
    const db = this.db!;
    const now = Date.now();
    const s: Session = {
      id: genId('sess'),
      name,
      agentId,
      createdAt: now,
      updatedAt: now,
    };
    await new Promise<void>((resolve, reject) => {
      const t = tx(db, 'readwrite', ['sessions']);
      const req = t.objectStore('sessions').add(s);
      req.onsuccess = () => resolve();
      req.onerror = () => reject(req.error);
    });
    return s;
  }

  async listSessions(): Promise<Session[]> {
    await this.init();
    const db = this.db!;
    return new Promise((resolve, reject) => {
      const t = tx(db, 'readonly', ['sessions']);
      const req = t.objectStore('sessions').getAll();
      req.onsuccess = () => {
        const all = (req.result as Session[]).sort((a, b) => b.updatedAt - a.updatedAt);
        resolve(all);
      };
      req.onerror = () => reject(req.error);
    });
  }

  private async saveBlobIfNeeded(text?: string): Promise<string | undefined> {
    if (!text) return undefined;
    const db = this.db!;
    const id = genId('blob');
    const blob = new Blob([text], { type: 'text/plain' });
    await new Promise<void>((resolve, reject) => {
      const t = tx(db, 'readwrite', ['blobs']);
      const req = t.objectStore('blobs').put({ id, blob });
      req.onsuccess = () => resolve();
      req.onerror = () => reject(req.error);
    });
    return id;
  }

  async saveEntry(input: SaveEntryInput): Promise<EntryMeta> {
    await this.init();
    const db = this.db!;
    const session = input.sessionId
      ? await this.getSession(input.sessionId).then((r) => r || this.ensureDefaultSession())
      : await this.ensureDefaultSession(input.sessionName);

    const [requestBlobId, responseBlobId] = await Promise.all([
      this.saveBlobIfNeeded(input.requestText),
      this.saveBlobIfNeeded(input.responseText),
    ]);

    const now = Date.now();
    const promptPreview = (input.ideas && input.ideas.length)
      ? input.ideas.map(i => i.text).join(' • ').slice(0, 300)
      : (input.requestText || '').slice(0, 300);
    const full: EntryFull = {
      id: genId('ent'),
      sessionId: session.id,
      provider: input.provider,
      model: input.model,
      params: input.params,
      promptPreview,
      status: input.status || (input.error ? 'error' : 'ok'),
      createdAt: now,
      updatedAt: now,
      requestBlobId,
      responseBlobId,
      // Keep tiny payloads inline to simplify
      requestText: input.requestText && input.requestText.length <= 2000 ? input.requestText : undefined,
      responseText: input.responseText && input.responseText.length <= 8000 ? input.responseText : undefined,
      error: input.error,
      ideas: input.ideas,
    };

    await new Promise<void>((resolve, reject) => {
      const t = tx(db, 'readwrite', ['entries', 'sessions']);
      t.objectStore('entries').add(full);
      const sess = { ...session, updatedAt: now } as Session;
      t.objectStore('sessions').put(sess);
      t.oncomplete = () => resolve();
      t.onerror = () => reject(t.error);
      t.onabort = () => reject(t.error);
    });

    const meta: EntryMeta = {
      id: full.id,
      sessionId: full.sessionId,
      provider: full.provider,
      model: full.model,
      promptPreview: full.promptPreview,
      status: full.status,
      createdAt: full.createdAt,
      updatedAt: full.updatedAt,
      tokenCounts: full.tokenCounts,
    };
    return meta;
  }

  async listEntries(sessionId: string, opts?: { limit?: number; offset?: number }): Promise<EntryMeta[]> {
    await this.init();
    const db = this.db!;
    const { limit = 100, offset = 0 } = opts || {};
    return new Promise((resolve, reject) => {
      const t = tx(db, 'readonly', ['entries']);
      const idx = t.objectStore('entries').index('by_session');
      const req = idx.openCursor(IDBKeyRange.only(sessionId), 'prev');
      const acc: EntryMeta[] = [];
      let skipped = 0;
      req.onsuccess = () => {
        const cursor = req.result;
        if (!cursor) return resolve(acc);
        if (skipped < offset) {
          skipped++;
          return cursor.continue();
        }
        if (acc.length < limit) {
          const v = cursor.value as EntryFull;
          acc.push({
            id: v.id,
            sessionId: v.sessionId,
            provider: v.provider,
            model: v.model,
            promptPreview: v.promptPreview,
            status: v.status,
            createdAt: v.createdAt,
            updatedAt: v.updatedAt,
            tokenCounts: v.tokenCounts,
          });
          cursor.continue();
        } else {
          resolve(acc);
        }
      };
      req.onerror = () => reject(req.error);
    });
  }

  async getEntry(id: string): Promise<EntryFull | undefined> {
    await this.init();
    const db = this.db!;
    const base: EntryFull | undefined = await new Promise((resolve, reject) => {
      const t = tx(db, 'readonly', ['entries']);
      const req = t.objectStore('entries').get(id);
      req.onsuccess = () => resolve(req.result as EntryFull | undefined);
      req.onerror = () => reject(req.error);
    });
    if (!base) return undefined;
    const withPayloads: EntryFull = { ...base };
    // If inline exists, skip blob read
    if (!withPayloads.requestText && withPayloads.requestBlobId) {
      withPayloads.requestText = await this.readBlobText(withPayloads.requestBlobId).catch(() => undefined);
    }
    if (!withPayloads.responseText && withPayloads.responseBlobId) {
      withPayloads.responseText = await this.readBlobText(withPayloads.responseBlobId).catch(() => undefined);
    }
    return withPayloads;
  }

  private async readBlobText(id: string): Promise<string> {
    const db = this.db!;
    const rec: { id: string; blob: Blob } | undefined = await new Promise((resolve, reject) => {
      const t = tx(db, 'readonly', ['blobs']);
      const req = t.objectStore('blobs').get(id);
      req.onsuccess = () => resolve(req.result as any);
      req.onerror = () => reject(req.error);
    });
    if (!rec) throw new Error('Blob not found');
    return await rec.blob.text();
  }

  async deleteEntry(id: string): Promise<void> {
    await this.init();
    const db = this.db!;
    await new Promise<void>((resolve, reject) => {
      const t = tx(db, 'readwrite', ['entries']);
      const req = t.objectStore('entries').delete(id);
      req.onsuccess = () => resolve();
      req.onerror = () => reject(req.error);
    });
  }
}
