import { BotMessageSquare, ChevronsLeft, ChevronsRight, Moon, Sun } from 'lucide-react';
import React from 'react';
import { Theme } from '../../../types';
import { useTranslations } from '../../i18n/useTranslations';
import LanguageSelector from './LanguageSelector';

interface HeaderProps {
  theme: Theme;
  toggleTheme: () => void;
  onToggleMenu?: () => void;
  onToggleSidebarCollapse?: () => void;
  sidebarCollapsed?: boolean;
}

const Header: React.FC<HeaderProps> = ({ theme, toggleTheme, onToggleMenu, onToggleSidebarCollapse, sidebarCollapsed }) => {
  const { t } = useTranslations();
  const isLight = theme === 'light';
  const collapseLabel = sidebarCollapsed ? t('expandSidebar') : t('collapseSidebar');

  return (
    <header className="flex items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
      <div className="flex items-center gap-4">
        <button
          type="button"
          onClick={onToggleMenu}
          className="inline-flex h-10 w-10 items-center justify-center rounded-lg border border-[#e2e8f0] bg-white text-[#475569] shadow-sm transition hover:border-[#cbd5f5] hover:text-[#1f2937] focus:outline-none focus:ring-2 focus:ring-[#06b6d4]/40 dark:border-[#0b1220] dark:bg-[#0a1523] dark:text-[#94a3b8] dark:hover:border-[#13263a] lg:hidden"
          aria-label="Open navigation"
        >
          <span className="sr-only">Open navigation</span>
          ≡
        </button>
        <div
          className={`w-12 h-12 rounded-full flex items-center justify-center shadow-soft-card transition-colors duration-300 ${isLight
            ? 'bg-white border border-[#e2e8f0] text-[#06b6d4]'
            : 'bg-[#0a1523] border-2 border-[#06b6d4]/60'
            }`}
        >
          <BotMessageSquare
            size={24}
            className={`${isLight ? 'text-[#06b6d4]' : 'text-[#06b6d4]'}`}
          />
        </div>
        <div>
          <div className="flex items-center gap-2">
            <h1
              className={`text-3xl font-bold font-orbitron tracking-widest uppercase ${isLight ? 'text-[#111827]' : 'text-[#06b6d4]'
                }`}
            >
              Grompt
            </h1>
            {/*
            <!-- Não exibir versão por ora -->
            <span className="text-xs bg-gradient-to-r from-emerald-500 to-[#06b6d4] dark:from-[#16a34a] dark:to-[#38cde4] text-white dark:text-black px-2 py-1 rounded-full font-bold">
              v1.0.8
            </span> */}
          </div>
          <h2 className={`text-lg font-medium font-plex-mono ${isLight ? 'text-[#475569]' : 'text-[#94a3b8]'}`}>
            {t('promptCrafter')}
          </h2>
          {/*
          <!-- Desativado por ora... Não há landing page ainda, então essa está na raiz -->
          <a
            href="https://kubex.world"
            target="_blank"
            rel="noopener noreferrer"
            className="text-xs text-[#06b6d4] transition-colors duration-200 hover:text-[#0891b2] font-plex-mono dark:text-[#38cde4]/80 dark:hover:text-[#38cde4]"
          >
            ← Back to Kubex Ecosystem
          </a> */}
        </div>
      </div>
      <div className="flex items-center gap-3">
        {onToggleSidebarCollapse && (
          <button
            type="button"
            onClick={onToggleSidebarCollapse}
            className="hidden h-10 w-10 items-center justify-center rounded-lg border border-[#e2e8f0] bg-white text-[#475569] shadow-sm transition hover:border-[#cbd5f5] hover:text-[#1f2937] focus:outline-none focus:ring-2 focus:ring-[#06b6d4]/40 dark:border-[#0b1220] dark:bg-[#0a1523] dark:text-[#94a3b8] dark:hover:border-[#13263a] lg:inline-flex"
            title={collapseLabel}
            aria-label={collapseLabel}
          >
            {sidebarCollapsed ? <ChevronsRight size={20} /> : <ChevronsLeft size={20} />}
          </button>
        )}
        <LanguageSelector />
        <button
          type='button'
          title={t('toggleTheme', { theme: theme === 'light' ? 'dark' : 'light' })}
          onClick={toggleTheme}
          className={`p-2 rounded-full transition-colors duration-200 ${isLight
            ? 'bg-white border border-[#e2e8f0] text-[#475569] hover:text-[#0f172a] hover:border-[#cbd5f5] shadow-sm'
            : 'bg-[#0a1523] border border-[#13263a] text-[#94a3b8] hover:text-[#f5f3ff]'
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
