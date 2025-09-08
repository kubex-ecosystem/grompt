// import React from 'react'
import { BotMessageSquare, Moon, Sun } from 'lucide-react';
import { createContext, useContext, useEffect, useState } from 'react';
import PromptCrafter from './components/PromptCrafter';

type Theme = 'light' | 'dark';
type Language = 'en' | 'es' | 'zh';

const translations: Record<Language, Record<string, string>> = {
  en: {
    'promptCrafter': 'Prompt Crafter',
    'toggleTheme': 'Switch to {theme} mode',
    'poweredBy': 'Kubex Principles: Radical Simplicity, No Cages.',
    'motto': 'CODE FAST. OWN EVERYTHING.',
  },
  es: {
    'promptCrafter': 'Creador de Prompts',
    'toggleTheme': 'Cambiar a modo {theme}',
    'poweredBy': 'Principios de Kubex: Simplicidad Radical, Sin Jaulas.',
    'motto': 'CODIFICA RÁPIDO. SÉ DUEÑO DE TODO.',
  },
  zh: {
    'promptCrafter': '提示词构建器',
    'toggleTheme': '切换到{theme}模式',
    'poweredBy': 'Kubex 原则：极致简约，无拘无束。',
    'motto': '快速编码。拥有一切。',
  }
};

interface LanguageContextType {
  language: Language;
  setLanguage: (lang: Language) => void;
  t: (key: string, params?: Record<string, string>) => string;
}

export const LanguageContext = createContext<LanguageContextType>({
  language: 'en',
  setLanguage: () => { },
  t: (key) => key,
});

const LanguageSelector: React.FC = () => {
  const { language, setLanguage } = useContext(LanguageContext);

  const languages: { code: Language; flag: React.ReactNode }[] = [
    { code: 'en', flag: <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 72 48"><path fill="#6a7bff" d="M0 0h72v48H0z" /><path fill="#fff" d="M0 0h72v24H0z" /><path fill="#ff4b55" d="M0 0h72v12H0z" /><path fill="#fff" d="M0 0h36v24H0z" /><path fill="#1e2a3d" d="M0 0h24v24H0z" /><path fill="#fff" d="m13 13-1-3-1 3h2zm-3-1 3 1v-2l-3 1zm1 3-1 1h2l-1-1zm-1 3 1-3 1 3h-2zM8 9l1 3 1-3H8zm3 1-3 1v-2l3 1zm-1 3 1 1H9l1-1zm1 3-1-3-1 3h2zm5-10-1 3-1-3h2zm-3-1 3 1v-2l-3 1zm1 3-1 1h2l-1-1zm-1 3 1-3 1 3h-2zM8 2l1 3 1-3H8zm3 1-3 1V1l3 1zm-1 3 1 1H9l1-1zm1 3-1-3-1 3h2z" /></svg> },
    { code: 'es', flag: <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 72 48"><path fill="#ff4b55" d="M0 0h72v48H0z" /><path fill="#ffd500" d="M0 12h72v24H0z" /><path fill="#ff4b55" d="M30 18h2v12h-2zm-2 0h-2v12h2zm-2 1h-1v10h1v-4h2v-2h-2v-4z" /><path fill="#ffd500" d="M30 18h-2v12h2v-5h-1v-2h1v-5z" /><path fill="#ff4b55" d="M31 23h-6v1h6v-1z" /></svg> },
    { code: 'zh', flag: <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 72 48"><path fill="#ff4b55" d="M0 0h72v48H0z" /><path fill="#ffd500" d="m18 9 2.2 6.8h7.1L21 20l2.2 6.8L18 22.1l-5.3 4.7L15 20l-6.2-4.2h7.1zm12.5 1.5 1 .8-.5 1.3-1.2-.3.1-1.4zm-1 16 .3-1.3-1.2.3.9 1zm6.5-13.5.8-1-1.3.6.5.4zm-4 14-1.3-.5.5-1.3.8 1zm3.5-1.5-.7 1.1.9.8.3-1.4-.5-.5z" /></svg> },
  ];

  return (
    <div className="flex items-center gap-1 bg-slate-200/50 dark:bg-[#10151b] p-1 rounded-full">
      {languages.map(lang => (
        <button
          key={lang.code}
          onClick={() => setLanguage(lang.code)}
          className={`w-8 h-8 rounded-full flex items-center justify-center transition-all duration-300 overflow-hidden ${language === lang.code ? 'ring-2 ring-sky-500 dark:ring-[#00f0ff] scale-110' : 'hover:opacity-80'}`}
          aria-label={`Switch to ${lang.code}`}
        >
          <div className="w-full h-full">{lang.flag}</div>
        </button>
      ))}
    </div>
  );
};

const App: React.FC = () => {
  const [theme, setTheme] = useState<Theme>('dark');
  const [language, setLanguage] = useState<Language>('en');

  useEffect(() => {
    const savedTheme = localStorage.getItem('theme') as Theme | null;
    const prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
    const initialTheme = savedTheme || (prefersDark ? 'dark' : 'light');
    setTheme(initialTheme);

    const savedLang = localStorage.getItem('language') as Language | null;
    const browserLang = navigator.language.split('-')[0] as Language;
    const initialLang = savedLang || (['en', 'es', 'zh'].includes(browserLang) ? browserLang : 'en');
    setLanguage(initialLang);
  }, []);

  useEffect(() => {
    const root = window.document.documentElement;
    const body = window.document.body;
    root.classList.remove('light', 'dark');
    root.classList.add(theme);
    body.classList.remove('light-theme', 'dark-theme');
    body.classList.add(`${theme}-theme`);
    localStorage.setItem('theme', theme);
  }, [theme]);

  useEffect(() => {
    localStorage.setItem('language', language);
  }, [language]);

  const toggleTheme = () => {
    setTheme(prevTheme => (prevTheme === 'light' ? 'dark' : 'light'));
  };

  const t = (key: string, params?: Record<string, string>): string => {
    let translation = translations[language][key] || translations.en[key] || key;
    if (params) {
      Object.keys(params).forEach(paramKey => {
        translation = translation.replace(`{${paramKey}}`, params[paramKey]);
      });
    }
    return translation;
  };

  return (
    <LanguageContext.Provider value={{ language, setLanguage, t }}>
      <div className="min-h-screen text-slate-800 dark:text-[#e0f7fa] font-plex-mono p-4 sm:p-6 lg:p-8">
        <div className="max-w-7xl mx-auto">
          <header className="mb-8 flex items-center justify-between">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-white/50 dark:bg-[#10151b] border-2 border-sky-500 dark:border-[#00f0ff] rounded-full flex items-center justify-center shadow-lg dark:neon-border-cyan">
                <BotMessageSquare size={24} className="text-sky-500 dark:text-[#00f0ff] dark:neon-glow-cyan" />
              </div>
              <div>
                <h1 className="text-3xl font-bold font-orbitron text-sky-500 light-shadow-sky dark:text-[#00f0ff] dark:neon-glow-cyan tracking-widest uppercase">
                  KUBEX
                </h1>
                <h2 className="text-lg text-slate-500 dark:text-[#90a4ae] font-medium font-plex-mono">
                  {t('promptCrafter')}
                </h2>
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
          <main>
            <PromptCrafter />
          </main>
          <footer className="text-center mt-12 text-slate-500 dark:text-[#90a4ae] text-xs">
            <p>{t('poweredBy')}</p>
            <p className="mt-1 font-orbitron tracking-wider">{t('motto')}</p>
          </footer>
        </div>
      </div>
    </LanguageContext.Provider>
  );
};

export default App;
