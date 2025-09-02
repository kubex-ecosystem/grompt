'use client';

import { CheckCircle, Download, KeyRound, Lock, PlugZap, Trash2, Unlock, Upload, XCircle } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { ApiKeysVault, clearVault, exportStoredEnvelope, importStoredEnvelope, saveVault, unlockVault } from '../hooks/useApiKeys';
import { ProviderId, testProvider } from '../lib/providers';


var _window: Window | undefined = typeof window !== 'undefined' ? window : undefined;

type Props = {
  isOpen: boolean;
  onClose: () => void;
};

function getProcessPort() {
  if (_window && _window.location) {
    const m = _window.location.href.match(/^(https?:\/\/[^\/]+)/);
    if (m) {
      const url = new URL(m[1]);
      return url.port;
    }
  }
  return `${process.env.PORT || ''}`;
}

function getProcessHost() {
  if (_window && _window.location) {
    const m = _window.location.href.match(/^(https?:\/\/[^\/]+)/);
    if (m) {
      return m[1];
    }
  }
  return `${process.env.HOST || '127.0.0.1'}`;
}

function getProcessPath() {
  if (_window && _window.location) {
    const m = _window.location.href.match(/^(https?:\/\/[^\/]+)/);
    if (m) {
      const url = new URL(m[1]);
      return url.pathname;
    }
  }
  return `${process.env.PATH || ''}`;
}

function getBaseURL() {
  const port = getProcessPort();
  const host = getProcessHost();
  const path = getProcessPath();
  const baseUrlParts = [host, `:${port}`, path].filter(Boolean);
  try {
    const sanitizedUrl = new URL(baseUrlParts.join('').replaceAll('//', '/'));
    return sanitizedUrl.href;
  } catch {
    return '';
  }
}

