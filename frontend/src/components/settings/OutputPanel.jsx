import { Check, Copy } from 'lucide-react';
import DemoMode from '../config/demoMode.js';

const OutputPanel = ({
  generatedPrompt,
  copyToClipboard,
  copied,
  outputType,
  agentFramework,
  agentProvider,
  maxLength,
  mcpServers,
  currentTheme
}) => {
  return (
    <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg`}>
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold">
          {outputType === 'prompt' ? '🚀 Prompt Estruturado' : '🤖 Agent Gerado'}
        </h2>
        {generatedPrompt && (
          <div className="flex items-center gap-2">
            <span className={`text-xs px-2 py-1 rounded-full ${outputType === 'prompt'
                ? 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200'
                : 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
              }`}>
              {outputType === 'prompt' ? 'Prompt' : agentFramework} {DemoMode.isActive ? '🎪' : ''}
            </span>
            <button
              onClick={copyToClipboard}
              className={`flex items-center gap-2 px-3 py-2 rounded-lg ${currentTheme.buttonSecondary} hover:bg-opacity-80 transition-colors`}
            >
              {copied ? <Check size={16} /> : <Copy size={16} />}
              {copied ? 'Copiado!' : 'Copiar'}
            </button>
          </div>
        )}
      </div>

      {generatedPrompt ? (
        <div className="space-y-4">
          <div className={`text-xs ${currentTheme.textSecondary} flex justify-between items-center`}>
            <span>Caracteres: {generatedPrompt.length.toLocaleString()}</span>
            <div className="flex items-center gap-4">
              {outputType === 'agent' && (
                <span className="bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200 px-2 py-1 rounded-full">
                  {agentFramework} + {agentProvider} + MCP
                </span>
              )}
              {outputType === 'prompt' && (
                <span>Limite: {maxLength.toLocaleString()}</span>
              )}
            </div>
          </div>
          <div className={`max-h-96 overflow-y-auto p-4 rounded-lg border ${currentTheme.border} bg-opacity-50`}>
            <pre className="whitespace-pre-wrap text-sm font-mono">{generatedPrompt}</pre>
          </div>

          {outputType === 'agent' && (
            <div className={`p-3 rounded-lg border ${currentTheme.border} bg-gradient-to-r from-purple-50 to-blue-50 dark:from-purple-900/20 dark:to-blue-900/20`}>
              <p className="text-sm text-purple-800 dark:text-purple-200">
                🚀 <strong>Agent Avançado:</strong> Integração com {agentProvider} + MCP + Config TOML
                {mcpServers.length > 0 && (
                  <span className="block mt-1">
                    🔌 <strong>Servidores MCP:</strong> {mcpServers.slice(0, 3).join(', ')}
                    {mcpServers.length > 3 && ` +${mcpServers.length - 3} mais`}
                  </span>
                )}
              </p>
            </div>
          )}
        </div>
      ) : (
        <div className={`${currentTheme.textSecondary} text-center py-12`}>
          <div className="text-4xl mb-4">🎯</div>
          <p className="mb-2">
            {outputType === 'prompt'
              ? 'Seu prompt estruturado aparecerá aqui'
              : 'Seu agent gerado aparecerá aqui'
            }
          </p>
          <p className="text-sm">
            Adicione ideias e clique em "Criar {outputType === 'prompt' ? 'Prompt' : 'Agent'}" para começar
          </p>
        </div>
      )}
    </div>
  );
};

export default OutputPanel;
