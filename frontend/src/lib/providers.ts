import { ApiKeysVault, Provider } from '../hooks/useApiKeys';

export type ProviderId = Provider;

export const providerHeaders = (provider: ProviderId, vault: ApiKeysVault): Record<string, string> => {
  const headers: Record<string, string> = {};
  switch (provider) {
    case 'openai': {
      const key = vault.openai?.apiKey;
      if (key) headers['Authorization'] = `Bearer ${key}`;
      if (vault.openai?.org) headers['OpenAI-Organization'] = vault.openai.org;
      break;
    }
    case 'chatgpt': {
      const key = vault.chatgpt?.apiKey || vault.openai?.apiKey;
      if (key) headers['Authorization'] = `Bearer ${key}`;
      break;
    }
    case 'anthropic': {
      const key = vault.anthropic?.apiKey;
      if (key) headers['x-api-key'] = key;
      break;
    }
    case 'gemini': {
      const key = vault.gemini?.apiKey;
      if (key) headers['x-goog-api-key'] = key;
      break;
    }
    case 'deepseek': {
      const key = vault.deepseek?.apiKey;
      if (key) headers['Authorization'] = `Bearer ${key}`;
      break;
    }
    case 'ollama':
      // no key header
      break;
  }
  return headers;
};

export function guessProviderFromEndpoint(endpoint: string): ProviderId | undefined {
  if (/\/api\/openai\b/.test(endpoint)) return 'openai';
  if (/\/api\/chatgpt\b/.test(endpoint)) return 'chatgpt';
  if (/\/api\/claude\b/.test(endpoint)) return 'anthropic';
  if (/\/api\/gemini\b/.test(endpoint)) return 'gemini';
  if (/\/api\/deepseek\b/.test(endpoint)) return 'deepseek';
  if (/\/api\/ollama\b/.test(endpoint)) return 'ollama';
  return undefined;
}

export const hasByok = (provider: ProviderId, vault: ApiKeysVault): boolean => {
  switch (provider) {
    case 'openai': return !!vault.openai?.apiKey;
    case 'chatgpt': return !!(vault.chatgpt?.apiKey || vault.openai?.apiKey);
    case 'anthropic': return !!vault.anthropic?.apiKey;
    case 'gemini': return !!vault.gemini?.apiKey;
    case 'deepseek': return !!vault.deepseek?.apiKey;
    case 'ollama': return !!vault.ollama?.baseUrl;
  }
};

