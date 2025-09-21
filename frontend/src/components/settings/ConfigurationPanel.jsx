import DemoMode from '@/config/DemoMode';
import { Info, Loader2, Wifi, WifiOff } from 'lucide-react';
import * as React from 'react';

const ConfigurationPanel = ({
  outputType,
  setOutputType,
  agentFramework,
  setAgentFramework,
  agentProvider,
  setAgentProvider,
  agentRole,
  setAgentRole,
  agentTools,
  setAgentTools,
  mcpServers,
  setMcpServers,
  customMcpServer,
  setCustomMcpServer,
  purpose,
  setPurpose,
  customPurpose,
  setCustomPurpose,
  maxLength,
  setMaxLength,
  currentTheme,
  showEducation,
  handleFeatureClick,
  providers
}) => {
  const outputTypes = [
    { value: 'prompt', label: '📝 Prompt', icon: '📝' },
    { value: 'agent', label: '🤖 Agent', icon: '🤖' }
  ];

  const agentFrameworks = [
    { value: 'crewai', label: 'CrewAI' },
    { value: 'autogen', label: 'AutoGen' },
    { value: 'langchain', label: 'LangChain Agents' },
    { value: 'semantic-kernel', label: 'Semantic Kernel' },
    { value: 'custom', label: 'Agent Customizado' }
  ];

  // Use actual providers from API with fallback
  const availableProviders = providers?.providers || [];
  const providersLoading = providers?.loading || false;
  const providersError = providers?.error;

  const tools = [
    'web_search', 'file_handler', 'calculator', 'email_sender',
    'database', 'api_caller', 'code_executor', 'image_generator',
    'git_ops', 'docker_manager'
  ];

  const mcpServersList = [
    { name: 'filesystem', desc: '📁 Sistema de arquivos' },
    { name: 'database', desc: '🗄️ Banco de dados' },
    { name: 'web-scraper', desc: '🕷️ Web scraping' },
    { name: 'git', desc: '🔄 Controle de versão' },
    { name: 'docker', desc: '🐳 Containers' },
    { name: 'kubernetes', desc: '☸️ Kubernetes' },
    { name: 'slack', desc: '💬 Slack' },
    { name: 'github', desc: '🐙 GitHub' },
    { name: 'notion', desc: '📝 Notion' },
    { name: 'calendar', desc: '📅 Calendário' }
  ];

  const purposeOptions = outputType === 'prompt'
    ? ['Código', 'Imagem', 'Análise', 'Escrita', 'Outros']
    : ['Automação', 'Análise', 'Suporte', 'Pesquisa', 'Outros'];

  return (
    <div className="space-y-4">
      {/* Output Type Selector */}
      <div id="output-selector">
        <label className="block text-sm font-medium mb-2 text-white">Tipo de Saída</label>
        <div className="flex gap-2">
          {outputTypes.map((option) => (
            <button
              key={option.value}
              onClick={() => setOutputType(option.value)}
              className={`flex-1 px-4 py-3 rounded-lg text-sm border transition-all ${outputType === option.value
                ? 'bg-purple-600 text-white border-purple-600 shadow-lg'
                : 'bg-gray-700/80 border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white hover:border-purple-500/50'
                }`}
            >
              <div className="text-center">
                <div className="text-lg mb-1">{option.icon}</div>
                <div>{option.label.split(' ')[1]}</div>
              </div>
            </button>
          ))}
        </div>
      </div>

      {/* Agent Configuration */}
      {outputType === 'agent' && (
        <div className="space-y-4 p-4 rounded-lg border border-purple-500/20 bg-purple-500/5 backdrop-blur-sm" id="mcp-section">
          {/* Framework Selection */}
          <div>
            <label className="block text-sm font-medium mb-2 text-white">Framework do Agent</label>
            <select
              value={agentFramework}
              onChange={(e) => setAgentFramework(e.target.value)}
              className="w-full px-3 py-2 rounded-lg border border-gray-600 bg-gray-700/80 text-white focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
            >
              {agentFrameworks.map((framework) => (
                <option key={framework.value} value={framework.value}>
                  {DemoMode.getLabel(framework.value, framework.label)}
                </option>
              ))}
            </select>
          </div>

          {/* Provider Selection */}
          <div>
            <label className="block text-sm font-medium mb-2 flex items-center gap-2 text-white">
              🤖 Provider LLM
              {providersLoading && <Loader2 size={16} className="animate-spin text-blue-400" />}
              {DemoMode.isActive && (
                <button
                  onClick={() => showEducation('agents')}
                  className="text-blue-400 hover:text-blue-300"
                >
                  <Info size={16} />
                </button>
              )}
            </label>

            {/* Provider Status */}
            {providersError && (
              <div className="mb-2 p-2 bg-red-900/50 border border-red-700 rounded text-red-400 text-xs">
                Erro ao carregar providers: {providersError.message}
              </div>
            )}

            <select
              value={agentProvider}
              onChange={(e) => setAgentProvider(e.target.value)}
              className="w-full px-3 py-2 rounded-lg border border-gray-600 bg-gray-700/80 text-white focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
              disabled={providersLoading}
            >
              {availableProviders.length > 0 ? (
                availableProviders.map((provider) => (
                  <option key={provider.name} value={provider.name}>
                    {provider.available ? (
                      <span className="flex items-center gap-2">
                        <Wifi size={12} />
                        {provider.name} {provider.defaultModel ? `(${provider.defaultModel})` : ''}
                      </span>
                    ) : (
                      <span className="flex items-center gap-2">
                        <WifiOff size={12} />
                        {provider.name} (Indisponível)
                      </span>
                    )}
                  </option>
                ))
              ) : (
                <option value="">Carregando providers...</option>
              )}
            </select>

            {/* Provider Info */}
            {availableProviders.length > 0 && (
              <div className="mt-2 text-xs text-gray-400">
                {availableProviders.filter(p => p.available).length} de {availableProviders.length} providers disponíveis
              </div>
            )}
          </div>

          {/* Agent Role */}
          <div>
            <label className="block text-sm font-medium mb-2 text-white">Papel do Agent</label>
            <input
              type="text"
              value={agentRole}
              onChange={(e) => setAgentRole(e.target.value)}
              placeholder="Ex: Especialista em Marketing Digital, Analista de Dados..."
              className="w-full px-3 py-2 rounded-lg border border-gray-600 bg-gray-700/80 text-white placeholder-gray-400 focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
            />
          </div>

          {/* Traditional Tools */}
          <div>
            <label className="block text-sm font-medium mb-2 text-white">🔧 Ferramentas Tradicionais</label>
            <div className="flex flex-wrap gap-2 mb-2">
              {tools.map((tool) => (
                <button
                  key={tool}
                  onClick={() => {
                    setAgentTools(prev =>
                      prev.includes(tool)
                        ? prev.filter(t => t !== tool)
                        : [...prev, tool]
                    );
                  }}
                  className={`px-3 py-1 rounded-full text-xs border transition-colors ${agentTools.includes(tool)
                    ? 'bg-teal-600 text-white border-teal-600'
                    : 'bg-gray-700/80 border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white'
                    }`}
                >
                  {tool}
                </button>
              ))}
            </div>
          </div>

          {/* MCP Servers */}
          <div className="border-t border-blue-500/20 pt-4">
            <label className="block text-sm font-medium mb-2 flex items-center gap-2 text-white">
              🔌 Servidores MCP (Model Context Protocol)
              {DemoMode.isActive && (
                <button
                  onClick={() => showEducation('mcp')}
                  className="text-purple-400 hover:text-purple-300"
                >
                  <Info size={16} />
                </button>
              )}
            </label>
            <p className="text-xs text-purple-400 mb-3">
              Configure servidores MCP para estender as capacidades do agent
            </p>

            <div className="space-y-3">
              <div className="flex flex-wrap gap-2">
                {mcpServersList.map((server) => (
                  <button
                    key={server.name}
                    onClick={() => {
                      if (DemoMode.isActive) {
                        const demoResult = DemoMode.handleDemoCall('mcp_real');
                        alert('🔌 ' + server.desc + '\n\n' + demoResult.message + '\n\nETA: ' + demoResult.eta);
                        return;
                      }
                      setMcpServers(prev =>
                        prev.includes(server.name)
                          ? prev.filter(s => s !== server.name)
                          : [...prev, server.name]
                      );
                    }}
                    className={`px-3 py-2 rounded-lg text-xs border transition-colors ${mcpServers.includes(server.name)
                      ? 'bg-purple-600 text-white border-purple-600'
                      : 'bg-gray-700/80 border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white'
                      }`}
                    title={server.desc + ' (demo)'}
                  >
                    {server.desc} 🎪
                  </button>
                ))}
              </div>

              <div className="flex gap-2">
                <input
                  type="text"
                  value={customMcpServer}
                  onChange={(e) => setCustomMcpServer(e.target.value)}
                  placeholder="Servidor MCP customizado (ex: meu-servidor-personalizado)"
                  className="flex-1 px-3 py-2 rounded-lg border border-gray-600 bg-gray-700/80 text-white placeholder-gray-400 focus:ring-2 focus:ring-purple-500 focus:border-purple-500 text-xs"
                />
                <button
                  onClick={() => {
                    if (customMcpServer.trim()) {
                      if (DemoMode.isActive) {
                        const demoResult = DemoMode.handleDemoCall('mcp_real');
                        alert('🔌 Servidor MCP Customizado\n\n' + demoResult.message + '\n\nETA: ' + demoResult.eta);
                        return;
                      }
                      setMcpServers(prev => [...prev, customMcpServer.trim()]);
                      setCustomMcpServer('');
                    }
                  }}
                  className="px-3 py-2 rounded-lg bg-gray-700/80 border border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white text-xs transition-colors"
                >
                  + Adicionar 🎪
                </button>
              </div>

              {mcpServers.length > 0 && (
                <div className="bg-purple-50 dark:bg-purple-900/20 p-3 rounded-lg">
                  <p className="text-xs font-medium text-purple-800 dark:text-purple-200 mb-2">
                    Servidores MCP selecionados:
                  </p>
                  <div className="flex flex-wrap gap-1">
                    {mcpServers.map((server) => (
                      <span
                        key={server}
                        className="inline-flex items-center gap-1 px-2 py-1 bg-purple-600 text-white rounded-full text-xs"
                      >
                        {server} 🎪
                        <button
                          onClick={() => setMcpServers(prev => prev.filter(s => s !== server))}
                          className="hover:bg-purple-700 rounded-full w-4 h-4 flex items-center justify-center"
                        >
                          ×
                        </button>
                      </span>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Purpose Selection */}
      <div>
        <label className="block text-sm font-medium mb-2 text-white">
          {outputType === 'prompt' ? 'Propósito do Prompt' : 'Área de Atuação do Agent'}
        </label>
        <div className="space-y-2">
          <div className="flex gap-2 flex-wrap">
            {purposeOptions.map((option) => (
              <button
                key={option}
                onClick={() => setPurpose(option)}
                className={`px-3 py-2 rounded-lg text-sm border transition-colors ${purpose === option
                  ? 'bg-purple-600 text-white border-purple-600'
                  : 'bg-gray-700/80 border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white'
                  }`}
              >
                {option}
              </button>
            ))}
          </div>
          {purpose === 'Outros' && (
            <input
              type="text"
              value={customPurpose}
              onChange={(e) => setCustomPurpose(e.target.value)}
              placeholder={
                outputType === 'prompt'
                  ? "Descreva o objetivo do prompt..."
                  : "Descreva a área de atuação do agent..."
              }
              className="w-full px-3 py-2 rounded-lg border border-gray-600 bg-gray-700/80 text-white placeholder-gray-400 focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
            />
          )}
        </div>
      </div>

      {/* Max Length for Prompts */}
      {outputType === 'prompt' && (
        <div>
          <label className="block text-sm font-medium mb-2 text-white">
            Tamanho Máximo: {maxLength.toLocaleString()} caracteres
          </label>
          <input
            type="range"
            min="500"
            max="130000"
            step="500"
            value={maxLength}
            onChange={(e) => setMaxLength(parseInt(e.target.value))}
            className="w-full"
          />
        </div>
      )}
    </div>
  );
};

export default ConfigurationPanel;
