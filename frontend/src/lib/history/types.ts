export type ProviderId = 'claude' | 'openai' | 'deepseek' | 'ollama' | 'gemini' | 'chatgpt' | 'demo' | string;

export type HistoryStatus = 'ok' | 'error' | 'partial' | 'draft';

export interface Session {
  id: string;
  name: string;
  agentId?: string;
  tags?: string[];
  createdAt: number; // epoch ms
  updatedAt: number; // epoch ms
}

export interface EntryMeta {
  id: string;
  sessionId: string;
  provider: ProviderId;
  model?: string;
  promptPreview: string;
  status: HistoryStatus;
  createdAt: number;
  updatedAt: number;
  tokenCounts?: { input?: number; output?: number };
}

export interface EntryFull extends EntryMeta {
  params?: Record<string, any>;
  requestBlobId?: string;
  responseBlobId?: string;
  // Convenience fields when adapter chooses to keep small payloads inline
  requestText?: string;
  responseText?: string;
  error?: string;
  ideas?: { id: number; text: string }[];
}

export interface SaveEntryInput {
  sessionId?: string;
  sessionName?: string; // if no sessionId, will create/find by name
  provider: ProviderId;
  model?: string;
  params?: Record<string, any>;
  requestText?: string;
  responseText?: string;
  status?: HistoryStatus;
  error?: string;
  ideas?: { id: number; text: string }[];
}

export interface IHistoryAdapter {
  init(): Promise<void>;
  ensureDefaultSession(name?: string): Promise<Session>;
  createSession(name: string, agentId?: string): Promise<Session>;
  listSessions(): Promise<Session[]>;
  saveEntry(input: SaveEntryInput): Promise<EntryMeta>;
  listEntries(sessionId: string, opts?: { limit?: number; offset?: number }): Promise<EntryMeta[]>;
  getEntry(id: string): Promise<EntryFull | undefined>;
  deleteEntry(id: string): Promise<void>;
}