// Direct calls (no backend). Note: requires provider CORS support.
export async function directCall(provider: ProviderId, vault: ApiKeysVault, payload: any): Promise<Response> {
  switch (provider) {
    case 'openai':
    case 'chatgpt': {
      const key = (vault.chatgpt?.apiKey || vault.openai?.apiKey) as string;
      const messages = payload.messages || [{ role: 'user', content: payload.prompt }];
      const model = payload.model || vault.chatgpt?.defaultModel || vault.openai?.defaultModel || 'gpt-3.5-turbo';
      return fetch('https://api.openai.com/v1/chat/completions', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${key}`,
          ...(vault.openai?.org ? { 'OpenAI-Organization': vault.openai.org } : {}),
        },
        body: JSON.stringify({ model, messages, temperature: payload.temperature ?? 0.2, max_tokens: payload.max_tokens ?? 256 })
      });
    }
    case 'anthropic': {
      const key = vault.anthropic?.apiKey as string;
      const model = payload.model || vault.anthropic?.defaultModel || 'claude-3-opus-20240229';
      const messages = payload.messages || [{ role: 'user', content: payload.prompt }];
      const content = messages.map((m: any) => ({ role: m.role, content: [{ type: 'text', text: typeof m.content === 'string' ? m.content : m.content?.[0]?.text || '' }] }));
      return fetch('https://api.anthropic.com/v1/messages', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'x-api-key': key,
          'anthropic-version': '2023-06-01',
        },
        body: JSON.stringify({ model, max_tokens: payload.max_tokens ?? 256, messages: content })
      });
    }
    case 'gemini': {
      const key = vault.gemini?.apiKey as string;
      const model = encodeURIComponent(payload.model || vault.gemini?.defaultModel || 'gemini-1.5-flash');
      const url = `https://generativelanguage.googleapis.com/v1beta/models/${model}:generateContent?key=${encodeURIComponent(key)}`;
      const contents = payload.contents || [{ parts: [{ text: payload.prompt }] }];
      return fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ contents })
      });
    }
    case 'deepseek': {
      const key = vault.deepseek?.apiKey as string;
      const messages = payload.messages || [{ role: 'user', content: payload.prompt }];
      const model = payload.model || vault.deepseek?.defaultModel || 'deepseek-chat';
      return fetch('https://api.deepseek.com/chat/completions', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${key}`,
        },
        body: JSON.stringify({ model, messages, temperature: payload.temperature ?? 0.2, max_tokens: payload.max_tokens ?? 256 })
      });
    }
    case 'ollama': {
      const base = vault.ollama?.baseUrl || 'http://localhost:11434';
      const url = `${base.replace(/\/$/, '')}/api/generate`;
      return fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ model: payload.model || vault.ollama?.defaultModel || 'llama3.1', prompt: payload.prompt, stream: false })
      });
    }
  }
}

// Streaming helpers
function textDecoderStream(reader: ReadableStreamDefaultReader<Uint8Array>, onChunk: (s: string) => void) {
  const td = new TextDecoder();
  let buffer = '';
  return (async () => {
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      buffer += td.decode(value, { stream: true });
      onChunk(buffer);
      buffer = '';
    }
  })();
}

export async function streamDirectCall(
  provider: ProviderId,
  vault: ApiKeysVault,
  payload: any,
  onDelta: (text: string) => void
): Promise<void> {
  if (provider === 'openai' || provider === 'chatgpt' || provider === 'deepseek') {
    const isDeepseek = provider === 'deepseek';
    const base = isDeepseek ? 'https://api.deepseek.com' : 'https://api.openai.com';
    const key = isDeepseek ? vault.deepseek?.apiKey : (vault.chatgpt?.apiKey || vault.openai?.apiKey);
    const messages = payload.messages || [{ role: 'user', content: payload.prompt }];
    const model = payload.model || (isDeepseek ? (vault.deepseek?.defaultModel || 'deepseek-chat') : (vault.chatgpt?.defaultModel || vault.openai?.defaultModel || 'gpt-3.5-turbo'));
    const res = await fetch(`${base}/v1/chat/completions`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${key}`,
        ...(provider !== 'deepseek' && vault.openai?.org ? { 'OpenAI-Organization': vault.openai.org } : {}),
      },
      body: JSON.stringify({ model, messages, stream: true, temperature: payload.temperature ?? 0.2 })
    });
    if (!res.ok || !res.body) throw new Error(`HTTP ${res.status}`);
    const reader = res.body.getReader();
    let buf = '';
    const td = new TextDecoder();
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      buf += td.decode(value, { stream: true });
      const lines = buf.split(/\r?\n/);
      buf = lines.pop() || '';
      for (const l of lines) {
        const m = l.match(/^data:\s*(.*)$/);
        if (!m) continue;
        const data = m[1];
        if (data === '[DONE]') continue;
        try {
          const j = JSON.parse(data);
          const delta = j.choices?.[0]?.delta?.content || '';
          if (delta) onDelta(delta);
        } catch {}
      }
    }
    return;
  }
  if (provider === 'anthropic') {
    const key = vault.anthropic?.apiKey as string;
    const model = payload.model || vault.anthropic?.defaultModel || 'claude-3-opus-20240229';
    const messages = payload.messages || [{ role: 'user', content: payload.prompt }];
    const content = messages.map((m: any) => ({ role: m.role, content: [{ type: 'text', text: typeof m.content === 'string' ? m.content : m.content?.[0]?.text || '' }] }));
    const res = await fetch('https://api.anthropic.com/v1/messages', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'x-api-key': key, 'anthropic-version': '2023-06-01' },
      body: JSON.stringify({ model, max_tokens: payload.max_tokens ?? 256, messages: content, stream: true })
    });
    if (!res.ok || !res.body) throw new Error(`HTTP ${res.status}`);
    const reader = res.body.getReader();
    let buf = '';
    const td = new TextDecoder();
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      buf += td.decode(value, { stream: true });
      const lines = buf.split(/\r?\n/);
      buf = lines.pop() || '';
      for (const l of lines) {
        const m = l.match(/^data:\s*(.*)$/);
        if (!m) continue;
        const data = m[1];
        if (data === '[DONE]') continue;
        try {
          const j = JSON.parse(data);
          if (j.type === 'content_block_delta' && j.delta?.text) {
            onDelta(j.delta.text);
          }
        } catch {}
      }
    }
    return;
  }
  if (provider === 'ollama') {
    const base = (vault.ollama?.baseUrl || 'http://localhost:11434').replace(/\/$/, '');
    const res = await fetch(`${base}/api/generate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ model: payload.model || vault.ollama?.defaultModel || 'llama3.1', prompt: payload.prompt, stream: true })
    });
    if (!res.ok || !res.body) throw new Error(`HTTP ${res.status}`);
    const reader = res.body.getReader();
    const td = new TextDecoder();
    let buf = '';
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      buf += td.decode(value, { stream: true });
      const lines = buf.split(/\r?\n/);
      buf = lines.pop() || '';
      for (const l of lines) {
        try {
          const j = JSON.parse(l);
          if (j.response) onDelta(j.response);
        } catch {}
      }
    }
    return;
  }
  throw new Error('Streaming não suportado para este provider');
}

export async function testProvider(
  provider: ProviderId,
  vault: ApiKeysVault,
  baseURL = ''
): Promise<{ ok: boolean; status: number; data?: any; error?: string }>
{
  // Prefer BYOK direct call when available to avoid exposing keys to backend
  if (hasByok(provider, vault)) {
    try {
      const payload: any = { prompt: 'ping' };
      if (provider === 'openai') payload.model = 'gpt-3.5-turbo';
      if (provider === 'chatgpt') payload.model = 'gpt-4o-mini';
      if (provider === 'gemini') payload.model = 'gemini-1.5-pro';
      if (provider === 'deepseek') payload.model = 'deepseek-chat';
      if (provider === 'ollama') payload.model = 'llama3.1';
      const res = await directCall(provider, vault, payload);
      const contentType = res.headers.get('content-type') || '';
      const data = contentType.includes('application/json') ? await res.json() : await res.text();
      return { ok: res.ok, status: res.status, data };
    } catch (e: any) {
      return { ok: false, status: 0, error: e?.message || 'Network error' };
    }
  }

  // Fallback to backend (no key headers injected)
  let endpoint = '';
  let body: Record<string, any> = {};
  switch (provider) {
    case 'openai': endpoint = '/api/openai'; body = { prompt: 'ping', model: 'gpt-3.5-turbo' }; break;
    case 'chatgpt': endpoint = '/api/chatgpt'; body = { prompt: 'ping', model: 'gpt-4o-mini' }; break;
    case 'anthropic': endpoint = '/api/claude'; body = { prompt: 'ping' }; break;
    case 'gemini': endpoint = '/api/gemini'; body = { prompt: 'ping', model: 'gemini-1.5-pro' }; break;
    case 'deepseek': endpoint = '/api/deepseek'; body = { prompt: 'ping', model: 'deepseek-chat' }; break;
    case 'ollama': endpoint = '/api/ollama'; body = { prompt: 'ping', model: 'llama3.1' }; break;
  }
  const url = `${baseURL}${endpoint}`;
  try {
    const res = await fetch(url, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) });
    const contentType = res.headers.get('content-type') || '';
    const data = contentType.includes('application/json') ? await res.json() : await res.text();
    return { ok: res.ok, status: res.status, data };
  } catch (e: any) {
    // Provável CORS quando BYOK tenta direto e o provedor bloqueia
    return { ok: false, status: 0, error: 'Falha de rede possivelmente por CORS no navegador. Sua chave permanece local e segura.' };
  }
}
