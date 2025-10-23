import { Info, Settings, X } from 'lucide-react';
import * as React from 'react';
import IdeasInput from '../components/ideas/IdeasInput';
import IdeasList from '../components/ideas/IdeasList';
import OutputPanel from '../components/settings/OutputPanel';
import { DemoMode } from '../config/demoMode';
import { useGromptAPI } from '../hooks/useGromptAPI';
import { AgentFramework, Purpose } from '../hooks/usePromptCrafter';

interface Theme {
  [key: string]: string;
}

const {
  // State
  generatePrompt,
  providers,
  health,
  rateLimit
} = useGromptAPI();

interface AgentManagementScreenProps {
  // Ideas state
  currentInput: { id: string; text: string; timestamp: Date };
  setCurrentInput: (value: { id: string; text: string; timestamp: Date }) => void;
  ideas: Array<{ id: string; text: string; timestamp: Date }>;
  editingId: string | null;
  editingText: string;
  setEditingText: (value: string) => void;

  // Agent configuration
  agentFramework: AgentFramework;
  setAgentFramework: (value: AgentFramework) => void;
  agentProvider: string;
  setAgentProvider: (value: string) => void;
  agentRole: string;
  setAgentRole: (value: string) => void;
  agentTools: string[];
  setAgentTools: (value: string[] | ((prev: string[]) => string[])) => void;
  mcpServers: string[];
  setMcpServers: (value: string[] | ((prev: string[]) => string[])) => void;
  customMcpServer: string;
  setCustomMcpServer: (value: string) => void;
  purpose: Purpose;
  setPurpose: (value: Purpose) => void;
  customPurpose: string;
  setCustomPurpose: (value: string) => void;

  // Generation state
  generatedPrompt: string;
  isGenerating: boolean;
  copied: boolean;

  // Actions
  addIdea: (text: { id: string; text: string; timestamp: Date }) => void;
  removeIdea: (id: string) => void;
  startEditing: (id: string, text: string) => void;
  saveEdit: () => void;
  cancelEdit: () => void;
  generatePrompt: () => void;
  copyToClipboard: () => void;
  showEducation: (topic: string) => void;
  handleFeatureClick: (feature: string) => boolean;

  // Theme, API and providers
  currentTheme: Theme;
  apiGenerateState: typeof generatePrompt;
  apiRateLimitState: typeof rateLimit;
  apiHealthState: typeof health;
  apiRateLimit?: typeof rateLimit;
  apiHealth?: typeof health;
  apiGenerate?: typeof generatePrompt;
  apiProviders?: typeof providers;
}

