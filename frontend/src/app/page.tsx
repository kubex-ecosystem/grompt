'use client';

import { Sparkles, Users } from 'lucide-react';
import Link from 'next/link';
import { useTranslation } from 'react-i18next';
 

export default function Home() {
  const { t } = useTranslation();
  

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

  const currentTheme = theme.dark;

  return (
    <div className={`min-h-screen transition-colors duration-300 ${currentTheme.bg} ${currentTheme.text}`}>
      <main className="py-4">
        {/* Welcome Section */}
        <div className={`${currentTheme.cardBg} rounded-lg shadow-lg p-8 mb-8 ${currentTheme.border} border`}>
          <div className="text-center">
            <div className="mb-4">
              <div className="inline-flex p-4 bg-gradient-to-r from-blue-500 to-purple-600 rounded-full">
                <Sparkles className="h-12 w-12 text-white" />
              </div>
            </div>
            <h1 className="text-4xl font-bold mb-4 bg-gradient-to-r from-blue-400 to-purple-500 bg-clip-text text-transparent">
              {t('welcome.title')}
            </h1>
            <p className="text-xl text-gray-500 mb-8 max-w-2xl mx-auto">
              {t('welcome.description')}
            </p>

            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Link
                href="/prompt"
                className={`inline-flex items-center px-6 py-3 ${currentTheme.button} text-white rounded-lg font-medium transition-colors`}
              >
                <Sparkles className="h-5 w-5 mr-2" />
                {t('welcome.start_crafting')}
              </Link>

              <Link
                href="/agents"
                className={`inline-flex items-center px-6 py-3 border-2 ${currentTheme.border} ${currentTheme.text} rounded-lg font-medium hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors`}
              >
                <Users className="h-5 w-5 mr-2" />
                {t('welcome.manage_agents')}
              </Link>
            </div>
          </div>
        </div>

        {/* Features Grid */}
        <div className="grid md:grid-cols-3 gap-6">
          <div className={`${currentTheme.cardBg} rounded-lg shadow p-6 ${currentTheme.border} border`}>
            <div className="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg w-fit mb-4">
              <Sparkles className="h-6 w-6 text-blue-600" />
            </div>
            <h3 className="text-lg font-semibold mb-2">{t('features.prompt_engineering.title')}</h3>
            <p className="text-gray-500">{t('features.prompt_engineering.description')}</p>
          </div>

          <div className={`${currentTheme.cardBg} rounded-lg shadow p-6 ${currentTheme.border} border`}>
            <div className="p-3 bg-purple-100 dark:bg-purple-900 rounded-lg w-fit mb-4">
              <Users className="h-6 w-6 text-purple-600" />
            </div>
            <h3 className="text-lg font-semibold mb-2">{t('features.agent_management.title')}</h3>
            <p className="text-gray-500">{t('features.agent_management.description')}</p>
          </div>

          <div className={`${currentTheme.cardBg} rounded-lg shadow p-6 ${currentTheme.border} border`}>
            <div className="p-3 bg-green-100 dark:bg-green-900 rounded-lg w-fit mb-4">
              <Sparkles className="h-6 w-6 text-green-600" />
            </div>
            <h3 className="text-lg font-semibold mb-2">{t('features.multi_llm.title')}</h3>
            <p className="text-gray-500">{t('features.multi_llm.description')}</p>
          </div>
        </div>
      </main>
    </div>
  );
}
