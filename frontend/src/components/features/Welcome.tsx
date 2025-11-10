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
      <Card className="bg-gradient-to-br from-[#ffffff] via-[#f9fafb] to-[#ecfeff] dark:from-[#0a1523] dark:via-[#0a1523] dark:to-[#0a1523]/90">
        <div className="flex flex-col gap-6 lg:flex-row lg:items-center">
          <div className="flex-1 space-y-4">
            <p className="inline-flex items-center gap-2 text-xs font-medium uppercase tracking-[0.4em] text-[#06b6d4] dark:text-[#38cde4]">
              <Compass size={16} /> Kubex Ecosystem
            </p>
            <h1 className="text-3xl font-bold tracking-tight text-[#111827] dark:text-[#f5f3ff] sm:text-4xl md:text-5xl">
              Governança, Criatividade, Produtividade, Liberdade!
            </h1>
            <p className="text-base text-[#475569] dark:text-[#94a3b8]">
              O Grompt reúne os fluxos de Prompt Engineering, geração de conteúdo e análise de código em um único arquivo, com frontend
              estável e leve, gateway coalescente, CLI dinâmica e intuitiva, multi-provider com mais de 100 APIs. Personalize idiomas, temas e provedores de IA na mesma interface.
            </p>
            <div className="flex flex-col items-start gap-3 sm:flex-row">
              <button
                type="button"
                onClick={onGetStarted}
                className="inline-flex items-center gap-2 rounded-full border border-[#06b6d4] bg-[#06b6d4] px-6 py-2 text-sm font-semibold text-white shadow-soft-card transition hover:scale-[1.01] hover:bg-[#0891b2] focus:outline-none focus:ring-2 focus:ring-[#06b6d4]/40 dark:border-[#06b6d4] dark:bg-[#06b6d4] dark:text-[#0a1523]"
              >
                Começar agora
                <ArrowRight size={16} />
              </button>
              <div className="rounded-full border border-[#e2e8f0] bg-white/70 px-4 py-2 text-xs text-[#475569] shadow-sm dark:border-[#13263a] dark:bg-[#0a1523]/80 dark:text-[#94a3b8]">
                Build: front-end React + Vite 7 + Tailwind 3
              </div>
            </div>
          </div>
          <div className="flex-1 rounded-2xl border border-[#e2e8f0] bg-white/95 p-6 shadow-soft-card dark:border-[#13263a] dark:bg-[#0a1523]/75">
            <p className="text-xs uppercase tracking-[0.45em] text-[#94a3b8] dark:text-[#64748b]">Stack unificada</p>
            <ul className="mt-4 space-y-3 text-sm text-[#475569] dark:text-[#cbd5f5]">
              <li className="flex items-start gap-3">
                <Sparkles className="mt-1 h-4 w-4 text-[#06b6d4]" />
                Multi-provedores (Gemini, OpenAI, Anthropic) com wrapper unificado.
              </li>
              <li className="flex items-start gap-3">
                <Sparkles className="mt-1 h-4 w-4 text-[#a855f7]" />
                Persistência híbrida usando IndexedDB com fallback automágico para localStorage.
              </li>
              <li className="flex items-start gap-3">
                <Sparkles className="mt-1 h-4 w-4 text-[#d946ef]" />
                Componentes responsivos com Tailwind e animações do Framer Motion.
              </li>
            </ul>
          </div>
        </div>
      </Card>

      <div className="grid gap-6 md:grid-cols-3">
        {featureHighlights.map((feature) => (
          <Card key={feature.title} title={feature.title} description={feature.description}>
            <div className="flex items-center justify-between pt-1 text-xs text-[#64748b] dark:text-[#94a3b8]">
              <span>Conectado ao fluxo Kubex</span>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
};

export default Welcome;
