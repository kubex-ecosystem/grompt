import { ArrowRight, Compass, Sparkles } from 'lucide-react';
import React from 'react';
import Card from '../ui/Card';
import { useTranslations } from '../../i18n/useTranslations';

interface WelcomeProps {
  onGetStarted: () => void;
}

const Welcome: React.FC<WelcomeProps> = ({ onGetStarted }) => {
  const { t } = useTranslations();
  const featureHighlights = [
    {
      phase: t('welcomePhaseCreation'),
      title: t('welcomePromptTitle'),
      description: t('welcomePromptDescription'),
    },
    {
      phase: t('welcomePhaseAnalysis'),
      title: t('welcomeChatTitle'),
      description: t('welcomeChatDescription'),
    },
    {
      phase: t('welcomePhaseConsolidation'),
      title: t('welcomeSummaryTitle'),
      description: t('welcomeSummaryDescription'),
    },
  ];

  return (
    <div className="space-y-8">
      <Card className="bg-gradient-to-br from-[#ffffff] via-[#f9fafb] to-[#ecfeff] dark:from-[#0a1523] dark:via-[#0a1523] dark:to-[#0a1523]/90">
        <div className="flex flex-col gap-6 lg:flex-row lg:items-center">
          <div className="flex-1 space-y-4">
            <p className="inline-flex items-center gap-2 text-xs font-medium uppercase tracking-[0.4em] text-[#06b6d4] dark:text-[#38cde4]">
              <Compass size={16} /> {t('welcomeKicker')}
            </p>
            <h1 className="text-3xl font-bold tracking-tight text-[#111827] dark:text-[#f5f3ff] sm:text-4xl md:text-5xl">
              {t('welcomeHeadline')}
            </h1>
            <p className="text-base text-[#475569] dark:text-[#94a3b8]">
              {t('welcomeSubheadline')}
            </p>
            <div className="flex flex-col items-start gap-3 sm:flex-row">
              <button
                type="button"
                onClick={onGetStarted}
                className="inline-flex items-center gap-2 rounded-full border border-[#06b6d4] bg-[#06b6d4] px-6 py-2 text-sm font-semibold text-white shadow-soft-card transition hover:scale-[1.01] hover:bg-[#0891b2] focus:outline-none focus:ring-2 focus:ring-[#06b6d4]/40 dark:border-[#06b6d4] dark:bg-[#06b6d4] dark:text-[#0a1523]"
              >
                {t('welcomeCta')}
                <ArrowRight size={16} />
              </button>
            </div>
          </div>
          <div className="flex-1 rounded-2xl border border-[#e2e8f0] bg-white/95 p-6 shadow-soft-card dark:border-[#13263a] dark:bg-[#0a1523]/75">
            <p className="text-xs uppercase tracking-[0.45em] text-[#94a3b8] dark:text-[#64748b]">{t('welcomeStackTitle')}</p>
            <ul className="mt-4 space-y-3 text-sm text-[#475569] dark:text-[#cbd5f5]">
              <li className="flex items-start gap-3">
                <Sparkles className="mt-1 h-4 w-4 text-[#06b6d4]" />
                {t('welcomeStackItemOne')}
              </li>
              <li className="flex items-start gap-3">
                <Sparkles className="mt-1 h-4 w-4 text-[#a855f7]" />
                {t('welcomeStackItemTwo')}
              </li>
              <li className="flex items-start gap-3">
                <Sparkles className="mt-1 h-4 w-4 text-[#d946ef]" />
                {t('welcomeStackItemThree')}
              </li>
            </ul>
          </div>
        </div>
      </Card>

      <div className="grid gap-6 md:grid-cols-3">
        {featureHighlights.map((feature) => (
          <Card key={feature.title} title={feature.title} description={feature.description}>
            <div className="text-xs font-semibold uppercase tracking-[0.35em] text-[#06b6d4] dark:text-[#38cde4]">
              {feature.phase}
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
};

export default Welcome;
