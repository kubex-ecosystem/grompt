import React, { useState, useEffect } from 'react';
import { Trash2, Edit3, Plus, Wand2, Sun, Moon, Copy, Check, AlertCircle } from 'lucide-react';

const PromptCrafter = () => {
  const [darkMode, setDarkMode] = useState(true);
  const [currentInput, setCurrentInput] = useState('');
  const [ideas, setIdeas] = useState([]);
  const [editingId, setEditingId] = useState(null);
  const [editingText, setEditingText] = useState('');
  const [purpose, setPurpose] = useState('Outros');
  const [customPurpose, setCustomPurpose] = useState('');
  const [maxLength, setMaxLength] = useState(5000);
  const [generatedPrompt, setGeneratedPrompt] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [copied, setCopied] = useState(false);
  const [apiProvider, setApiProvider] = useState('demo');
  const [selectedModel, setSelectedModel] = useState('');
  const [availableAPIs, setAvailableAPIs] = useState({
    claude_available: false,
    openai_available: false,
    deepseek_available: false,
    ollama_available: false,
    demo_mode: true,
    available_models: {
      openai: [],
      deepseek: [],
      claude: [],
      ollama: []
    }
  });
  const [connectionStatus, setConnectionStatus] = useState('checking');
  const [serverInfo, setServerInfo] = useState(null);

  // =========================================
  // CONFIGURAÇÃO DE URL BASE
  // =========================================
  const getBaseURL = () => {
    // Se estamos em desenvolvimento (npm start), usar proxy ou porta específica
    if (process.env.NODE_ENV === 'development') {
      return 'http://localhost:8080'; // Servidor Go
    }
    // Se estamos em produção (servido pelo Go), usar URL relativa
    return '';
  };

  const apiCall = async (endpoint, options = {}) => {
    const baseURL = getBaseURL();
    const url = `${baseURL}${endpoint}`;
    
    console.log(`🔗 Fazendo requisição para: ${url}`);
    
    const defaultOptions = {
      headers: {
        'Content-Type': 'application/json',
      },
      ...options
    };

    try {
      const response = await fetch(url, defaultOptions);
      return response;
    } catch (error) {
      console.error(`❌ Erro na requisição para ${url}:`, error);
      throw error;
    }
  };

  useEffect(() => {
    document.documentElement.className = darkMode ? 'dark' : '';
  }, [darkMode]);

  // Verificar configuração e APIs disponíveis na inicialização
  useEffect(() => {
    checkAPIAvailability();
  }, []);

  const checkAPIAvailability = async () => {
    try {
      console.log('🔍 Verificando disponibilidade das APIs...');
      
      // Primeiro, verificar se o servidor Go está rodando
      const healthResponse = await apiCall('/api/health');
      
      if (healthResponse.ok) {
        const healthData = await healthResponse.json();
        setServerInfo(healthData);
        console.log('✅ Servidor Go conectado:', healthData);
      }

      // Verificar configuração das APIs
      const configResponse = await apiCall('/api/config');
      
      if (configResponse.ok) {
        const config = await configResponse.json();
        setAvailableAPIs(config);
        setConnectionStatus('connected');
        
        console.log('📋 Configuração recebida:', config);
        
        // Definir provider padrão baseado na disponibilidade
        if (config.claude_available) {
          setApiProvider('claude');
        } else if (config.openai_available) {
          setApiProvider('openai');
        } else if (config.deepseek_available) {
          setApiProvider('deepseek');
        } else if (config.ollama_available) {
          setApiProvider('ollama');
        } else {
          setApiProvider('demo');
        }
      } else {
        throw new Error(`Servidor retornou status ${configResponse.status}`);
      }
    } catch (error) {
      console.error('❌ Erro ao verificar APIs:', error);
      setConnectionStatus('offline');
      setAvailableAPIs({ demo_mode: true });
      setApiProvider('demo');
      
      // Se estivermos em desenvolvimento, mostrar dica
      if (process.env.NODE_ENV === 'development') {
        console.log('💡 Dica: Certifique-se de que o servidor Go está rodando na porta 8080');
        console.log('🔧 Execute: go run . ou make run');
      }
    }
  };

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

  const generateDemoPrompt = () => {
    const purposeText = purpose === 'Outros' && customPurpose 
      ? customPurpose 
      : purpose;

    return `# Prompt Estruturado - ${purposeText}

## 🎯 Contexto
Você é um assistente especializado em **${purposeText.toLowerCase()}** com conhecimento profundo na área.

## 📝 Ideias do Usuário Organizadas:
${ideas.map((idea, index) => `**${index + 1}.** ${idea.text}`).join('\n')}

## 🔧 Instruções Específicas
- Analise cuidadosamente todas as ideias apresentadas acima
- Identifique o objetivo principal e objetivos secundários
- Forneça uma resposta estruturada e bem organizada
- Mantenha o foco no propósito definido: **${purposeText}**
- Use exemplos práticos quando apropriado
- Seja específico e actionável

## 📋 Formato de Resposta Esperado
1. **Análise Inicial**: Resumo do que foi solicitado
2. **Desenvolvimento**: Resposta detalhada seguindo as ideias
3. **Conclusão**: Próximos passos ou considerações finais

## ⚙️ Configurações Técnicas
- Máximo de caracteres: ${maxLength.toLocaleString()}
- Propósito: ${purposeText}
- Total de ideias processadas: ${ideas.length}
- Modo: ${ connectionStatus === 'connected' ? 'Demo (servidor Go conectado)' : 'Demo (modo offline)'}

---
*Prompt gerado automaticamente pelo Prompt Crafter v1.0*
*${connectionStatus === 'connected' ? 'Configure uma API key para funcionalidade completa' : 'Servidor Go offline - usando modo demo'}*`;
  };

  const generatePrompt = async () => {
    if (ideas.length === 0) return;
    
    setIsGenerating(true);
    
    const purposeText = purpose === 'Outros' && customPurpose 
      ? customPurpose 
      : purpose;
    
    const engineeringPrompt = `
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

    try {
      let response;
      
      if (apiProvider === 'demo' || connectionStatus === 'offline') {
        // Simular delay para parecer real
        await new Promise(resolve => setTimeout(resolve, 2000));
        response = generateDemoPrompt();
      } else if (apiProvider === 'claude') {
        console.log('🤖 Enviando para Claude API...');
        const result = await apiCall('/api/claude', {
          method: 'POST',
          body: JSON.stringify({
            prompt: engineeringPrompt,
            max_tokens: maxLength
          })
        });
        
        if (!result.ok) {
          const errorText = await result.text();
          throw new Error(`Erro HTTP ${result.status}: ${errorText}`);
        }
        
        const data = await result.json();
        response = data.response || data.content || 'Resposta vazia do servidor';
        console.log('✅ Resposta recebida do Claude');
        
      } else if (apiProvider === 'openai') {
        console.log('🧠 Enviando para OpenAI API...');
        const result = await apiCall('/api/openai', {
          method: 'POST',
          body: JSON.stringify({
            prompt: engineeringPrompt,
            max_tokens: maxLength,
            model: selectedModel || 'gpt-3.5-turbo'
          })
        });
        
        if (!result.ok) {
          const errorText = await result.text();
          throw new Error(`Erro HTTP ${result.status}: ${errorText}`);
        }
        
        const data = await result.json();
        response = data.response || 'Resposta vazia do OpenAI';
        console.log('✅ Resposta recebida do OpenAI');
        
      } else if (apiProvider === 'deepseek') {
        console.log('🔍 Enviando para DeepSeek API...');
        const result = await apiCall('/api/deepseek', {
          method: 'POST',
          body: JSON.stringify({
            prompt: engineeringPrompt,
            max_tokens: maxLength,
            model: selectedModel || 'deepseek-chat'
          })
        });
        
        if (!result.ok) {
          const errorText = await result.text();
          throw new Error(`Erro HTTP ${result.status}: ${errorText}`);
        }
        
        const data = await result.json();
        response = data.response || 'Resposta vazia do DeepSeek';
        console.log('✅ Resposta recebida do DeepSeek');
      } else if (apiProvider === 'ollama') {
        console.log('🦙 Enviando para Ollama...');
        const result = await apiCall('/api/ollama', {
          method: 'POST',
          body: JSON.stringify({
            model: selectedModel || 'llama2',
            prompt: engineeringPrompt,
            stream: false
          })
        });
        
        if (!result.ok) {
          const errorText = await result.text();
          throw new Error(`Erro HTTP ${result.status}: ${errorText}`);
        }
        
        const data = await result.json();
        response = data.response || 'Resposta vazia do Ollama';
        console.log('✅ Resposta recebida do Ollama');
      }
      
      setGeneratedPrompt(response);
    } catch (error) {
      console.error('❌ Erro ao gerar prompt:', error);
      setGeneratedPrompt(`# Erro ao Gerar Prompt

**Erro:** ${error.message}

**Detalhes:** Não foi possível conectar com a API selecionada.

## 🔍 Verificações:
- **Status do servidor:** ${connectionStatus}
- **Modo atual:** ${process.env.NODE_ENV || 'production'}
- **Provider selecionado:** ${apiProvider}
- **Base URL:** ${getBaseURL()}

## 💡 Soluções:
1. **Se em desenvolvimento:** Certifique-se de que o servidor Go está rodando na porta 8080
2. **Se em produção:** Verifique se as APIs estão configuradas corretamente
3. **Tente usar o modo demo** como alternativa

**Comando para iniciar servidor Go:**
\`\`\`
go run .
# ou
make run
\`\`\`
`);
    }
    
    setIsGenerating(false);
  };

  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(generatedPrompt);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (error) {
      console.error('Erro ao copiar:', error);
      // Fallback para navegadores mais antigos
      const textArea = document.createElement('textarea');
      textArea.value = generatedPrompt;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
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

  const getConnectionStatusColor = () => {
    switch (connectionStatus) {
      case 'connected': return 'text-green-500';
      case 'offline': return 'text-red-500';
      default: return 'text-yellow-500';
    }
  };

  const getConnectionStatusText = () => {
    switch (connectionStatus) {
      case 'connected': return 'Conectado ao Go Server';
      case 'offline': return process.env.NODE_ENV === 'development' ? 'Servidor Go Offline' : 'Modo Offline';
      default: return 'Verificando conexão...';
    }
  };

  return (
    <div className={`min-h-screen ${currentTheme.bg} ${currentTheme.text} p-4 transition-colors duration-300`}>
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-4xl font-bold mb-2">
              <span className={currentTheme.accent}>Prompt</span> Crafter
            </h1>
            <p className={currentTheme.textSecondary}>
              Transforme suas ideias brutas em prompts estruturados e profissionais
            </p>
            {/* Debug info em desenvolvimento */}
            {process.env.NODE_ENV === 'development' && (
              <p className="text-xs text-yellow-400 mt-1">
                🔧 Modo desenvolvimento | Base URL: {getBaseURL()} | Status: {connectionStatus}
              </p>
            )}
          </div>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <div className={`h-2 w-2 rounded-full ${connectionStatus === 'connected' ? 'bg-green-500' : connectionStatus === 'offline' ? 'bg-red-500' : 'bg-yellow-500'}`}></div>
              <span className={`text-sm ${getConnectionStatusColor()}`}>
                {getConnectionStatusText()}
              </span>
            </div>
            <select 
              value={apiProvider}
              onChange={(e) => {
                setApiProvider(e.target.value);
                setSelectedModel(''); // Reset model when changing provider
              }}
              className={`px-3 py-2 rounded-lg ${currentTheme.input} border focus:ring-2 focus:ring-blue-500`}
            >
              {availableAPIs.claude_available && (
                <option value="claude">Claude API</option>
              )}
              {availableAPIs.openai_available && (
                <option value="openai">OpenAI API</option>
              )}
              {availableAPIs.deepseek_available && (
                <option value="deepseek">DeepSeek API</option>
              )}
              {availableAPIs.ollama_available && (
                <option value="ollama">Ollama Local</option>
              )}
              <option value="demo">Modo Demo</option>
            </select>
            
            {/* Model Selection */}
            {apiProvider !== 'demo' && availableAPIs.available_models && availableAPIs.available_models[apiProvider] && (
              <select
                value={selectedModel}
                onChange={(e) => setSelectedModel(e.target.value)}
                className={`px-3 py-2 rounded-lg ${currentTheme.input} border focus:ring-2 focus:ring-blue-500`}
              >
                <option value="">Modelo padrão</option>
                {availableAPIs.available_models[apiProvider].map((model) => (
                  <option key={model} value={model}>{model}</option>
                ))}
              </select>
            )}
            <button
              onClick={() => setDarkMode(!darkMode)}
              className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
            >
              {darkMode ? <Sun size={20} /> : <Moon size={20} />}
            </button>
          </div>
        </div>

        {/* Status Alert */}
        {connectionStatus === 'offline' && (
          <div className="mb-6 p-4 bg-yellow-900 border border-yellow-600 rounded-lg flex items-center gap-3">
            <AlertCircle className="text-yellow-400" size={20} />
            <div className="text-yellow-100">
              <strong>Modo Offline:</strong> 
              {process.env.NODE_ENV === 'development' 
                ? ' Servidor Go não está respondendo. Certifique-se de executar "go run ." ou "make run"' 
                : ' Executando em modo demo. Configure APIs para funcionalidade completa.'
              }
            </div>
          </div>
        )}

        {/* Server Info (em desenvolvimento) */}
        {process.env.NODE_ENV === 'development' && serverInfo && (
          <div className="mb-6 p-4 bg-blue-900 border border-blue-600 rounded-lg">
            <p className="text-blue-100">
              <strong>🔧 Info do Servidor:</strong> v{serverInfo.version} | 
              APIs: {serverInfo.apis?.claude ? '✅' : '❌'} Claude, {serverInfo.apis?.ollama ? '✅' : '❌'} Ollama
            </p>
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Input Section */}
          <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg`}>
            <h2 className="text-xl font-semibold mb-4">📝 Adicionar Ideias</h2>
            <div className="space-y-4">
              <textarea
                value={currentInput}
                onChange={(e) => setCurrentInput(e.target.value)}
                placeholder="Cole suas notas, ideias brutas ou pensamentos desorganizados aqui..."
                className={`w-full h-32 px-4 py-3 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-blue-500 resize-none`}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' && e.ctrlKey) {
                    addIdea();
                  }
                }}
              />
              <button
                onClick={addIdea}
                disabled={!currentInput.trim()}
                className={`w-full flex items-center justify-center gap-2 px-4 py-3 rounded-lg ${currentTheme.button} disabled:opacity-50 disabled:cursor-not-allowed transition-all`}
              >
                <Plus size={20} />
                Incluir (Ctrl+Enter)
              </button>
            </div>

            {/* Configuration */}
            <div className="mt-6 space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2">Propósito do Prompt</label>
                <div className="space-y-2">
                  <div className="flex gap-2">
                    {['Código', 'Imagem', 'Outros'].map((option) => (
                      <button
                        key={option}
                        onClick={() => setPurpose(option)}
                        className={`px-3 py-2 rounded-lg text-sm border transition-colors ${
                          purpose === option 
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
                      placeholder="Descreva o objetivo do prompt..."
                      className={`w-full px-3 py-2 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-blue-500`}
                    />
                  )}
                </div>
              </div>

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
                  className="w-full h-2 bg-gray-300 rounded-lg appearance-none cursor-pointer slider"
                />
              </div>
            </div>
          </div>

          {/* Ideas List */}
          <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg`}>
            <h2 className="text-xl font-semibold mb-4">💡 Suas Ideias ({ideas.length})</h2>
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
                onClick={generatePrompt}
                disabled={isGenerating}
                className={`w-full mt-4 flex items-center justify-center gap-2 px-4 py-3 rounded-lg bg-gradient-to-r from-purple-600 to-blue-600 text-white hover:from-purple-700 hover:to-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all transform hover:scale-105`}
              >
                <Wand2 size={20} className={isGenerating ? 'animate-spin' : ''} />
                {isGenerating ? 'Gerando...' : 'Me ajude, engenheiro?!'}
              </button>
            )}
          </div>

          {/* Generated Prompt */}
          <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg ${generatedPrompt ? 'lg:col-span-1' : ''}`}>
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-semibold">🚀 Prompt Estruturado</h2>
              {generatedPrompt && (
                <button
                  onClick={copyToClipboard}
                  className={`flex items-center gap-2 px-3 py-2 rounded-lg ${currentTheme.buttonSecondary} hover:bg-opacity-80 transition-colors`}
                >
                  {copied ? <Check size={16} /> : <Copy size={16} />}
                  {copied ? 'Copiado!' : 'Copiar'}
                </button>
              )}
            </div>
            
            {generatedPrompt ? (
              <div className="space-y-4">
                <div className={`text-xs ${currentTheme.textSecondary} flex justify-between`}>
                  <span>Caracteres: {generatedPrompt.length}</span>
                  <span>Limite: {maxLength.toLocaleString()}</span>
                </div>
                <div className={`max-h-96 overflow-y-auto p-4 rounded-lg border ${currentTheme.border} bg-opacity-50`}>
                  <pre className="whitespace-pre-wrap text-sm font-mono">{generatedPrompt}</pre>
                </div>
              </div>
            ) : (
              <div className={`${currentTheme.textSecondary} text-center py-12`}>
                <Wand2 size={48} className="mx-auto mb-4 opacity-50" />
                <p>Seu prompt estruturado aparecerá aqui</p>
                <p className="text-sm mt-2">Adicione ideias e clique em "Me ajude, engenheiro?!"</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default PromptCrafter;
