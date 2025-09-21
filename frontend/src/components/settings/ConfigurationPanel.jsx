import { Info } from 'lucide-react';
import DemoMode from '../config/demoMode.js';

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
  handleFeatureClick
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

  const agentProviders = [
    { value: 'claude', label: '🎭 Claude (Anthropic) ✅', available: true },
    { value: 'codex', label: '💻 Codex (OpenAI)', feature: 'openai' },
    { value: 'gpt4', label: '🧠 GPT-4 (OpenAI)', feature: 'openai' },
    { value: 'gemini', label: '💎 Gemini (Google)', feature: 'gemini' },
    { value: 'copilot', label: '🚁 GitHub Copilot', feature: 'copilot' },
    { value: 'ollama', label: '🦙 Ollama (Local)', feature: 'ollama' }
  ];

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
        <label className="block text-sm font-medium mb-2">Tipo de Saída</label>
        <div className="flex gap-2">
          {outputTypes.map((option) => (
            <button
              key={option.value}
              onClick={() => setOutputType(option.value)}
              className={`flex-1 px-4 py-3 rounded-lg text-sm border transition-all ${outputType === option.value
                  ? 'bg-blue-600 text-white border-blue-600 shadow-lg'
                  : `${currentTheme.buttonSecondary} ${currentTheme.border}`
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
        <div className="space-y-4 p-4 rounded-lg border border-blue-500/20 bg-blue-500/5" id="mcp-section">
          {/* Framework Selection */}
          <div>
            <label className="block text-sm font-medium mb-2">Framework do Agent</label>
            <select
              value={agentFramework}
              onChange={(e) => setAgentFramework(e.target.value)}
              className={`w-full px-3 py-2 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-blue-500`}
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
            <label className="block text-sm font-medium mb-2 flex items-center gap-2">
              🤖 Provider LLM
              {DemoMode.isActive && (
                <button
                  onClick={() => showEducation('agents')}
                  className="text-blue-500 hover:text-blue-600"
                >
                  <Info size={16} />
                </button>
              )}
            </label>
            <select
              value={agentProvider}
              onChange={(e) => {
                if (e.target.value !== 'claude' && DemoMode.isActive) {
                  handleFeatureClick(e.target.value);
                  return;
                }
                setAgentProvider(e.target.value);
              }}
              className={`w-full px-3 py-2 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-blue-500`}
            >
              {agentProviders.map((provider) => (
                <option key={provider.value} value={provider.value}>
                  {provider.available ? provider.label : DemoMode.getLabel(provider.feature, provider.label)}
                </option>
              ))}
            </select>
          </div>

          {/* Agent Role */}
          <div>
            <label className="block text-sm font-medium mb-2">Papel do Agent</label>
            <input
              type="text"
              value={agentRole}
              onChange={(e) => setAgentRole(e.target.value)}
              placeholder="Ex: Especialista em Marketing Digital, Analista de Dados..."
              className={`w-full px-3 py-2 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-blue-500`}
            />
          </div>

          {/* Traditional Tools */}
          <div>
            <label className="block text-sm font-medium mb-2">🔧 Ferramentas Tradicionais</label>
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
                      ? 'bg-green-600 text-white border-green-600'
                      : `${currentTheme.buttonSecondary} ${currentTheme.border}`
                    }`}
                >
                  {tool}
                </button>
              ))}
            </div>
          </div>

          {/* MCP Servers */}
          <div className="border-t border-blue-500/20 pt-4">
            <label className="block text-sm font-medium mb-2 flex items-center gap-2">
              🔌 Servidores MCP (Model Context Protocol)
              {DemoMode.isActive && (
                <button
                  onClick={() => showEducation('mcp')}
                  className="text-blue-500 hover:text-blue-600"
                >
                  <Info size={16} />
                </button>
              )}
            </label>
            <p className="text-xs text-blue-600 dark:text-blue-400 mb-3">
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
                        : `${currentTheme.buttonSecondary} ${currentTheme.border}`
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
                  className={`flex-1 px-3 py-2 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-blue-500 text-xs`}
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
                  className={`px-3 py-2 rounded-lg ${currentTheme.buttonSecondary} text-xs`}
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
        <label className="block text-sm font-medium mb-2">
          {outputType === 'prompt' ? 'Propósito do Prompt' : 'Área de Atuação do Agent'}
        </label>
        <div className="space-y-2">
          <div className="flex gap-2 flex-wrap">
            {purposeOptions.map((option) => (
              <button
                key={option}
                onClick={() => setPurpose(option)}
                className={`px-3 py-2 rounded-lg text-sm border transition-colors ${purpose === option
                    ? 'bg-blue-600 text-white border-blue-600'
                    : `${currentTheme.buttonSecondary} ${currentTheme.border}`
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
              className={`w-full px-3 py-2 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-blue-500`}
            />
          )}
        </div>
      </div>

      {/* Max Length for Prompts */}
      {outputType === 'prompt' && (
        <div>
          <label className="block text-sm font-medium mb-2">
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
