import { Braces, ClipboardCheck, ClipboardCopy, Loader2, Sparkles } from 'lucide-react';
import React, { useState } from 'react';
import Card from '../ui/Card';

interface CodeGeneratorProps {
  onGenerate?: (spec: {
    stack: string;
    goal: string;
    constraints: string[];
    extras: string;
  }) => Promise<string>;
}

const stacks = ['Go + Fiber', 'TypeScript + React', 'Python + FastAPI', 'Rust + Axum'];

const CodeGenerator: React.FC<CodeGeneratorProps> = ({ onGenerate }) => {
  const [stack, setStack] = useState(stacks[0]);
  const [goal, setGoal] = useState('');
  const [constraints, setConstraints] = useState<string[]>([]);
  const [extraNotes, setExtraNotes] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [result, setResult] = useState<string>('');
  const [copied, setCopied] = useState(false);

  const toggleConstraint = (value: string) => {
    setConstraints((prev) =>
      prev.includes(value) ? prev.filter((item) => item !== value) : [...prev, value]
    );
  };

  const handleGenerate = async () => {
    if (!goal.trim() || !onGenerate) return;
    setIsGenerating(true);
    setResult('');
    try {
      const code = await onGenerate({
        stack,
        goal: goal.trim(),
        constraints,
        extras: extraNotes.trim(),
      });
      setResult(code);
    } catch (error) {
      setResult(error instanceof Error ? error.message : 'Não foi possível gerar o código.');
    } finally {
      setIsGenerating(false);
    }
  };

  const handleCopy = async () => {
    if (!result) return;
    await navigator.clipboard.writeText(result);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const disabled = !goal.trim() || isGenerating;

  return (
    <div className="space-y-6">
      <Card title="Code Generator" description="Gere scaffolds de código com constraints Kubex-ready.">
        <div className="grid gap-6 lg:grid-cols-2">
          <div className="space-y-4">
            <label className="block text-xs font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
              Stack alvo
            </label>
            <div className="grid grid-cols-2 gap-2">
              {stacks.map((option) => (
                <button
                  key={option}
                  type="button"
                  onClick={() => setStack(option)}
                  className={`rounded-xl border px-4 py-3 text-sm font-semibold transition ${
                    stack === option
                      ? 'border-slate-900 bg-slate-900 text-white shadow-[0_20px_45px_-35px_rgba(15,23,42,0.8)] dark:border-[#00f0ff] dark:bg-[#00f0ff] dark:text-[#010409]'
                      : 'border-slate-200 bg-white text-slate-600 hover:border-slate-300 dark:border-slate-700 dark:bg-[#0f172a] dark:text-slate-300'
                  }`}
                >
                  {option}
                </button>
              ))}
            </div>

            <div>
              <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
                Objetivo principal
              </label>
              <textarea
                value={goal}
                onChange={(event) => setGoal(event.target.value)}
                placeholder="Descreva o módulo, endpoint ou fluxo que deseja gerar."
                rows={6}
                className="w-full resize-none rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-200 dark:border-slate-700 dark:bg-[#0f172a] dark:text-slate-200"
              />
            </div>

            <div className="space-y-2">
              <p className="text-xs font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
                Constraints recomendadas
              </p>
              {['Testes unitários incluídos', 'Sem frameworks proprietários', 'Documentação inline'].map((constraint) => (
                <label key={constraint} className="flex items-center gap-2 text-sm text-slate-600 dark:text-slate-300">
                  <input
                    type="checkbox"
                    checked={constraints.includes(constraint)}
                    onChange={() => toggleConstraint(constraint)}
                    className="h-4 w-4 rounded border-slate-300 text-slate-900 focus:ring-slate-400 dark:border-slate-700 dark:bg-[#0f172a]"
                  />
                  {constraint}
                </label>
              ))}
            </div>

            <div>
              <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
                Observações extras
              </label>
              <textarea
                value={extraNotes}
                onChange={(event) => setExtraNotes(event.target.value)}
                placeholder="Dependências preferidas, integrações, padrões arquiteturais..."
                rows={4}
                className="w-full resize-none rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-700 shadow-inner transition focus:border-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-200 dark:border-slate-700 dark:bg-[#0f172a] dark:text-slate-200"
              />
            </div>

            <button
              type="button"
              onClick={handleGenerate}
              disabled={disabled}
              className="inline-flex items-center gap-2 rounded-full border border-slate-900 bg-slate-900 px-6 py-2 text-sm font-semibold text-white shadow-[0_20px_45px_-35px_rgba(15,23,42,0.8)] transition disabled:cursor-not-allowed disabled:opacity-60 dark:border-[#00f0ff] dark:bg-[#00f0ff] dark:text-[#010409]"
            >
              {isGenerating ? <Loader2 className="h-4 w-4 animate-spin" /> : <Sparkles className="h-4 w-4" />}
              {isGenerating ? 'Gerando blueprint...' : 'Gerar código'}
            </button>
          </div>

          <div className="flex h-full flex-col rounded-2xl border border-slate-200/80 bg-white/80 p-5 shadow-inner dark:border-slate-800/60 dark:bg-[#0a0f14]/70">
            <div className="flex items-center justify-between text-xs font-semibold uppercase tracking-[0.3em] text-slate-500 dark:text-slate-400">
              <span className="inline-flex items-center gap-2"><Braces className="h-4 w-4" /> Saída</span>
              <button
                type="button"
                onClick={handleCopy}
                disabled={!result}
                className="inline-flex items-center gap-1 rounded-full border border-slate-200 bg-white px-3 py-1 text-[11px] font-semibold text-slate-600 transition disabled:cursor-not-allowed disabled:opacity-50 dark:border-slate-700 dark:bg-[#0f172a] dark:text-slate-300"
              >
                {copied ? <ClipboardCheck className="h-3 w-3" /> : <ClipboardCopy className="h-3 w-3" />}
                {copied ? 'Copiado' : 'Copiar'}
              </button>
            </div>
            <div className="mt-4 flex-1 overflow-auto rounded-xl border border-dashed border-slate-200/80 bg-white/80 p-4 text-sm text-slate-600 dark:border-slate-700/80 dark:bg-[#0f172a]/80 dark:text-slate-300">
              {isGenerating && (
                <div className="flex items-center gap-2 text-sm text-slate-500 dark:text-slate-400">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  Gerando snippet idiomático...
                </div>
              )}
              {!isGenerating && result && (
                <pre className="whitespace-pre-wrap break-words text-sm leading-relaxed">{result}</pre>
              )}
              {!isGenerating && !result && (
                <p>O código e instruções aparecerão aqui após a geração.</p>
              )}
            </div>
          </div>
        </div>
      </Card>
    </div>
  );
};

export default CodeGenerator;
