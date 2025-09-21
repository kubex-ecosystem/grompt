import React from 'react';
import { Check, Copy, Download, Loader2, AlertCircle } from 'lucide-react';
import { DemoMode } from '../../config/demoMode';
import { OutputType, AgentFramework } from '../../hooks/usePromptCrafter';
import { UseGeneratePromptState } from '../../hooks/useGromptAPI';

interface Theme {
  [key: string]: string;
}

interface OutputPanelProps {
  generatedPrompt: string;
  copyToClipboard: () => void;
  copied: boolean;
  outputType: OutputType;
  agentFramework: AgentFramework;
  agentProvider: string;
  maxLength: number;
  mcpServers: string[];
  currentTheme: Theme;
  apiGenerateState?: UseGeneratePromptState;
}

const OutputPanel: React.FC<OutputPanelProps> = ({
  generatedPrompt,
  copyToClipboard,
  copied,
  outputType,
  agentFramework,
  agentProvider,
  maxLength,
  mcpServers,
  currentTheme,
  apiGenerateState
}) => {
  // Show API prompt if available, fallback to legacy prompt
  const promptToShow = apiGenerateState?.data?.prompt ||
                       apiGenerateState?.progress?.content ||
                       generatedPrompt;

  const isAPILoading = apiGenerateState?.loading;
  const isAPIStreaming = apiGenerateState?.progress?.isStreaming;
  const apiError = apiGenerateState?.error;
  const apiUsage = apiGenerateState?.data?.usage || apiGenerateState?.progress?.usage;

  const handleDownload = () => {
    const blob = new Blob([promptToShow], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `grompt-${outputType}-${Date.now()}.txt`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold text-white">
          {outputType === 'prompt' ? '🚀 Prompt Estruturado' : '🤖 Agent Gerado'}
        </h2>
        {promptToShow && (
          <div className="flex items-center gap-2">
            <span className={`text-xs px-2 py-1 rounded-full ${outputType === 'prompt'
                ? 'bg-purple-900/50 text-purple-200 border border-purple-700'
                : 'bg-green-900/50 text-green-200 border border-green-700'
              }`}>
              {outputType === 'prompt' ? 'Prompt' : agentFramework} {DemoMode.isActive ? '🎪' : ''}
            </span>

            {/* Loading indicator */}
            {(isAPILoading || isAPIStreaming) && (
              <div className="flex items-center gap-1 text-purple-400">
                <Loader2 size={14} className="animate-spin" />
                <span className="text-xs">{isAPIStreaming ? 'Streaming...' : 'Gerando...'}</span>
              </div>
            )}

            <button
              onClick={copyToClipboard}
              disabled={!promptToShow}
              className="flex items-center gap-2 px-3 py-2 rounded-lg bg-gray-700/80 border border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {copied ? <Check size={16} /> : <Copy size={16} />}
              {copied ? 'Copiado!' : 'Copiar'}
            </button>
          </div>
        )}
      </div>

      {/* Error display */}
      {apiError && (
        <div className="mb-4 p-3 bg-red-900/50 border border-red-700 rounded-lg flex items-center gap-2">
          <AlertCircle size={16} className="text-red-400" />
          <span className="text-red-400 text-sm">
            Erro na geração: {apiError.message}
          </span>
        </div>
      )}

      {promptToShow ? (
        <div className="space-y-4">
          {/* Stats and metadata */}
          <div className="text-xs text-gray-400 flex justify-between items-center">
            <span>Caracteres: {promptToShow.length.toLocaleString()}</span>
            <div className="flex items-center gap-4">
              {apiUsage && (
                <span className="bg-blue-900/50 text-blue-200 px-2 py-1 rounded-full border border-blue-700">
                  Tokens: {apiUsage.tokens} | Custo: ${apiUsage.costUSD?.toFixed(6) || '0.000000'}
                </span>
              )}
              {outputType === 'agent' && (
                <span className="bg-teal-900/50 text-teal-200 px-2 py-1 rounded-full border border-teal-700">
                  {agentFramework} + {agentProvider} + MCP
                </span>
              )}
              {outputType === 'prompt' && (
                <span className="text-gray-500">Limite: {maxLength.toLocaleString()}</span>
              )}
            </div>
          </div>

          {/* Content display */}
          <div className="max-h-96 overflow-y-auto p-4 rounded-lg border border-gray-600 bg-gray-800/50 backdrop-blur-sm">
            <pre className="whitespace-pre-wrap text-sm font-mono text-gray-300">{promptToShow}</pre>
          </div>

          {/* Agent info */}
          {outputType === 'agent' && (
            <div className="p-3 rounded-lg border border-purple-500/20 bg-gradient-to-r from-purple-900/20 to-blue-900/20 backdrop-blur-sm">
              <p className="text-sm text-purple-200">
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

          {/* Download button for large prompts */}
          {promptToShow.length > 1000 && (
            <button
              onClick={handleDownload}
              className="w-full flex items-center justify-center gap-2 px-4 py-2 rounded-lg bg-gray-700/80 border border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white transition-colors"
            >
              <Download size={16} />
              Baixar como arquivo
            </button>
          )}
        </div>
      ) : (
        <div className="text-gray-400 text-center py-12">
          <div className="text-4xl mb-4">🎯</div>
          <p className="mb-2 text-gray-300">
            {outputType === 'prompt'
              ? 'Seu prompt estruturado aparecerá aqui'
              : 'Seu agent gerado aparecerá aqui'
            }
          </p>
          <p className="text-sm text-gray-500">
            Adicione ideias e clique em "Criar {outputType === 'prompt' ? 'Prompt' : 'Agent'}" para começar
          </p>
        </div>
      )}
    </div>
  );
};

export default OutputPanel;