import { Menu, Moon, Sun } from 'lucide-react';
import React from 'react';
import { Theme } from '@/types';
import LanguageSelector from './LanguageSelector';

interface HeaderProps {
  theme: Theme;
  onToggleTheme: () => void;
  onToggleSidebar: () => void;
  collapsed?: boolean;
}

const Header: React.FC<HeaderProps> = ({ theme, onToggleTheme, onToggleSidebar, collapsed = false }) => {
  return (
    <div className="flex items-center justify-between px-4 py-4">
      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={onToggleSidebar}
          className="flex items-center justify-center rounded-full border border-[#e2e8f0] bg-white/80 p-2 text-[#475569] transition hover:bg-[#ecfeff] dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#94a3b8]"
          aria-label="Toggle sidebar"
        >
          <Menu size={20} />
        </button>
        {!collapsed && (
          <div className="hidden sm:flex flex-col">
            <h2 className="text-sm font-semibold text-[#111827] dark:text-[#e5f2f2]">
              Grompt Workspace
            </h2>
            <p className="text-xs text-[#64748b] dark:text-[#94a3b8]">
              AI-powered prompt engineering
            </p>
          </div>
        )}
      </div>
      <div className="flex items-center gap-3">
        <LanguageSelector />
        <button
          type="button"
          onClick={onToggleTheme}
          className="flex items-center justify-center rounded-full border border-[#e2e8f0] bg-white/80 p-2 text-[#475569] transition hover:bg-[#ecfeff] dark:border-[#13263a] dark:bg-[#0a1523] dark:text-[#94a3b8]"
          aria-label="Toggle theme"
        >
          {theme === 'dark' ? <Sun size={20} /> : <Moon size={20} />}
        </button>
      </div>
    </div>
  );
};

export default Header;
