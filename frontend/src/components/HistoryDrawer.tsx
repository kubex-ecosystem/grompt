import { Clock, FolderPlus, History as HistoryIcon, Loader2, X } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { history } from '../lib/history/store';
import { EntryFull, EntryMeta, Session } from '../lib/history/types';
import ProviderBadge from './ProviderBadge';

type Props = {
  isOpen: boolean;
  onClose: () => void;
};

export default function HistoryDrawer({ isOpen, onClose }: Props) {
  const { t } = useTranslation();
  const [sessions, setSessions] = useState<Session[]>([]);
  const [selectedSessionId, setSelectedSessionId] = useState<string | null>(null);
  const [entries, setEntries] = useState<EntryMeta[]>([]);
  const [selectedEntry, setSelectedEntry] = useState<EntryFull | null>(null);
  const [busy, setBusy] = useState(false);
  const [migrated, setMigrated] = useState<null | { created: number }>(null);
  const [query, setQuery] = useState('');
  const [providerFilter, setProviderFilter] = useState<string>('');
  const [modelFilter, setModelFilter] = useState<string>('');
  const [showTech, setShowTech] = useState(false);
  const [dismissed, setDismissed] = useState<boolean>(false);

  const loadSessions = async () => {
    await history.init();
    const s = await history.listSessions();
    setSessions(s);
    if (!selectedSessionId && s[0]) {
      setSelectedSessionId(s[0].id);
    }
  };

  const loadEntries = async (sessionId: string) => {
    const es = await history.listEntries(sessionId, { limit: 200 });
    setEntries(es);
    setSelectedEntry(null);
  };

  useEffect(() => {
    if (!isOpen) return;
    (async () => {
      setBusy(true);
      try { setDismissed(localStorage.getItem('grompt.history.migrationDismissed') === '1'); } catch {}
      await loadSessions();
      // Try migration once per open if we haven't before
      if (!migrated) {
        const res = await history.migrateFromLocalStorage();
        if (res) {
          setMigrated(res);
          await loadSessions();
        }
      }
      setBusy(false);
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isOpen]);

  useEffect(() => {
    if (selectedSessionId) loadEntries(selectedSessionId);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedSessionId]);

  const theme = useMemo(() => ({
    bg: 'bg-gray-900',
    cardBg: 'bg-gray-800',
    text: 'text-white',
    border: 'border-gray-700',
    button: 'bg-blue-600 hover:bg-blue-700',
    buttonSecondary: 'bg-gray-700 hover:bg-gray-600',
  }), []);

  const providerOptions = useMemo(() => {
    const set = new Set<string>();
    entries.forEach(e => e.provider && set.add(String(e.provider)));
    return Array.from(set).sort();
  }, [entries]);

  const modelOptions = useMemo(() => {
    const set = new Set<string>();
    entries.forEach(e => {
      if (providerFilter && e.provider !== providerFilter) return;
      if (e.model) set.add(String(e.model));
    });
    return Array.from(set).sort();
  }, [entries, providerFilter]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex">
      <div className="flex-1 bg-black/40" onClick={onClose} />
      <div className={`w-full max-w-6xl h-full ${theme.bg} ${theme.text} shadow-xl p-0 overflow-hidden border-l ${theme.border}`}>
        {/* Header */}
        <div className={`flex items-center justify-between p-4 border-b ${theme.border}`}>
          <div className="flex items-center gap-2">
            <HistoryIcon className="h-5 w-5" />
            <h2 className="font-semibold">{t('history.title')}</h2>
            {busy && <Loader2 className="h-4 w-4 animate-spin" />}
            {migrated && (
              <span className="text-xs text-green-400 ml-2">{t('history.migrated', { count: migrated.created })}</span>
            )}
          </div>
          <button title={t('common.close')} onClick={onClose} className="p-2 rounded hover:bg-gray-800"><X className="h-4 w-4" /></button>
        </div>

        {/* Migrated banner (one-shot) */}
        {migrated && !dismissed && (
          <div className="px-4 py-2 text-xs bg-blue-900 text-blue-100 border-b border-blue-700 flex items-center justify-between">
            <span>{t('history.migratedBanner')}</span>
            <button
              className="px-2 py-1 bg-blue-700 hover:bg-blue-600 rounded"
              onClick={() => { try { localStorage.setItem('grompt.history.migrationDismissed', '1'); } catch {}; setDismissed(true); }}
            >{t('history.gotIt')}</button>
          </div>
        )}

        {/* Content */}
        <div className="flex h-full">
          {/* Sessions */}
          <aside className={`w-64 h-full border-r ${theme.border} p-3`}>
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm text-gray-300">{t('history.sessions')}</span>
              <div className="flex items-center gap-2">
                {selectedSessionId && (
                  <>
                    <button
                      className={`p-1 rounded ${theme.buttonSecondary}`}
                      title={t('history.clearSession')}
                      onClick={async () => {
                        if (!confirm(t('history.clearConfirm'))) return;
                        await history.clearSession(selectedSessionId);
                        await loadEntries(selectedSessionId);
                      }}
                    >
                      🧹
                    </button>
                    <button
                      className={`p-1 rounded ${theme.buttonSecondary}`}
                      title={t('history.deleteSession')}
                      onClick={async () => {
                        if (!confirm(t('history.deleteConfirm'))) return;
                        const sid = selectedSessionId;
                        await history.deleteSession(sid);
                        await loadSessions();
                        const first = (await history.listSessions())[0];
                        setSelectedSessionId(first ? first.id : null);
                        if (first) await loadEntries(first.id);
                        else setEntries([]);
                      }}
                    >
                      🗑️
                    </button>
                  </>
                )}
                <button
                  className={`p-1 rounded ${theme.buttonSecondary}`}
                  title={t('history.newSession')}
                  onClick={async () => {
                    const name = prompt(t('history.promptName'));
                    if (!name) return;
                    const s = await history.createSession(name);
                    await loadSessions();
                    setSelectedSessionId(s.id);
                    await loadEntries(s.id);
                  }}
                >
                  <FolderPlus className="h-4 w-4" />
                </button>
              </div>
            </div>
            <div className="space-y-1 overflow-auto max-h-[calc(100vh-8rem)]">
              {sessions.length === 0 && (
                <div className="text-xs text-gray-400">{t('history.emptySessions')}</div>
              )}
              {sessions.map((s) => (
                <button
                  key={s.id}
                  onClick={() => setSelectedSessionId(s.id)}
                  className={`w-full text-left px-2 py-2 rounded border ${theme.border} ${selectedSessionId === s.id ? 'bg-gray-800' : ''}`}
                >
                  <div className="text-sm font-medium truncate">{s.name}</div>
                  <div className="text-[10px] text-gray-400 flex items-center gap-1"><Clock className="h-3 w-3" /> {new Date(s.updatedAt).toLocaleString()}</div>
                </button>
              ))}
            </div>
          </aside>

          {/* Entries list */}
          <section className="flex-1 h-full flex">
            <div className="w-96 border-r border-gray-800 p-3 overflow-auto">
              <div className="flex items-center gap-2 mb-2">
                <div className="text-sm text-gray-300">{t('history.entries')}</div>
                <input
                  placeholder={t('history.searchPlaceholder')}
                  value={query}
                  onChange={(e) => setQuery(e.target.value)}
                  className="flex-1 px-2 py-1 text-sm rounded bg-gray-900 border border-gray-700"
                />
                {query && (
                  <button
                    className={`px-2 py-1 text-xs rounded ${theme.buttonSecondary}`}
                    onClick={() => setQuery('')}
                    title={t('history.clearSearch')}
                  >
                    ✕
                  </button>
                )}
              </div>
              <div className="flex items-center gap-2 mb-3">
                <select
                  title={t('history.filterProvider')}
                  value={providerFilter}
                  onChange={(e) => { setProviderFilter(e.target.value); setModelFilter(''); }}
                  className="px-2 py-1 text-sm rounded bg-gray-900 border border-gray-700"
                >
                  <option value="">{t('history.allProviders')}</option>
                  {providerOptions.map(p => (
                    <option key={p} value={p}>{p}</option>
                  ))}
                </select>
                <select
                  title={t('history.filterModel')}
                  value={modelFilter}
                  onChange={(e) => setModelFilter(e.target.value)}
                  className="px-2 py-1 text-sm rounded bg-gray-900 border border-gray-700"
                  disabled={!providerFilter}
                >
                  <option value="">{t('history.allModels')}</option>
                  {modelOptions.map(m => (
                    <option key={m} value={m}>{m}</option>
                  ))}
                </select>
              </div>
              {entries.length === 0 && (
                <div className="text-xs text-gray-400">{t('history.emptyEntries')}</div>
              )}
              <div className="space-y-2">
                {entries
                  .filter((e) => {
                    if (providerFilter && e.provider !== providerFilter) return false;
                    if (modelFilter && e.model !== modelFilter) return false;
                    const q = query.trim().toLowerCase();
                    if (!q) return true;
                    return (
                      (e.promptPreview || '').toLowerCase().includes(q) ||
                      (e.provider || '').toLowerCase().includes(q) ||
                      (e.model || '').toLowerCase().includes(q)
                    );
                  })
                  .map((e) => (
                    <button
                      key={e.id}
                      onClick={async () => setSelectedEntry((await history.getEntry(e.id)) || null)}
                      className={`w-full text-left p-2 rounded border ${theme.border} hover:bg-gray-800`}
                      title={`${e.provider} ${e.model || ''}`}
                    >
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <ProviderBadge provider={e.provider} />
                          {e.model && <span className="text-[10px] text-gray-400">{e.model}</span>}
                        </div>
                        <div className="text-[10px] text-gray-500">{new Date(e.createdAt).toLocaleString()}</div>
                      </div>
                      <div className="text-sm truncate mt-1">{e.promptPreview}</div>
                    </button>
                  ))}
              </div>
            </div>

            {/* Detail */}
            <div className="flex-1 p-3 overflow-auto">
              {!selectedEntry ? (
                <div className="text-sm text-gray-400">{t('history.selectToView')}</div>
              ) : (
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <ProviderBadge provider={selectedEntry.provider} showLabel={true} />
                      {selectedEntry.model && <span className="text-xs text-gray-400">{selectedEntry.model}</span>}
                    </div>
                    <div className="flex items-center">
                      <button
                        className={`px-2 py-1 rounded text-xs border ${theme.border} hover:bg-gray-800`}
                        onClick={() => setShowTech(v => !v)}
                        title={t('history.toggleTech')}
                      >
                        {showTech ? t('history.hideTech') : t('history.showTech')}
                      </button>
                      <button
                        className={`px-2 py-1 rounded text-xs border ${theme.border} hover:bg-gray-800`}
                        onClick={() => {
                          try {
                            window.dispatchEvent(new CustomEvent('grompt:history-load', { detail: selectedEntry } as any));
                            onClose();
                          } catch { }
                        }}
                        title={t('history.loadTitle')}
                      >
                        {t('history.load')}
                      </button>
                      <button
                        className={`px-2 py-1 rounded text-xs border ${theme.border} hover:bg-gray-800`}
                        onClick={() => {
                          try {
                            window.dispatchEvent(new CustomEvent('grompt:history-load-draft', { detail: selectedEntry } as any));
                            onClose();
                          } catch { }
                        }}
                        title={t('history.draftTitle')}
                      >
                        {t('history.draft')}
                      </button>
                      <button
                        className={`px-2 py-1 rounded text-xs border ${theme.border} hover:bg-gray-800`}
                        onClick={() => {
                          try {
                            window.dispatchEvent(new CustomEvent('grompt:history-reexec', { detail: selectedEntry } as any));
                            onClose();
                          } catch { }
                        }}
                        title={t('history.reexecTitle')}
                      >
                        {t('history.reexec')}
                      </button>
                      <button
                        className={`px-2 py-1 rounded text-xs border ${theme.border} hover:bg-gray-800`}
                        onClick={() => {
                          try {
                            window.dispatchEvent(new CustomEvent('grompt:history-edit-reexec', { detail: selectedEntry } as any));
                            onClose();
                          } catch { }
                        }}
                        title={t('history.editReexecTitle')}
                      >
                        {t('history.editReexec')}
                      </button>
                      <button
                        className={`px-2 py-1 rounded text-xs border ${theme.border} hover:bg-gray-800`}
                        onClick={() => {
                          const blob = new Blob([JSON.stringify(selectedEntry, null, 2)], { type: 'application/json' });
                          const url = URL.createObjectURL(blob);
                          const a = document.createElement('a');
                          a.href = url;
                          a.download = `grompt-entry-${selectedEntry.id}.json`;
                          a.click();
                          URL.revokeObjectURL(url);
                        }}
                        title={t('history.exportTitle')}
                      >
                        {/* <Download className="h-3 w-3 inline mr-1" /> */} {t('history.export')}
                      </button>
                      <button
                        className={`px-2 py-1 rounded text-xs border ${theme.border} hover:bg-gray-800 text-red-400`}
                        onClick={async () => {
                          if (!confirm(t('history.deleteConfirmEntry'))) return;
                          await history.deleteEntry(selectedEntry.id);
                          if (selectedSessionId) await loadEntries(selectedSessionId);
                          setSelectedEntry(null);
                        }}
                        title={t('history.deleteTitle')}
                      >
                        {/* <Trash2 className="h-3 w-3 inline mr-1" /> */} {t('history.delete')}
                      </button>
                    </div>
                  </div>
                  <div>
                    <div className="text-xs text-gray-400 mb-1">{t('history.ideas')}</div>
                    {Array.isArray(selectedEntry.ideas) && selectedEntry.ideas.length > 0 ? (
                      <ul className="list-disc pl-5 text-xs space-y-1">
                        {selectedEntry.ideas.map((it, idx) => (
                          <li key={idx}>{it.text}</li>
                        ))}
                      </ul>
                    ) : (
                      <div className="text-xs text-gray-600">{t('history.ideasUnavailable')}</div>
                    )}
                  </div>
                  <div>
                    <div className="text-xs text-gray-400 mb-1">{t('history.result')}</div>
                    <pre className="whitespace-pre-wrap text-xs p-2 rounded border border-gray-800 bg-gray-900 max-h-[50vh] overflow-auto">{selectedEntry.responseText || '(payload externo)'}</pre>
                  </div>
                  {showTech && (
                    <div className="space-y-2">
                      <div className="text-xs text-gray-400">{t('history.technicalDetails')}</div>
                      <pre className="whitespace-pre-wrap text-[10px] p-2 rounded border border-gray-800 bg-gray-950 max-h-48 overflow-auto">{selectedEntry.requestText || '(payload externo)'}</pre>
                    </div>
                  )}
                </div>
              )}
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}
