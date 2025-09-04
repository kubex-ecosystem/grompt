'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useEffect, useState } from 'react';
import { Menu, Sparkles, X, KeyRound, Sun, Moon, Users } from 'lucide-react';
import LanguageSelector from './LanguageSelector';
import { useTranslation } from 'react-i18next';
import ApiKeysDrawer from './ApiKeysDrawer';
import HistoryDrawer from './HistoryDrawer';

type Props = { children: React.ReactNode };

export default function AppShell({ children }: Props) {
  const { t } = useTranslation();
  const pathname = usePathname();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [showKeysDrawer, setShowKeysDrawer] = useState(false);
  const [showHistoryDrawer, setShowHistoryDrawer] = useState(false);
  const [darkMode, setDarkMode] = useState(true);

  useEffect(() => {
    document.documentElement.className = darkMode ? 'dark' : '';
  }, [darkMode]);

  const isActive = (path: string) => pathname === path;

  const theme = {
    bg: 'bg-gray-900',
    cardBg: 'bg-gray-800',
    text: 'text-white',
    accent: 'text-blue-400',
    button: 'bg-blue-600 hover:bg-blue-700',
    buttonSecondary: 'bg-gray-700 hover:bg-gray-600',
    border: 'border-gray-700'
  };

  return (
    <div className={`min-h-screen ${theme.bg} ${theme.text}`}>
      {/* Topbar */}
      <header className={`sticky top-0 z-40 ${theme.cardBg} border-b ${theme.border}`}>
        <div className="max-w-7xl mx-auto px-3 sm:px-6 lg:px-8 h-14 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <button
              aria-label="Abrir menu"
              className={`p-2 rounded md:hidden ${theme.buttonSecondary}`}
              onClick={() => setSidebarOpen(true)}
            >
              <Menu className="h-5 w-5" />
            </button>
            <Link href="/" className="flex items-center gap-2">
              <div className="p-2 bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg">
                <Sparkles className="h-5 w-5 text-white" />
              </div>
              <span className="font-semibold">Grompt</span>
            </Link>
            <nav className="hidden md:flex items-center gap-4 ml-4">
              <Link href="/" className={`${isActive('/') ? theme.accent : 'text-gray-400 hover:text-white'} text-sm`}>Home</Link>
              <Link href="/prompt" className={`${isActive('/prompt') ? theme.accent : 'text-gray-400 hover:text-white'} text-sm`}>Prompt</Link>
              <Link href="/agents" className={`${isActive('/agents') ? theme.accent : 'text-gray-400 hover:text-white'} text-sm flex items-center gap-1`}><Users className="h-4 w-4"/>Agents</Link>
            </nav>
          </div>
          <div className="flex items-center gap-2">
            <LanguageSelector currentTheme={{ buttonSecondary: theme.buttonSecondary }} />
            <button
              onClick={() => setShowHistoryDrawer(true)}
              className={`p-2 rounded ${theme.buttonSecondary}`}
              title={t('history.title')}
            >
              {/* Simple clock/history icon via emoji to avoid extra deps */}
              <span role="img" aria-label="hist">🕘</span>
            </button>
            <button
              onClick={() => setShowKeysDrawer(true)}
              className={`p-2 rounded ${theme.buttonSecondary}`}
              title="Gerenciar Chaves (BYOK)"
            >
              <KeyRound className="h-4 w-4" />
            </button>
            <button
              onClick={() => setDarkMode(!darkMode)}
              className={`p-2 rounded ${theme.buttonSecondary}`}
              title={darkMode ? 'Tema claro' : 'Tema escuro'}
            >
              {darkMode ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
            </button>
          </div>
        </div>
      </header>

      {/* Sidebar (overlay on mobile) */}
      {sidebarOpen && (
        <div className="fixed inset-0 z-50 flex md:hidden">
          <div className="flex-1 bg-black/40" onClick={() => setSidebarOpen(false)} />
          <aside className={`w-64 h-full ${theme.cardBg} border-l ${theme.border} p-4`}>
            <div className="flex justify-between items-center mb-4">
              <span className="font-semibold">Navegação</span>
              <button onClick={() => setSidebarOpen(false)} className={`p-2 rounded ${theme.buttonSecondary}`} aria-label="Fechar menu">
                <X className="h-4 w-4" />
              </button>
            </div>
            <nav className="flex flex-col gap-2">
              <Link href="/" onClick={() => setSidebarOpen(false)} className={`${isActive('/') ? theme.accent : 'text-gray-300 hover:text-white'} text-sm`}>Home</Link>
              <Link href="/prompt" onClick={() => setSidebarOpen(false)} className={`${isActive('/prompt') ? theme.accent : 'text-gray-300 hover:text-white'} text-sm`}>Prompt</Link>
              <Link href="/agents" onClick={() => setSidebarOpen(false)} className={`${isActive('/agents') ? theme.accent : 'text-gray-300 hover:text-white'} text-sm`}>Agents</Link>
            </nav>
          </aside>
        </div>
      )}

      {/* Content */}
      <main className="max-w-7xl mx-auto px-3 sm:px-6 lg:px-8 py-6">
        {children}
      </main>

      {/* BYOK Drawer Global */}
      <ApiKeysDrawer isOpen={showKeysDrawer} onClose={() => setShowKeysDrawer(false)} />
      {/* History Drawer Global */}
      <HistoryDrawer isOpen={showHistoryDrawer} onClose={() => setShowHistoryDrawer(false)} />
    </div>
  );
}
