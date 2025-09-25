/**
 * Onboarding Steps Configuration - Integrated Bundle Version
 * Simple JavaScript object to avoid external module issues
 */

export interface OnboardingStep {
  id: string;
  title: string;
  content: string;
  target: string;
}

export const onboardingSteps: OnboardingStep[] = [
  {
    id: 'welcome',
    title: 'Bem-vindo ao Grompt! 🎉',
    content: 'Esta ferramenta transforma suas ideias em prompts profissionais e agents inteligentes.',
    target: 'header'
  },
  {
    id: 'ideas',
    title: 'Comece adicionando suas ideias 💡',
    content: 'Cole notas, pensamentos ou requisitos. A IA organizará tudo para você!',
    target: 'ideas-input'
  },
  {
    id: 'output-type',
    title: 'Escolha o que criar 🎯',
    content: 'Prompt = Instruções estruturadas | Agent = Código Python funcional',
    target: 'output-selector'
  },
  {
    id: 'mcp',
    title: 'Poder do MCP 🔌',
    content: 'Model Context Protocol conecta IA com ferramentas reais. Revolucionário!',
    target: 'mcp-section'
  }
];

// Default export for compatibility
export default onboardingSteps;