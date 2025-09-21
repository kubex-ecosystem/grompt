import { useCallback, useEffect, useState } from 'react';
import DemoMode from '../config/DemoMode';
import OnboardingSteps from '../constants/onboardingSteps.js';

const usePromptCrafter = ({ apiGenerate }) => {
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

  // Generation Logic using new API
  const generatePrompt = useCallback(async () => {
    if (ideas.length === 0) return;

    setIsGenerating(true);

    try {
      // Clear previous results
      setGeneratedPrompt('');

      const purposeText = purpose === 'Outros' && customPurpose
        ? customPurpose
        : purpose;

      // Prepare request for the API
      const generateRequest = {
        provider: agentProvider,
        ideas: ideas.map(idea => idea.text),
        purpose: purposeText.toLowerCase(),
        temperature: 0.7,
        maxTokens: maxLength,
        context: {
          outputType,
          agentFramework: outputType === 'agent' ? agentFramework : undefined,
          agentRole: outputType === 'agent' ? agentRole : undefined,
          agentTools: outputType === 'agent' ? agentTools : undefined,
          mcpServers: outputType === 'agent' ? mcpServers : undefined,
          maxLength: outputType === 'prompt' ? maxLength : undefined
        }
      };

      // Use streaming by default for better UX
      if (apiGenerate?.generateStream) {
        await apiGenerate.generateStream(generateRequest);
      } else if (apiGenerate?.generateSync) {
        const result = await apiGenerate.generateSync(generateRequest);
        setGeneratedPrompt(result.prompt);
      } else {
        // Fallback to legacy generation for demo mode
        await generatePromptLegacy(purposeText);
      }

    } catch (error) {
      console.error('Erro ao gerar:', error);
      setGeneratedPrompt('Erro ao gerar o ' + (outputType === 'prompt' ? 'prompt' : 'agent') + '. ' + error.message);
    }

    setIsGenerating(false);
  }, [ideas, purpose, customPurpose, outputType, agentProvider, agentFramework, agentRole, agentTools, mcpServers, maxLength, apiGenerate]);

  // Legacy generation for fallback and demo mode
  const generatePromptLegacy = async (purposeText) => {
    if (agentProvider !== 'claude' && DemoMode.isActive) {
      const demoResult = DemoMode.handleDemoCall(agentProvider);
      setGeneratedPrompt('# 🎪 Demo Mode\n\n' + demoResult.message + '\n\n**ETA:** ' + demoResult.eta + '\n\n---\n\n*Configurações salvas:*\n- Framework: ' + agentFramework + '\n- Provider: ' + agentProvider + '\n- Ferramentas: ' + (agentTools.join(', ') || 'Nenhuma') + '\n- Servidores MCP: ' + (mcpServers.join(', ') || 'Nenhum') + '\n\nEssas configurações serão aplicadas quando o provider estiver disponível!');
      return;
    }

    // If no API available, show fallback message
    setGeneratedPrompt('API integration not available. Please check the backend connection.');
  };

  // Clipboard functionality
  const copyToClipboard = useCallback(async () => {
    try {
      // Get content from API state or fallback to legacy state
      const contentToCopy = apiGenerate?.data?.prompt ||
        apiGenerate?.progress?.content ||
        generatedPrompt;

      if (!contentToCopy) {
        console.warn('No content to copy');
        return;
      }

      await navigator.clipboard.writeText(contentToCopy);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (error) {
      console.error('Erro ao copiar:', error);
    }
  }, [apiGenerate, generatedPrompt]);

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
