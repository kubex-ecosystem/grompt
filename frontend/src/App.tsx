import { AlertTriangle, BrainCircuit, Bot, Image as ImageIcon, MessageCircle, NotebookPen, Sparkles, Workflow } from 'lucide-react';
import React, { useEffect, useMemo, useState } from 'react';
import ChatInterface, { type ChatMessage } from './components/features/ChatInterface';
import CodeGenerator from './components/features/CodeGenerator';
import ContentSummarizer from './components/features/ContentSummarizer';
import ImageGenerator from './components/features/ImageGenerator';
import AgentsGenerator from './components/features/AgentsGenerator';
import PromptCrafter from './components/features/PromptCrafter';
import Welcome from './components/features/Welcome';
import Footer from './components/layout/Footer';
import Header from './components/layout/Header';
import Layout from './components/layout/Layout';
import Sidebar from './components/layout/Sidebar';
import ProjectExtractor from './components/projectsFiles/Extractor';
import { LanguageContext } from './context/LanguageContext';
import { useLanguage } from './hooks/useLanguage';
import { useTheme } from './hooks/useTheme';
import { useAnalytics } from './services/analytics';
import { initStorage } from './services/storageService';
import { unifiedAIService } from './services/unifiedAIService';

type SectionKey =
  | 'welcome'
  | 'prompt'
  | 'agents'
  | 'chat'
  | 'summaries'
  | 'code'
  | 'images'
  | 'projects';

