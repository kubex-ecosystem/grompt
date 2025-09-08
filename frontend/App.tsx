import React, { useState, useEffect } from 'react';
import PromptCrafter from './components/PromptCrafter';
import { BotMessageSquare, Sun, Moon } from 'lucide-react';

type Theme = 'light' | 'dark';

const App: React.FC = () => {
  const [theme, setTheme] = useState<Theme>('dark');

  useEffect(() => {
    const savedTheme = localStorage.getItem('theme') as Theme | null;
    const prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
    const initialTheme = savedTheme || (prefersDark ? 'dark' : 'light');
    setTheme(initialTheme);
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

  const toggleTheme = () => {
    setTheme(prevTheme => (prevTheme === 'light' ? 'dark' : 'light'));
  };

  return (
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
                Prompt Crafter
              </h2>
            </div>
          </div>
          <button
            onClick={toggleTheme}
            className="p-2 rounded-full bg-slate-200/50 dark:bg-[#10151b] text-slate-600 dark:text-[#90a4ae] hover:text-sky-500 dark:hover:text-[#00f0ff] transition-colors duration-200"
            aria-label={`Switch to ${theme === 'light' ? 'dark' : 'light'} mode`}
          >
            {theme === 'light' ? <Moon size={24} /> : <Sun size={24} />}
          </button>
        </header>
        <main>
          <PromptCrafter />
        </main>
        <footer className="text-center mt-12 text-slate-500 dark:text-[#90a4ae] text-xs">
          <p>Powered by Gemini API. Adhering to the Kubex Principles: Radical Simplicity, No Cages.</p>
          <p className="mt-1 font-orbitron tracking-wider">CODE FAST. OWN EVERYTHING.</p>
        </footer>
      </div>
    </div>
  );
};

export default App;