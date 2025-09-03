export type ProviderKey = 'openai' | 'anthropic' | 'claude' | 'gemini' | 'deepseek' | 'ollama' | 'chatgpt' | 'demo' | string;

type Meta = { key: ProviderKey; label: string; color: string; text: string };

const MAP: Record<string, Meta> = {
  openai:    { key: 'openai',    label: 'OpenAI',    color: '#10a37f', text: '#063a2f' }, // green
  chatgpt:   { key: 'chatgpt',   label: 'ChatGPT',   color: '#10a37f', text: '#063a2f' },
  anthropic: { key: 'anthropic', label: 'Claude',    color: '#efb417', text: '#5a4300' }, // amber
  claude:    { key: 'claude',    label: 'Claude',    color: '#efb417', text: '#5a4300' },
  gemini:    { key: 'gemini',    label: 'Gemini',    color: '#4285F4', text: '#0b3ea8' }, // google blue
  deepseek:  { key: 'deepseek',  label: 'DeepSeek',  color: '#f97316', text: '#7c2d12' }, // orange
  ollama:    { key: 'ollama',    label: 'Ollama',    color: '#9CA3AF', text: '#111827' }, // grey
  demo:      { key: 'demo',      label: 'Demo',      color: '#8b5cf6', text: '#3b0764' }, // violet
};

export function getProviderMeta(provider: string | undefined): Meta {
  if (!provider) return MAP.demo;
  const key = String(provider).toLowerCase();
  if (MAP[key]) return MAP[key];
  if (key === 'anthropic' || key === 'claude') return MAP.anthropic;
  return { key, label: provider, color: '#6b7280', text: '#111827' };
}

