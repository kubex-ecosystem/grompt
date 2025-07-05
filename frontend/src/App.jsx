import React, { useState } from 'react';
import { Routes, Route, Link, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Sparkles, Users, Sun, Moon } from 'lucide-react';
import PromptCrafter from './PromptCrafter';
import AgentsDashboard from './components/AgentsDashboard';
import AgentForm from './components/AgentForm';
import AgentView from './components/AgentView';
import LanguageSelector from './components/LanguageSelector';

const App = () => {
  const { t } = useTranslation();
  const location = useLocation();
  const [darkMode, setDarkMode] = useState(true);

  // Temas
  const theme = {
    dark: {
      bg: 'bg-gray-900',
      cardBg: 'bg-gray-800',
      text: 'text-gray-100',
      textSecondary: 'text-gray-300',
      border: 'border-gray-700',
      button: 'bg-blue-600 hover:bg-blue-700 text-white',
      buttonSecondary: 'bg-gray-700 hover:bg-gray-600 text-gray-200',
      accent: 'text-blue-400'
    },
    light: {
      bg: 'bg-gray-50',
      cardBg: 'bg-white',
      text: 'text-gray-900',
      textSecondary: 'text-gray-600',
      border: 'border-gray-300',
      button: 'bg-blue-600 hover:bg-blue-700 text-white',
      buttonSecondary: 'bg-gray-200 hover:bg-gray-300 text-gray-700',
      accent: 'text-blue-600'
    }
  };

  const currentTheme = darkMode ? theme.dark : theme.light;

  const isActivePath = (path) => {
    if (path === '/') return location.pathname === '/';
    return location.pathname.startsWith(path);
  };

  const getNavLinkClass = (path) => {
    const baseClass = 'flex items-center gap-3 px-4 py-2 rounded-lg transition-all duration-200 font-medium';
    const activeClass = 'bg-blue-600 text-white shadow-lg';
    const inactiveClass = `${currentTheme.textSecondary} hover:bg-blue-50 dark:hover:bg-gray-700 hover:text-blue-600 dark:hover:text-blue-400`;
    
    return `${baseClass} ${isActivePath(path) ? activeClass : inactiveClass}`;
  };

  return (
    <div className={`min-h-screen ${currentTheme.bg} ${currentTheme.text}`}>
      {/* Navigation Header */}
      <nav className={`${currentTheme.cardBg} ${currentTheme.border} border-b sticky top-0 z-50`}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            {/* Logo e Navigation Links */}
            <div className="flex items-center gap-8">
              <div className="flex items-center gap-3">
                <div className="bg-gradient-to-r from-blue-500 to-purple-600 p-2 rounded-lg">
                  <Sparkles className="h-6 w-6 text-white" />
                </div>
                <h1 className="text-xl font-bold bg-gradient-to-r from-blue-500 to-purple-600 bg-clip-text text-transparent">
                  Grompt
                </h1>
              </div>
              
              <div className="hidden md:flex items-center gap-2">
                <Link to="/" className={getNavLinkClass('/')}>
                  <Sparkles size={18} />
                  {t('nav.promptCrafter')}
                </Link>
                
                <Link to="/agents" className={getNavLinkClass('/agents')}>
                  <Users size={18} />
                  {t('nav.agents')}
                </Link>
              </div>
            </div>

            {/* Controls */}
            <div className="flex items-center gap-4">
              <LanguageSelector currentTheme={currentTheme} />
              
              <button
                onClick={() => setDarkMode(!darkMode)}
                className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
                title={darkMode ? t('theme.light') : t('theme.dark')}
              >
                {darkMode ? <Sun size={20} /> : <Moon size={20} />}
              </button>
            </div>
          </div>
          
          {/* Mobile Navigation */}
          <div className="md:hidden pb-4">
            <div className="flex flex-col gap-2">
              <Link to="/" className={getNavLinkClass('/')}>
                <Sparkles size={18} />
                {t('nav.promptCrafter')}
              </Link>
              
              <Link to="/agents" className={getNavLinkClass('/agents')}>
                <Users size={18} />
                {t('nav.agents')}
              </Link>
            </div>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main>
        <Routes>
          <Route path="/" element={<PromptCrafter />} />
          <Route path="/agents" element={<AgentsDashboard />} />
          <Route path="/agents/new" element={<AgentForm />} />
          <Route path="/agents/:id" element={<AgentView />} />
          <Route path="/agents/:id/edit" element={<AgentForm />} />
        </Routes>
      </main>
    </div>
  );
};

export default App;
