import { BotMessageSquare, Moon, Sun } from 'lucide-react';
import React from 'react';
import { useTranslations } from '../../i18n/useTranslations';
import { Theme } from '../../types';
import LanguageSelector from './LanguageSelector';

interface HeaderProps {
  theme: Theme;
  toggleTheme: () => void;
}

const Header: React.FC<HeaderProps> = ({ theme, toggleTheme }) => {
  const { t } = useTranslations();

  return (
    <header className="mb-8 flex items-center justify-between">
      <div className="flex items-center gap-4">
        <div className="w-12 h-12 bg-white/50 dark:bg-[#10151b] border-2 border-sky-500 dark:border-[#00f0ff] rounded-full flex items-center justify-center shadow-lg dark:neon-border-cyan">
          <BotMessageSquare size={24} className="text-sky-500 dark:text-[#00f0ff] dark:neon-glow-cyan" />
        </div>
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-3xl font-bold font-orbitron text-sky-500 light-shadow-sky dark:text-[#00f0ff] dark:neon-glow-cyan tracking-widest uppercase">
              Grompt
            </h1>
            <span className="text-xs bg-gradient-to-r from-emerald-500 to-sky-500 dark:from-[#00e676] dark:to-[#00f0ff] text-white dark:text-black px-2 py-1 rounded-full font-bold">
              v2.0
            </span>
          </div>
          <h2 className="text-lg text-slate-500 dark:text-[#90a4ae] font-medium font-plex-mono">
            {t('promptCrafter')}
          </h2>
          <a
            href="https://kubex.world"
            target="_blank"
            rel="noopener noreferrer"
            className="text-xs text-sky-600 dark:text-[#00f0ff]/80 hover:text-sky-500 dark:hover:text-[#00f0ff] transition-colors duration-200 font-plex-mono"
          >
            ← Back to Kubex Ecosystem
          </a>
        </div>
      </div>
      <div className="flex items-center gap-3">
        <LanguageSelector />
        <button
          onClick={toggleTheme}
          className="p-2 rounded-full bg-slate-200/50 dark:bg-[#10151b] text-slate-600 dark:text-[#90a4ae] hover:text-sky-500 dark:hover:text-[#00f0ff] transition-colors duration-200"
          aria-label={t('toggleTheme', { theme: theme === 'light' ? 'dark' : 'light' })}
        >
          {theme === 'light' ? <Moon size={24} /> : <Sun size={24} />}
        </button>
      </div>
    </header>
  );
};

export default Header;