export default function ApiKeysDrawer({ isOpen, onClose }: Props) {
  const [passphrase, setPassphrase] = useState('');
  const [locked, setLocked] = useState<boolean>(false);
  const [enc, setEnc] = useState<boolean>(false);
  const [vault, setVault] = useState<ApiKeysVault>({});
  const [busy, setBusy] = useState(false);
  const [testResult, setTestResult] = useState<Record<string, { status: 'idle' | 'ok' | 'fail'; detail?: string }>>({});

  const baseURL = useMemo(() => getBaseURL(), []);

  useEffect(() => {
    if (!isOpen) return;
    (async () => {
      const res = await unlockVault();
      setEnc(res.enc);
      if (res.locked) {
        setLocked(true);
        setVault({});
      } else if (res.ok && res.vault) {
        setLocked(false);
        setVault(res.vault);
      } else {
        setVault({});
      }
    })();
  }, [isOpen]);

  const handleUnlock = async () => {
    const res = await unlockVault(passphrase);
    setEnc(res.enc);
    if (res.ok && res.vault) {
      setLocked(false);
      setVault(res.vault);
      try { window.dispatchEvent(new CustomEvent('grompt:apikeys-updated')); } catch { }
    } else {
      setLocked(true);
      alert(res.error || 'Falha ao desbloquear');
    }
  };

  const handleSave = async () => {
    setBusy(true);
    try {
      await saveVault(vault, passphrase.trim() ? passphrase : undefined);
      const res = await unlockVault(passphrase.trim() ? passphrase : undefined);
      setEnc(res.enc);
      setLocked(!!res.locked);
      alert('Chaves salvas no navegador.');
      try { window.dispatchEvent(new CustomEvent('grompt:apikeys-updated')); } catch { }
    } catch (e) {
      alert('Erro ao salvar chaves');
    } finally {
      setBusy(false);
    }
  };

  const handleClear = async () => {
    if (!confirm('Limpar chaves salvas no navegador?')) return;
    clearVault();
    setVault({});
    setLocked(false);
    setEnc(false);
    setPassphrase('');
    try { window.dispatchEvent(new CustomEvent('grompt:apikeys-updated')); } catch { }
  };

  const handleExport = () => {
    const blob = new Blob([exportStoredEnvelope()], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'grompt.apikeys.v1.json';
    a.click();
    URL.revokeObjectURL(url);
  };

  const handleImport = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = () => {
      try {
        importStoredEnvelope(String(reader.result));
        alert('Importado com sucesso. Se criptografado, informe a passphrase e Desbloqueie.');
        try { window.dispatchEvent(new CustomEvent('grompt:apikeys-updated')); } catch { }
      } catch (err) {
        alert('JSON inválido para importação');
      }
    };
    reader.readAsText(file);
  };

  const test = async (provider: ProviderId) => {
    setTestResult((s) => ({ ...s, [provider]: { status: 'idle' } }));
    const res = await testProvider(provider, vault, baseURL);
    const humanDetail = (() => {
      if (res.ok) {
        return typeof res.data === 'string' ? String(res.data).slice(0, 200) : JSON.stringify(res.data).slice(0, 200);
      }
      // Non-ok: try to extract provider message
      if (typeof res.data === 'object' && res.data) {
        const j: any = res.data;
        const msg = j.error?.message || j.message || j.error || JSON.stringify(j);
        return String(msg).slice(0, 240);
      }
      if (res.status === 0) return 'Possível bloqueio CORS pelo navegador. Sua chave permanece local.';
      return res.error || `HTTP ${res.status}`;
    })();

    setTestResult((s) => ({
      ...s,
      [provider]: res.ok
        ? { status: 'ok', detail: humanDetail }
        : { status: 'fail', detail: humanDetail }
    }));
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex">
      <div className="flex-1 bg-black/40" onClick={onClose} />
      <div className="w-full max-w-md h-full bg-gray-900 text-white shadow-xl p-5 overflow-y-auto">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <KeyRound size={20} />
            <h2 className="text-lg font-semibold">Gerenciar Chaves (BYOK)</h2>
          </div>
          <button onClick={onClose} className="text-gray-400 hover:text-white">✕</button>
        </div>

        <p className="text-xs text-gray-400 mb-4">As chaves ficam no seu navegador; não armazenamos no servidor.</p>

        {/* Passphrase / lock */}
        <div className="mb-4">
          <label className="block text-sm text-gray-300 mb-1">Passphrase (opcional para criptografar)</label>
          <div className="flex gap-2">
            <input
              type="password"
              className="flex-1 px-3 py-2 bg-gray-800 border border-gray-700 rounded"
              placeholder="••••••••"
              value={passphrase}
              onChange={(e) => setPassphrase(e.target.value)}
            />
            {enc && locked ? (
              <button onClick={handleUnlock} className="px-3 py-2 bg-blue-600 rounded flex items-center gap-1"><Unlock size={16} /> Desbloquear</button>
            ) : (
              <div className="px-3 py-2 border border-gray-700 rounded flex items-center gap-1 text-gray-300">
                {enc ? <Lock size={16} /> : <Unlock size={16} />} {enc ? 'Criptografado' : 'Sem criptografia'}
              </div>
            )}
          </div>
        </div>

        {/* Providers fields */}
        <div className="space-y-4">
          <section className="p-3 bg-gray-800 rounded border border-gray-700">
            <h3 className="font-medium mb-2">OpenAI</h3>
            <input
              type="text" placeholder="sk-..."
              className="w-full mb-2 px-3 py-2 bg-gray-900 border border-gray-700 rounded"
              value={vault.openai?.apiKey || ''}
              onChange={(e) => setVault({ ...vault, openai: { ...(vault.openai || {}), apiKey: e.target.value } })}
            />
            <input
              type="text" placeholder="Default model (ex.: gpt-4o-mini)"
              className="w-full mb-2 px-3 py-2 bg-gray-900 border border-gray-700 rounded"
              value={(vault.openai?.defaultModel || '') as string}
              onChange={(e) =>
                setVault({
                  ...vault,
                  openai: {
                    apiKey: vault.openai?.apiKey ?? '',
                    org: vault.openai?.org,
                    defaultModel: e.target.value
                  }
                })
              }
            />
            <input
              type="text" placeholder="OpenAI-Organization (opcional)"
              className="w-full mb-2 px-3 py-2 bg-gray-900 border border-gray-700 rounded"
              value={vault.openai?.org || ''}
              onChange={(e) => setVault({ ...vault, openai: { apiKey: vault.openai?.apiKey ?? '', org: e.target.value, defaultModel: vault.openai?.defaultModel } })}
            />
            <button onClick={() => test('openai')} className="text-sm px-3 py-1 bg-gray-700 rounded flex items-center gap-2">
              <PlugZap size={16} /> Testar Conexão
              {testResult['openai']?.status === 'ok' && <CheckCircle size={16} className="text-green-400" />}
              {testResult['openai']?.status === 'fail' && <XCircle size={16} className="text-red-400" />}
            </button>
            {testResult['openai']?.status === 'fail' && (
              <p className="mt-2 text-xs text-yellow-400">{testResult['openai']?.detail || 'Falha na conexão. Se for CORS, use proxy/backend sem enviar a chave.'}</p>
            )}
          </section>

          <section className="p-3 bg-gray-800 rounded border border-gray-700">
            <h3 className="font-medium mb-2">Anthropic (Claude)</h3>
            <input
              type="text" placeholder="sk-ant-..."
              className="w-full mb-2 px-3 py-2 bg-gray-900 border border-gray-700 rounded"
              value={vault.anthropic?.apiKey || ''}
              onChange={(e) => setVault({ ...vault, anthropic: { apiKey: e.target.value } })}
            />
            <input
              type="text" placeholder="Default model (ex.: claude-3-haiku)"
              className="w-full mb-2 px-3 py-2 bg-gray-900 border border-gray-700 rounded"
              value={vault.anthropic?.defaultModel || ''}
              onChange={(e) => setVault({ ...vault, anthropic: { apiKey: vault.anthropic?.apiKey ?? '', defaultModel: e.target.value } })}
            />
            <button onClick={() => test('anthropic')} className="text-sm px-3 py-1 bg-gray-700 rounded flex items-center gap-2">
              <PlugZap size={16} /> Testar Conexão
              {testResult['anthropic']?.status === 'ok' && <CheckCircle size={16} className="text-green-400" />}
              {testResult['anthropic']?.status === 'fail' && <XCircle size={16} className="text-red-400" />}
            </button>
            {testResult['anthropic']?.status === 'fail' && (
              <p className="mt-2 text-xs text-yellow-400">{testResult['anthropic']?.detail || 'Falha na conexão. Se for CORS, use proxy/backend sem enviar a chave.'}</p>
            )}
          </section>

          <section className="p-3 bg-gray-800 rounded border border-gray-700">
            <h3 className="font-medium mb-2">Gemini</h3>
            <input
              type="text" placeholder="AIza..."
              className="w-full mb-2 px-3 py-2 bg-gray-900 border border-gray-700 rounded"
              value={vault.gemini?.apiKey || ''}
              onChange={(e) => setVault({ ...vault, gemini: { apiKey: e.target.value } })}
            />
            <input
              type="text" placeholder="Default model (ex.: gemini-2.0-flash)"
              className="w-full mb-2 px-3 py-2 bg-gray-900 border border-gray-700 rounded"
              value={vault.gemini?.defaultModel || ''}
              onChange={(e) => setVault({ ...vault, gemini: { apiKey: vault.gemini?.apiKey ?? '', defaultModel: e.target.value } })}
            />
            <button onClick={() => test('gemini')} className="text-sm px-3 py-1 bg-gray-700 rounded flex items-center gap-2">
              <PlugZap size={16} /> Testar Conexão
              {testResult['gemini']?.status === 'ok' && <CheckCircle size={16} className="text-green-400" />}
              {testResult['gemini']?.status === 'fail' && <XCircle size={16} className="text-red-400" />}
            </button>
            {testResult['gemini']?.status === 'fail' && (
              <p className="mt-2 text-xs text-yellow-400">{testResult['gemini']?.detail || 'Falha na conexão. Se for CORS, use proxy/backend sem enviar a chave.'}</p>
            )}
          </section>

          <section className="p-3 bg-gray-800 rounded border border-gray-700">
            <h3 className="font-medium mb-2">DeepSeek</h3>
            <input
              type="text" placeholder="sk-..."
              className="w-full mb-2 px-3 py-2 bg-gray-900 border border-gray-700 rounded"
              value={vault.deepseek?.apiKey || ''}
              onChange={(e) => setVault({ ...vault, deepseek: { apiKey: e.target.value } })}
            />
            <input
              type="text" placeholder="Default model (ex.: deepseek-chat)"
              className="w-full mb-2 px-3 py-2 bg-gray-900 border border-gray-700 rounded"
              value={vault.deepseek?.defaultModel || ''}
              onChange={(e) => setVault({ ...vault, deepseek: { apiKey: vault.deepseek?.apiKey ?? '', defaultModel: e.target.value } })}
            />
            <button onClick={() => test('deepseek')} className="text-sm px-3 py-1 bg-gray-700 rounded flex items-center gap-2">
              <PlugZap size={16} /> Testar Conexão
              {testResult['deepseek']?.status === 'ok' && <CheckCircle size={16} className="text-green-400" />}
              {testResult['deepseek']?.status === 'fail' && <XCircle size={16} className="text-red-400" />}
            </button>
            {testResult['deepseek']?.status === 'fail' && (
              <p className="mt-2 text-xs text-yellow-400">{testResult['deepseek']?.detail || 'Falha na conexão. Se for CORS, use proxy/backend sem enviar a chave.'}</p>
            )}
          </section>

          <section className="p-3 bg-gray-800 rounded border border-gray-700">
            <h3 className="font-medium mb-2">Ollama</h3>
            <input
              type="text" placeholder="http://localhost:11434"
              className="w-full mb-2 px-3 py-2 bg-gray-900 border border-gray-700 rounded"
              value={vault.ollama?.baseUrl || ''}
              onChange={(e) => setVault({ ...vault, ollama: { baseUrl: e.target.value } })}
            />
            <input
              type="text" placeholder="Default model (ex.: llama3.1)"
              className="w-full mb-2 px-3 py-2 bg-gray-900 border border-gray-700 rounded"
              value={vault.ollama?.defaultModel || ''}
              onChange={(e) => setVault({ ...vault, ollama: { ...(vault.ollama || {}), defaultModel: e.target.value } })}
            />
            <button onClick={() => test('ollama')} className="text-sm px-3 py-1 bg-gray-700 rounded flex items-center gap-2">
              <PlugZap size={16} /> Testar Conexão
              {testResult['ollama']?.status === 'ok' && <CheckCircle size={16} className="text-green-400" />}
              {testResult['ollama']?.status === 'fail' && <XCircle size={16} className="text-red-400" />}
            </button>
            {testResult['ollama']?.status === 'fail' && (
              <p className="mt-2 text-xs text-yellow-400">{testResult['ollama']?.detail || 'Falha na conexão. Verifique URL do Ollama e CORS.'}</p>
            )}
          </section>
        </div>

        {/* Actions */}
        <div className="mt-6 flex flex-wrap gap-2">
          <button disabled={busy} onClick={handleSave} className="px-3 py-2 bg-blue-600 rounded flex items-center gap-2 disabled:opacity-60">
            <Lock size={16} /> Salvar (local)
          </button>
          <label className="px-3 py-2 bg-gray-700 rounded flex items-center gap-2 cursor-pointer">
            <Upload size={16} /> Importar JSON
            <input type="file" accept="application/json" className="hidden" onChange={handleImport} />
          </label>
          <button onClick={handleExport} className="px-3 py-2 bg-gray-700 rounded flex items-center gap-2">
            <Download size={16} /> Exportar JSON
          </button>
          <button onClick={handleClear} className="px-3 py-2 bg-red-700 rounded flex items-center gap-2">
            <Trash2 size={16} /> Limpar
          </button>
        </div>

        <div className="mt-4 text-xs text-gray-400">
          Headers por provider serão injetados automaticamente nas requisições.
        </div>
      </div>
    </div>
  );
}
