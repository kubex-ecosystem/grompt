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
            <label className="block text-xs font-semibold uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
              Conteúdo base
            </label>
            <textarea
              value={input}
              onChange={(event) => setInput(event.target.value)}
              placeholder="Cole aqui atas, relatórios ou mensagens extensas."
              rows={12}
              className="w-full resize-none rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm text-[#475569] shadow-inner transition focus:border-[#06b6d4] focus:outline-none focus:ring-2 focus:ring-[#06b6d4]/20 dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#e5f2f2] dark:focus:border-[#38cde4] dark:focus:ring-[#38cde4]/20"
            />

            <div className="flex flex-wrap gap-3">
              {tonePresets.map((preset) => (
                <button
                  key={preset.id}
                  type="button"
                  onClick={() => setTone(preset.id)}
                  className={`rounded-full border px-4 py-2 text-xs font-semibold transition ${
                    tone === preset.id
                      ? 'border-[#06b6d4] bg-[#06b6d4] text-white shadow-soft-card dark:border-[#06b6d4] dark:bg-[#06b6d4] dark:text-[#0a1523]'
                      : 'border-slate-200 bg-white text-[#475569] hover:border-[#bae6fd] dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#94a3b8]'
                  }`}
                >
                  {preset.label}
                </button>
              ))}
            </div>

            <div className="flex items-center gap-3">
              <label htmlFor="max-words" className="text-xs uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
                Limite de palavras
              </label>
              <input
                id="max-words"
                type="number"
                min={100}
                max={600}
                value={maxWords}
                onChange={(event) => setMaxWords(Number(event.target.value))}
                className="w-24 rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm text-[#475569] shadow-sm focus:border-[#06b6d4] focus:outline-none focus:ring-2 focus:ring-[#06b6d4]/20 dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#e5f2f2]"
              />
            </div>
            <button
              type="button"
              disabled={disabled}
              onClick={handleSummarize}
              className="inline-flex items-center gap-2 rounded-full border border-[#06b6d4] bg-[#06b6d4] px-6 py-2 text-sm font-semibold text-white shadow-soft-card transition hover:bg-[#0891b2] disabled:cursor-not-allowed disabled:opacity-60 focus:outline-none focus:ring-2 focus:ring-[#06b6d4]/30 dark:border-[#06b6d4] dark:bg-[#06b6d4] dark:text-[#0a1523]"
            >
              {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <NotebookPen className="h-4 w-4" />}
              {isLoading ? 'Gerando resumo...' : 'Gerar resumo'}
            </button>
            <p className="text-[11px] uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
              O resumo respeita as diretrizes Kubex e mantém a origem do conteúdo.
            </p>
          </div>

          <div className="flex h-full flex-col rounded-2xl border border-slate-200/80 bg-white/85 p-5 shadow-soft-card dark:border-[#13263a]/70 dark:bg-[#0a1523]/70">
            <div className="flex items-center gap-3 text-xs font-semibold uppercase tracking-[0.3em] text-[#94a3b8] dark:text-[#64748b]">
              <Wand2 className="h-4 w-4" /> Entregável
            </div>
            <div className="mt-4 flex-1 overflow-auto rounded-xl border border-dashed border-slate-200/80 bg-white/90 p-4 text-sm text-[#475569] dark:border-[#13263a]/80 dark:bg-[#0a1523]/75 dark:text-[#e5f2f2]">
              {summary && <pre className="whitespace-pre-wrap break-words text-sm leading-relaxed">{summary}</pre>}
              {!summary && !error && !isLoading && (
                <p>O resultado aparecerá aqui em Markdown pronto para compartilhar.</p>
              )}
              {error && (
                <p className="text-sm text-red-500 dark:text-red-400">{error}</p>
              )}
              {isLoading && (
                <div className="flex items-center gap-2 text-sm text-[#94a3b8] dark:text-[#64748b]">
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
