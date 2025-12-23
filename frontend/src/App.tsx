import React, { useEffect, useRef, useState } from 'react';
import AgentsGenerator from './components/features/AgentsGenerator';
import ChatInterface from './components/features/ChatInterface';
import CodeGenerator from './components/features/CodeGenerator';
import ContentSummarizer from './components/features/ContentSummarizer';
import ImageGenerator from './components/features/ImageGenerator';
import PromptCrafter from './components/features/PromptCrafter';
import Welcome from './components/features/Welcome';
import Footer from './components/layout/Footer';
import Header from './components/layout/Header';
import Layout from './components/layout/Layout';
import Sidebar, { SidebarSection } from './components/layout/Sidebar';
import { LanguageContext } from './context/LanguageContext';
import { configService } from './services/configService';
import { Language, Theme } from './types';

// Translation strings
const translations: Record<Language, Record<string, string>> = {
  en: {
    poweredBy: 'Powered by Kubex Ecosystem',
    motto: 'Governance, Creativity, Productivity, Freedom',
    lang_en: 'English',
    lang_pt: 'Portuguese',
    headerTitle: 'Grompt Workspace',
    headerTagline: 'Engineering freedom to evolve real systems',
    sidebarTitle: 'Workspace',
    sidebarSuiteLabel: 'Kubex Suite',
    sidebarHubTitle: 'Grompt Hub',
    sidebarDescription: 'Navigate through engineering modules for creation, analysis, and delivery.',
    sidebarFooterLineOne: 'Build long-lived systems with deliberate tooling.',
    sidebarFooterLineTwo: 'Portable, explicit, and human-guided.',
    sectionWelcomeLabel: 'Welcome',
    sectionWelcomeDescription: 'Getting started',
    sectionPromptLabel: 'Prompt Crafter',
    sectionPromptDescription: 'Build prompts',
    sectionAgentsLabel: 'Agents',
    sectionAgentsDescription: 'AI Agents',
    sectionChatLabel: 'Chat',
    sectionChatDescription: 'Conversations',
    sectionSummarizerLabel: 'Summarizer',
    sectionSummarizerDescription: 'Content summary',
    sectionCodeLabel: 'Code',
    sectionCodeDescription: 'Code generation',
    sectionImagesLabel: 'Images',
    sectionImagesDescription: 'Image generation',
    welcomeKicker: 'Kubex Ecosystem',
    welcomeHeadline: 'Autonomous, agnostic workspace for engineering.',
    welcomeSubheadline:
      'Grompt organizes creation, analysis, and delivery of software systems with explicit tools, including AI when needed. The focus is architecture, operation, and continuous evolution.',
    welcomeCta: 'Explore the workspace',
    welcomeStackTitle: 'Technical workspace',
    welcomeStackItemOne: 'External providers and tools connected as optional resources.',
    welcomeStackItemTwo: 'Local persistence for continuity and technical context.',
    welcomeStackItemThree: 'Responsive interface organized around work flows.',
    welcomePhaseCreation: 'Creation',
    welcomePhaseAnalysis: 'Analysis',
    welcomePhaseConsolidation: 'Consolidation',
    welcomePromptTitle: 'Prompt Crafter',
    welcomePromptDescription: 'Structure context, objectives, and constraints to start with precision and traceability.',
    welcomeChatTitle: 'Assisted Conversations',
    welcomeChatDescription: 'Investigate decisions, hypotheses, and trade-offs with contextual dialogue.',
    welcomeSummaryTitle: 'Summaries & Deliverables',
    welcomeSummaryDescription: 'Consolidate results into technical summaries and clear deliverables.',
  },
  pt: {
    poweredBy: 'Desenvolvido pelo Kubex Ecosystem',
    motto: 'Governança, Criatividade, Produtividade, Liberdade',
    lang_en: 'English',
    lang_pt: 'Português',
    headerTitle: 'Grompt Workspace',
    headerTagline: 'Liberdade de engenharia para evoluir sistemas reais',
    sidebarTitle: 'Workspace',
    sidebarSuiteLabel: 'Kubex Suite',
    sidebarHubTitle: 'Grompt Hub',
    sidebarDescription: 'Navegue por módulos de engenharia para criação, análise e entrega.',
    sidebarFooterLineOne: 'Construa sistemas duradouros com ferramentas explícitas.',
    sidebarFooterLineTwo: 'Portátil, transparente e guiado por pessoas.',
    sectionWelcomeLabel: 'Boas-vindas',
    sectionWelcomeDescription: 'Primeiros passos',
    sectionPromptLabel: 'Prompt Crafter',
    sectionPromptDescription: 'Construir prompts',
    sectionAgentsLabel: 'Agentes',
    sectionAgentsDescription: 'Agentes de IA',
    sectionChatLabel: 'Chat',
    sectionChatDescription: 'Conversas',
    sectionSummarizerLabel: 'Sumarizador',
    sectionSummarizerDescription: 'Resumo de conteúdo',
    sectionCodeLabel: 'Código',
    sectionCodeDescription: 'Geração de código',
    sectionImagesLabel: 'Imagens',
    sectionImagesDescription: 'Geração de imagens',
    welcomeKicker: 'Kubex Ecosystem',
    welcomeHeadline: 'Workspace autônomo e agnóstico para engenharia.',
    welcomeSubheadline:
      'O Grompt organiza criação, análise e entrega de sistemas de software com ferramentas explícitas, incluindo IA quando necessário. O foco é arquitetura, operação e evolução contínua.',
    welcomeCta: 'Explorar o workspace',
    welcomeStackTitle: 'Workspace técnico',
    welcomeStackItemOne: 'Provedores e ferramentas externos conectados como recursos opcionais.',
    welcomeStackItemTwo: 'Persistência local para continuidade e contexto técnico.',
    welcomeStackItemThree: 'Interface responsiva organizada por fluxos de trabalho.',
    welcomePhaseCreation: 'Criação',
    welcomePhaseAnalysis: 'Análise',
    welcomePhaseConsolidation: 'Consolidação',
    welcomePromptTitle: 'Prompt Crafter',
    welcomePromptDescription: 'Estruture contexto, objetivos e restrições para iniciar trabalhos com precisão e rastreabilidade.',
    welcomeChatTitle: 'Conversas Assistidas',
    welcomeChatDescription: 'Investigue decisões, hipóteses e trade-offs com conversa contextual.',
    welcomeSummaryTitle: 'Sumarização e Entregáveis',
    welcomeSummaryDescription: 'Consolide resultados em resumos técnicos e entregáveis claros.',
  },
};

