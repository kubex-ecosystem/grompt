import React, { useState, useEffect } from 'react';
import { Trash2, Edit3, Plus, Wand2, Sun, Moon, Copy, Check, AlertCircle, ChevronDown, ChevronUp, RefreshCw } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import LanguageSelector from './components/LanguageSelector';

const PromptCrafter = () => {
  const { t } = useTranslation();
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

  const [isOutputCollapsed, setIsOutputCollapsed] = useState(true);
  const [isInputCollapsed, setIsInputCollapsed] = useState(false);

  // Controlar collapse automático e reativo de input/output
  useEffect(() => {
    if (generatedPrompt) {
      // Quando há prompt gerado: minimizar input, expandir output
      setIsInputCollapsed(true);
      setIsOutputCollapsed(false);
    } else {
      // Quando não há prompt: expandir input, minimizar output
      setIsInputCollapsed(false);
      setIsOutputCollapsed(true);
    }
  }, [generatedPrompt]);

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
  // eslint-disable-next-line react-hooks/exhaustive-deps
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
    if (editingId === id) {
      setEditingId(null);
      setEditingText('');
    }
    if (ideas.length === 0) {
      setIsInputCollapsed(false); // Expandir input se não houver ideias
      setIsOutputCollapsed(true); // Colapsar output se não houver ideias
    }
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

    setIsInputCollapsed(false); // Expandir input ao salvar edição (A ideia foi alterada)
    setIsOutputCollapsed(true); // Colapsar output ao salvar edição (O prompt gerado será alterado)
  };

  const cancelEdit = () => {
    setEditingId(null);
    setEditingText('');
  };

  // Função para limpar tudo e reiniciar
  const clearAll = () => {
    // Limpar todas as ideias
    setIdeas([]);
    // Limpar input atual
    setCurrentInput('');
    // Limpar prompt gerado
    setGeneratedPrompt('');
    // Cancelar qualquer edição em andamento
    setEditingId(null);
    setEditingText('');
    // Resetar estados de collapse para o inicial
    setIsInputCollapsed(false);
    setIsOutputCollapsed(true);
    // Resetar configurações para padrão
    setPurpose('Outros');
    setCustomPurpose('');
    setMaxLength(5000);
    // Feedback visual
    console.log('🧹 Interface limpa - pronto para começar!');
  };

  const generateDemoPrompt = () => {
    const purposeText = purpose === 'Outros' && customPurpose 
      ? customPurpose 
      : purpose;

    if (ideas.length === 0) {
      return `# ${t('output.title')} - ${purposeText}
        ## 🎯 ${t('demo.context')}
        ${t('demo.contextDesc', { purpose: purposeText.toLowerCase() })}
      `;
    }

    return `# ${t('output.title')} - ${purposeText}

      ## 🎯 ${t('demo.context')}

      ${t('demo.contextDesc', { purpose: purposeText.toLowerCase() })}

      ## 📝 ${t('demo.ideasTitle')}

      ${ideas.map((idea, index) => `**${index + 1}.** ${idea.text}`).join('\n')}

      ## 🔧 ${t('demo.instructions')}

      ${t('demo.instructionsList', { returnObjects: true }).map(instruction => `- ${instruction}`).join('\n')}

      ## 📋 ${t('demo.responseFormat')}

      ${t('demo.responseSteps', { returnObjects: true }).map((step, index) => `${index + 1}. ${step}`).join('\n')}

      ## ⚙️ ${t('demo.technicalConfig')}

      - ${t('demo.maxChars')}: ${maxLength.toLocaleString()}
      - ${t('demo.purpose')}: ${purposeText}
      - ${t('demo.totalIdeas')}: ${ideas.length}
      - ${t('demo.mode')}: ${ connectionStatus === 'connected' ? t('demo.modeConnected') : t('demo.modeOffline')}

      ---

      *${t('demo.footer')}*
      *${connectionStatus === 'connected' ? t('demo.footerConnected') : t('demo.footerOffline')}*`
        .replaceAll('      ', '')
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
      case 'connected': return t('connection.connected');
      case 'offline': return process.env.NODE_ENV === 'development' ? t('connection.offline') : t('connection.offlineProduction');
      default: return t('connection.checking');
    }
  };

  const inputDiv = (
    <div className={`${currentTheme.cardBg} rounded-xl border ${currentTheme.border} shadow-lg transition-all duration-500 ease-in-out hover:shadow-xl transform hover:scale-[1.02] will-change-transform ${
      isInputCollapsed ? 'h-20' : 'h-auto'
    }`}>

      <div className="items-center p-6 pb-0 gap-4">
        <div className="text-xl font-semibold mb-4 flex items-center">
          <h2>
            📝 {t('input.title')}
          </h2>
          <div className={`h-px flex-1 bg-gradient-to-r from-blue-500/20 to-transparent ml-4 transition-all duration-300 ${
            isInputCollapsed ? 'opacity-0' : 'opacity-100'
          }`}></div>
        </div>
      </div>

      {/* Conteúdo colapsável com animação super suave */}
      <div className={`overflow-hidden transition-all duration-700 ease-in-out transform origin-top p-6 ${
        isInputCollapsed 
          ? 'max-h-0 opacity-0 scale-y-0 -translate-y-4 pointer-events-none' 
          : 'max-h-[800px] opacity-100 scale-y-100 translate-y-0 pointer-events-auto'
      }`}>
        <div className={`transition-all duration-500 delay-100 ${
          isInputCollapsed ? 'opacity-0' : 'opacity-100'
        }`}>
        {/* Inputs editáveis */}
        <div className="space-y-4">
          <textarea
            value={currentInput}
            onChange={(e) => setCurrentInput(e.target.value)}
            placeholder={t('input.placeholder')}
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
            {t('input.addButton')}
          </button>
        </div>

        {/* Configuration */}
        <div className="mt-6 space-y-4">
          <div>
            <label className="block text-sm font-medium mb-2">{t('config.purpose')}</label>
            <div className="space-y-2">
              <div className="flex gap-2">
                {[
                  { key: 'purposeCode', value: 'Código' },
                  { key: 'purposeImage', value: 'Imagem' },
                  { key: 'purposeOthers', value: 'Outros' }
                ].map((option) => (
                  <button
                    key={option.value}
                    onClick={() => setPurpose(option.value)}
                    className={`px-3 py-2 rounded-lg text-sm border transition-colors ${
                      purpose === option.value 
                        ? 'bg-blue-600 text-white border-blue-600' 
                        : `${currentTheme.buttonSecondary} ${currentTheme.border}`
                    }`}
                  >
                    {t(`config.${option.key}`)}
                  </button>
                ))}
              </div>
              {purpose === 'Outros' && (
                <input
                  type="text"
                  value={customPurpose}
                  onChange={(e) => setCustomPurpose(e.target.value)}
                  placeholder={t('config.customPurpose')}
                  className={`w-full px-3 py-2 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-blue-500`}
                />
              )}
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">
              {t('config.maxLength')}: {maxLength.toLocaleString()} {t('config.characters')}
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
      </div>
    </div>
  );

  const outputDiv = (
    <div className={`${currentTheme.cardBg} rounded-xl border ${currentTheme.border} shadow-lg transition-all duration-500 ease-in-out hover:shadow-xl transform hover:scale-[1.02] will-change-transform ${
      isOutputCollapsed ? 'h-20' : 'h-auto'
    }`}>
      
      {/* Header sempre visível */}
      <div className="flex justify-between items-center p-6 pb-0">
        <div className="flex items-center gap-3">
          <h2 className="text-xl font-semibold">🚀 {t('output.title')}</h2>
          {generatedPrompt && (
            <div className={`px-2 py-1 rounded-full text-xs ${currentTheme.textSecondary} bg-opacity-50 ${currentTheme.cardBg} transition-all duration-300 ${
              isOutputCollapsed ? 'opacity-0 scale-0' : 'opacity-100 scale-100'
            }`}>
              {generatedPrompt.length.toLocaleString()} chars
            </div>
          )}
        </div>
        
        <div className="flex items-center gap-2">
          {/* Botão de copiar quando há conteúdo - com animação suave */}
          {generatedPrompt && !isOutputCollapsed && (
            <button
              onClick={copyToClipboard}
              className={`flex items-center gap-2 px-3 py-2 rounded-lg ${currentTheme.buttonSecondary} hover:bg-opacity-80 transition-all duration-300 transform hover:scale-105 active:scale-95`}
            >
              <div className="transition-all duration-200">
                {copied ? <Check size={16} /> : <Copy size={16} />}
              </div>
              {copied ? t('output.copied') : t('output.copy')}
            </button>
          )}
        </div>
      </div>
      
      {/* Conteúdo colapsável com animação super suave */}
      <div className={`overflow-hidden transition-all duration-700 ease-in-out transform origin-top ${
        isOutputCollapsed 
          ? 'max-h-0 opacity-0 scale-y-0 -translate-y-4 pointer-events-none' 
          : 'max-h-[800px] opacity-100 scale-y-100 translate-y-0 pointer-events-auto'
      }`}>
        <div className={`p-6 pt-4 transition-all duration-500 delay-100 ${
          isOutputCollapsed ? 'opacity-0' : 'opacity-100'
        }`}>
          {generatedPrompt ? (
            <div className="space-y-4">
              <div className={`text-xs ${currentTheme.textSecondary} flex justify-between items-center`}>
                <span>{t('output.characters')}: {generatedPrompt.length.toLocaleString()}</span>
                <span>{t('output.limit')}: {maxLength.toLocaleString()}</span>
                <div className={`w-24 h-1 rounded-full ${currentTheme.border} overflow-hidden`}>
                  <div 
                    className="h-full bg-gradient-to-r from-blue-500 to-purple-500 transition-all duration-300"
                    style={{ width: `${Math.min((generatedPrompt.length / maxLength) * 100, 100)}%` }}
                  />
                </div>
              </div>
              <div className={`max-h-96 overflow-y-auto p-4 rounded-lg border ${currentTheme.border} bg-opacity-50 ${currentTheme.cardBg}`}>
                <pre className="whitespace-pre-wrap text-sm font-mono leading-relaxed">{generatedPrompt}</pre>
              </div>
            </div>
          ) : (
            <div className={`${currentTheme.textSecondary} text-center py-16`}>
              <div className="flex flex-col items-center space-y-4">
                <Wand2 size={64} className="opacity-30 animate-pulse" />
                <div>
                  <p className="text-lg font-medium">{t('output.emptyTitle')}</p>
                  <p className="text-sm mt-2 opacity-75">{t('output.emptySubtitle')}</p>
                </div>
                <div className="flex items-center gap-2 mt-4 text-xs opacity-50">
                  <ChevronUp size={16} />
                  <span>Clique no ícone acima para minimizar esta seção</span>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
      
      {/* Indicador visual quando colapsado */}
      {isOutputCollapsed && (
        <div className="px-6 pb-3">
          <div className={`text-center ${currentTheme.textSecondary}`}>
            <div className="flex items-center justify-center gap-2 text-sm opacity-60">
              {generatedPrompt ? (
                <>
                  <span>{t('output.ready')} ({generatedPrompt.length.toLocaleString()} chars)</span>
                  <ChevronDown size={16} className="animate-bounce" />
                </>
              ) : (
                <>
                  <span>{t('output.minimized')}</span>
                  <ChevronDown size={16} className="animate-bounce" />
                </>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );

  return (
    <div className={`min-h-screen ${currentTheme.bg} ${currentTheme.text} p-4 transition-colors duration-300`}>
      <div className="max-w-[90%] mx-auto">{/* Expandido de max-w-7xl para 90% */}
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-4xl font-bold mb-2">
              <span className={currentTheme.accent}>{t('header.title')}</span>
            </h1>
            <p className={currentTheme.textSecondary}>
              {t('header.subtitle')}
            </p>
            {/* Debug info em desenvolvimento */}
            {process.env.NODE_ENV === 'development' && (
              <p className="text-xs text-yellow-400 mt-1">
                🔧 {t('header.debugMode')} | {t('header.baseUrl')}: {getBaseURL()} | {t('header.status')}: {connectionStatus}
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
                <option value="claude">{t('providers.claude')}</option>
              )}
              {availableAPIs.openai_available && (
                <option value="openai">{t('providers.openai')}</option>
              )}
              {availableAPIs.deepseek_available && (
                <option value="deepseek">{t('providers.deepseek')}</option>
              )}
              {availableAPIs.ollama_available && (
                <option value="ollama">{t('providers.ollama')}</option>
              )}
              <option value="demo">{t('providers.demo')}</option>
            </select>
            
            {/* Model Selection */}
            {apiProvider !== 'demo' && availableAPIs.available_models && availableAPIs.available_models[apiProvider] && (
              <select
                value={selectedModel}
                onChange={(e) => setSelectedModel(e.target.value)}
                className={`px-3 py-2 rounded-lg ${currentTheme.input} border focus:ring-2 focus:ring-blue-500`}
              >
                <option value="">{t('providers.defaultModel')}</option>
                {availableAPIs.available_models[apiProvider].map((model) => (
                  <option key={model} value={model}>{model}</option>
                ))}
              </select>
            )}
            
            <LanguageSelector currentTheme={currentTheme} />
            
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
              <strong>{t('alerts.offlineTitle')}</strong> 
              {process.env.NODE_ENV === 'development' 
                ? t('alerts.offlineDev')
                : t('alerts.offlineProduction')
              }
            </div>
          </div>
        )}

        {/* Server Info (em desenvolvimento) */}
        {process.env.NODE_ENV === 'development' && serverInfo && (
          <div className="mb-6 p-4 bg-blue-900 border border-blue-600 rounded-lg">
            <p className="text-blue-100">
              <strong>🔧 {t('alerts.serverInfo')}:</strong> v{serverInfo.version} | 
              APIs: {serverInfo.apis?.claude ? '✅' : '❌'} Claude, {serverInfo.apis?.ollama ? '✅' : '❌'} Ollama
            </p>
          </div>
        )}

        <div className={`grid gap-6 transition-all duration-700 ease-in-out ${
          // Grid que se adapta fluida e dinamicamente
          isInputCollapsed && isOutputCollapsed
            ? 'grid-cols-1' // Apenas Ideas quando ambos colapsados
            : isInputCollapsed && !isOutputCollapsed 
            ? 'grid-cols-1 lg:grid-cols-2' // Ideas + Output
            : !isInputCollapsed && isOutputCollapsed
            ? 'grid-cols-1 lg:grid-cols-2' // Input + Ideas
            : 'grid-cols-1 lg:grid-cols-3' // Todos expandidos: Input + Ideas + Output
        }`}>

          {/* Input Section - Sempre no DOM, mas com transições ultra-fluidas */}
          <div className={`transition-all duration-700 ease-out transform origin-center will-change-transform ${
            isInputCollapsed 
              ? 'opacity-0 scale-95 translate-y-4 pointer-events-none filter blur-sm' 
              : 'opacity-100 scale-100 translate-y-0 pointer-events-auto filter blur-0'
          } ${isInputCollapsed ? 'lg:hidden' : ''}`}>
            {inputDiv}
          </div>

          {/* Ideas List - Sempre visível e responsivo */}
          <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg hover:shadow-xl transition-all duration-500 ease-in-out transform hover:scale-[1.02] ${
            // Expande quando é o único card visível
            isInputCollapsed && isOutputCollapsed ? 'lg:col-span-1' : ''
          }`}>
            <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
              💡 {t('ideas.title')} 
              <span className={`px-2 py-1 rounded-full text-sm ${currentTheme.textSecondary} bg-blue-500/10 border border-blue-500/20`}>
                {ideas.length}
              </span>
              <div className="h-px flex-1 bg-gradient-to-r from-purple-500/20 to-transparent ml-4"></div>
              
              {/* Botão Limpar Tudo - aparece quando há conteúdo */}
              {(ideas.length > 0 || generatedPrompt || currentInput.trim()) && (
                <button
                  onClick={() => {
                    if (window.confirm(t('ideas.clearConfirm'))) {
                      clearAll();
                    }
                  }}
                  className={`group flex items-center gap-2 px-3 py-2 rounded-lg ${currentTheme.buttonSecondary} hover:bg-opacity-80 transition-all duration-300 transform hover:scale-105 active:scale-95 hover:bg-red-600 hover:text-white`}
                  title={t('ideas.clearAll')}
                >
                  <RefreshCw size={16} className="transition-transform duration-300 group-hover:rotate-180" />
                  <span className="text-sm font-medium">{t('ideas.clearAll')}</span>
                </button>
              )}
            </h2>
            <div className="space-y-3 max-h-96 overflow-y-auto">
              {ideas.length === 0 ? (
                <p className={`${currentTheme.textSecondary} text-center py-8`}>
                  {t('input.emptyState')}
                </p>
              ) : (
                ideas.map((idea, index) => (
                  <div 
                    key={idea.id} 
                    className={`p-3 rounded-lg border ${currentTheme.border} bg-opacity-50 transition-all duration-300 ease-out transform hover:scale-[1.02] hover:shadow-md animate-in slide-in-from-bottom-4 fade-in duration-500`}
                    style={{ animationDelay: `${index * 50}ms` }}
                  >
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
                            className="px-2 py-1 bg-green-600 text-white rounded text-xs hover:bg-green-700 transition-all duration-200 transform hover:scale-105 active:scale-95"
                          >
                            {t('ideas.save')}
                          </button>
                          <button
                            onClick={cancelEdit}
                            className={`px-2 py-1 rounded text-xs ${currentTheme.buttonSecondary} transition-all duration-200 transform hover:scale-105 active:scale-95`}
                          >
                            {t('ideas.cancel')}
                          </button>
                        </div>
                      </div>
                    ) : (
                      <>
                        <p className="text-sm mb-2">{idea.text}</p>
                        <div className="flex justify-end gap-1">
                          <button
                            onClick={() => startEditing(idea.id, idea.text)}
                            className={`p-1 rounded ${currentTheme.buttonSecondary} hover:bg-opacity-80 transition-all duration-200 transform hover:scale-110 active:scale-90`}
                            title={t('ideas.edit')}
                          >
                            <Edit3 size={14} />
                          </button>
                          <button
                            onClick={() => removeIdea(idea.id)}
                            className="p-1 rounded bg-red-600 text-white hover:bg-red-700 transition-all duration-200 transform hover:scale-110 active:scale-90"
                            title={t('ideas.delete')}
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
              <div className="items-bottom justify-between border-opacity-0 h-auto">
                <div className="mt-6 pt-4 border-t border-opacity-20 border-gradient-to-r from-purple-500 to-blue-500">
                  <button
                    onClick={generatePrompt}
                    disabled={isGenerating}
                    className={`
                      w-full 
                      flex 
                      items-center 
                      justify-center 
                      gap-3 
                      px-6 
                      py-4 
                      rounded-xl 
                      bg-gradient-to-r 
                      from-purple-600 
                      via-blue-600 
                      to-indigo-600 
                      text-white 
                      hover:from-purple-700 
                      hover:via-blue-700 
                      hover:to-indigo-700 
                      disabled:opacity-50 
                      disabled:cursor-not-allowed 
                      transition-all 
                      duration-500
                      transform 
                      hover:scale-105 
                      active:scale-95
                      hover:shadow-xl
                      hover:shadow-purple-500/25
                      font-medium 
                      text-lg 
                      relative 
                      overflow-hidden 
                      group
                      ${isGenerating ? 'animate-pulse' : ''}
                    `}
                  >
                    {/* Efeito de brilho no hover mais suave */}
                    <div className="absolute inset-0 bg-gradient-to-r from-white/0 via-white/20 to-white/0 translate-x-[-200%] group-hover:translate-x-[200%] transition-transform duration-1500 ease-in-out"></div>
                    
                    <div className={`transition-transform duration-300 ${isGenerating ? 'animate-spin' : 'group-hover:animate-pulse'}`}>
                      <Wand2 size={24} />
                    </div>
                    <span className="relative z-10 transition-all duration-300">
                      {isGenerating ? t('ideas.generating') : t('ideas.generateButton')}
                    </span>
                    
                    {!isGenerating && (
                      <div className="absolute right-4 opacity-60 group-hover:opacity-100 transition-all duration-300 group-hover:animate-bounce">
                        ✨
                      </div>
                    )}
                  </button>
                </div>
              </div>
            )}
          </div>

          {/* Generated Prompt - Sempre no DOM, mas com transições ultra-fluidas */}
          <div className={`transition-all duration-700 ease-out transform origin-center will-change-transform ${
            isOutputCollapsed 
              ? 'opacity-0 scale-95 translate-y-4 pointer-events-none filter blur-sm' 
              : 'opacity-100 scale-100 translate-y-0 pointer-events-auto filter blur-0'
          } ${isOutputCollapsed ? 'lg:hidden' : ''}`}>
            {outputDiv}
          </div>
        </div>
      </div>
    </div>
  );
};

export default PromptCrafter;
