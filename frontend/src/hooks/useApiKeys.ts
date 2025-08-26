// BYOK storage + optional AES-GCM encryption
// localStorage key: grompt.apikeys.v1

export type Provider = 'openai' | 'chatgpt' | 'anthropic' | 'gemini' | 'deepseek' | 'ollama';

export type ApiKeysVault = {
  openai?: { apiKey: string; org?: string; defaultModel?: string };
  chatgpt?: { apiKey: string; defaultModel?: string }; // alias to OpenAI headers
  anthropic?: { apiKey: string; defaultModel?: string };
  gemini?: { apiKey: string; defaultModel?: string };
  deepseek?: { apiKey: string; defaultModel?: string };
  ollama?: { baseUrl?: string; defaultModel?: string };
};

type StoredEnvelope = {
  v: 1;
  enc: boolean;
  salt?: string; // base64
  iv?: string; // base64
  data: string; // base64 (plaintext JSON if enc=false; ciphertext if enc=true)
};

export type UnlockResult = {
  ok: boolean;
  locked?: boolean;
  vault?: ApiKeysVault;
  error?: string;
  enc: boolean;
};

const STORAGE_KEY = 'grompt.apikeys.v1';

// --- helpers: base64 <-> bytes ---
const enc = new TextEncoder();
const dec = new TextDecoder();

function bytesToBase64(bytes: Uint8Array): string {
  // Browser-safe base64
  let binary = '';
  const len = bytes.byteLength;
  for (let i = 0; i < len; i++) binary += String.fromCharCode(bytes[i]);
  return btoa(binary);
}

function base64ToBytes(base64: string): Uint8Array {
  const binary = atob(base64);
  const len = binary.length;
  const bytes = new Uint8Array(len);
  for (let i = 0; i < len; i++) bytes[i] = binary.charCodeAt(i);
  return bytes;
}

async function deriveAesGcmKey(passphrase: string, salt: Uint8Array) {
  const keyMaterial = await crypto.subtle.importKey(
    'raw',
    enc.encode(passphrase),
    'PBKDF2',
    false,
    ['deriveKey']
  );
  return crypto.subtle.deriveKey(
    {
      name: 'PBKDF2',
      salt,
      iterations: 100_000,
      hash: 'SHA-256'
    },
    keyMaterial,
    { name: 'AES-GCM', length: 256 },
    false,
    ['encrypt', 'decrypt']
  );
}

export function getStoredEnvelope(): StoredEnvelope | null {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw) as StoredEnvelope;
    if (!parsed || parsed.v !== 1 || typeof parsed.enc !== 'boolean' || typeof parsed.data !== 'string') {
      return null;
    }
    return parsed;
  } catch {
    return null;
  }
}

export async function saveVault(vault: ApiKeysVault, passphrase?: string): Promise<void> {
  const json = JSON.stringify(vault);
  if (passphrase && passphrase.trim().length > 0) {
    const salt = crypto.getRandomValues(new Uint8Array(16));
    const iv = crypto.getRandomValues(new Uint8Array(12));
    const key = await deriveAesGcmKey(passphrase, salt);
    const cipher = await crypto.subtle.encrypt({ name: 'AES-GCM', iv }, key, enc.encode(json));
    const envelope: StoredEnvelope = {
      v: 1,
      enc: true,
      salt: bytesToBase64(salt),
      iv: bytesToBase64(iv),
      data: bytesToBase64(new Uint8Array(cipher))
    };
    localStorage.setItem(STORAGE_KEY, JSON.stringify(envelope));
  } else {
    const envelope: StoredEnvelope = {
      v: 1,
      enc: false,
      data: bytesToBase64(enc.encode(json))
    };
    localStorage.setItem(STORAGE_KEY, JSON.stringify(envelope));
  }
}

export async function unlockVault(passphrase?: string): Promise<UnlockResult> {
  const stored = getStoredEnvelope();
  if (!stored) return { ok: true, enc: false, vault: {} };

  if (!stored.enc) {
    try {
      const json = dec.decode(base64ToBytes(stored.data));
      return { ok: true, enc: false, vault: JSON.parse(json) };
    } catch (e) {
      return { ok: false, enc: false, error: 'Dados inválidos no armazenamento' };
    }
  }

  if (!passphrase || passphrase.trim().length === 0) {
    return { ok: true, enc: true, locked: true };
  }

  try {
    const salt = base64ToBytes(stored.salt || '');
    const iv = base64ToBytes(stored.iv || '');
    const key = await deriveAesGcmKey(passphrase, salt);
    const plainBuf = await crypto.subtle.decrypt({ name: 'AES-GCM', iv }, key, base64ToBytes(stored.data));
    const json = dec.decode(new Uint8Array(plainBuf));
    return { ok: true, enc: true, vault: JSON.parse(json) };
  } catch (e) {
    return { ok: false, enc: true, locked: true, error: 'Passphrase incorreta' };
  }
}

export function clearVault() {
  localStorage.removeItem(STORAGE_KEY);
}

export function exportStoredEnvelope(): string {
  const stored = getStoredEnvelope();
  return JSON.stringify(stored ?? { v: 1, enc: false, data: bytesToBase64(enc.encode(JSON.stringify({}))) }, null, 2);
}

export function importStoredEnvelope(json: string) {
  const parsed = JSON.parse(json) as StoredEnvelope;
  if (!parsed || parsed.v !== 1 || typeof parsed.enc !== 'boolean' || typeof parsed.data !== 'string') {
    throw new Error('Formato inválido para importação');
  }
  localStorage.setItem(STORAGE_KEY, JSON.stringify(parsed));
}

export function useApiKeys() {
  // Simple utility wrapper for components
  return {
    STORAGE_KEY,
    getStoredEnvelope,
    save: saveVault,
    unlock: unlockVault,
    clear: clearVault,
    exportJSON: exportStoredEnvelope,
    importJSON: importStoredEnvelope,
  };
}