const SIDEBAR_COLLAPSED_KEY = 'grompt.sidebar.collapsed';
const SIDEBAR_AUTO_EXPAND_SEEN_KEY = 'grompt.sidebar.autoExpandSeen';

const SECTION_IDS = ['welcome', 'prompt', 'agents', 'chat', 'summarizer', 'code', 'images'] as const;
type SectionId = (typeof SECTION_IDS)[number];
const ACTIVE_SECTION_KEY = 'grompt.activeSection';

const getSectionFromHash = (hash: string): SectionId | null => {
  if (hash.startsWith('#prompt=')) return 'prompt';
  if (hash.startsWith('#section=')) {
    const raw = hash.replace('#section=', '');
    const candidate = decodeURIComponent(raw) as SectionId;
    if (SECTION_IDS.includes(candidate)) return candidate;
  }
  return null;
};

const App: React.FC = () => {
  // State management
  const [theme, setTheme] = useState<Theme>('dark');
  const [language, setLanguage] = useState<Language>('en');
  const [activeSection, setActiveSection] = useState<SectionId>(() => {
    if (typeof window === 'undefined') return 'welcome';
    const fromHash = getSectionFromHash(window.location.hash);
    if (fromHash) return fromHash;
    const saved = localStorage.getItem(ACTIVE_SECTION_KEY) as SectionId | null;
    if (saved && SECTION_IDS.includes(saved)) return saved;
    return 'welcome';
  });
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [sidebarCollapsed, setSidebarCollapsed] = useState<boolean>(() => {
    if (typeof window === 'undefined') return true;
    const saved = localStorage.getItem(SIDEBAR_COLLAPSED_KEY);
    if (saved === 'true') return true;
    if (saved === 'false') return false;
    return true;
  });
  const [autoExpandSeen, setAutoExpandSeen] = useState<boolean>(() => {
    if (typeof window === 'undefined') return false;
    return localStorage.getItem(SIDEBAR_AUTO_EXPAND_SEEN_KEY) === 'true';
  });
  const userInteractedRef = useRef(false);
  const activeSectionRef = useRef(activeSection);
  const [demoMode, setDemoMode] = useState<boolean>(true); // Default to demo mode

  // Load theme from localStorage
  useEffect(() => {
    const savedTheme = localStorage.getItem('theme') as Theme;
    if (savedTheme && (savedTheme === 'light' || savedTheme === 'dark')) {
      setTheme(savedTheme);
      if (savedTheme === 'dark') {
        document.documentElement.classList.add('dark');
      } else {
        document.documentElement.classList.remove('dark');
      }
    }
  }, []);

  // Load language from localStorage
  useEffect(() => {
    const savedLanguage = localStorage.getItem('language') as Language;
    if (savedLanguage && ['en', 'pt'].includes(savedLanguage)) {
      setLanguage(savedLanguage);
      return;
    }
    if (savedLanguage) {
      localStorage.removeItem('language');
    }
  }, []);

  useEffect(() => {
    activeSectionRef.current = activeSection;
    localStorage.setItem(ACTIVE_SECTION_KEY, activeSection);
  }, [activeSection]);

  useEffect(() => {
    const handleHashChange = () => {
      const nextSection = getSectionFromHash(window.location.hash);
      if (nextSection && nextSection !== activeSectionRef.current) {
        setActiveSection(nextSection);
      }
    };
    window.addEventListener('hashchange', handleHashChange);
    handleHashChange();
    return () => window.removeEventListener('hashchange', handleHashChange);
  }, []);

  useEffect(() => {
    localStorage.setItem(SIDEBAR_AUTO_EXPAND_SEEN_KEY, String(autoExpandSeen));
  }, [autoExpandSeen]);

  useEffect(() => {
    if (typeof window === 'undefined') return;
    if (window.location.hash.startsWith('#prompt=')) return;
    const nextHash = `#section=${encodeURIComponent(activeSection)}`;
    if (window.location.hash !== nextHash) {
      window.history.replaceState(null, document.title, `${window.location.pathname}${window.location.search}${nextHash}`);
    }
  }, [activeSection]);

  useEffect(() => {
    if (autoExpandSeen) return;
    if (!sidebarCollapsed) {
      setAutoExpandSeen(true);
      return;
    }

    let timer = window.setTimeout(() => {
      if (userInteractedRef.current) return;
      setSidebarCollapsed(false);
      setAutoExpandSeen(true);
    }, 3500);

    const resetTimer = () => {
      if (userInteractedRef.current) return;
      window.clearTimeout(timer);
      timer = window.setTimeout(() => {
        if (userInteractedRef.current) return;
        setSidebarCollapsed(false);
        setAutoExpandSeen(true);
      }, 3500);
    };

    window.addEventListener('scroll', resetTimer, { passive: true });
    window.addEventListener('mousemove', resetTimer);
    window.addEventListener('keydown', resetTimer);
    window.addEventListener('touchstart', resetTimer);

    return () => {
      window.clearTimeout(timer);
      window.removeEventListener('scroll', resetTimer);
      window.removeEventListener('mousemove', resetTimer);
      window.removeEventListener('keydown', resetTimer);
      window.removeEventListener('touchstart', resetTimer);
    };
  }, [autoExpandSeen, sidebarCollapsed]);

  // Fetch backend config to check demo mode
  useEffect(() => {
    const fetchConfig = async () => {
      try {
        const isDemo = await configService.isDemoMode();
        setDemoMode(isDemo);
        console.log('[App] Demo mode status from backend:', isDemo);
      } catch (error) {
        console.error('[App] Failed to fetch config:', error);
        // Keep default demo mode = true on error
      }
    };
    fetchConfig();
  }, []);

  // Theme toggle handler
  const handleToggleTheme = () => {
    const newTheme: Theme = theme === 'dark' ? 'light' : 'dark';
    setTheme(newTheme);
    localStorage.setItem('theme', newTheme);
    if (newTheme === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  };

  // Language change handler
  const handleSetLanguage = (lang: Language) => {
    setLanguage(lang);
    localStorage.setItem('language', lang);
  };

  const markSidebarInteraction = () => {
    if (!userInteractedRef.current) {
      userInteractedRef.current = true;
      setAutoExpandSeen(true);
    }
  };

  const setSidebarCollapsedPreference = (next: boolean) => {
    localStorage.setItem(SIDEBAR_COLLAPSED_KEY, String(next));
    setSidebarCollapsed(next);
  };

  const handleSidebarToggle = () => {
    markSidebarInteraction();
    if (window.matchMedia('(min-width: 1024px)').matches) {
      setSidebarCollapsedPreference(!sidebarCollapsed);
      return;
    }
    setSidebarOpen(prev => !prev);
  };

  // Translation function
  const t = (key: string, params?: Record<string, string>): string => {
    const localeStrings = translations[language] || translations.en;
    let translation = localeStrings[key] || translations.en[key] || key;
    if (params) {
      Object.keys(params).forEach(paramKey => {
        translation = translation.replace(`{${paramKey}}`, params[paramKey]);
      });
    }
    return translation;
  };

  // Sidebar sections configuration
  const sections: SidebarSection[] = [
    {
      id: 'welcome',
      label: t('sectionWelcomeLabel'),
      description: t('sectionWelcomeDescription'),
    },
    {
      id: 'prompt',
      label: t('sectionPromptLabel'),
      description: t('sectionPromptDescription'),
    },
    {
      id: 'agents',
      label: t('sectionAgentsLabel'),
      description: t('sectionAgentsDescription'),
    },
    {
      id: 'chat',
      label: t('sectionChatLabel'),
      description: t('sectionChatDescription'),
    },
    {
      id: 'summarizer',
      label: t('sectionSummarizerLabel'),
      description: t('sectionSummarizerDescription'),
    },
    {
      id: 'code',
      label: t('sectionCodeLabel'),
      description: t('sectionCodeDescription'),
    },
    {
      id: 'images',
      label: t('sectionImagesLabel'),
      description: t('sectionImagesDescription'),
    },
  ];

  // Render active section content
  const renderContent = () => {
    switch (activeSection) {
      case 'welcome':
        return <Welcome onGetStarted={() => setActiveSection('prompt')} />;
      case 'prompt':
        return <PromptCrafter theme={theme} isApiKeyMissing={demoMode} />;
      case 'agents':
        return <AgentsGenerator theme={theme} isApiKeyMissing={demoMode} />;
      case 'chat':
        return <ChatInterface theme={theme} isApiKeyMissing={demoMode} />;
      case 'summarizer':
        return <ContentSummarizer theme={theme} isApiKeyMissing={demoMode} />;
      case 'code':
        return <CodeGenerator theme={theme} isApiKeyMissing={demoMode} />;
      case 'images':
        return <ImageGenerator theme={theme} isApiKeyMissing={demoMode} />;
      default:
        return <Welcome onGetStarted={() => setActiveSection('prompt')} />;
    }
  };

  const handleSectionChange = (section: string) => {
    if (SECTION_IDS.includes(section as SectionId)) {
      setActiveSection(section as SectionId);
      return;
    }
    setActiveSection('welcome');
  };

  return (
    <LanguageContext.Provider value={{ language, setLanguage: handleSetLanguage, t }}>
      <Layout
        sidebar={
          <Sidebar
            sections={sections}
            activeSection={activeSection}
            onSectionChange={handleSectionChange}
            onClose={() => setSidebarOpen(false)}
            collapsed={sidebarCollapsed}
          />
        }
        header={
          <Header
            theme={theme}
            onToggleTheme={handleToggleTheme}
            onToggleSidebar={handleSidebarToggle}
            collapsed={sidebarCollapsed}
          />
        }
        footer={<Footer />}
        sidebarOpen={sidebarOpen}
        onSidebarClose={() => {
          markSidebarInteraction();
          setSidebarOpen(false);
        }}
        sidebarCollapsed={sidebarCollapsed}
      >
        {renderContent()}
      </Layout>
    </LanguageContext.Provider>
  );
};

export default App;