const App: React.FC = () => {
  const [theme, toggleTheme] = useTheme();
  const { language, setLanguage, t } = useLanguage();
  const [isApiKeyMissing, setIsApiKeyMissing] = useState(false);
  const [activeSection, setActiveSection] = useState<SectionKey>('welcome');
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

  const languageContextValue = useMemo(
    () => ({ language, setLanguage, t }),
    [language, setLanguage, t]
  );

  useAnalytics();

  useEffect(() => {
    initStorage();
    if (!process.env.API_KEY) {
      console.log('Running in demo mode - API key not configured');
      setIsApiKeyMissing(true);
    }
    if (typeof window !== 'undefined') {
      const stored = window.localStorage.getItem('grompt.sidebar.collapsed');
      if (stored !== null) {
        setSidebarCollapsed(stored === '1');
      }
    }
  }, []);

  const handleToggleSidebarCollapse = () => {
    setSidebarCollapsed((prev) => {
      const next = !prev;
      if (typeof window !== 'undefined') {
        window.localStorage.setItem('grompt.sidebar.collapsed', next ? '1' : '0');
      }
      return next;
    });
  };

  const sections = useMemo(
    () => [
      {
        id: 'welcome',
        label: 'Visão Geral',
        description: 'Tour rápido pelas ferramentas Kubex',
        icon: Sparkles,
      },
      {
        id: 'prompt',
        label: 'Prompt Crafter',
        description: 'Estruture ideias e governe entregas',
        icon: BrainCircuit,
      },
      {
        id: 'agents',
        label: 'Agents Squad',
        description: 'Modele squads IA alinhados ao Kubex',
        icon: Bot,
      },
      {
        id: 'chat',
        label: 'Chat Assistido',
        description: 'Converse com seu provedor principal',
        icon: MessageCircle,
      },
      {
        id: 'summaries',
        label: 'Summarizer',
        description: 'Resumos executivos e planos de ação',
        icon: NotebookPen,
      },
      {
        id: 'code',
        label: 'Code Generator',
        description: 'Gere scaffolds idiomáticos',
        icon: Workflow,
      },
      {
        id: 'images',
        label: 'Visual Prompts',
        description: 'Briefings para modelos de imagem',
        icon: ImageIcon,
      },
      {
        id: 'projects',
        label: 'Project Extractor',
        description: 'Explore artefatos e arquivos',
        icon: NotebookPen,
      },
    ],
    []
  );

  const handleChatSend = async (
    history: ChatMessage[],
    input: string,
    apiKey?: string
  ): Promise<{ content: string; provider?: string } | null> => {
    const transcript = history
      .map((message) => `${message.role === 'user' ? 'User' : 'Assistant'}: ${message.content}`)
      .join('\n');

    const systemInstruction = `You are Kubex Copilot, a pragmatic assistant that answers in concise, actionable blocks.
Always respect the Kubex principles: Modular Thinking, Practicality, Interoperability.`;

    const composedPrompt = `${systemInstruction}

Conversation so far:
${transcript}

Reply to the latest user message with a helpful, structured answer. Use Markdown when it improves comprehension.`;

    try {
      // BYOK Support: Pass external API key if provided
      const response = await unifiedAIService.generateDirectPrompt(composedPrompt, undefined, undefined, 700, apiKey);
      return {
        content: response.response.trim(),
        provider: response.provider,
      };
    } catch (error) {
      console.error('Chat generation failed', error);
      return {
        content: 'Não foi possível contactar o provedor. Estamos operando em modo demo — revise o backend /api/unified.',
      };
    }
  };

  const handleSummarize = async (input: string, tone: string, maxWords: number, apiKey?: string): Promise<string> => {
    const prompt = `You are an expert summarizer for the Kubex ecosystem.
Tone requested: ${tone}.
Target length: up to ${maxWords} words.

Summarize the following content into a structured Markdown deliverable with sections for Context, Key Insights, and Recommended Next Steps.

---
${input}
---

Return only the formatted summary.`;
    try {
      // BYOK Support: Pass external API key if provided
      const response = await unifiedAIService.generateDirectPrompt(
        prompt,
        undefined,
        undefined,
        Math.min(maxWords * 3, 1200),
        apiKey
      );
      return response.response.trim();
    } catch (error) {
      console.error('Summarization failed', error);
      throw new Error('Falha ao gerar o resumo (verifique a rota /api/unified).');
    }
  };

  const handleCodeGenerate = async (spec: {
    stack: string;
    goal: string;
    constraints: string[];
    extras: string;
  }, apiKey?: string): Promise<string> => {
    const constraintList = spec.constraints.length > 0 ? spec.constraints.join('; ') : 'Sem constraints adicionais.';
    const prompt = `You are an experienced engineer generating boilerplates aligned with Kubex guidelines.

Tech stack: ${spec.stack}
Goal: ${spec.goal}
Constraints: ${constraintList}
Additional notes: ${spec.extras || 'Nenhuma.'}

Produce a concise Markdown response containing:
1. Implementation outline (bullet list).
2. Key code snippet(s) with explanations.
3. Checklist of follow-up tasks.

Use fenced code blocks with the appropriate language identifier and avoid proprietary dependencies.`;

    try {
      // BYOK Support: Pass external API key if provided
      const response = await unifiedAIService.generateDirectPrompt(prompt, undefined, undefined, 1200, apiKey);
      return response.response.trim();
    } catch (error) {
      console.error('Code generation failed', error);
      throw new Error('Não foi possível gerar o blueprint. Confirme se o backend unificado está operacional.');
    }
  };

  const handleImageBrief = async (payload: {
    subject: string;
    mood: string;
    style: string;
    details: string;
  }, apiKey?: string): Promise<string> => {
    const prompt = `Craft a single prompt for an image generation model.

Subject: ${payload.subject}
Mood: ${payload.mood}
Visual style: ${payload.style}
Additional details: ${payload.details || 'Nenhum detalhe extra informado.'}

Return the final prompt in Markdown with sections Persona, Composition, Style Notes, and Output Requirements. Keep it under 180 words.`;

    try {
      // BYOK Support: Pass external API key if provided
      const response = await unifiedAIService.generateDirectPrompt(prompt, undefined, undefined, 400, apiKey);
      return response.response.trim();
    } catch (error) {
      console.error('Image briefing failed', error);
      throw new Error('Erro ao criar o briefing visual. Forneça um provider compatível em /api/unified.');
    }
  };

  const renderActiveSection = () => {
    switch (activeSection) {
      case 'welcome':
        return <Welcome onGetStarted={() => setActiveSection('prompt')} />;
      case 'prompt':
        return <PromptCrafter theme={theme} isApiKeyMissing={isApiKeyMissing} />;
      case 'agents':
        return <AgentsGenerator />;
      case 'chat':
        return <ChatInterface onSend={handleChatSend} />;
      case 'summaries':
        return <ContentSummarizer onSummarize={handleSummarize} />;
      case 'code':
        return <CodeGenerator onGenerate={handleCodeGenerate} />;
      case 'images':
        return <ImageGenerator onCraftPrompt={handleImageBrief} />;
      case 'projects':
        return (
          <ProjectExtractor
            projectFile="frontend"
            projectName="Grompt Frontend"
            description="Explore a árvore de arquivos carregada diretamente do servidor."
          />
        );
      default:
        return null;
    }
  };

  return (
    <LanguageContext.Provider value={languageContextValue}>
      <Layout
        sidebar={
          <Sidebar
            sections={sections}
            activeSection={activeSection}
            onSectionChange={(section) => {
              setActiveSection(section as SectionKey);
              setSidebarOpen(false);
            }}
            onClose={() => setSidebarOpen(false)}
            collapsed={sidebarCollapsed}
          />
        }
        header={
          <Header
            theme={theme}
            toggleTheme={toggleTheme}
            onToggleMenu={() => setSidebarOpen(true)}
            onToggleSidebarCollapse={handleToggleSidebarCollapse}
            sidebarCollapsed={sidebarCollapsed}
          />
        }
        footer={<Footer />}
        sidebarOpen={sidebarOpen}
        sidebarCollapsed={sidebarCollapsed}
        onSidebarClose={() => setSidebarOpen(false)}
      >
        {isApiKeyMissing && (
          <div
            className="mb-6 flex items-start gap-3 rounded-2xl border border-blue-200/70 bg-blue-50/80 p-4 text-sm text-slate-600 shadow-sm dark:border-blue-400/40 dark:bg-blue-900/20 dark:text-slate-100"
            role="alert"
          >
            <AlertTriangle className="mt-1 h-5 w-5 text-blue-500 dark:text-blue-300" />
            <div>
              <p className="font-semibold">{t('apiKeyMissingTitle')}</p>
              <p className="text-sm">{t('apiKeyMissingText')}</p>
            </div>
          </div>
        )}

        {renderActiveSection()}
      </Layout>
    </LanguageContext.Provider>
  );
};

export default App;
