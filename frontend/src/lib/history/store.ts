import { IHistoryAdapter, SaveEntryInput, Session, EntryMeta, EntryFull } from './types';
import { IndexedDBAdapter } from './indexeddb';
import { LocalStorageAdapter } from './localstorage';

class HistoryStore implements IHistoryAdapter {
  private adapter: IHistoryAdapter;
  private inited = false;

  constructor() {
    // Choose IndexedDB if available, else fallback to LocalStorage
    const idbSupported = typeof indexedDB !== 'undefined';
    this.adapter = idbSupported ? new IndexedDBAdapter() : new LocalStorageAdapter();
  }

  async init(): Promise<void> {
    if (this.inited) return;
    await this.adapter.init();
    this.inited = true;
  }

  ensureDefaultSession(name?: string): Promise<Session> {
    return this.adapter.ensureDefaultSession(name);
  }

  createSession(name: string, agentId?: string): Promise<Session> {
    return this.adapter.createSession(name, agentId);
  }

  listSessions(): Promise<Session[]> {
    return this.adapter.listSessions();
  }

  saveEntry(input: SaveEntryInput): Promise<EntryMeta> {
    return this.adapter.saveEntry(input);
  }

  listEntries(sessionId: string, opts?: { limit?: number; offset?: number }): Promise<EntryMeta[]> {
    return this.adapter.listEntries(sessionId, opts);
  }

  getEntry(id: string): Promise<EntryFull | undefined> {
    return this.adapter.getEntry(id);
  }

  deleteEntry(id: string): Promise<void> {
    return this.adapter.deleteEntry(id);
  }

  clearSession(sessionId: string): Promise<number> {
    return this.adapter.clearSession(sessionId);
  }

  deleteSession(sessionId: string): Promise<void> {
    return this.adapter.deleteSession(sessionId);
  }

  async migrateFromLocalStorage(): Promise<{ created: number } | null> {
    try {
      await this.init();
      // Read previous UI state keys and turn into a single entry
      const ideasRaw = localStorage.getItem('grompt.prompt.ideas');
      const gen = localStorage.getItem('grompt.prompt.generated');
      const purpose = localStorage.getItem('grompt.prompt.purpose');
      const customPurpose = localStorage.getItem('grompt.prompt.customPurpose');
      const maxLength = localStorage.getItem('grompt.prompt.maxLength');
      if (!ideasRaw && !gen) return null;
      const ideas: { id: number; text: string }[] = ideasRaw ? JSON.parse(ideasRaw) : [];
      const engineeringPrompt = `Migração localStorage -> IndexedDB\n\nNotas:\n${ideas.map((i, idx) => `${idx + 1}. ${i.text}`).join('\n')}`;
      const session = await this.ensureDefaultSession();
      await this.saveEntry({
        sessionId: session.id,
        provider: 'demo',
        model: 'migration',
        params: { purpose, customPurpose, maxLength },
        requestText: engineeringPrompt,
        responseText: gen ? JSON.parse(gen) : '',
        status: 'ok',
        ideas,
      });
      // Clean keys only related to generated prompt (leave UI prefs)
      localStorage.removeItem('grompt.prompt.generated');
      return { created: 1 };
    } catch {
      return null;
    }
  }
}

export const history = new HistoryStore();
