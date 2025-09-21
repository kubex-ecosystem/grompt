import { BookOpen, Moon, Play, Sun } from 'lucide-react';

import DemoMode from '@/config/DemoMode';

const Header = ({
  darkMode,
  setDarkMode,
  currentTheme,
  startOnboarding,
  showEducation
}) => {
  return (
    <div className="flex justify-between items-center mb-8" id="header">
      <div>
        <h1 className="text-4xl font-bold mb-2">
          <span className={currentTheme.accent}>Agent</span> & <span className={currentTheme.accent}>Prompt</span> Crafter
          <span className="text-lg ml-2 px-2 py-1 bg-gradient-to-r from-purple-500 to-blue-500 text-white rounded-full">
            +MCP
          </span>
          {DemoMode.isActive && (
            <span className="text-xs ml-2 px-2 py-1 bg-blue-500 text-blue-100 rounded-full">
              DEMO v1.0.0
            </span>
          )}
        </h1>
        <p className={currentTheme.textSecondary}>
          Crie prompts profissionais e agents inteligentes com Model Context Protocol
        </p>
      </div>
      <div className="flex items-center gap-4">
        {DemoMode.isActive && (
          <div className="flex gap-2">
            <button
              onClick={startOnboarding}
              className="px-3 py-2 rounded-lg bg-green-600 text-white hover:bg-green-700 flex items-center gap-2 text-sm"
            >
              <Play size={16} />
              Tour
            </button>
            <button
              onClick={() => showEducation('mcp')}
              className="px-3 py-2 rounded-lg bg-purple-600 text-white hover:bg-purple-700 flex items-center gap-2 text-sm"
            >
              <BookOpen size={16} />
              O que é MCP?
            </button>
          </div>
        )}
        <select
          value="claude"
          className={`px-3 py-2 rounded-lg ${currentTheme.input} border focus:ring-2 focus:ring-blue-500`}
        >
          <option value="claude">Claude API ✅</option>
          <option disabled>Outros providers em breve...</option>
        </select>
        <button
          onClick={() => setDarkMode(!darkMode)}
          className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
        >
          {darkMode ? <Sun size={20} /> : <Moon size={20} />}
        </button>
      </div>
    </div>
  );
};

export default Header;
