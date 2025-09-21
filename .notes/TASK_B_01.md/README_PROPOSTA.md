# Agent & Prompt Crafter

Esta é uma SUGESTÃO DE APRIMORAMENTOS PARA a aplicação React para criação de prompts profissionais e agents inteligentes com Model Context Protocol (MCP).

## 📁 Estrutura de Arquivos

```plaintext
frontend/
├── components/              # Componentes React reutilizáveis
│   ├── Header.jsx          # Cabeçalho com título e controles
│   ├── OnboardingModal.jsx # Modal de tutorial
│   ├── EducationalModal.jsx # Modal educacional sobre MCP/Agents
│   ├── IdeasInput.jsx      # Input para adicionar ideias
│   ├── ConfigurationPanel.jsx # Painel de configurações
│   ├── IdeasList.jsx       # Lista de ideias adicionadas
│   ├── OutputPanel.jsx     # Painel de resultado/output
│   └── DemoStatusFooter.jsx # Footer com informações do demo
├── config/                 # Configurações globais
│   └── demoMode.js        # Configuração do modo demo
├── constants/              # Constantes e dados estáticos
│   ├── onboardingSteps.js # Passos do tutorial
│   └── themes.js          # Definições de tema claro/escuro
├── hooks/                  # Hooks customizados
│   └── usePromptCrafter.js # Hook principal com toda a lógica
├── utils/                  # Utilitários (futuro)
├── PromptCrafter.jsx      # Componente principal
├── index.js               # Ponto de entrada
└── Untitled-1             # Arquivo original (backup)
```

## 🎯 Componentes

### PromptCrafter.jsx

Componente principal que orquestra toda a aplicação, importando e organizando todos os sub-componentes.

### Header.jsx

- Título da aplicação
- Botões de tour e educação (quando em modo demo)
- Seletor de provider (Claude API)
- Toggle de tema claro/escuro

### OnboardingModal.jsx

Modal que apresenta o tutorial inicial para novos usuários.

### EducationalModal.jsx

Modal educacional que explica conceitos como MCP e Agents de IA.

### IdeasInput.jsx

Área para adicionar novas ideias/notas, com textarea e botão de envio.

### ConfigurationPanel.jsx

Painel de configurações que inclui:

- Seletor de tipo de output (Prompt vs Agent)
- Configurações específicas de agents (framework, provider, role, tools)
- Configurações de servidores MCP
- Seleção de propósito
- Controle de tamanho máximo (para prompts)

### IdeasList.jsx

Lista das ideias adicionadas com:

- Edição inline
- Remoção de itens
- Botão de geração

### OutputPanel.jsx

Painel que mostra:

- Prompt/Agent gerado
- Informações de caracteres e configurações
- Botão de copiar
- Estado vazio quando não há output

### DemoStatusFooter.jsx

Footer informativo sobre o status do modo demo e funcionalidades.

## 🔧 Configuração e Estado

### config/demoMode.js

Controlador central do modo demo com:

- Status de funcionalidades
- Modos de demonstração
- Conteúdo educacional
- Handlers para features não implementadas

### hooks/usePromptCrafter.js

Hook principal que centraliza:

- Todo o estado da aplicação
- Lógica de manipulação de ideias
- Geração de prompts/agents via API
- Gestão de onboarding e modais
- Funcionalidades de clipboard

### constants/

- `themes.js`: Definições de cores e estilos para modo claro/escuro
- `onboardingSteps.js`: Dados dos passos do tutorial

## 🚀 Como Usar

```jsx
import PromptCrafter from './grompt';

function App() {
  return <PromptCrafter />;
}
```

## 🎪 Modo Demo

A aplicação funciona em modo demo com:

- ✅ Claude API funcional
- 🎪 Outras funcionalidades simuladas
- Onboarding interativo
- Conteúdo educacional sobre MCP

## 📝 Dependências

- React + Hooks
- Lucide React (ícones)
- Tailwind CSS (estilos)
- Claude API (Anthropic)
