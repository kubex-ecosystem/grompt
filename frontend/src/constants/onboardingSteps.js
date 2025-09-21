// =============================================================================
// 🎓 ONBOARDING STEPS
// =============================================================================

const OnboardingSteps = [
  {
    id: 'welcome',
    title: 'Bem-vindo ao Agent Crafter! 🎉',
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

export default OnboardingSteps;
