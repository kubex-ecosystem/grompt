'use client';

import { KeyRound, Moon, Settings, Sparkles, Sun, Users } from 'lucide-react';
import { usePathname } from 'next/navigation';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
// import { Link } from 'react-router-dom';
import Link from 'next/link';
import ApiKeysDrawer from './ApiKeysDrawer';
import ConfigModal from './ConfigModal';
import LanguageSelector from './LanguageSelector';

export default function Header() {
  const { t } = useTranslation();
  const pathname = usePathname();
  const [darkMode, setDarkMode] = useState<boolean>(true);
  const [showConfigModal, setShowConfigModal] = useState<boolean>(false);
  const [showKeysDrawer, setShowKeysDrawer] = useState<boolean>(false);

  // Temas
  const theme = {
    dark: {
      bg: 'bg-gray-900',
      cardBg: 'bg-gray-800',
      text: 'text-white',
      accent: 'text-blue-400',
      button: 'bg-blue-600 hover:bg-blue-700',
      buttonSecondary: 'bg-gray-700 hover:bg-gray-600',
      border: 'border-gray-700'
    },
    light: {
      bg: 'bg-gray-50',
      cardBg: 'bg-white',
      text: 'text-gray-900',
      accent: 'text-blue-600',
      button: 'bg-blue-600 hover:bg-blue-700',
      buttonSecondary: 'bg-gray-200 hover:bg-gray-300',
      border: 'border-gray-200'
    }
  };

  const currentTheme = darkMode ? theme.dark : theme.light;

  const isActive = (path: string) => pathname === path;

  return (
    <header className={`${currentTheme.cardBg} shadow-lg border-b ${currentTheme.border}`}>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg">
              <Sparkles className="h-6 w-6 text-white" />
            </div>
            <div>
              <h1 className="text-xl font-bold bg-gradient-to-r from-blue-400 to-purple-500 bg-clip-text text-transparent">
                Grompt
              </h1>
              <p className="text-xs text-gray-500">{t('common.subtitle')}</p>
            </div>
          </div>

          {/* Navigation */}
          <nav className="hidden md:flex space-x-8">
            <Link
              href="/"
              className={`flex items-center space-x-2 px-3 py-2 rounded-md text-sm font-medium transition-colors ${isActive('/') ? currentTheme.accent : 'text-gray-500 hover:text-gray-300'
                }`}
            >
              <Sparkles className="h-4 w-4" />
              <span>{t('nav.home')}</span>
            </Link>

            <Link
              href="/prompt"
              className={`flex items-center space-x-2 px-3 py-2 rounded-md text-sm font-medium transition-colors ${isActive('/prompt') ? currentTheme.accent : 'text-gray-500 hover:text-gray-300'
                }`}
            >
              <Sparkles className="h-4 w-4" />
              <span>{t('nav.prompt_crafter')}</span>
            </Link>

            <Link
              href="/agents"
              className={`flex items-center space-x-2 px-3 py-2 rounded-md text-sm font-medium transition-colors ${isActive('/agents') ? currentTheme.accent : 'text-gray-500 hover:text-gray-300'
                }`}
            >
              <Users className="h-4 w-4" />
              <span>{t('nav.agents')}</span>
            </Link>
          </nav>

          {/* Actions */}
          <div className="flex items-center space-x-4">
            <LanguageSelector currentTheme={currentTheme} />

            <button
              onClick={() => setShowKeysDrawer(true)}
              className={`p-2 rounded-md ${currentTheme.buttonSecondary} transition-colors`}
              title="Gerenciar Chaves (BYOK)"
            >
              <KeyRound className="h-4 w-4" />
            </button>

            <button
              onClick={() => setShowConfigModal(true)}
              className={`p-2 rounded-md ${currentTheme.buttonSecondary} transition-colors`}
              title="Configurações"
            >
              <Settings className="h-4 w-4" />
            </button>

            <button
              onClick={() => setDarkMode(!darkMode)}
              className={`p-2 rounded-md ${currentTheme.button} transition-colors`}
              title={darkMode ? t('common.light_mode') : t('common.dark_mode')}
            >
              {darkMode ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
            </button>
          </div>
        </div>
      </div>

      {/* Config Modal */}
      <ConfigModal
        isOpen={showConfigModal}
        onClose={() => setShowConfigModal(false)}
        onSave={() => {
          // Reload page to refresh API configuration
          window.location.reload();
        }}
      />

      {/* BYOK Drawer */}
      <ApiKeysDrawer isOpen={showKeysDrawer} onClose={() => setShowKeysDrawer(false)} />
    </header>
  );
}
