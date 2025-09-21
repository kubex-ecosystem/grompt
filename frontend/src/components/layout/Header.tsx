import { BookOpen, Moon, Play, Sun } from 'lucide-react';
import * as React from 'react';
import { DemoMode } from '../../config/demoMode';
import { Theme } from '../../constants/themes';
import { UseHealthState, UseProvidersState } from '../../hooks/useGromptAPI';

interface HeaderProps {
  darkMode: boolean;
  setDarkMode: (value: boolean) => void;
  currentTheme: Theme;
  startOnboarding: () => void;
  showEducation: (topic: string) => void;
  providers?: UseProvidersState;
  health?: UseHealthState;
  isHealthy: boolean;
}

const Header: React.FC<HeaderProps> = ({
  darkMode,
  setDarkMode,
  currentTheme,
  startOnboarding,
  showEducation,
  providers,
  health,
  isHealthy
}) => {
  const healthyProviders = providers?.providers?.filter(p => p.available) || [];
  const healthStatus = isHealthy ? 'healthy' : 'degraded';

  return (
    <div className="flex justify-between items-center mb-8" id="header">
      <div>
        <h1 className="text-4xl font-bold mb-2 text-white">
          <span className="text-purple-400">Grompt</span>{' '}
          <span className="text-teal-400">AI</span>
          <span className="text-lg ml-2 px-2 py-1 bg-gradient-to-r from-purple-500 to-blue-500 text-white rounded-full">
            +API
          </span>
          {DemoMode.isActive && (
            <span className="text-xs ml-2 px-2 py-1 bg-blue-500 text-blue-100 rounded-full">
              DEMO v1.0.0
            </span>
          )}
        </h1>
        <p className="text-gray-300">
          Crie prompts profissionais e agents inteligentes com Multi-Provider API
        </p>
        {providers && (
          <div className="mt-2 flex items-center gap-2 text-xs">
            <span className={`px-2 py-1 rounded-full ${healthStatus === 'healthy'
              ? 'bg-green-900/50 text-green-200 border border-green-700'
              : 'bg-yellow-900/50 text-yellow-200 border border-yellow-700'
              }`}>
              {healthyProviders.length} providers ativos
            </span>
            {healthyProviders.length > 0 && (
              <span className="text-gray-400">
                ({healthyProviders.map(p => p.name).join(', ')})
              </span>
            )}
          </div>
        )}
      </div>
      <div className="flex items-center gap-4">
        {DemoMode.isActive && (
          <div className="flex gap-2">
            <button
              onClick={startOnboarding}
              className="px-3 py-2 rounded-lg bg-green-600 text-white hover:bg-green-700 flex items-center gap-2 text-sm transition-colors"
            >
              <Play size={16} />
              Tour
            </button>
            <button
              onClick={() => showEducation('mcp')}
              className="px-3 py-2 rounded-lg bg-purple-600 text-white hover:bg-purple-700 flex items-center gap-2 text-sm transition-colors"
            >
              <BookOpen size={16} />
              O que é MCP?
            </button>
          </div>
        )}

        {/* Provider selector */}
        <select
          title='Selecione o provedor de IA'
          className="px-3 py-2 rounded-lg border border-gray-600 bg-gray-700/80 text-white focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
          defaultValue={healthyProviders[0]?.name || 'claude'}
        >
          {healthyProviders.length > 0 ? (
            healthyProviders.map(provider => (
              <option key={provider.name} value={provider.name}>
                {provider.name} ✅
              </option>
            ))
          ) : (
            <option value="claude">Claude API (configurar)</option>
          )}
          {providers?.providers?.filter(p => !p.available).map(provider => (
            <option key={provider.name} value={provider.name} disabled>
              {provider.name} (indisponível)
            </option>
          ))}
        </select>

        <button
          onClick={() => setDarkMode(!darkMode)}
          className="p-2 rounded-lg bg-gray-700/80 border border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white transition-colors"
        >
          {darkMode ? <Sun size={20} /> : <Moon size={20} />}
        </button>
      </div>
    </div>
  );
};

export default Header;
