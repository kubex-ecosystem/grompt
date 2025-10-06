import { BotMessageSquare, Moon, Sun } from 'lucide-react';
import React from 'react';
import { Theme } from '../../../types';
import { useTranslations } from '../../i18n/useTranslations';
import LanguageSelector from './LanguageSelector';

interface HeaderProps {
  theme: Theme;
  toggleTheme: () => void;
  onToggleMenu?: () => void;
}

const Header: React.FC<HeaderProps> = ({ theme, toggleTheme, onToggleMenu }) => {
  const { t } = useTranslations();
  const isLight = theme === 'light';

  return (
    <header className="flex items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
      <div className="flex items-center gap-4">
        <button
          type="button"
          onClick={onToggleMenu}
          className="inline-flex h-10 w-10 items-center justify-center rounded-lg border border-slate-200/80 bg-white text-slate-600 shadow-sm transition hover:border-slate-300 hover:text-slate-900 focus:outline-none focus:ring-2 focus:ring-slate-400 dark:border-slate-800/80 dark:bg-[#0a0f14] dark:text-slate-300 dark:hover:border-slate-600 lg:hidden"
          aria-label="Open navigation"
        >
          <span className="sr-only">Open navigation</span>
          ≡
        </button>
        <div
          className={`w-12 h-12 rounded-full flex items-center justify-center shadow-md transition-colors duration-300 ${isLight
            ? 'bg-white border border-slate-200 text-sky-600 shadow-[0_16px_40px_-30px_rgba(15,23,42,0.45)]'
            : 'bg-[#10151b] border-2 border-[#00f0ff] dark:neon-border-cyan'
            }`}
        >
          <BotMessageSquare
            size={24}
            className={`${isLight ? 'text-sky-600' : 'text-[#00f0ff] dark:neon-glow-cyan'}`}
          />
        </div>
        <div>
          <div className="flex items-center gap-2">
            <h1
              className={`text-3xl font-bold font-orbitron tracking-widest uppercase ${isLight ? 'text-slate-900' : 'text-[#00f0ff] dark:neon-glow-cyan'
                }`}
            >
              Grompt
            </h1>
            {/*
            <!-- Não exibir versão por ora -->
            <span className="text-xs bg-gradient-to-r from-emerald-500 to-sky-500 dark:from-[#00e676] dark:to-[#00f0ff] text-white dark:text-black px-2 py-1 rounded-full font-bold">
              v1.0.8
            </span> */}
          </div>
          <h2 className={`text-lg font-medium font-plex-mono ${isLight ? 'text-slate-500' : 'text-[#90a4ae]'}`}>
            {t('promptCrafter')}
          </h2>
          {/*
          <!-- Desativado por ora... Não há landing page ainda, então essa está na raiz -->
          <a
            href="https://kubex.world"
            target="_blank"
            rel="noopener noreferrer"
            className="text-xs text-sky-600 dark:text-[#00f0ff]/80 hover:text-sky-500 dark:hover:text-[#00f0ff] transition-colors duration-200 font-plex-mono"
          >
            ← Back to Kubex Ecosystem
          </a> */}
        </div>
      </div>
      <div className="flex items-center gap-3">
        <LanguageSelector />
        <button
          type='button'
          title={t('toggleTheme', { theme: theme === 'light' ? 'dark' : 'light' })}
          onClick={toggleTheme}
          className={`p-2 rounded-full transition-colors duration-200 ${isLight
            ? 'bg-white border border-slate-200 text-slate-600 hover:text-slate-900 hover:border-slate-300 shadow-sm'
            : 'bg-[#10151b] text-[#90a4ae] hover:text-[#00f0ff]'
            }`}
          aria-label={t('toggleTheme', { theme: theme === 'light' ? 'dark' : 'light' })}
        >
          {theme === 'light' ? <Moon size={24} /> : <Sun size={24} />}
        </button>
      </div>
    </header>
  );
};

export default Header;
