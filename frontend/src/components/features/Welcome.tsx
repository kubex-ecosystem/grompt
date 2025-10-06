import { ArrowRight, Compass, Sparkles } from 'lucide-react';
import React from 'react';
import Card from '../ui/Card';

interface WelcomeProps {
  onGetStarted: () => void;
}

const featureHighlights = [
  {
    title: 'Prompt Crafter',
    description: 'Estruture ideias em prompts profissionais com histórico, compartilhamento e otimizações de tom.',
  },
  {
    title: 'Conversas Assistidas',
    description: 'Use chat contextual com múltiplos provedores de IA e alternância rápida entre modelos.',
  },
  {
    title: 'Sumarização & Entregáveis',
    description: 'Gere resumos executivos, planos de ação e mensagens focadas a partir de conteúdos extensos.',
  },
];

const Welcome: React.FC<WelcomeProps> = ({ onGetStarted }) => {
  return (
    <div className="space-y-8">
      <Card
        className="bg-gradient-to-br from-white via-white to-slate-100/80 dark:from-[#0a0f14] dark:via-[#0a0f14] dark:to-[#0a0f14]/90"
      >
        <div className="flex flex-col gap-6 lg:flex-row lg:items-center">
          <div className="flex-1 space-y-4">
            <p className="inline-flex items-center gap-2 text-xs font-medium uppercase tracking-[0.4em] text-slate-500 dark:text-slate-400">
              <Compass size={16} /> Kubex Ecosystem
            </p>
            <h1 className="text-3xl font-bold tracking-tight text-slate-900 dark:text-[#e0f7fa] sm:text-4xl md:text-5xl">
              Governança, Criatividade, Produtividade, Liberdade!
            </h1>
            <p className="text-base text-slate-600 dark:text-slate-300">
              O Grompt reúne os fluxos de Prompt Engineering, geração de conteúdo e análise de código em um único arquivo, com frontend
              estável e leve, gateway coalescente, CLI dinâmica e intuitiva, multi-provider com mais de 100 APIs. Personalize idiomas, temas e provedores de IA na mesma interface.
            </p>
            <div className="flex flex-col items-start gap-3 sm:flex-row">
              <button
                type="button"
                onClick={onGetStarted}
                className="inline-flex items-center gap-2 rounded-full border border-slate-900 bg-slate-900 px-6 py-2 text-sm font-semibold text-white shadow-[0_20px_45px_-35px_rgba(15,23,42,0.8)] transition hover:scale-[1.01] focus:outline-none focus:ring-2 focus:ring-slate-400 dark:border-[#00f0ff] dark:bg-[#00f0ff] dark:text-[#010409]"
              >
                Começar agora
                <ArrowRight size={16} />
              </button>
              <div className="rounded-full border border-slate-200/70 bg-white/60 px-4 py-2 text-xs text-slate-500 shadow-sm dark:border-slate-700 dark:bg-[#0a0f14]/80 dark:text-slate-300">
                Build: front-end React + Vite 7 + Tailwind 3
              </div>
            </div>
          </div>
          <div className="flex-1 rounded-2xl border border-slate-200/80 bg-white/80 p-6 shadow-inner dark:border-slate-800/60 dark:bg-[#0f172a]/60">
            <p className="text-xs uppercase tracking-[0.45em] text-slate-500 dark:text-slate-400">Stack unificada</p>
            <ul className="mt-4 space-y-3 text-sm text-slate-600 dark:text-slate-300">
              <li className="flex items-start gap-3">
                <Sparkles className="mt-1 h-4 w-4 text-emerald-500" />
                Multi-provedores (Gemini, OpenAI, Anthropic) com wrapper unificado.
              </li>
              <li className="flex items-start gap-3">
                <Sparkles className="mt-1 h-4 w-4 text-sky-500" />
                Persistência híbrida usando IndexedDB com fallback automágico para localStorage.
              </li>
              <li className="flex items-start gap-3">
                <Sparkles className="mt-1 h-4 w-4 text-purple-500" />
                Componentes responsivos com Tailwind e animações do Framer Motion.
              </li>
            </ul>
          </div>
        </div>
      </Card>

      <div className="grid gap-6 md:grid-cols-3">
        {featureHighlights.map((feature) => (
          <Card key={feature.title} title={feature.title} description={feature.description}>
            <div className="flex items-center justify-between pt-1 text-xs text-slate-500 dark:text-slate-400">
              <span>Conectado ao fluxo Kubex</span>
              <span>Sem lock-in</span>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
};

export default Welcome;
