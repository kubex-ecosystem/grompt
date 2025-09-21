<BookOpen size={16} />
                    Sobre Prompts
                  </button >
                </div >
              )}
<button
  onClick={() => setDarkMode(!darkMode)}
  className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
>
  {darkMode ? <Sun size={20} /> : <Moon size={20} />}
</button>
            </div >
          </div >

  <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
    {/* Input Section */}
    <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg`} id="ideas-input">
      <h2 className="text-xl font-semibold mb-4">📝 Suas Ideias</h2>
      <div className="space-y-4">
        <textarea
          value={currentInput}
          onChange={(e) => setCurrentInput(e.target.value)}
          placeholder="Digite suas ideias, requisitos ou pensamentos aqui..."
          className={`w-full h-32 px-4 py-3 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-purple-500 resize-none`}
          onKeyDown={(e) => {
            if (e.key === 'Enter' && e.ctrlKey) {
              addIdea();
            }
          }}
        />
        <button
          onClick={addIdea}
          disabled={!currentInput.trim()}
          className={`w-full flex items-center justify-center gap-2 px-4 py-3 rounded-lg bg-purple-600 hover:bg-purple-700 text-white disabled:opacity-50 disabled:cursor-not-allowed transition-all`}
        >
          <Plus size={20} />
          Adicionar Ideia (Ctrl+Enter)
        </button>
      </div>

      {/* Purpose Section */}
      <div className="mt-6" id="purpose-section">
        <label className="block text-sm font-medium mb-2">Propósito do Prompt</label>
        <div className="space-y-2">
          <div className="flex gap-2 flex-wrap">
            {['Código', 'Imagem', 'Análise', 'Escrita', 'Educacional', 'Outros'].map((option) => (
              <button
                key={option}
                onClick={() => setPromptPurpose(option)}
                className={`px-3 py-2 rounded-lg text-sm border transition-colors ${promptPurpose === option
                  ? 'bg-purple-600 text-white border-purple-600'
                  : `${currentTheme.buttonSecondary} ${currentTheme.border}`
                  }`}
              >
                {option}
              </button>
            ))}
          </div>
          {promptPurpose === 'Outros' && (
            <input
              type="text"
              value={customPromptPurpose}
              onChange={(e) => setCustomPromptPurpose(e.target.value)}
              placeholder="Descreva o propósito específico..."
              className={`w-full px-3 py-2 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-purple-500`}
            />
          )}
        </div>
      </div>

      {/* Length Control */}
      <div className="mt-6">
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
          className="w-full h-2 bg-gray-300 rounded-lg appearance-none cursor-pointer slider"
        />
      </div>
    </div>

    {/* Ideas List */}
    <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg`}>
      <h2 className="text-xl font-semibold mb-4">💡 Lista de Ideias ({ideas.length})</h2>
      <div className="space-y-3 max-h-96 overflow-y-auto">
        {ideas.length === 0 ? (
          <p className={`${currentTheme.textSecondary} text-center py-8`}>
            Adicione suas primeiras ideias ao lado ←
          </p>
        ) : (
          ideas.map((idea) => (
            <div key={idea.id} className={`p-3 rounded-lg border ${currentTheme.border} bg-opacity-50`}>
              {editingId === idea.id ? (
                <div className="space-y-2">
                  <textarea
                    value={editingText}
                    onChange={(e) => setEditingText(e.target.value)}
                    className={`w-full px-2 py-1 rounded border ${currentTheme.input} text-sm`}
                    rows="2"
                  />
                  <div className="flex gap-1">
                    <button
                      onClick={saveEdit}
                      className="px-2 py-1 bg-green-600 text-white rounded text-xs hover:bg-green-700"
                    >
                      Salvar
                    </button>
                    <button
                      onClick={cancelEdit}
                      className={`px-2 py-1 rounded text-xs ${currentTheme.buttonSecondary}`}
                    >
                      Cancelar
                    </button>
                  </div>
                </div>
              ) : (
                <>
                  <p className="text-sm mb-2">{idea.text}</p>
                  <div className="flex justify-end gap-1">
                    <button
                      onClick={() => startEditing(idea.id, idea.text)}
                      className={`p-1 rounded ${currentTheme.buttonSecondary} hover:bg-opacity-80`}
                    >
                      <Edit3 size={14} />
                    </button>
                    <button
                      onClick={() => removeIdea(idea.id)}
                      className="p-1 rounded bg-red-600 text-white hover:bg-red-700"
                    >
                      <Trash2 size={14} />
                    </button>
                  </div>
                </>
              )}
            </div>
          ))
        )}
      </div>

      {ideas.length > 0 && (
        <button
          onClick={generateContent}
          disabled={isGenerating}
          className={`w-full mt-4 flex items-center justify-center gap-2 px-4 py-3 rounded-lg bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700 text-white disabled:opacity-50 disabled:cursor-not-allowed transition-all transform hover:scale-105`}
        >
          <Wand2 size={20} className={isGenerating ? 'animate-spin' : ''} />
          {isGenerating ? 'Gerando prompt...' : 'Gerar Prompt Profissional 🚀'}
        </button>
      )}
    </div>

    {/* Generated Prompt */}
    <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg`}>
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold">🚀 Prompt Estruturado</h2>
        {generatedContent && (
          <div className="flex items-center gap-2">
            <span className="text-xs px-2 py-1 rounded-full bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200">
              Prompt Profissional
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

      {generatedContent ? (
        <div className="space-y-4">
          <div className={`text-xs ${currentTheme.textSecondary} flex justify-between items-center`}>
            <span>Caracteres: {generatedContent.length.toLocaleString()}</span>
            <span>Limite: {maxLength.toLocaleString()}</span>
          </div>
          <div className={`max-h-96 overflow-y-auto p-4 rounded-lg border ${currentTheme.border} bg-opacity-50`}>
            <pre className="whitespace-pre-wrap text-sm font-mono">{generatedContent}</pre>
          </div>
        </div>
      ) : (
        <div className={`${currentTheme.textSecondary} text-center py-12`}>
          <Wand2 size={48} className="mx-auto mb-4 opacity-50" />
          <p>Seu prompt estruturado aparecerá aqui</p>
          <p className="text-sm mt-2">Adicione ideias e clique em "Gerar Prompt Profissional 🚀"</p>
        </div>
      )}
    </div>
  </div>
        </div >

  {/* Onboarding Modal */ }
{
  showOnboarding && (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
      <div className={`${currentTheme.cardBg} rounded-xl p-6 max-w-md border ${currentTheme.border} shadow-xl`}>
        <h3 className="text-xl font-bold mb-4">{OnboardingSteps.prompt[currentStep].title}</h3>
        <p className={`${currentTheme.textSecondary} mb-6`}>{OnboardingSteps.prompt[currentStep].content}</p>
        <div className="flex justify-between">
          <span className="text-sm text-gray-500">
            {currentStep + 1} de {OnboardingSteps.prompt.length}
          </span>
          <button
            onClick={nextOnboardingStep}
            className={`px-4 py-2 rounded-lg ${currentTheme.button}`}
          >
            {currentStep < OnboardingSteps.prompt.length - 1 ? 'Próximo' : 'Finalizar'}
          </button>
        </div>
      </div>
    </div>
  )
}

{/* Educational Modal */ }
{
  showEducational && educationalTopic && (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
      <div className={`${currentTheme.cardBg} rounded-xl p-6 max-w-lg border ${currentTheme.border} shadow-xl`}>
        <h3 className="text-xl font-bold mb-4">{DemoMode.education[educationalTopic].title}</h3>
        <p className={`${currentTheme.textSecondary} mb-4`}>{DemoMode.education[educationalTopic].description}</p>
        <div className="mb-6">
          <h4 className="font-semibold mb-2">Benefícios:</h4>
          <ul className="space-y-1">
            {DemoMode.education[educationalTopic].benefits.map((benefit, index) => (
              <li key={index} className={currentTheme.textSecondary}>{benefit}</li>
            ))}
          </ul>
        </div>
        <button
          onClick={() => setShowEducational(false)}
          className={`px-4 py-2 rounded-lg ${currentTheme.button} w-full`}
        >
          Entendi!
        </button>
      </div>
    </div>
  )
}
      </div >
    );
  }

// AGENT CRAFTER SCREEN
if (currentScreen === 'agent') {
  return (
    <div className={`min-h-screen ${currentTheme.bg} ${currentTheme.text} p-4 transition-colors duration-300`}>
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-8" id="header">
          <div className="flex items-center gap-4">
            <button
              onClick={() => navigateToScreen('home')}
              className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
            >
              <ArrowLeft size={20} />
            </button>
            <div>
              <h1 className="text-4xl font-bold mb-2 flex items-center gap-3">
                <Bot className="text-green-500" size={40} />
                <span className="text-green-500">Agent</span> Crafter
              </h1>
              <p className={currentTheme.textSecondary}>
                Crie agents inteligentes com MCP e configurações profissionais
              </p>
            </div>
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
                  onClick={() => showEducation('agents')}
                  className="px-3 py-2 rounded-lg bg-purple-600 text-white hover:bg-purple-700 flex items-center gap-2 text-sm"
                >
                  <BookOpen size={16} />
                  Sobre Agents
                </button>
                <button
                  onClick={() => showEducation('mcp')}
                  className="px-3 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 flex items-center gap-2 text-sm"
                >
                  <Info size={16} />
                  O que é MCP?
                </button>
              </div>
            )}
            <button
              onClick={() => setDarkMode(!darkMode)}
              className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
            >
              {darkMode ? <Sun size={20} /> : <Moon size={20} />}
            </button>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Configuration Section */}
          <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg`} id="framework-section">
            <h2 className="text-xl font-semibold mb-4">⚙️ Configuração do Agent</h2>

            {/* Framework */}
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2">Framework</label>
                <select
                  value={agentFramework}
                  onChange={(e) => setAgentFramework(e.target.value)}
                  className={`w-full px-3 py-2 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-green-500`}
                >
                  <option value="crewai">{DemoMode.getLabel('crewai', 'CrewAI')}</option>
                  <option value="autogen">{DemoMode.getLabel('autogen', 'AutoGen')}</option>
                  <option value="langchain">{DemoMode.getLabel('langchain', 'LangChain Agents')}</option>
                  <option value="semantic-kernel">{DemoMode.getLabel('semantic-kernel', 'Semantic Kernel')}</option>
                  <option value="custom">{DemoMode.getLabel('custom', 'Agent Customizado')}</option>
                </select>
              </div>

              {/* Provider */}
              <div>
                <label className="block text-sm font-medium mb-2">Provider LLM</label>
                <select
                  value={agentProvider}
                  onChange={(e) => {
                    if (e.target.value !== 'claude' && DemoMode.isActive) {
                      handleFeatureClick(e.target.value);
                      return;
                    }
                    setAgentProvider(e.target.value);
                  }}
                  className={`w-full px-3 py-2 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-green-500`}
                >
                  <option value="claude">🎭 Claude (Anthropic) ✅</option>
                  <option value="codex">{DemoMode.getLabel('openai', '💻 Codex (OpenAI)')}</option>
                  <option value="gpt4">{DemoMode.getLabel('openai', '🧠 GPT-4 (OpenAI)')}</option>
                  <option value="gemini">{DemoMode.getLabel('gemini', '💎 Gemini (Google)')}</option>
                  <option value="copilot">{DemoMode.getLabel('copilot', '🚁 GitHub Copilot')}</option>
                  <option value="ollama">{DemoMode.getLabel('ollama', '🦙 Ollama (Local)')}</option>
                </select>
              </div>

              {/* Role */}
              <div>
                <label className="block text-sm font-medium mb-2">Papel do Agent</label>
                <input
                  type="text"
                  value={agentRole}
                  onChange={(e) => setAgentRole(e.target.value)}
                  placeholder="Ex: Especialista em DevOps, Analista de Dados..."
                  className={`w-full px-3 py-2 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-green-500`}
                />
              </div>

              {/* Purpose */}
              <div>
                <label className="block text-sm font-medium mb-2">Área de Atuação</label>
                <div className="space-y-2">
                  <div className="flex gap-2 flex-wrap">
                    {['Automação', 'Análise', 'Suporte', 'Pesquisa', 'DevOps', 'Outros'].map((option) => (
                      <button
                        key={option}
                        onClick={() => setAgentPurpose(option)}
                        className={`px-3 py-2 rounded-lg text-sm border transition-colors ${agentPurpose === option
                          ? 'bg-green-600 text-white border-green-600'
                          : `${currentTheme.buttonSecondary} ${currentTheme.border}`
                          }`}
                      >
                        {option}
                      </button>
                    ))}
                  </div>
                  {agentPurpose === 'Outros' && (
                    <input
                      type="text"
                      value={customAgentPurpose}
                      onChange={(e) => setCustomAgentPurpose(e.target.value)}
                      placeholder="Descreva a área de atuação..."
                      className={`w-full px-3 py-2 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-green-500`}
                    />
                  )}
                </div>
              </div>

              {/* Tools */}
              <div>
                <label className="block text-sm font-medium mb-2">🔧 Ferramentas</label>
                <div className="flex flex-wrap gap-2 mb-2">
                  {['web_search', 'file_handler', 'calculator', 'email_sender', 'database', 'api_caller', 'code_executor', 'docker_manager'].map((tool) => (
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
            </div>
          </div>

          {/* Ideas + MCP Section */}
          <div className="space-y-6">
            {/* Ideas */}
            <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg`} id="ideas-input">
              <h2 className="text-xl font-semibold mb-4">💡 Requisitos do Agent</h2>
              <div className="space-y-4">
                <textarea
                  value={currentInput}
                  onChange={(e) => setCurrentInput(e.target.value)}
                  placeholder="Descreva o que o agent deve fazer, suas responsabilidades..."
                  className={`w-full h-24 px-4 py-3 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-green-500 resize-none`}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter' && e.ctrlKey) {
                      addIdea();
                    }
                  }}
                />
                <button
                  onClick={addIdea}
                  disabled={!currentInput.trim()}
                  className={`w-full flex items-center justify-center gap-2 px-4 py-2 rounded-lg bg-green-600 hover:bg-green-700 text-white disabled:opacity-50 disabled:cursor-not-allowed transition-all`}
                >
                  <Plus size={18} />
                  Adicionar (Ctrl+Enter)
                </button>
              </div>

              {/* Ideas List */}
              <div className="space-y-2 mt-4 max-h-48 overflow-y-auto">
                {ideas.length === 0 ? (
                  <p className={`${currentTheme.textSecondary} text-center py-4 text-sm`}>
                    Adicione requisitos para o agent
                  </p>
                ) : (
                  ideas.map((idea) => (
                    <div key={idea.id} className={`p-2 rounded border ${currentTheme.border} bg-opacity-50`}>
                      <p className="text-sm">{idea.text}</p>
                      <div className="flex justify-end mt-1">
                        <button
                          onClick={() => removeIdea(idea.id)}
                          className="text-red-500 hover:text-red-700 p-1"
                        >
                          <Trash2 size={12} />
                        </button>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>

            {/* MCP Servers */}
            <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg`} id="mcp-section">
              <h2 className="text-xl font-semibold mb-4">🔌 Servidores MCP</h2>
              <p className="text-xs text-blue-600 dark:text-blue-400 mb-3">
                Model Context Protocol - conecte com sistemas externos
              </p>

              <div className="space-y-3">
                <div className="flex flex-wrap gap-2">
                  {[
                    { name: 'filesystem', desc: '📁 Arquivos' },
                    { name: 'database', desc: '🗄️ Database' },
                    { name: 'git', desc: '🔄 Git' },
                    { name: 'docker', desc: '🐳 Docker' },
                    { name: 'slack', desc: '💬 Slack' },
                    { name: 'github', desc: '🐙 GitHub' }
                  ].map((server) => (
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
                    >
                      {server.desc} 🎪
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </div>

          {/* Generated Agent */}
          <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg`}>
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-semibold">🤖 Agent Gerado</h2>
              {generatedContent && (
                <div className="flex items-center gap-2">
                  <span className="text-xs px-2 py-1 rounded-full bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200">
                    {agentFramework} + MCP
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

            {generatedContent ? (
              <div className="space-y-4">
                <div className={`text-xs ${currentTheme.textSecondary} flex justify-between items-center`}>
                  <span>Caracteres: {generatedContent.length.toLocaleString()}</span>
                  <span className="bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200 px-2 py-1 rounded-full">
                    Production Ready
                  </span>
                </div>
                <div className={`max-h-96 overflow-y-auto p-4 rounded-lg border ${currentTheme.border} bg-opacity-50`}>
                  <pre className="whitespace-pre-wrap text-sm font-mono">{generatedContent}</pre>
                </div>

                <div className={`p-3 rounded-lg border ${currentTheme.border} bg-gradient-to-r from-green-50 to-blue-50 dark:from-green-900/20 dark:to-blue-900/20`}>
                  <p className="text-sm text-green-800 dark:text-green-200">
                    🚀 <strong>Agent Profissional:</strong> {agentProvider} + {agentFramework} + MCP + Config TOML
                    {DemoMode.isActive && (
                      <span className="block mt-1 text-blue-700 dark:text-blue-300">
                        🎪 <strong>Demo:</strong> Arquivos completos gerados com configurações de produção
                      </span>
                    )}
                  </p>
                  <div className="mt-2 text-xs text-gray-600 dark:text-gray-400">
                    <div className="bg-gray-200 dark:bg-gray-700 p-2 rounded font-mono text-xs">
                      <div>📁 <strong>Arquivos inclusos:</strong></div>
                      <div>├── config.toml (configuração TOML)</div>
                      <div>├── agent.py (código do agent)</div>
                      <div>├── requirements.txt (dependências)</div>
                      <div>└── README.md (documentação)</div>
                    </div>
                  </div>
                </div>
              </div>
            ) : (
              <div className={`${currentTheme.textSecondary} text-center py-12`}>
                <Bot size={48} className="mx-auto mb-4 opacity-50" />
                <p>Seu agent será gerado aqui</p>
                <p className="text-sm mt-2">Configure o agent e gere o código</p>
                {DemoMode.isActive && (
                  <p className="text-xs mt-4 text-green-600 dark:text-green-400">
                    🎪 Versão Demo: Gera código completo e profissional
                  </p>
                )}
              </div>
            )}

            {/* Generate Button */}
            {ideas.length > 0 && (
              <button
                onClick={generateContent}
                disabled={isGenerating}
                className={`w-full mt-4 flex items-center justify-center gap-2 px-4 py-3 rounded-lg bg-gradient-to-r from-green-600 to-blue-600 hover:from-green-700 hover:to-blue-700 text-white disabled:opacity-50 disabled:cursor-not-allowed transition-all transform hover:scale-105`}
              >
                <Wand2 size={20} className={isGenerating ? 'animate-spin' : ''} />
                {isGenerating ? 'Gerando agent...' : 'Criar Agent Profissional 🚀'}
              </button>
            )}
          </div>
        </div>
      </div>

      {/* Onboarding Modal */}
      {showOnboarding && (
        <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
          <div className={`${currentTheme.cardBg} rounded-xl p-6 max-w-md border ${currentTheme.border} shadow-xl`}>
            <h3 className="text-xl font-bold mb-4">{OnboardingSteps.agent[currentStep].title}</h3>
            <p className={`${currentTheme.textSecondary} mb-6`}>{OnboardingSteps.agent[currentStep].content}</p>
            <div className="flex justify-between">
              <span className="text-sm text-gray-500">
                {currentStep + 1} de {OnboardingSteps.agent.length}
              </span>
              <button
                onClick={nextOnboardingStep}
                className={`px-4 py-2 rounded-lg ${currentTheme.button}`}
              >
                {currentStep < OnboardingSteps.agent.length - 1 ? 'Próximo' : 'Finalizar'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Educational Modal */}
      {showEducational && educationalTopic && (
        <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
          <div className={`${currentTheme.cardBg} rounded-xl p-6 max-w-lg border ${currentTheme.border} shadow-xl`}>
            <h3 className="text-xl font-bold mb-4">{DemoMode.education[educationalTopic].title}</h3>
            <p className={`${currentTheme.textSecondary} mb-4`}>{DemoMode.education[educationalTopic].description}</p>
            <div className="mb-6">
              <h4 className="font-semibold mb-2">Benefícios:</h4>
              <ul className="space-y-1">
                {DemoMode.education[educationalTopic].benefits.map((benefit, index) => (
                  <li key={index} className={currentTheme.textSecondary}>{benefit}</li>
                ))}
              </ul>
            </div>
            <button
              onClick={() => setShowEducational(false)}
              className={`px-4 py-2 rounded-lg ${currentTheme.button} w-full`}
            >
              Entendi!
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

return null;
};

export default PromptCrafter; import React, { useState, useEffect } from 'react';
import { Trash2, Edit3, Plus, Wand2, Sun, Moon, Copy, Check, Info, Play, BookOpen, Bot, FileText, ArrowLeft, BookOpenCheckIcon, BookOpenCheck } from 'lucide-react';

// =============================================================================
// 🎯 DEMO MODE CONTROLLER - Single Source of Truth
// =============================================================================

const DemoMode = {
  isActive: true,

  modes: {
    SIMPLE: 'simple',
    ONBOARDING: 'onboarding',
    EDUCATIONAL: 'educational',
    PREVIEW: 'preview'
  },

  currentMode: 'onboarding',

  features: {
    ollama: { ready: false, eta: 'Q2 2024' },
    openai: { ready: false, eta: 'Q1 2024' },
    gemini: { ready: false, eta: 'Q2 2024' },
    mcp_real: { ready: false, eta: 'Q1 2024' },
    agent_execution: { ready: false, eta: 'Q3 2024' },
    copilot: { ready: false, eta: 'Q2 2024' }
  },

  education: {
    mcp: {
      title: "Model Context Protocol (MCP)",
      description: "Protocolo que permite que modelos de IA se conectem com sistemas externos de forma padronizada",
      benefits: [
        "🔌 Conecta IA com ferramentas reais",
        "🛡️ Segurança e controle de acesso",
        "🔄 Reutilização entre diferentes modelos",
        "⚡ Performance otimizada"
      ]
    },
    agents: {
      title: "Agents de IA",
      description: "Sistemas autônomos que podem usar ferramentas e tomar decisões para completar tarefas",
      benefits: [
        "🤖 Automação inteligente",
        "🧠 Tomada de decisão contextual",
        "🔧 Uso de ferramentas múltiplas",
        "📈 Escalabilidade de tarefas"
      ]
    },
    prompts: {
      title: "Engenharia de Prompts",
      description: "Arte e ciência de criar instruções eficazes para modelos de linguagem",
      benefits: [
        "🎯 Resultados mais precisos",
        "⚡ Respostas consistentes",
        "🔧 Controle sobre outputs",
        "📊 Melhor performance"
      ]
    }
  },

  getLabel: function (feature, defaultLabel) {
    if (!this.isActive) return defaultLabel;

    const featureStatus = this.features[feature];

    switch (this.currentMode) {
      case 'simple':
        return defaultLabel + ' 🎪';
      case 'onboarding':
        return featureStatus ? defaultLabel + ' (' + featureStatus.eta + ')' : defaultLabel + ' 🎪';
      case 'educational':
        return featureStatus ? defaultLabel + ' - Chegando em ' + featureStatus.eta : defaultLabel + ' 🎪';
      case 'preview':
        return defaultLabel + ' - Preview';
      default:
        return defaultLabel;
    }
  },

  handleDemoCall: function (feature, action) {
    if (!this.isActive) return null;

    const messages = {
      ollama: '🦙 Ollama será integrado na versão completa! Conecte modelos locais diretamente.',
      openai: '🧠 OpenAI GPT-4 em breve! Múltiplos providers em um só lugar.',
      gemini: '💎 Google Gemini chegando! Diversidade de modelos para diferentes tarefas.',
      mcp_real: '🔌 Servidores MCP reais em desenvolvimento! Conecte com qualquer sistema.',
      copilot: '🚁 GitHub Copilot API será integrada! Agents com capacidades de código avançadas.'
    };

    return {
      success: false,
      message: messages[feature] || 'Feature "' + feature + '" em modo demo',
      eta: this.features[feature] ? this.features[feature].eta : 'Em breve'
    };
  }
};

// =============================================================================
// 🎓 ONBOARDING STEPS
// =============================================================================

const OnboardingSteps = {
  prompt: [
    {
      id: 'welcome-prompt',
      title: 'Bem-vindo ao Prompt Engineer! 📝',
      content: 'Transforme suas ideias brutas em prompts profissionais e estruturados.',
      target: 'header'
    },
    {
      id: 'ideas-prompt',
      title: 'Adicione suas ideias 💡',
      content: 'Cole notas, pensamentos ou requisitos. A IA organizará tudo em um prompt perfeito!',
      target: 'ideas-input'
    },
    {
      id: 'purpose-prompt',
      title: 'Defina o propósito 🎯',
      content: 'Especifique o objetivo do seu prompt para melhor estruturação.',
      target: 'purpose-section'
    }
  ],
  agent: [
    {
      id: 'welcome-agent',
      title: 'Bem-vindo ao Agent Crafter! 🤖',
      content: 'Crie agents inteligentes com MCP, frameworks modernos e configurações profissionais.',
      target: 'header'
    },
    {
      id: 'framework-agent',
      title: 'Escolha o Framework 🏗️',
      content: 'CrewAI, AutoGen, LangChain - cada framework tem suas especialidades.',
      target: 'framework-section'
    },
    {
      id: 'mcp-agent',
      title: 'Poder do MCP 🔌',
      content: 'Model Context Protocol conecta seu agent com ferramentas e sistemas externos.',
      target: 'mcp-section'
    }
  ]
};

// =============================================================================
// 🎨 MAIN COMPONENT - Multi-Screen Application
// =============================================================================

const PromptCrafter = () => {
  // Navigation state
  const [currentScreen, setCurrentScreen] = useState('home'); // 'home', 'prompt', 'agent'

  // Common state
  const [darkMode, setDarkMode] = useState(true);
  const [currentInput, setCurrentInput] = useState('');
  const [ideas, setIdeas] = useState([]);
  const [editingId, setEditingId] = useState(null);
  const [editingText, setEditingText] = useState('');
  const [generatedContent, setGeneratedContent] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [copied, setCopied] = useState(false);

  // Prompt-specific state
  const [promptPurpose, setPromptPurpose] = useState('Outros');
  const [customPromptPurpose, setCustomPromptPurpose] = useState('');
  const [maxLength, setMaxLength] = useState(5000);

  // Agent-specific state
  const [agentFramework, setAgentFramework] = useState('crewai');
  const [agentRole, setAgentRole] = useState('');
  const [agentTools, setAgentTools] = useState([]);
  const [agentProvider, setAgentProvider] = useState('claude');
  const [mcpServers, setMcpServers] = useState([]);
  const [customMcpServer, setCustomMcpServer] = useState('');
  const [agentPurpose, setAgentPurpose] = useState('Automação');
  const [customAgentPurpose, setCustomAgentPurpose] = useState('');

  // Onboarding state
  const [showOnboarding, setShowOnboarding] = useState(false);
  const [currentStep, setCurrentStep] = useState(0);
  const [showEducational, setShowEducational] = useState(false);
  const [educationalTopic, setEducationalTopic] = useState(null);

  useEffect(() => {
    document.documentElement.className = darkMode ? 'dark' : '';
  }, [darkMode]);

  const addIdea = () => {
    if (currentInput.trim()) {
      const newIdea = {
        id: Date.now(),
        text: currentInput.trim()
      };
      setIdeas([...ideas, newIdea]);
      setCurrentInput('');
    }
  };

  const removeIdea = (id) => {
    setIdeas(ideas.filter(idea => idea.id !== id));
  };

  const startEditing = (id, text) => {
    setEditingId(id);
    setEditingText(text);
  };

  const saveEdit = () => {
    setIdeas(ideas.map(idea =>
      idea.id === editingId
        ? { ...idea, text: editingText }
        : idea
    ));
    setEditingId(null);
    setEditingText('');
  };

  const cancelEdit = () => {
    setEditingId(null);
    setEditingText('');
  };

  const generateContent = async () => {
    if (ideas.length === 0) return;

    setIsGenerating(true);

    let engineeringPrompt = '';

    if (currentScreen === 'prompt') {
      const purposeText = promptPurpose === 'Outros' && customPromptPurpose
        ? customPromptPurpose
        : promptPurpose;

      engineeringPrompt = `
Você é um especialista em engenharia de prompts com conhecimento profundo em técnicas de prompt engineering. Sua tarefa é transformar ideias brutas e desorganizadas em um prompt estruturado, profissional e eficaz.

CONTEXTO: O usuário inseriu as seguintes notas/ideias brutas:
${ideas.map((idea, index) => `${index + 1}. "${idea.text}"`).join('\n')}

PROPÓSITO DO PROMPT: ${purposeText}
TAMANHO MÁXIMO: ${maxLength} caracteres

INSTRUÇÕES PARA ESTRUTURAÇÃO:
1. Analise todas as ideias e identifique o objetivo principal
2. Organize as informações de forma lógica e hierárquica
3. Aplique técnicas de engenharia de prompt como:
   - Definição clara de contexto e papel
   - Instruções específicas e mensuráveis
   - Exemplos quando apropriado
   - Formato de saída bem definido
   - Chain-of-thought se necessário
4. Use markdown para estruturação clara
5. Seja preciso, objetivo e profissional
6. Mantenha o escopo dentro do limite de caracteres

IMPORTANTE: Responda APENAS com o prompt estruturado em markdown, sem explicações adicionais ou texto introdutório. O prompt deve ser completo e pronto para uso.
`;
    } else if (currentScreen === 'agent') {
      const purposeText = agentPurpose === 'Outros' && customAgentPurpose
        ? customAgentPurpose
        : agentPurpose;
      const toolsList = agentTools.length > 0 ? agentTools.join(', ') : 'ferramentas padrão';
      const mcpServersList = mcpServers.length > 0 ? mcpServers.join(', ') : 'nenhum servidor MCP configurado';

      engineeringPrompt = `
Você é um especialista em desenvolvimento de agents de IA com conhecimento avançado em Model Context Protocol (MCP), arquitetura de sistemas multi-agent e integração com diversos provedores de LLM.

CONTEXTO: O usuário inseriu as seguintes notas/ideias brutas para o agent:
${ideas.map((idea, index) => `${index + 1}. "${idea.text}"`).join('\n')}

CONFIGURAÇÕES DO AGENT:
- Framework: ${agentFramework}
- Provider LLM: ${agentProvider}
- Papel/Role: ${agentRole || 'A ser definido baseado nas ideias'}
- Ferramentas Tradicionais: ${toolsList}
- Servidores MCP: ${mcpServersList}
- Propósito: ${purposeText}

INSTRUÇÕES PARA CRIAÇÃO DO AGENT COM MCP E CONFIG TOML:
1. Analise as ideias e defina claramente o papel e objetivo do agent
2. Crie um agent ${agentFramework} completo e funcional
3. Configure integração com ${agentProvider} via API ou MCP
4. Inclua configurações MCP detalhadas e arquivo config.toml profissional
5. Use configurações baseadas em produção:
   - Context scoping com tokens limitados
   - Guards contra comandos perigosos
   - Summarizers específicos por tipo
   - Goal-driven context management
   - Fail-fast behaviors
6. Gere TODOS os arquivos necessários:
   - config.toml (configuração principal)
   - agent.py (código do agent)
   - requirements.txt (dependências)
   - README.md (documentação)

ESTRUTURA ESPERADA:
\`\`\`toml
# config.toml - Configuração principal do agent
[settings]
model_reasoning_summary = "concise"
user_intent_summary = "detailed"
# ... resto da configuração profissional
\`\`\`

\`\`\`python
# agent.py - Implementação do agent
# Framework: ${agentFramework}
# Provider: ${agentProvider}
\`\`\`

IMPORTANTE: Responda com código estruturado e pronto para uso, incluindo config.toml profissional.
`;
    }

    try {
      // Only Claude is functional - others trigger demo mode
      if (agentProvider !== 'claude' && DemoMode.isActive && currentScreen === 'agent') {
        const demoResult = DemoMode.handleDemoCall(agentProvider);
        setGeneratedContent('# 🎪 Demo Mode\n\n' + demoResult.message + '\n\n**ETA:** ' + demoResult.eta + '\n\n---\n\n*Configurações salvas:*\n- Framework: ' + agentFramework + '\n- Provider: ' + agentProvider + '\n- Ferramentas: ' + (agentTools.join(', ') || 'Nenhuma') + '\n- Servidores MCP: ' + (mcpServers.join(', ') || 'Nenhum') + '\n\nEssas configurações serão aplicadas quando o provider estiver disponível!');
        setIsGenerating(false);
        return;
      }

      // Real Claude API call
      const response = await fetch("https://api.anthropic.com/v1/messages", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          model: "claude-sonnet-4-20250514",
          max_tokens: 4000,
          messages: [{ role: "user", content: engineeringPrompt }]
        })
      });

      if (!response.ok) {
        throw new Error('API request failed: ' + response.status);
      }

      const data = await response.json();
      const result = data.content[0].text;

      setGeneratedContent(result);
    } catch (error) {
      console.error('Erro ao gerar:', error);
      setGeneratedContent('Erro ao gerar o conteúdo. ' + error.message);
    }

    setIsGenerating(false);
  };

  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(generatedContent);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (error) {
      console.error('Erro ao copiar:', error);
    }
  };

  const handleFeatureClick = (feature) => {
    if (DemoMode.isActive && (!DemoMode.features[feature] || !DemoMode.features[feature].ready)) {
      const demoResult = DemoMode.handleDemoCall(feature);
      alert(demoResult.message + '\n\nETA: ' + demoResult.eta);
      return false;
    }
    return true;
  };

  const startOnboarding = () => {
    setShowOnboarding(true);
    setCurrentStep(0);
  };

  const nextOnboardingStep = () => {
    const steps = currentScreen === 'prompt' ? OnboardingSteps.prompt : OnboardingSteps.agent;
    if (currentStep < steps.length - 1) {
      setCurrentStep(currentStep + 1);
    } else {
      setShowOnboarding(false);
    }
  };

  const showEducation = (topic) => {
    setEducationalTopic(topic);
    setShowEducational(true);
  };

  const resetState = () => {
    setIdeas([]);
    setCurrentInput('');
    setGeneratedContent('');
    setEditingId(null);
    setEditingText('');
  };

  const navigateToScreen = (screen) => {
    setCurrentScreen(screen);
    resetState();
  };

  const theme = {
    dark: {
      bg: 'bg-gray-900',
      cardBg: 'bg-gray-800',
      text: 'text-gray-100',
      textSecondary: 'text-gray-300',
      border: 'border-gray-700',
      input: 'bg-gray-700 border-gray-600 text-gray-100',
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
      input: 'bg-white border-gray-300 text-gray-900',
      button: 'bg-blue-600 hover:bg-blue-700 text-white',
      buttonSecondary: 'bg-gray-200 hover:bg-gray-300 text-gray-700',
      accent: 'text-blue-600'
    }
  };

  const currentTheme = darkMode ? theme.dark : theme.light;

  // HOME SCREEN
  if (currentScreen === 'home') {
    return (
      <div className={`min-h-screen ${currentTheme.bg} ${currentTheme.text} p-4 transition-colors duration-300`}>
        <div className="max-w-6xl mx-auto">
          {/* Header */}
          <div className="flex justify-between items-center mb-12" id="header">
            <div>
              <h1 className="text-5xl font-bold mb-2">
                <span className={currentTheme.accent}>AI</span> Craft <span className={currentTheme.accent}>Studio</span>
                <span className="text-lg ml-3 px-3 py-1 bg-gradient-to-r from-purple-500 to-blue-500 text-white rounded-full">
                  +MCP
                </span>
                {DemoMode.isActive && (
                  <span className="text-xs ml-2 px-2 py-1 bg-blue-500 text-blue-100 rounded-full">
                    DEMO v1.0.0
                  </span>
                )}
              </h1>
              <p className={`${currentTheme.textSecondary} text-xl`}>
                Plataforma completa para criação de prompts e agents inteligentes
              </p>
            </div>
            <div className="flex items-center gap-4">
              {DemoMode.isActive && (
                <button
                  onClick={() => showEducation('prompts')}
                  className="px-3 py-2 rounded-lg bg-purple-600 text-white hover:bg-purple-700 flex items-center gap-2 text-sm"
                >
                  <BookOpen size={16} />
                  Docs
                </button>
              )}
              <button
                onClick={() => setDarkMode(!darkMode)}
                className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
              >
                {darkMode ? <Sun size={20} /> : <Moon size={20} />}
              </button>
            </div>
          </div>

          {/* Tools Selection */}
          <div className="grid md:grid-cols-2 gap-8 mb-12">
            {/* Prompt Engineer */}
            <div
              onClick={() => navigateToScreen('prompt')}
              className={`${currentTheme.cardBg} rounded-xl p-8 border ${currentTheme.border} shadow-lg hover:shadow-xl transition-all cursor-pointer hover:scale-105 group`}
            >
              <div className="flex items-center mb-4">
                <div className="p-3 rounded-lg bg-purple-500 text-white mr-4 group-hover:scale-110 transition-transform">
                  <FileText size={32} />
                </div>
                <div>
                  <h2 className="text-2xl font-bold text-purple-500">Prompt Engineer</h2>
                  <p className={currentTheme.textSecondary}>Engenharia de Prompts Profissional</p>
                </div>
              </div>
              <p className={`${currentTheme.textSecondary} mb-6`}>
                Transforme ideias brutas em prompts estruturados, otimizados e profissionais.
                Técnicas avançadas de prompt engineering para resultados consistentes.
              </p>
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <span className="text-green-500">✓</span>
                  <span className="text-sm">Estruturação automática</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-green-500">✓</span>
                  <span className="text-sm">Chain-of-thought integration</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-green-500">✓</span>
                  <span className="text-sm">Markdown formatting</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-green-500">✓</span>
                  <span className="text-sm">Token optimization</span>
                </div>
              </div>
            </div>

            {/* Agent Crafter */}
            <div
              onClick={() => navigateToScreen('agent')}
              className={`${currentTheme.cardBg} rounded-xl p-8 border ${currentTheme.border} shadow-lg hover:shadow-xl transition-all cursor-pointer hover:scale-105 group`}
            >
              <div className="flex items-center mb-4">
                <div className="p-3 rounded-lg bg-green-500 text-white mr-4 group-hover:scale-110 transition-transform">
                  <Bot size={32} />
                </div>
                <div>
                  <h2 className="text-2xl font-bold text-green-500">Agent Crafter</h2>
                  <p className={currentTheme.textSecondary}>Criação de Agents Inteligentes</p>
                </div>
              </div>
              <p className={`${currentTheme.textSecondary} mb-6`}>
                Crie agents autônomos com MCP, frameworks modernos e configurações profissionais.
                Integração com múltiplos providers e ferramentas externas.
              </p>
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <span className="text-green-500">✓</span>
                  <span className="text-sm">Model Context Protocol (MCP)</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-green-500">✓</span>
                  <span className="text-sm">Multi-framework support</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-green-500">✓</span>
                  <span className="text-sm">Professional TOML config</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-green-500">✓</span>
                  <span className="text-sm">Production-ready code</span>
                </div>
              </div>
            </div>
          </div>

          {/* Features Grid */}
          <div className="grid md:grid-cols-3 gap-6 mb-12">
            <div className={`${currentTheme.cardBg} rounded-lg p-6 border ${currentTheme.border}`}>
              <div className="text-2xl mb-3">🔌</div>
              <h3 className="font-bold mb-2">Model Context Protocol</h3>
              <p className={`${currentTheme.textSecondary} text-sm`}>
                Conecte agents com sistemas externos de forma padronizada e segura.
              </p>
            </div>
            <div className={`${currentTheme.cardBg} rounded-lg p-6 border ${currentTheme.border}`}>
              <div className="text-2xl mb-3">🏗️</div>
              <h3 className="font-bold mb-2">Multi-Framework</h3>
              <p className={`${currentTheme.textSecondary} text-sm`}>
                Suporte a CrewAI, AutoGen, LangChain e frameworks customizados.
              </p>
            </div>
            <div className={`${currentTheme.cardBg} rounded-lg p-6 border ${currentTheme.border}`}>
              <div className="text-2xl mb-3">⚡</div>
              <h3 className="font-bold mb-2">Production Ready</h3>
              <p className={`${currentTheme.textSecondary} text-sm`}>
                Configurações profissionais, guards de segurança e otimizações.
              </p>
            </div>
          </div>

          {/* Demo Footer */}
          {DemoMode.isActive && (
            <div className="p-4 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800">
              <div className="flex items-start gap-3">
                <span className="text-2xl">🎪</span>
                <div>
                  <h3 className="font-semibold text-blue-800 dark:text-blue-200 mb-2">
                    AI Craft Studio - Powered by Grompt Engine
                  </h3>
                  <p className="text-blue-700 dark:text-blue-300 text-sm">
                    💡 <strong>Inspirado no Grompt CLI:</strong> Esta plataforma web é uma evolução do Grompt CLI,
                    mantendo a filosofia Kubex de simplicidade radical e anti-lock-in.
                  </p>
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Educational Modal */}
        {showEducational && educationalTopic && (
          <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
            <div className={`${currentTheme.cardBg} rounded-xl p-6 max-w-lg border ${currentTheme.border} shadow-xl`}>
              <h3 className="text-xl font-bold mb-4">{DemoMode.education[educationalTopic].title}</h3>
              <p className={`${currentTheme.textSecondary} mb-4`}>{DemoMode.education[educationalTopic].description}</p>
              <div className="mb-6">
                <h4 className="font-semibold mb-2">Benefícios:</h4>
                <ul className="space-y-1">
                  {DemoMode.education[educationalTopic].benefits.map((benefit, index) => (
                    <li key={index} className={currentTheme.textSecondary}>{benefit}</li>
                  ))}
                </ul>
              </div>
              <button
                onClick={() => setShowEducational(false)}
                className={`px-4 py-2 rounded-lg ${currentTheme.button} w-full`}
              >
                Entendi!
              </button>
            </div>
          </div>
        )}
      </div>
    );
  }

  // PROMPT ENGINEER SCREEN
  if (currentScreen === 'prompt') {
    return (
      <div className={`min-h-screen ${currentTheme.bg} ${currentTheme.text} p-4 transition-colors duration-300`}>
        <div className="max-w-7xl mx-auto">
          {/* Header */}
          <div className="flex justify-between items-center mb-8" id="header">
            <div className="flex items-center gap-4">
              <button
                onClick={() => navigateToScreen('home')}
                className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
              >
                <ArrowLeft size={20} />
              </button>
              <div>
                <h1 className="text-4xl font-bold mb-2 flex items-center gap-3">
                  <FileText className="text-purple-500" size={40} />
                  <span className="text-purple-500">Prompt</span> Engineer
                </h1>
                <p className={currentTheme.textSecondary}>
                  Transforme ideias em prompts estruturados e profissionais
                </p>
              </div>
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
                    onClick={() => showEducation('prompts')}
                    className="px-3 py-2 rounded-lg bg-purple-600 text-white hover:bg-purple-700 flex items-center gap-2 text-sm"
                  >
                    <BookOpenCheckIcon size={16} />
                    Docs
                  </button>
                </div>
              )}
              <button
                onClick={() => setDarkMode(!darkMode)}
                className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
              >
                {darkMode ? <Sun size={20} /> : <Moon size={20} />}
              </button>
            </div>
          </div>

          <div className="grid md:grid-cols-3 gap-8">
            {/* Left Panel - Ideas Input */}
            <div className="md:col-span-1">
              <h2 className="text-2xl font-bold mb-4">Suas Ideias 💡</h2>
              <div className="mb-4" >
                <textarea
                  value={userInput}
                  onChange={(e) => setUserInput(e.target.value)}
                  className={`w-full p-2 border rounded-lg ${currentTheme.border} ${currentTheme.bg} ${currentTheme.text}`}
                  rows={5}
                  placeholder="Digite suas ideias aqui..."
                />
              </div>
              <button
                onClick={addIdea}
                className={`w-full px-4 py-2 rounded-lg ${currentTheme.button} transition-transform hover:scale-105`}
              >
                <Plus size={16} className="inline-block mr-2" />
                Adicionar Ideia
              </button>

              <div className="mt-6">
                {ideas.length === 0 ? (
                  <p className={currentTheme.textSecondary}>Nenhuma ideia adicionada ainda.</p>
                ) : (
                  <ul className="space-y-4 max-h-[400px] overflow-y-auto pr-2">
                    {ideas.map((idea) => (
                      <li key={idea.id} className={`${currentTheme.cardBg} p-3 rounded-lg border ${currentTheme.border} shadow-sm`}>
                        {editingId === idea.id ? (
                          <div>
                            <textarea
                              value={editingText}
                              onChange={(e) => setAgentFramework(e.target.value)}
                              className={`w-full p-2 border rounded-lg ${currentTheme.border} ${currentTheme.bg} ${currentTheme.text}`}
                              rows={3}
                            />
                            <div className="flex justify-end gap-2 mt-2">
                              <button
                                onClick={saveEdit}
                                className={`px-3 py-1 rounded-lg ${currentTheme.button} text-sm`}
                              >
                                Salvar
                              </button>
                              <button
                                onClick={cancelEdit}
                                className={`px-3 py-1 rounded-lg ${currentTheme.buttonSecondary} text-sm`}
                              >
                                Cancelar
                              </button>
                            </div>
                          </div>
                        ) : (
                          <div className="flex justify-between items-start">
                            <p className={currentTheme.textSecondary}>{idea.text}</p>
                            <div className="flex gap-2">
                              <button
                                onClick={() => startEditing(idea.id, idea.text)}
                                className={`p-1 rounded-lg ${currentTheme.buttonSecondary} hover:bg-gray-600`}
                              >
                                <Edit3 size={16} />
                              </button>
                              <button
                                onClick={() => removeIdea(idea.id)}
                                className={`p-1 rounded-lg ${currentTheme.buttonSecondary} hover:bg-gray-600`}
                              >
                                <Trash2 size={16} />
                              </button>
                            </div>
                          </div>
                        )}
                      </li>
                    ))}
                  </ul>
                )}
              </div>
            </div>

            {/* Middle Panel - Configuration */}
            <div className="md:col-span-1">
              <h2 className="text-2xl font-bold mb-4">Configurações do Prompt ⚙️</h2>

              <div className="mb-6" id="purpose-section">
                <label className="block font-semibold mb-2">Propósito do {promptPurpose}</label>
                <textarea
                  value={promptPurpose}
                  onChange={(e) => setPromptPurpose(e.target.value)}
                  className={`w-full p-2 border rounded-lg ${currentTheme.border} ${currentTheme.bg} ${currentTheme.text}`}
                  rows={3}
                />
              </div>
              <div className="mb-6">
                <label className="block font-semibold mb-2">Tamanho Máximo (caracteres)</label>
                <input
                  type="number"
                  value={maxLength}
                  onChange={(e) => setMaxLength(parseInt(e.target.value) || 5000)}
                  className={`w-full p-2 border rounded-lg ${currentTheme.border} ${currentTheme.bg} ${currentTheme.text}`}
                  min={100}
                  max={10000}
                />
              </div>

              <button
                onClick={generateContent}
                disabled={isGenerating || ideas.length === 0}
                className={`w-full px-4 py-2 rounded-lg ${isGenerating || ideas.length === 0 ? currentTheme.buttonSecondary : currentTheme.button} transition-transform hover:scale-105`}
              >
                {isGenerating ? (
                  <span>Gerando...</span>
                ) : (
                  <>
                    <Wand2 size={16} className="inline-block mr-2" />
                    Gerar Prompt
                  </>
                )}
              </button>
            </div>

            {/* Right Panel - Generated Content */}
            <div className="md:col-span-1">
              <h2 className="text-2xl font-bold mb-4">Prompt Gerado 🏗️</h2>
              <div className={`${currentTheme.cardBg} p-4 rounded-lg border ${currentTheme.border} shadow-sm min-h-[400px] flex flex-col`}>
                {generatedContent ? (
                  <>
                    <div className="mb-4 flex justify-end">
                      <button
                        onClick={copyToClipboard}
                        className={`px-3 py-1 rounded-lg ${currentTheme.buttonSecondary} hover:bg-gray-600 text-sm`}
                      >
                        {copied ? 'Copiado!' : 'Copiar'}
                      </button>
                    </div>
                    <pre className="whitespace-pre-wrap break-words overflow-y-auto max-h-[350px]">
                      {generatedContent}
                    </pre>
                  </>
                ) : (
                  <p className={currentTheme.textSecondary}>O prompt gerado aparecerá aqui.</p>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Onboarding Modal */}
        {showOnboarding && (
          <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
            <div className={`${currentTheme.cardBg} rounded-xl p-6 max-w-lg border ${currentTheme.border} shadow-xl`}>
              <h3 className="text-xl font-bold mb-4">{OnboardingSteps.prompt[currentStep].title}</h3>
              <p className={`${currentTheme.textSecondary} mb-6`}>{OnboardingSteps.prompt[currentStep].content}</p>
              <button
                onClick={nextOnboardingStep}
                className={`px-4 py-2 rounded-lg ${currentTheme.button} w-full`}
              >
                {currentStep < OnboardingSteps.prompt.length - 1 ? 'Próximo' : 'Começar'}
              </button>
            </div>
          </div>
        )}
      </div>
    );
  }

  // AGENT CRAFTER SCREEN
  if (currentScreen === 'agent') {
    return (
      <div className={`min-h-screen ${currentTheme.bg} ${currentTheme.text} p-4 transition-colors duration-300`}>
        <div className="max-w-7xl mx-auto">
          <h1 className="text-3xl font-bold mb-6">Agent Crafter 🤖</h1>
