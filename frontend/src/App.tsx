import { AlertTriangle } from 'lucide-react';
import { useEffect, useState } from 'react';
import Footer from './components/layout/Footer';
import Header from './components/layout/Header';
import PromptCrafter from './components/prompt-crafter/PromptCrafter';
import { LanguageContext } from './context/LanguageContext';
import { useLanguage } from './hooks/useLanguage';
import { useTheme } from './hooks/useTheme';
import { useAnalytics } from './services/analytics';
import { initStorage } from './services/storageService';

const App: React.FC = () => {
  const [theme, toggleTheme] = useTheme();
  const { language, setLanguage, t } = useLanguage();
  const [isApiKeyMissing, setIsApiKeyMissing] = useState(false);

  // Initialize analytics
  useAnalytics();

  // Initialize the storage service on app load
  useEffect(() => {
    initStorage();
    // Check for the API key availability on mount
    if (!process.env.API_KEY) {
      console.log("Running in demo mode - API key not configured");
      setIsApiKeyMissing(true);
    }
  }, []);

  return (
    <LanguageContext.Provider value={{ language, setLanguage, t }}>
      <div className="min-h-screen text-slate-800 dark:text-[#e0f7fa] font-plex-mono p-4 sm:p-6 lg:p-8">
        <div className="max-w-7xl mx-auto">
          <Header theme={theme} toggleTheme={toggleTheme} />
          {isApiKeyMissing && (
            <div className="bg-blue-100 border-l-4 border-blue-500 text-blue-700 p-4 rounded-md mb-6 dark:bg-blue-900/30 dark:text-blue-300 dark:border-blue-600 flex items-start gap-3" role="alert">
              <AlertTriangle className="h-6 w-6 text-blue-500 dark:text-blue-400 flex-shrink-0 mt-0.5" />
              <div>
                <p className="font-bold">{t('apiKeyMissingTitle')}</p>
                <p>{t('apiKeyMissingText')}</p>
              </div>
            </div>
          )}
          <main>
            <PromptCrafter theme={theme} isApiKeyMissing={isApiKeyMissing} />
          </main>
          <Footer />
        </div>
      </div>
    </LanguageContext.Provider>
  );
};

export default App;
