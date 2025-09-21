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

export default DemoMode;
