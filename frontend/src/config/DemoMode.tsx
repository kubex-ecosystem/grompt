// =============================================================================
// 🎯 DEMO MODE CONTROLLER - Single Source of Truth
// =============================================================================

// Tipos para as features disponíveis
export type FeatureKey = 'ollama' | 'openai' | 'gemini' | 'mcp_real' | 'agent_execution' | 'copilot';

// Tipos para os modos disponíveis
export type ModeKey = 'simple' | 'onboarding' | 'educational' | 'preview';

// Interface para status de feature
export interface FeatureStatus {
  ready: boolean;
  eta: string;
}

// Interface para conteúdo educacional
export interface EducationContent {
  title: string;
  description: string;
  benefits: string[];
}

// Interface para resposta de demo call
export interface DemoCallResponse {
  success: boolean;
  message: string;
  eta: string;
}

// Interface principal do DemoMode
export interface DemoModeConfig {
  isActive: boolean;
  modes: Record<string, ModeKey>;
  currentMode: ModeKey;
  features: Record<FeatureKey, FeatureStatus>;
  education: Record<string, EducationContent>;
  getLabel: (feature: FeatureKey, defaultLabel: string) => string;
  handleDemoCall: (feature: FeatureKey, action?: string) => DemoCallResponse | null;
}

const features: Record<FeatureKey, FeatureStatus> = {
  ollama: { ready: false, eta: 'Q2 2024' },
  openai: { ready: false, eta: 'Q1 2024' },
  gemini: { ready: false, eta: 'Q2 2024' },
  mcp_real: { ready: false, eta: 'Q1 2024' },
  agent_execution: { ready: false, eta: 'Q3 2024' },
  copilot: { ready: false, eta: 'Q2 2024' }
}

const modes: Record<string, ModeKey> = {
  SIMPLE: 'simple',
  ONBOARDING: 'onboarding',
  EDUCATIONAL: 'educational',
  PREVIEW: 'preview'
}

const education: Record<string, EducationContent> = {
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
}

const getLabel = (feature: FeatureKey, defaultLabel: string): string => {
  if (!DemoMode.isActive) return defaultLabel;

  const featureStatus = DemoMode.features[feature];

  switch (DemoMode.currentMode) {
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
}

const handleDemoCall = function (this: DemoModeConfig, feature: FeatureKey, action?: string): DemoCallResponse | null {
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
    message: messages[feature as keyof typeof messages] || 'Feature "' + feature + '" em modo demo',
    eta: this.features[feature] ? this.features[feature].eta : 'Em breve'
  };
}

const DemoMode: DemoModeConfig = {
  isActive: true,
  modes,
  currentMode: 'onboarding',
  features,
  education,
  getLabel,
  handleDemoCall
};

export default DemoMode;

// Exportações adicionais para uso externo
export { DemoMode };
export type { DemoModeConfig as DemoModeType };

