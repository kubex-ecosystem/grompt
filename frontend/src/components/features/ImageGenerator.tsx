import { Image as ImageIcon, Loader2, Sparkles } from 'lucide-react';
import React, { useState } from 'react';
import Card from '../ui/Card';

interface ImageGeneratorProps {
  onCraftPrompt?: (payload: { subject: string; mood: string; style: string; details: string }) => Promise<string>;
}

const moods = ['Vibrante', 'Minimalista', 'Futurista', 'Orgânico'];
const styles = ['Ilustração digital', 'Fotorrealista', 'Flat design', 'Isométrico'];

const ImageGenerator: React.FC<ImageGeneratorProps> = ({ onCraftPrompt }) => {
  const { t } = useTranslation();
  const [subject, setSubject] = useState('');
  const [mood, setMood] = useState(moods[0]);
  const [style, setStyle] = useState(styles[0]);
  const [details, setDetails] = useState('');
  const [prompt, setPrompt] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);

  const disabled = !subject.trim() || isLoading;

  const handleGenerate = async () => {
    if (!onCraftPrompt || !subject.trim()) return;
    setIsLoading(true);
    setPrompt('');
    try {
      const crafted = await onCraftPrompt({
        subject: subject.trim(),
        mood,
        style,
        details: details.trim(),
      });
      setPrompt(crafted);
    } catch (error) {
      setPrompt(error instanceof Error ? error.message : 'Não foi possível criar o prompt.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      <Card title="Prompt para imagens" description="Defina briefing visual e gere instruções padronizadas para modelos de imagem.">
        <div className="grid gap-6 lg:grid-cols-[1.25fr,1fr]">
          <div className="space-y-4">
            <div>
              <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
                Tema principal
              </label>
              <input
                value={subject}
                onChange={(event) => setSubject(event.target.value)}
                placeholder="Ex.: Aplicativo de finanças para universitários"
                className="w-full rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-700 shadow-inner focus:border-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-200 dark:border-slate-700 dark:bg-[#0f172a] dark:text-slate-200"
              />
            </div>
            <div className="grid gap-4 sm:grid-cols-2">
              <div>
                <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
                  Clima
                </label>
                <select
                  title={t("Selecione o clima ou tom desejado para a imagem")}
                  value={mood}
                  onChange={(event) => setMood(event.target.value)}
                  className="w-full rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 shadow-inner focus:border-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-200 dark:border-slate-700 dark:bg-[#0f172a] dark:text-slate-200"
                >
                  {moods.map((option) => (
                    <option key={option} value={option}>
                      {option}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
                  Estilo
                </label>
                <select
                  title={t("Selecione o estilo visual desejado para a imagem")}
                  value={style}
                  onChange={(event) => setStyle(event.target.value)}
                  className="w-full rounded-xl border border-slate-200 bg-white px-3 py-2 text-sm text-slate-700 shadow-inner focus:border-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-200 dark:border-slate-700 dark:bg-[#0f172a] dark:text-slate-200"
                >
                  {styles.map((option) => (
                    <option key={option} value={option}>
                      {option}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            <div>
              <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
                Detalhes extras
              </label>
              <textarea
                value={details}
                onChange={(event) => setDetails(event.target.value)}
                rows={5}
                placeholder="Paleta desejada, elementos obrigatórios, texto de UI, referências visuais..."
                className="w-full resize-none rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-200 dark:border-slate-700 dark:bg-[#0f172a] dark:text-slate-200"
              />
            </div>
            <button
              type="button"
              onClick={handleGenerate}
              disabled={disabled}
              className="inline-flex items-center gap-2 rounded-full border border-slate-900 bg-slate-900 px-6 py-2 text-sm font-semibold text-white shadow-[0_20px_45px_-35px_rgba(15,23,42,0.8)] transition disabled:cursor-not-allowed disabled:opacity-60 dark:border-[#00f0ff] dark:bg-[#00f0ff] dark:text-[#010409]"
            >
              {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Sparkles className="h-4 w-4" />}
              {isLoading ? 'Gerando briefing...' : 'Gerar prompt'}
            </button>
          </div>

          <div className="flex h-full flex-col rounded-2xl border border-slate-200/80 bg-white/90 p-6 shadow-inner dark:border-slate-800/60 dark:bg-[#0a0f14]/70">
            <div className="flex items-center gap-2 text-xs font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
              <ImageIcon className="h-4 w-4" /> Prompt sugerido
            </div>
            <div className="mt-4 flex-1 overflow-auto rounded-xl border border-dashed border-slate-200/80 bg-white/80 p-4 text-sm text-slate-600 dark:border-slate-700/80 dark:bg-[#0f172a]/80 dark:text-slate-300">
              {isLoading && (
                <div className="flex items-center gap-2 text-sm text-slate-500 dark:text-slate-400">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  Montando instruções visuais...
                </div>
              )}
              {!isLoading && prompt && <pre className="whitespace-pre-wrap break-words text-sm leading-relaxed">{prompt}</pre>}
              {!isLoading && !prompt && (
                <p>O texto pronto para IA de imagens aparecerá aqui.</p>
              )}
            </div>
          </div>
        </div>
      </Card>
    </div>
  );
};

export default ImageGenerator;
function useTranslation(): { t: any; } {
  throw new Error('Function not implemented.');
}

