import { EntryFull, EntryMeta, IHistoryAdapter, SaveEntryInput, Session } from './types';
import { genId } from './db';

// Minimal LocalStorage adapter as fallback (size-limited)

type Persisted = {
  sessions: Session[];
  entries: EntryFull[];
};

const LS_KEY = 'grompt.history.v1';

function load(): Persisted {
  try {
    const raw = localStorage.getItem(LS_KEY);
    if (!raw) return { sessions: [], entries: [] };
    const parsed = JSON.parse(raw) as Persisted;
    return parsed;
  } catch {
    return { sessions: [], entries: [] };
  }
}

function save(data: Persisted) {
  try {
    localStorage.setItem(LS_KEY, JSON.stringify(data));
  } catch {}
}

export class LocalStorageAdapter implements IHistoryAdapter {
  private defaultSessionId?: string;

  async init(): Promise<void> {
    // no-op
  }

  async ensureDefaultSession(name = 'Padrão'): Promise<Session> {
    const data = load();
    if (this.defaultSessionId) {
      const s = data.sessions.find((s) => s.id === this.defaultSessionId);
      if (s) return s;
    }
    let s = data.sessions.find((s) => s.name === name);
    if (!s) {
      s = await this.createSession(name);
    }
    this.defaultSessionId = s.id;
    return s;
  }

  async createSession(name: string, agentId?: string): Promise<Session> {
    const data = load();
    const now = Date.now();
    const s: Session = { id: genId('sess'), name, agentId, createdAt: now, updatedAt: now };
    data.sessions.push(s);
    save(data);
    return s;
  }

  async listSessions(): Promise<Session[]> {
    const data = load();
    return data.sessions.sort((a, b) => b.updatedAt - a.updatedAt);
  }

  async saveEntry(input: SaveEntryInput): Promise<EntryMeta> {
    const data = load();
    const sess = input.sessionId
      ? data.sessions.find((s) => s.id === input.sessionId) || (await this.ensureDefaultSession())
      : await this.ensureDefaultSession(input.sessionName);
    const now = Date.now();
    const promptPreview = (input.ideas && input.ideas.length)
      ? input.ideas.map(i => i.text).join(' • ').slice(0, 300)
      : (input.requestText || '').slice(0, 300);
    const full: EntryFull = {
      id: genId('ent'),
      sessionId: sess.id,
      provider: input.provider,
      model: input.model,
      params: input.params,
      promptPreview,
      status: input.status || (input.error ? 'error' : 'ok'),
      createdAt: now,
      updatedAt: now,
      requestText: input.requestText,
      responseText: input.responseText,
      error: input.error,
      ideas: input.ideas,
    };
    data.entries.push(full);
    // Update session
    const idx = data.sessions.findIndex((s) => s.id === sess.id);
    if (idx >= 0) data.sessions[idx] = { ...sess, updatedAt: now };
    save(data);
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
    const { limit = 100, offset = 0 } = opts || {};
    const data = load();
    return data.entries
      .filter((e) => e.sessionId === sessionId)
      .sort((a, b) => b.createdAt - a.createdAt)
      .slice(offset, offset + limit)
      .map((v) => ({
        id: v.id,
        sessionId: v.sessionId,
        provider: v.provider,
        model: v.model,
        promptPreview: v.promptPreview,
        status: v.status,
        createdAt: v.createdAt,
        updatedAt: v.updatedAt,
        tokenCounts: v.tokenCounts,
      }));
  }

  async getEntry(id: string): Promise<EntryFull | undefined> {
    const data = load();
    return data.entries.find((e) => e.id === id);
  }

  async deleteEntry(id: string): Promise<void> {
    const data = load();
    data.entries = data.entries.filter((e) => e.id !== id);
    save(data);
  }
}
