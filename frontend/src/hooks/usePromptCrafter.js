import { useEffect, useState } from 'react';
import DemoMode from '../config/demoMode.js';
import OnboardingSteps from '../constants/onboardingSteps.js';

const usePromptCrafter = () => {
  // UI State
  const [darkMode, setDarkMode] = useState(true);

  // Ideas State
  const [currentInput, setCurrentInput] = useState('');
  const [ideas, setIdeas] = useState([]);
  const [editingId, setEditingId] = useState(null);
  const [editingText, setEditingText] = useState('');

  // Configuration State
  const [purpose, setPurpose] = useState('Outros');
  const [customPurpose, setCustomPurpose] = useState('');
  const [maxLength, setMaxLength] = useState(5000);
  const [outputType, setOutputType] = useState('prompt');

  // Agent Configuration
  const [agentFramework, setAgentFramework] = useState('crewai');
  const [agentRole, setAgentRole] = useState('');
  const [agentTools, setAgentTools] = useState([]);
  const [agentProvider, setAgentProvider] = useState('claude');
  const [mcpServers, setMcpServers] = useState([]);
  const [customMcpServer, setCustomMcpServer] = useState('');

  // Output State
  const [generatedPrompt, setGeneratedPrompt] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [copied, setCopied] = useState(false);

  // Onboarding State
  const [showOnboarding, setShowOnboarding] = useState(false);
  const [currentStep, setCurrentStep] = useState(0);
  const [showEducational, setShowEducational] = useState(false);
  const [educationalTopic, setEducationalTopic] = useState(null);

  // Dark mode effect
  useEffect(() => {
    document.documentElement.className = darkMode ? 'dark' : '';
  }, [darkMode]);

  // Ideas Management
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

  // Generation Logic
  const generatePrompt = async () => {
    if (ideas.length === 0) return;

    setIsGenerating(true);

    const purposeText = purpose === 'Outros' && customPurpose
      ? customPurpose
      : purpose;

    let engineeringPrompt = '';

    if (outputType === 'prompt') {
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
    } else if (outputType === 'agent') {
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
      if (agentProvider !== 'claude' && DemoMode.isActive) {
        const demoResult = DemoMode.handleDemoCall(agentProvider);
        setGeneratedPrompt('# 🎪 Demo Mode\n\n' + demoResult.message + '\n\n**ETA:** ' + demoResult.eta + '\n\n---\n\n*Configurações salvas:*\n- Framework: ' + agentFramework + '\n- Provider: ' + agentProvider + '\n- Ferramentas: ' + (agentTools.join(', ') || 'Nenhuma') + '\n- Servidores MCP: ' + (mcpServers.join(', ') || 'Nenhum') + '\n\nEssas configurações serão aplicadas quando o provider estiver disponível!');
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

      setGeneratedPrompt(result);
    } catch (error) {
      console.error('Erro ao gerar:', error);
      setGeneratedPrompt('Erro ao gerar o ' + (outputType === 'prompt' ? 'prompt' : 'agent') + '. ' + error.message);
    }

    setIsGenerating(false);
  };

  // Clipboard functionality
  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(generatedPrompt);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (error) {
      console.error('Erro ao copiar:', error);
    }
  };

  // Feature handling
  const handleFeatureClick = (feature) => {
    if (DemoMode.isActive && (!DemoMode.features[feature] || !DemoMode.features[feature].ready)) {
      const demoResult = DemoMode.handleDemoCall(feature);
      alert(demoResult.message + '\n\nETA: ' + demoResult.eta);
      return false;
    }
    return true;
  };

  // Onboarding management
  const startOnboarding = () => {
    setShowOnboarding(true);
    setCurrentStep(0);
  };

  const nextOnboardingStep = () => {
    if (currentStep < OnboardingSteps.length - 1) {
      setCurrentStep(currentStep + 1);
    } else {
      setShowOnboarding(false);
    }
  };

  // Educational modal management
  const showEducation = (topic) => {
    setEducationalTopic(topic);
    setShowEducational(true);
  };

  return {
    // State
    darkMode,
    currentInput,
    ideas,
    editingId,
    editingText,
    purpose,
    customPurpose,
    maxLength,
    generatedPrompt,
    isGenerating,
    copied,
    outputType,
    agentFramework,
    agentRole,
    agentTools,
    agentProvider,
    mcpServers,
    customMcpServer,
    showOnboarding,
    currentStep,
    showEducational,
    educationalTopic,

    // Setters
    setDarkMode,
    setCurrentInput,
    setEditingText,
    setPurpose,
    setCustomPurpose,
    setMaxLength,
    setOutputType,
    setAgentFramework,
    setAgentRole,
    setAgentTools,
    setAgentProvider,
    setMcpServers,
    setCustomMcpServer,
    setShowEducational,

    // Actions
    addIdea,
    removeIdea,
    startEditing,
    saveEdit,
    cancelEdit,
    generatePrompt,
    copyToClipboard,
    handleFeatureClick,
    startOnboarding,
    nextOnboardingStep,
    showEducation
  };
};

export default usePromptCrafter;
