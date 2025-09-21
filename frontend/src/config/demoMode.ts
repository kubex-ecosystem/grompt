/**
 * Demo Mode Configuration - Integrated Bundle Version
 * All demo functionality consolidated to avoid external module issues
 */

export interface FeatureStatus {
  ready: boolean;
  eta: string;
}

export interface EducationContent {
  title: string;
  description: string;
  benefits: string[];
}

export interface DemoCallResponse {
  success: boolean;
  message: string;
  eta: string;
}

export type FeatureKey = 'autogen' | 'crewai' | 'ollama' |
  'openai' | 'gemini' | 'mcp_real' |
  'agent_execution' | 'copilot' | 'langchain' |
  'semantic-kernel' | 'custom_models' | 'custom' |
  'vectorstores' | 'tools' | 'agents' |
  'memory' | 'chat_models' | 'embedding_models';

// Demo configuration consolidated in bundle
const DEMO_CONFIG = {
  isActive: false, // Disable demo mode by default for production
  currentMode: 'simple' as const,

  features: {
    ollama: { ready: false, eta: 'Q2 2024' },
    openai: { ready: true, eta: 'Available' }, // OpenAI should work via backend
    gemini: { ready: false, eta: 'Q2 2024' },
    mcp_real: { ready: false, eta: 'Q1 2024' },
    agent_execution: { ready: false, eta: 'Q3 2024' },
    copilot: { ready: false, eta: 'Q2 2024' }
  } as Record<FeatureKey, FeatureStatus>,

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
  } as Record<string, EducationContent>,

  demoMessages: {
    ollama: '🦙 Ollama será integrado na versão completa! Conecte modelos locais diretamente.',
    openai: '🧠 OpenAI GPT-4 disponível via backend! Verifique a configuração da API.',
    gemini: '💎 Google Gemini chegando! Diversidade de modelos para diferentes tarefas.',
    mcp_real: '🔌 Servidores MCP reais em desenvolvimento! Conecte com qualquer sistema.',
    copilot: '🚁 GitHub Copilot API será integrada! Agents com capacidades de código avançadas.'
  } as Record<string, string>
};

/**
 * Demo Mode Controller - Simplified for bundling
 */
export class DemoMode {
  static get isActive(): boolean {
    return DEMO_CONFIG.isActive;
  }

  static get features(): Record<FeatureKey, FeatureStatus> {
    return DEMO_CONFIG.features;
  }

  static get education(): Record<string, EducationContent> {
    return DEMO_CONFIG.education;
  }

  static getLabel(feature: FeatureKey, defaultLabel: string): string {
    if (!this.isActive) {
      return defaultLabel;
    }

    const featureStatus = this.features[feature];

    if (featureStatus?.ready) {
      return defaultLabel;
    }

    return `${defaultLabel} 🎪`;
  }

  static handleDemoCall(feature: FeatureKey, action?: string): DemoCallResponse {
    if (!this.isActive) {
      return {
        success: true,
        message: 'Feature available',
        eta: 'Now'
      };
    }

    const message = DEMO_CONFIG.demoMessages[feature] || `Feature "${feature}" em modo demo`;
    const eta = this.features[feature]?.eta || 'Em breve';

    return {
      success: false,
      message,
      eta
    };
  }

  static enable(): void {
    DEMO_CONFIG.isActive = true;
  }

  static disable(): void {
    DEMO_CONFIG.isActive = false;
  }

  static setFeatureReady(feature: FeatureKey, ready: boolean = true): void {
    if (DEMO_CONFIG.features[feature]) {
      DEMO_CONFIG.features[feature].ready = ready;
      DEMO_CONFIG.features[feature].eta = ready ? 'Available' : DEMO_CONFIG.features[feature].eta;
    }
  }
}

// Default export for compatibility
export default DemoMode;
