import * as React from 'react';

export type Screen = 'prompts' | 'agents';

interface ScreenNavigationProps {
  currentScreen: Screen;
  setCurrentScreen: (screen: Screen) => void;
}

const ScreenNavigation: React.FC<ScreenNavigationProps> = ({
  currentScreen,
  setCurrentScreen
}) => {
  const screens = [
    { id: 'prompts' as const, label: 'Prompts', icon: '📝', description: 'Geração de prompts' },
    { id: 'agents' as const, label: 'Agents', icon: '🤖', description: 'Gestão de agents' }
  ];

  return (
    <div className="bg-gray-800/50 border border-gray-700/80 rounded-xl p-4 backdrop-blur-sm mb-8">
      <div className="flex gap-2">
        {screens.map((screen) => (
          <button
            key={screen.id}
            onClick={() => setCurrentScreen(screen.id)}
            className={`flex-1 px-6 py-4 rounded-lg text-center border transition-all duration-300 ${
              currentScreen === screen.id
                ? 'bg-purple-600 text-white border-purple-600 shadow-lg shadow-purple-600/25'
                : 'bg-gray-700/80 border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white hover:border-purple-500/50'
            }`}
          >
            <div className="text-2xl mb-2">{screen.icon}</div>
            <div className="font-semibold">{screen.label}</div>
            <div className="text-sm opacity-75 mt-1">{screen.description}</div>
          </button>
        ))}
      </div>
    </div>
  );
};

export default ScreenNavigation;