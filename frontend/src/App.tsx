import React, { useEffect, useState } from 'react';
import { Language, Theme } from './types';
import { LanguageContext } from './context/LanguageContext';
import Layout from './components/layout/Layout';
import Header from './components/layout/Header';
import Sidebar, { SidebarSection } from './components/layout/Sidebar';
import Footer from './components/layout/Footer';
import Welcome from './components/features/Welcome';
import PromptCrafter from './components/features/PromptCrafter';
import AgentsGenerator from './components/features/AgentsGenerator';
import ChatInterface from './components/features/ChatInterface';
import ContentSummarizer from './components/features/ContentSummarizer';
import CodeGenerator from './components/features/CodeGenerator';
import ImageGenerator from './components/features/ImageGenerator';
import { configService } from './services/configService';

// Translation strings
const translations: Record<Language, Record<string, string>> = {
  en: {
    poweredBy: 'Powered by Kubex Ecosystem',
    motto: 'Governance, Creativity, Productivity, Freedom',
    lang_en: 'English',
    lang_es: 'Español',
    lang_pt: 'Português',
    lang_zh: '中文',
  },
  es: {
    poweredBy: 'Impulsado por Kubex Ecosystem',
    motto: 'Gobernanza, Creatividad, Productividad, Libertad',
    lang_en: 'English',
    lang_es: 'Español',
    lang_pt: 'Português',
    lang_zh: '中文',
  },
  pt: {
    poweredBy: 'Desenvolvido pelo Kubex Ecosystem',
    motto: 'Governança, Criatividade, Produtividade, Liberdade',
    lang_en: 'English',
    lang_es: 'Español',
    lang_pt: 'Português',
    lang_zh: '中文',
  },
  zh: {
    poweredBy: '由 Kubex 生态系统提供支持',
    motto: '治理、创造力、生产力、自由',
    lang_en: 'English',
    lang_es: 'Español',
    lang_pt: 'Português',
    lang_zh: '中文',
  },
};

const App: React.FC = () => {
  // State management
  const [theme, setTheme] = useState<Theme>('dark');
  const [language, setLanguage] = useState<Language>('en');
  const [activeSection, setActiveSection] = useState<string>('welcome');
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
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
    if (savedLanguage && ['en', 'es', 'pt', 'zh'].includes(savedLanguage)) {
      setLanguage(savedLanguage);
    }
  }, []);

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

  // Translation function
  const t = (key: string, params?: Record<string, string>): string => {
    let translation = translations[language][key] || translations['en'][key] || key;
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
      label: 'Welcome',
      description: 'Getting started',
    },
    {
      id: 'prompt',
      label: 'Prompt Crafter',
      description: 'Build prompts',
    },
    {
      id: 'agents',
      label: 'Agents',
      description: 'AI Agents',
    },
    {
      id: 'chat',
      label: 'Chat',
      description: 'Conversations',
    },
    {
      id: 'summarizer',
      label: 'Summarizer',
      description: 'Content summary',
    },
    {
      id: 'code',
      label: 'Code',
      description: 'Code generation',
    },
    {
      id: 'images',
      label: 'Images',
      description: 'Image generation',
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
        return <AgentsGenerator theme={theme} />;
      case 'chat':
        return <ChatInterface theme={theme} />;
      case 'summarizer':
        return <ContentSummarizer theme={theme} />;
      case 'code':
        return <CodeGenerator theme={theme} />;
      case 'images':
        return <ImageGenerator theme={theme} />;
      default:
        return <Welcome onGetStarted={() => setActiveSection('prompt')} />;
    }
  };

  return (
    <LanguageContext.Provider value={{ language, setLanguage: handleSetLanguage, t }}>
      <Layout
        sidebar={
          <Sidebar
            sections={sections}
            activeSection={activeSection}
            onSectionChange={setActiveSection}
            onClose={() => setSidebarOpen(false)}
            collapsed={sidebarCollapsed}
          />
        }
        header={
          <Header
            theme={theme}
            onToggleTheme={handleToggleTheme}
            onToggleSidebar={() => setSidebarOpen(!sidebarOpen)}
            collapsed={sidebarCollapsed}
          />
        }
        footer={<Footer />}
        sidebarOpen={sidebarOpen}
        onSidebarClose={() => setSidebarOpen(false)}
        sidebarCollapsed={sidebarCollapsed}
      >
        {renderContent()}
      </Layout>
    </LanguageContext.Provider>
  );
};

export default App;
