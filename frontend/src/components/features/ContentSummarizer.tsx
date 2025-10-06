import { Loader2, NotebookPen, Wand2 } from 'lucide-react';
import React, { useState } from 'react';
import Card from '../ui/Card';

interface ContentSummarizerProps {
  onSummarize?: (input: string, tone: string, maxWords: number) => Promise<string>;
}

const tonePresets = [
  { id: 'executive', label: 'Executivo' },
  { id: 'technical', label: 'Técnico' },
  { id: 'casual', label: 'Casual' },
];

const ContentSummarizer: React.FC<ContentSummarizerProps> = ({ onSummarize }) => {
  const [input, setInput] = useState('');
  const [tone, setTone] = useState('executive');
  const [maxWords, setMaxWords] = useState(220);
  const [isLoading, setIsLoading] = useState(false);
  const [summary, setSummary] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleSummarize = async () => {
    if (!input.trim() || !onSummarize) return;
    setIsLoading(true);
    setError(null);
    setSummary(null);
    try {
      const result = await onSummarize(input.trim(), tone, maxWords);
      setSummary(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Não foi possível gerar o resumo.');
    } finally {
      setIsLoading(false);
    }
  };

  const disabled = !input.trim() || isLoading;

  return (
    <div className="space-y-6">
      <Card title="Summarizer" description="Transforme briefings longos em entregáveis prontos para stakeholders.">
        <div className="grid gap-6 lg:grid-cols-2">
          <div className="space-y-4">
            <label className="block text-xs font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
              Conteúdo base
            </label>
            <textarea
              value={input}
              onChange={(event) => setInput(event.target.value)}
              placeholder="Cole aqui atas, relatórios ou mensagens extensas."
              rows={12}
              className="w-full resize-none rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-200 dark:border-slate-700 dark:bg-[#0f172a] dark:text-slate-200 dark:focus:border-slate-500 dark:focus:ring-slate-700/40"
            />

            <div className="flex flex-wrap gap-3">
              {tonePresets.map((preset) => (
                <button
                  key={preset.id}
                  type="button"
                  onClick={() => setTone(preset.id)}
                  className={`rounded-full border px-4 py-2 text-xs font-semibold transition ${
                    tone === preset.id
                      ? 'border-slate-900 bg-slate-900 text-white shadow-[0_20px_45px_-35px_rgba(15,23,42,0.8)] dark:border-[#00f0ff] dark:bg-[#00f0ff] dark:text-[#010409]'
                      : 'border-slate-200 bg-white text-slate-600 hover:border-slate-300 dark:border-slate-700 dark:bg-[#0f172a] dark:text-slate-300'
                  }`}
                >
                  {preset.label}
                </button>
              ))}
            </div>

            <div className="flex items-center gap-3">
              <label htmlFor="max-words" className="text-xs uppercase tracking-[0.3em] text-slate-400 dark:text-slate-500">
                Limite de palavras
              </label>
              <input
                id="max-words"
                type="number"
                min={100}
                max={600}
                value={maxWords}
                onChange={(event) => setMaxWords(Number(event.target.value))}
                className="w-24 rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 shadow-sm focus:border-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-200 dark:border-slate-700 dark:bg-[#0f172a] dark:text-slate-200"
              />
            </div>
            <button
              type="button"
              disabled={disabled}
              onClick={handleSummarize}
              className="inline-flex items-center gap-2 rounded-full border border-slate-900 bg-slate-900 px-6 py-2 text-sm font-semibold text-white shadow-[0_20px_45px_-35px_rgba(15,23,42,0.8)] transition disabled:cursor-not-allowed disabled:opacity-60 dark:border-[#00f0ff] dark:bg-[#00f0ff] dark:text-[#010409]"
            >
              {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <NotebookPen className="h-4 w-4" />}
              {isLoading ? 'Gerando resumo...' : 'Gerar resumo'}
            </button>
            <p className="text-[11px] uppercase tracking-[0.3em] text-slate-400 dark:text-slate-500">
              O resumo respeita as diretrizes Kubex e mantém a origem do conteúdo.
            </p>
          </div>

          <div className="flex h-full flex-col rounded-2xl border border-slate-200/80 bg-white/70 p-5 shadow-inner dark:border-slate-800/60 dark:bg-[#0a0f14]/70">
            <div className="flex items-center gap-3 text-xs font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
              <Wand2 className="h-4 w-4" /> Entregável
            </div>
            <div className="mt-4 flex-1 overflow-auto rounded-xl border border-dashed border-slate-200/80 bg-white/70 p-4 text-sm text-slate-600 dark:border-slate-700/80 dark:bg-[#0f172a]/80 dark:text-slate-300">
              {summary && <pre className="whitespace-pre-wrap break-words text-sm leading-relaxed">{summary}</pre>}
              {!summary && !error && !isLoading && (
                <p>O resultado aparecerá aqui em Markdown pronto para compartilhar.</p>
              )}
              {error && (
                <p className="text-sm text-red-500 dark:text-red-400">{error}</p>
              )}
              {isLoading && (
                <div className="flex items-center gap-2 text-sm text-slate-500 dark:text-slate-400">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  Preparando síntese...
                </div>
              )}
            </div>
          </div>
        </div>
      </Card>
    </div>
  );
};

export default ContentSummarizer;