const AgentManagementScreen: React.FC<AgentManagementScreenProps> = ({
  currentInput,
  setCurrentInput,
  ideas,
  editingId,
  editingText,
  setEditingText,
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
  generatedPrompt,
  isGenerating,
  copied,
  addIdea,
  removeIdea,
  startEditing,
  saveEdit,
  cancelEdit,
  generatePrompt,
  copyToClipboard,
  showEducation,
  handleFeatureClick,
  currentTheme,
  apiGenerateState,
  apiRateLimitState,
  apiHealthState
}) => {
  // Modal states
  const [showConfigModal, setShowConfigModal] = React.useState(false);
  const [showToolsModal, setShowToolsModal] = React.useState(false);
  const [showMcpModal, setShowMcpModal] = React.useState(false);

  const agentFrameworks: { value: AgentFramework; label: string }[] = [
    { value: 'crewai' as const, label: 'CrewAI' },
    { value: 'autogen' as const, label: 'AutoGen' },
    { value: 'langchain' as const, label: 'LangChain Agents' },
    { value: 'semantic-kernel' as const, label: 'Semantic Kernel' },
    { value: 'custom' as const, label: 'Agent Customizado' }
  ];

  const purposeOptions: Purpose[] = ['Automação', 'Análise', 'Suporte', 'Pesquisa', 'Outros'];

  const availableProviders = providers?.providers || [];
  const providersLoading = providers?.loading || false;
  const providersError = providers?.error;

  const tools: string[] = [
    'web_search', 'file_handler', 'calculator', 'email_sender',
    'database', 'api_caller', 'code_executor', 'image_generator',
    'git_ops', 'docker_manager'
  ];

  const mcpServersList: { name: string; desc: string }[] = [
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

  const handleToolToggle = (tool: string) => {
    setAgentTools(prev =>
      prev.includes(tool)
        ? prev.filter(t => t !== tool)
        : [...prev, tool]
    );
  };

  const handleMcpServerToggle = (serverName: string) => {
    if (DemoMode.isActive) {
      const demoResult = DemoMode.handleDemoCall('mcp_real');
      alert('🔌 ' + serverName + '\n\n' + demoResult.message + '\n\nETA: ' + demoResult.eta);
      return;
    }
    setMcpServers(prev =>
      prev.includes(serverName)
        ? prev.filter(s => s !== serverName)
        : [...prev, serverName]
    );
  };

  const handleCustomMcpAdd = () => {
    if (customMcpServer.trim()) {
      if (DemoMode.isActive) {
        const demoResult = DemoMode.handleDemoCall('mcp_real');
        alert('🔌 Servidor MCP Customizado\n\n' + demoResult.message + '\n\nETA: ' + demoResult.eta);
        return;
      }
      setMcpServers(prev => [...prev, customMcpServer.trim()]);
      setCustomMcpServer('');
    }
  };

  const handleMcpServerRemove = (server: string) => {
    setMcpServers(prev => prev.filter(s => s !== server));
  };

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="text-center">
        <h1 className="text-3xl font-bold mb-2 text-white">🤖 Gestão de Agents</h1>
        <p className="text-gray-400">Configure e gerencie seus agents de IA especializados</p>
      </div>

      {/* Quick Configuration Panel */}
      <div className="bg-gray-800/50 border border-gray-700/80 rounded-xl p-6 backdrop-blur-sm">
        <div className="flex flex-wrap gap-4 items-center justify-between">
          <div className="flex flex-wrap gap-4 items-center">
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-white">Framework:</span>
              <span className="px-3 py-1 rounded-lg bg-purple-600 text-white text-sm">
                {agentFrameworks.find(f => f.value === agentFramework)?.label}
              </span>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-white">Provider:</span>
              <span className="px-3 py-1 rounded-lg bg-blue-600 text-white text-sm">
                {agentProvider || 'Não selecionado'}
              </span>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-white">Tools:</span>
              <span className="px-3 py-1 rounded-lg bg-teal-600 text-white text-sm">
                {agentTools.length}
              </span>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-white">MCP:</span>
              <span className="px-3 py-1 rounded-lg bg-orange-600 text-white text-sm">
                {mcpServers.length}
              </span>
            </div>
          </div>
          <div className="flex gap-2">
            <button
              onClick={() => setShowConfigModal(true)}
              className="px-4 py-2 rounded-lg bg-gray-700/80 border border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white transition-colors flex items-center gap-2"
            >
              <Settings size={16} />
              Configurar
            </button>
          </div>
        </div>
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">

        {/* Left Column: Input and Role */}
        <div className="space-y-6">
          {/* Ideas Input Card */}
          <div className="bg-gray-800/50 border border-gray-700/80 rounded-xl p-6 backdrop-blur-sm transition-all duration-300 hover:border-purple-500/50">
            <IdeasInput
              currentInput={currentInput}
              setCurrentInput={setCurrentInput}
              addIdea={addIdea}
              currentTheme={currentTheme}
            />
          </div>

          {/* Agent Role Configuration */}
          <div className="bg-gray-800/50 border border-gray-700/80 rounded-xl p-6 backdrop-blur-sm transition-all duration-300 hover:border-purple-500/50">
            <h2 className="text-xl font-semibold mb-4 text-white">👤 Papel do Agent</h2>
            <div className="space-y-4">
              <div>
                <label className="text-sm font-medium mb-2 text-white flex items-center">Especialização</label>
                <input
                  type="text"
                  value={agentRole}
                  onChange={(e) => setAgentRole(e.target.value)}
                  placeholder="Ex: Especialista em Marketing Digital, Analista de Dados..."
                  className="w-full px-3 py-2 rounded-lg border border-gray-600 bg-gray-700/80 text-white placeholder-gray-400 focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
                />
              </div>

              {/* Purpose Selection */}
              <div>
                <label className="text-sm font-medium mb-2 text-white flex items-center">
                  Área de Atuação
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
                      placeholder="Descreva a área de atuação do agent..."
                      className="w-full px-3 py-2 rounded-lg border border-gray-600 bg-gray-700/80 text-white placeholder-gray-400 focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
                    />
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Center Column: Ideas List */}
        <div className="bg-gray-800/50 border border-gray-700/80 rounded-xl p-6 backdrop-blur-sm transition-all duration-300 hover:border-purple-500/50">
          <IdeasList
            ideas={ideas}
            editingId={editingId}
            editingText={editingText}
            setEditingText={setEditingText}
            startEditing={startEditing}
            saveEdit={saveEdit}
            cancelEdit={cancelEdit}
            removeIdea={removeIdea}
            generatePrompt={generatePrompt}
            isGenerating={isGenerating}
            outputType="agent"
            currentTheme={currentTheme}
            apiGenerateState={apiGenerateState}
          />
        </div>

        {/* Right Column: Output Panel */}
        <div className="bg-gray-800/50 border border-gray-700/80 rounded-xl p-6 backdrop-blur-sm transition-all duration-300 hover:border-purple-500/50">
          <OutputPanel
            generatedPrompt={generatedPrompt}
            copyToClipboard={copyToClipboard}
            copied={copied}
            outputType="agent"
            agentFramework={agentFramework}
            agentProvider={agentProvider}
            maxLength={0}
            mcpServers={mcpServers}
            currentTheme={currentTheme}
            apiGenerateState={apiGenerateState}
          />
        </div>
      </div>

      {/* Configuration Modal */}
      {showConfigModal && (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-gray-800 rounded-xl border border-gray-700 w-full max-w-4xl max-h-[90vh] overflow-y-auto">
            <div className="sticky top-0 bg-gray-800 p-6 border-b border-gray-700 flex items-center justify-between">
              <h2 className="text-xl font-bold text-white">⚙️ Configurações do Agent</h2>
              <button
                title='Close Configuration'
                onClick={() => setShowConfigModal(false)}
                className="p-2 hover:bg-gray-700 rounded-lg transition-colors"
              >
                <X size={20} className="text-gray-400" />
              </button>
            </div>

            <div className="p-6 space-y-6">
              {/* Framework and Provider */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium mb-2 text-white flex items-center">Framework do Agent</label>
                  <select
                    title='Selecione o framework do agent'
                    value={agentFramework}
                    onChange={(e) => setAgentFramework(e.target.value as AgentFramework)}
                    className="w-full px-3 py-2 rounded-lg border border-gray-600 bg-gray-700/80 text-white focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
                  >
                    {agentFrameworks.map((framework) => (
                      <option key={framework.value} value={framework.value}>
                        {DemoMode.getLabel(framework.value, framework.label)}
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <label className="text-sm font-medium mb-2 text-white flex items-center gap-2">
                    🤖 Provider LLM
                    {providersLoading && <div className="animate-spin w-4 h-4 border-2 border-blue-400 border-t-transparent rounded-full" />}
                  </label>

                  {providersError && (
                    <div className="mb-2 p-2 bg-red-900/50 border border-red-700 rounded text-red-400 text-xs">
                      Erro ao carregar providers: {providersError.message}
                    </div>
                  )}

                  <select
                    title='Selecione o provider do agent'
                    value={agentProvider}
                    onChange={(e) => setAgentProvider(e.target.value)}
                    className="w-full px-3 py-2 rounded-lg border border-gray-600 bg-gray-700/80 text-white focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
                    disabled={providersLoading}
                  >
                    {availableProviders.length > 0 ? (
                      availableProviders.map((provider) => (
                        <option key={provider.name} value={provider.name}>
                          {provider.available
                            ? `${provider.name} ${provider.defaultModel ? `(${provider.defaultModel})` : ''}`
                            : `${provider.name} (Indisponível)`
                          }
                        </option>
                      ))
                    ) : (
                      <option value="">Carregando providers...</option>
                    )}
                  </select>
                </div>
              </div>

              {/* Tools and MCP sections */}
              <div className="flex gap-4">
                <button
                  onClick={() => setShowToolsModal(true)}
                  className="flex-1 px-4 py-3 rounded-lg bg-teal-600/20 border border-teal-600/50 text-teal-300 hover:bg-teal-600/30 transition-colors flex items-center gap-2 justify-center"
                >
                  <span>🔧</span>
                  Ferramentas ({agentTools.length})
                </button>
                <button
                  onClick={() => setShowMcpModal(true)}
                  className="flex-1 px-4 py-3 rounded-lg bg-purple-600/20 border border-purple-600/50 text-purple-300 hover:bg-purple-600/30 transition-colors flex items-center gap-2 justify-center"
                >
                  <span>🔌</span>
                  MCP Servers ({mcpServers.length})
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Tools Modal */}
      {showToolsModal && (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-gray-800 rounded-xl border border-gray-700 w-full max-w-2xl">
            <div className="p-6 border-b border-gray-700 flex items-center justify-between">
              <h3 className="text-lg font-bold text-white">🔧 Ferramentas Tradicionais</h3>
              <button
                title='Close Configuration'
                onClick={() => setShowToolsModal(false)}
                className="p-2 hover:bg-gray-700 rounded-lg transition-colors"
              >
                <X size={20} className="text-gray-400" />
              </button>
            </div>

            <div className="p-6">
              <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
                {tools.map((tool) => (
                  <button
                    key={tool}
                    onClick={() => handleToolToggle(tool)}
                    className={`px-3 py-2 rounded-lg text-sm border transition-colors ${agentTools.includes(tool)
                      ? 'bg-teal-600 text-white border-teal-600'
                      : 'bg-gray-700/80 border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white'
                      }`}
                  >
                    {tool}
                  </button>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* MCP Modal */}
      {showMcpModal && (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-gray-800 rounded-xl border border-gray-700 w-full max-w-3xl max-h-[80vh] overflow-y-auto">
            <div className="sticky top-0 bg-gray-800 p-6 border-b border-gray-700 flex items-center justify-between">
              <div>
                <h3 className="text-lg font-bold text-white flex items-center gap-2">
                  🔌 Servidores MCP (Model Context Protocol)
                  {DemoMode.isActive && (
                    <button
                      title='Learn more about MCP'
                      onClick={() => showEducation('mcp')}
                      className="text-purple-400 hover:text-purple-300"
                    >
                      <Info size={16} />
                    </button>
                  )}
                </h3>
                <p className="text-xs text-purple-400">
                  Configure servidores MCP para estender as capacidades do agent
                </p>
              </div>
              <button
                title='Close Configuration'
                onClick={() => setShowMcpModal(false)}
                className="p-2 hover:bg-gray-700 rounded-lg transition-colors"
              >
                <X size={20} className="text-gray-400" />
              </button>
            </div>

            <div className="p-6 space-y-4">
              {/* Predefined MCP Servers */}
              <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
                {mcpServersList.map((server) => (
                  <button
                    key={server.name}
                    onClick={() => handleMcpServerToggle(server.name)}
                    className={`px-3 py-2 rounded-lg text-sm border transition-colors ${mcpServers.includes(server.name)
                      ? 'bg-purple-600 text-white border-purple-600'
                      : 'bg-gray-700/80 border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white'
                      }`}
                    title={server.desc + ' (demo)'}
                  >
                    {server.desc} 🎪
                  </button>
                ))}
              </div>

              {/* Custom MCP Server */}
              <div className="border-t border-gray-700 pt-4">
                <label className="text-sm font-medium mb-2 text-white block">Servidor Customizado</label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={customMcpServer}
                    onChange={(e) => setCustomMcpServer(e.target.value)}
                    placeholder="Ex: meu-servidor-personalizado"
                    className="flex-1 px-3 py-2 rounded-lg border border-gray-600 bg-gray-700/80 text-white placeholder-gray-400 focus:ring-2 focus:ring-purple-500 focus:border-purple-500 text-sm"
                  />
                  <button
                    onClick={handleCustomMcpAdd}
                    className="px-4 py-2 rounded-lg bg-gray-700/80 border border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white text-sm transition-colors"
                  >
                    Adicionar 🎪
                  </button>
                </div>
              </div>

              {/* Selected MCP Servers */}
              {mcpServers.length > 0 && (
                <div className="bg-purple-900/20 p-4 rounded-lg border border-purple-600/30">
                  <p className="text-sm font-medium text-purple-200 mb-2">
                    Servidores MCP selecionados:
                  </p>
                  <div className="flex flex-wrap gap-2">
                    {mcpServers.map((server) => (
                      <span
                        key={server}
                        className="inline-flex items-center gap-1 px-3 py-1 bg-purple-600 text-white rounded-full text-sm"
                      >
                        {server} 🎪
                        <button
                          onClick={() => handleMcpServerRemove(server)}
                          className="hover:bg-purple-700 rounded-full w-5 h-5 flex items-center justify-center ml-1"
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
    </div>
  );
};

export default AgentManagementScreen;
