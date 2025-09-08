import { Lightbulb, Loader, Wand2 } from 'lucide-react';
import React, { useCallback, useState } from 'react';
import { examples } from '../../constants/prompts';
import { useTranslations } from '../../i18n/useTranslations';
import { generateStructuredPrompt } from '../../services/geminiService';
import { HistoryItem, Idea, Theme } from '../../types';

// Import hooks
import { useAutosaveDraft } from '../../hooks/useAutosaveDraft';
import { usePromptHistory } from '../../hooks/usePromptHistory';
import { useUrlSharing } from '../../hooks/useUrlSharing';

// Import sub-components
import ApiKeyInput from './ApiKeyInput';
import GeneratedPrompt from './GeneratedPrompt';
import IdeaInput from './IdeaInput';
import IdeasList from './IdeasList';
import PromptHistoryDisplay from './PromptHistoryDisplay';
import PurposeSelector from './PurposeSelector';

interface PromptCrafterProps {
  theme: Theme;
  isApiKeyMissing: boolean;
}

const PromptCrafter: React.FC<PromptCrafterProps> = ({ theme, isApiKeyMissing }) => {
  const { t } = useTranslations();

  // State managed by hooks
  const { ideas, setIdeas, purpose, setPurpose } = useAutosaveDraft();
  const { promptHistory, setPromptHistory, deleteFromHistory, clearHistory } = usePromptHistory();

  // Local component state
  const [currentIdea, setCurrentIdea] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [generatedPrompt, setGeneratedPrompt] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [tokenUsage, setTokenUsage] = useState<{ input: number; output: number; total: number; } | null>(null);
  const [isExampleLoaded, setIsExampleLoaded] = useState(false);
  const [userApiKey, setUserApiKey] = useState<string>('');

  // Handler for API key changes
  const handleApiKeyChange = (apiKey: string) => {
    setUserApiKey(apiKey);
    // Force re-render to update UI state when API key is added/removed
    if (apiKey && isApiKeyMissing) {
      // We can now use the API instead of demo mode
      setError(null);
    }
  };

  // Hook for loading shared URL
  useUrlSharing({ setIdeas, setPurpose, setGeneratedPrompt, setTokenUsage, setError });

  // Create wrapped state setters to reset the example flag on user interaction
  const handleSetCurrentIdea = (value: string) => {
    // If an example is loaded, any typing signifies a modification.
    // Reset the flag to re-enable the generate button.
    if (isExampleLoaded) {
      setIsExampleLoaded(false);
    }
    setCurrentIdea(value);
  };

  const handleSetPurpose = (value: string) => {
    // If an example is loaded, changing the purpose signifies a modification.
    // Reset the flag to re-enable the generate button.
    if (isExampleLoaded) {
      setIsExampleLoaded(false);
    }
    setPurpose(value);
  };

  const handleAddIdea = useCallback(() => {
    if (currentIdea.trim() !== '') {
      setIdeas(prev => [...prev, { id: Date.now().toString(), text: currentIdea.trim() }]);
      setCurrentIdea('');
      // If an example is loaded, adding a new idea signifies a modification.
      // Reset the flag to re-enable the generate button.
      if (isExampleLoaded) {
        setIsExampleLoaded(false);
      }
    }
  }, [currentIdea, setIdeas, isExampleLoaded]);

  const handleRemoveIdea = useCallback((id: string) => {
    setIdeas(prev => prev.filter(idea => idea.id !== id));
    // If an example is loaded, removing an idea signifies a modification.
    // Reset the flag to re-enable the generate button.
    if (isExampleLoaded) {
      setIsExampleLoaded(false);
    }
  }, [setIdeas, isExampleLoaded]);

  const handleGenerate = async () => {
    if (ideas.length === 0) {
      setError(t('errorAddIdea'));
      return;
    }
    if (!purpose.trim()) {
      setError(t('errorSpecifyPurpose'));
      return;
    }
    setError(null);
    setIsLoading(true);
    setGeneratedPrompt('');
    setTokenUsage(null);
    try {
      const result = await generateStructuredPrompt(ideas, purpose);
      setGeneratedPrompt(result.prompt);
      const inputTokens = result.usageMetadata?.promptTokenCount ?? 0;
      const outputTokens = result.usageMetadata?.candidatesTokenCount ?? 0;
      const totalTokens = result.usageMetadata?.totalTokenCount ?? (inputTokens + outputTokens);
      setTokenUsage({ input: inputTokens, output: outputTokens, total: totalTokens });

      const newItem: HistoryItem = {
        id: Date.now().toString(),
        prompt: result.prompt,
        purpose,
        ideas: [...ideas],
        timestamp: Date.now(),
        inputTokens,
        outputTokens,
        totalTokens,
      };
      setPromptHistory(prev => [newItem, ...prev].slice(0, 20));
    } catch (err: any) {
      setError(err.message || 'An unknown error occurred.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleLoadExample = useCallback(() => {
    setError(null);
    setGeneratedPrompt('');
    setTokenUsage(null);
    const randomIndex = Math.floor(Math.random() * examples.length);
    const example = examples[randomIndex];
    const exampleIdeas: Idea[] = example.ideas.map((text, index) => ({
      id: `example-${Date.now()}-${index}`,
      text,
    }));
    setPurpose(example.purpose);
    setIdeas(exampleIdeas);
    setIsExampleLoaded(true);
  }, [setIdeas, setPurpose]);

  const handleLoadFromHistory = useCallback((item: HistoryItem) => {
    setIdeas(item.ideas);
    setPurpose(item.purpose);
    setGeneratedPrompt(item.prompt);
    setTokenUsage(item.inputTokens != null && item.outputTokens != null
      ? {
        input: item.inputTokens,
        output: item.outputTokens,
        total: item.totalTokens ?? (item.inputTokens + item.outputTokens)
      }
      : null);
    setError(null);
    setIsExampleLoaded(false); // Reset flag when loading from history
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }, [setIdeas, setPurpose]);

  return (
    <div className="grid lg:grid-cols-2 gap-8">
      {/* API Key Input Section - only show when no API key is configured */}
      {isApiKeyMissing && (
        <div className="lg:col-span-2">
          <ApiKeyInput
            onApiKeyChange={handleApiKeyChange}
            isVisible={isApiKeyMissing}
          />
        </div>
      )}
      {/* Input Section */}
      <div className="bg-white/60 dark:bg-[#10151b]/30 p-6 rounded-lg border-2 border-slate-200 dark:border-[#7c4dff]/30 backdrop-blur-sm shadow-2xl shadow-slate-500/10">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-2xl font-bold font-orbitron text-indigo-600 dark:text-[#7c4dff] tracking-wider">{t('inputIdeasTitle')}</h3>
          <button
            onClick={handleLoadExample}
            className="flex items-center gap-2 px-3 py-1 text-sm bg-slate-200/50 dark:bg-[#10151b] text-indigo-600 dark:text-[#7c4dff] rounded-full hover:bg-slate-300/50 dark:hover:bg-slate-800 transition-colors duration-200"
            aria-label={t('loadExample')}
          >
            <Lightbulb size={16} />
            <span>{t('loadExample')}</span>
          </button>
        </div>
        <IdeaInput
          currentIdea={currentIdea}
          setCurrentIdea={handleSetCurrentIdea}
          onAddIdea={handleAddIdea}
          disabled={isExampleLoaded}
        />
        <IdeasList ideas={ideas} onRemoveIdea={handleRemoveIdea} />
        <PurposeSelector purpose={purpose} setPurpose={handleSetPurpose} />
      </div>

      {/* Output Section */}
      <div className="flex flex-col">
        <GeneratedPrompt
          isLoading={isLoading}
          error={error}
          setError={setError}
          generatedPrompt={generatedPrompt}
          ideas={ideas}
          purpose={purpose}
          tokenUsage={tokenUsage}
          theme={theme}
        />
        <button
          onClick={handleGenerate}
          disabled={isLoading || ideas.length === 0 || !purpose.trim() || isExampleLoaded}
          className="w-full mt-6 bg-gradient-to-r from-emerald-500 to-sky-500 dark:from-[#00e676] dark:to-[#00f0ff] text-white dark:text-black font-bold font-orbitron text-lg p-4 rounded-lg flex items-center justify-center gap-3 hover:scale-105 disabled:scale-100 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300 shadow-lg shadow-sky-500/40 dark:shadow-[0_0_15px_rgba(0,230,118,0.5)] hover:shadow-xl hover:shadow-sky-500/50 dark:hover:shadow-[0_0_25px_rgba(0,240,255,0.7)]"
        >
          {isLoading ? <Loader className="animate-spin" size={28} /> : <Wand2 size={28} />}
          {isLoading ? t('generatingButton') : t('generateButton')}
        </button>
      </div>

      {/* Prompt History Section */}
      <PromptHistoryDisplay
        history={promptHistory}
        onLoad={handleLoadFromHistory}
        onDelete={deleteFromHistory}
        onClear={clearHistory}
      />
    </div>
  );
};

export default PromptCrafter;
