
import { BrainCircuit, Check, Clipboard, Copy, Download, History, Lightbulb, Plus, RefreshCw, Search, Trash2, Wand2, X } from 'lucide-react';
import React, { useEffect, useMemo, useRef, useState } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Loader } from './components/ui/Loader';
import { PURPOSES } from './constants';
import { getFullSourceCode } from './services/codebaseService';
import { dbService } from './services/dbService';
import { geminiService } from './services/geminiService';
import type { PromptHistoryItem } from './types';

const App: React.FC = () => {
  const [currentIdea, setCurrentIdea] = useState('');
  const [ideas, setIdeas] = useState<string[]>([]);
  const [purpose, setPurpose] = useState(PURPOSES[0]);

  const [generatedPrompt, setGeneratedPrompt] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const [history, setHistory] = useState<PromptHistoryItem[]>([]);
  const [historySearch, setHistorySearch] = useState('');

  // State for the Evolve Modal
  const [isEvolveModalOpen, setIsEvolveModalOpen] = useState(false);
  const [evolveStep, setEvolveStep] = useState<'idle' | 'generating_prompt' | 'refactoring' | 'done' | 'error'>('idle');
  const [refactoredCode, setRefactoredCode] = useState('');
  const [selfImprovementPrompt, setSelfImprovementPrompt] = useState('');
  const [evolveError, setEvolveError] = useState('');

  const ideaInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    loadHistory();
  }, []);

  const loadHistory = async () => {
    const items = await dbService.getPrompts();
    setHistory(items);
  };

  const handleAddIdea = () => {
    if (currentIdea.trim()) {
      setIdeas(prev => [...prev, currentIdea.trim()]);
      setCurrentIdea('');
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAddIdea();
    }
  };

  const handleRemoveIdea = (indexToRemove: number) => {
    setIdeas(ideas.filter((_, index) => index !== indexToRemove));
  };

  const handleGeneratePrompt = async () => {
    if (ideas.length === 0 || isLoading) return;

    setIsLoading(true);
    setError(null);
    setGeneratedPrompt('');
    setCopied(false);

    try {
      const result = await geminiService.craftPrompt({ ideas, purpose });
      setGeneratedPrompt(result);
      await dbService.addPrompt({
        prompt: result,
        timestamp: Date.now(),
        inputs: { ideas, purpose }
      });
      await loadHistory();
    } catch (e) {
      setError('Failed to generate prompt. Please try again.');
      console.error(e);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCopy = () => {
    if (!generatedPrompt) return;
    navigator.clipboard.writeText(generatedPrompt);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleCopyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    // You could add a temporary "Copied!" message here if desired
  };

  const handleExport = () => {
    if (!generatedPrompt) return;
    const blob = new Blob([generatedPrompt], { type: 'text/plain;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'grompt_prompt.txt';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const handleDeleteHistory = async (id: number) => {
    await dbService.deletePrompt(id);
    await loadHistory();
  };

  const handleRestoreHistory = (item: PromptHistoryItem) => {
    setIdeas(item.inputs.ideas);
    setPurpose(item.inputs.purpose);
    setGeneratedPrompt(item.prompt);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleLoadExample = () => {
    setIdeas([
      "The app needs a user profile page.",
      "It should display the user's name, email, and a profile picture.",
      "There must be an 'Edit' button.",
      "Include a section for 'Recent Activity'."
    ]);
    setPurpose("Technical Documentation");
    ideaInputRef.current?.focus();
  };

  const handleEvolve = async () => {
    setIsEvolveModalOpen(true);
    setEvolveStep('generating_prompt');
    setEvolveError('');
    setRefactoredCode('');
    setSelfImprovementPrompt('');

    try {
      // Step 1: Craft the self-improvement prompt
      const selfImprovementIdeas = [
        "Refactor the entire Grompt application codebase provided below.",
        "Improve code clarity, performance, and maintainability.",
        "Follow modern React (Hooks) and TypeScript best practices.",
        "Ensure components are well-structured, reusable, and self-contained.",
        "The application uses Tailwind CSS for styling, so maintain that approach.",
        "Do not change the external API of the services, but improve their internal implementation.",
        "Return ONLY the complete, refactored code for all files, formatted within a single markdown code block.",
        "Preserve the file structure and denote each file clearly (e.g., using comments like '// --- START OF FILE: App.tsx ---')."
      ];
      const prompt = await geminiService.craftPrompt({
        ideas: selfImprovementIdeas,
        purpose: 'Code Generation',
      });
      setSelfImprovementPrompt(prompt);

      // Step 2: Get the full source code
      const sourceCode = getFullSourceCode();

      // Step 3: Send it for refactoring
      setEvolveStep('refactoring');
      const result = await geminiService.refactorCode({
        systemPrompt: prompt,
        code: sourceCode
      });
      setRefactoredCode(result);
      setEvolveStep('done');

    } catch (e) {
      console.error("Evolution process failed:", e);
      setEvolveError("The evolution process failed. The AI may be overloaded. Please try again later.");
      setEvolveStep('error');
    }
  };

  const closeEvolveModal = () => {
    setIsEvolveModalOpen(false);
    setEvolveStep('idle');
  }

  const filteredHistory = useMemo(() => {
    return history.filter(item =>
      item.prompt.toLowerCase().includes(historySearch.toLowerCase()) ||
      item.inputs.purpose.toLowerCase().includes(historySearch.toLowerCase())
    );
  }, [history, historySearch]);

  return (
    <>
      <div className="min-h-screen flex flex-col items-center p-4 sm:p-6 lg:p-8 space-y-8 text-foreground">
        <header className="w-full max-w-7xl mx-auto text-center relative">
          <h1 className="text-4xl md:text-5xl font-display font-bold text-primary tracking-wider">Grompt</h1>
          <p className="mt-2 text-muted-foreground font-sans">AI Prompt Crafter</p>
          <button
            onClick={handleEvolve}
            className="absolute top-0 right-0 p-2 text-muted-foreground hover:text-primary transition-colors"
            title="Evolve Code (Virtuous Meta-Recursion)"
          >
            <BrainCircuit className="w-8 h-8" />
          </button>
        </header>

        <main className="w-full max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Left Column: Inputs */}
          <section className="bg-card border border-border rounded-lg p-6 backdrop-blur-sm flex flex-col gap-6">
            <div className="flex justify-between items-center">
              <h2 className="text-xl font-display text-primary">1. INPUT IDEAS</h2>
              <button onClick={handleLoadExample} className="flex items-center gap-2 text-sm text-muted-foreground hover:text-primary transition-colors">
                <Lightbulb className="w-4 h-4" /> Load Example
              </button>
            </div>
            <div className="flex gap-2">
              <input
                ref={ideaInputRef}
                type="text"
                value={currentIdea}
                onChange={(e) => setCurrentIdea(e.target.value)}
                onKeyDown={handleKeyPress}
                placeholder="Enter a raw idea, concept, or requirement..."
                className="w-full bg-input border border-border rounded-md p-3 focus:ring-2 focus:ring-primary focus:outline-none"
              />
              <button title='Add Idea' onClick={handleAddIdea} disabled={!currentIdea.trim()} className="p-3 bg-secondary hover:bg-secondary-hover rounded-md disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
                <Plus className="w-5 h-5 text-secondary-foreground" />
              </button>
            </div>
            <div className="space-y-2 flex-grow overflow-y-auto max-h-48 pr-2">
              {ideas.map((idea, index) => (
                <div key={index} className="flex items-center justify-between bg-input p-2 rounded-md border border-transparent hover:border-border transition-colors">
                  <span className="text-sm flex-1">{idea}</span>
                  <button onClick={() => handleRemoveIdea(index)} className="p-1 text-muted-foreground hover:text-red-500" aria-label={`Remove idea: ${idea}`}>
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              ))}
            </div>

            <div className="space-y-3">
              <h3 className="font-display text-primary">Purpose</h3>
              <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
                {PURPOSES.map(p => (
                  <button key={p} onClick={() => setPurpose(p)} className={`p-2 text-sm rounded-md transition-all border ${purpose === p ? 'bg-accent text-accent-foreground font-bold border-accent' : 'bg-input hover:bg-gray-800 border-border'}`}>
                    {p}
                  </button>
                ))}
              </div>
              <div className="p-2 bg-input border border-border rounded-md text-sm">
                <span className="font-bold text-accent">{purpose}</span>
              </div>
            </div>
          </section>

          {/* Right Column: Output */}
          <section className="bg-card border border-border rounded-lg p-6 backdrop-blur-sm flex flex-col">
            <h2 className="text-xl font-display text-primary mb-4">2. GENERATED PROMPT</h2>
            <div className="flex-grow w-full bg-input rounded-lg flex p-4 border border-border overflow-y-auto min-h-[300px]">
              {isLoading ? (
                <div className="m-auto text-center space-y-2 text-muted-foreground">
                  <Loader size="md" />
                  <p>Crafting your prompt...</p>
                </div>
              ) : error ? (
                <p className="m-auto text-red-400">{error}</p>
              ) : generatedPrompt ? (
                // Fix: The className prop is not supported on ReactMarkdown. Wrap it in a div and apply classes there.
                <div className="prose prose-sm prose-invert max-w-none w-full">
                  <ReactMarkdown remarkPlugins={[remarkGfm]}>{generatedPrompt}</ReactMarkdown>
                </div>
              ) : (
                <div className="m-auto text-center text-muted-foreground">
                  <Wand2 className="w-10 h-10 mx-auto mb-2" />
                  <p>Your professional prompt will appear here.</p>
                </div>
              )}
            </div>
            <div className="flex items-center gap-2 mt-4">
              <button onClick={handleGeneratePrompt} disabled={isLoading || ideas.length === 0} className="w-full py-3 px-4 rounded-lg bg-accent hover:bg-accent-hover text-accent-foreground text-base font-bold flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed transition-shadow shadow-md hover:shadow-neon-green">
                <Wand2 className="w-5 h-5" />
                <span>GENERATE PROMPT</span>
              </button>
              <button onClick={handleCopy} disabled={!generatedPrompt} className="p-3 bg-secondary hover:bg-secondary-hover rounded-md text-secondary-foreground disabled:opacity-50"><span className="sr-only">Copy</span>{copied ? <Check className="w-5 h-5 text-green-400" /> : <Clipboard className="w-5 h-5" />}</button>
              <button onClick={handleExport} disabled={!generatedPrompt} className="p-3 bg-secondary hover:bg-secondary-hover rounded-md text-secondary-foreground disabled:opacity-50"><span className="sr-only">Export</span><Download className="w-5 h-5" /></button>
            </div>
          </section>
        </main>

        {/* History Section */}
        <section className="w-full max-w-7xl mx-auto bg-card border border-border rounded-lg p-6 backdrop-blur-sm">
          <h2 className="text-xl font-display text-primary mb-4">3. PROMPT HISTORY</h2>
          <div className="relative mb-4">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground" />
            <input
              type="text"
              value={historySearch}
              onChange={(e) => setHistorySearch(e.target.value)}
              placeholder="Search history..."
              className="w-full bg-input border border-border rounded-md p-3 pl-10 focus:ring-2 focus:ring-primary focus:outline-none"
            />
          </div>
          <div className="space-y-2 max-h-64 overflow-y-auto pr-2">
            {filteredHistory.length > 0 ? filteredHistory.map(item => (
              <div key={item.id} className="bg-input p-3 rounded-md flex justify-between items-start gap-4 border border-border/50">
                <div className="flex-1">
                  <p className="text-sm text-foreground truncate">{item.prompt}</p>
                  <div className="text-xs text-muted-foreground mt-1 flex items-center gap-4">
                    <span>{new Date(item.timestamp).toLocaleString()}</span>
                    <span className="px-2 py-0.5 bg-secondary text-secondary-foreground rounded-full">{item.inputs.purpose}</span>
                  </div>
                </div>
                <div className="flex gap-2">
                  <button onClick={() => handleRestoreHistory(item)} title="Restore" className="p-1.5 text-muted-foreground hover:text-primary"><RefreshCw className="w-4 h-4" /></button>
                  <button onClick={() => handleDeleteHistory(item.id)} title="Delete" className="p-1.5 text-muted-foreground hover:text-red-500"><Trash2 className="w-4 h-4" /></button>
                </div>
              </div>
            )) : (
              <div className="text-center py-8 text-muted-foreground">
                <History className="w-8 h-8 mx-auto mb-2" />
                <p>Your generated prompts will be saved here.</p>
              </div>
            )}
          </div>
        </section>
      </div>

      {/* Evolve Modal */}
      {isEvolveModalOpen && (
        <div className="fixed inset-0 bg-background/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-card border border-border rounded-lg shadow-neon-cyan w-full max-w-4xl h-[90vh] flex flex-col">
            <header className="flex items-center justify-between p-4 border-b border-border">
              <h2 className="text-xl font-display text-primary flex items-center gap-2">
                <BrainCircuit className="w-6 h-6" />
                Virtuous Meta-Recursion
              </h2>
              <button title='Close' onClick={closeEvolveModal} className="p-1 text-muted-foreground hover:text-primary"><X className="w-5 h-5" /></button>
            </header>
            <main className="p-6 flex-1 overflow-y-auto">
              {evolveStep === 'generating_prompt' && (
                <div className="flex flex-col items-center justify-center h-full text-center">
                  <Loader size="lg" />
                  <p className="mt-4 text-lg text-muted-foreground">Step 1/2: Generating self-improvement prompt...</p>
                </div>
              )}
              {evolveStep === 'refactoring' && (
                <div className="flex flex-col items-center justify-center h-full text-center">
                  <Loader size="lg" />
                  <p className="mt-4 text-lg text-muted-foreground">Step 2/2: Code evolution in progress...</p>
                  <p className="text-sm text-muted-foreground/70">Sending full source code to Gemini for refactoring. This may take a moment.</p>
                </div>
              )}
              {evolveStep === 'error' && (
                <div className="flex flex-col items-center justify-center h-full text-center text-red-400">
                  <X className="w-12 h-12 mb-4" />
                  <h3 className="text-xl">Evolution Failed</h3>
                  <p>{evolveError}</p>
                </div>
              )}
              {evolveStep === 'done' && (
                <div className="h-full flex flex-col">
                  <h3 className="text-lg font-bold text-accent mb-2">Evolution Complete!</h3>
                  <p className="text-sm text-muted-foreground mb-4">The AI has returned the refactored code. Review the changes below.</p>
                  <div className="relative flex-1">
                    <pre className="bg-input p-4 rounded-md h-full overflow-auto text-sm w-full border border-border">
                      <code>{refactoredCode}</code>
                    </pre>
                    <button onClick={() => handleCopyToClipboard(refactoredCode)} className="absolute top-2 right-2 p-2 bg-secondary hover:bg-secondary-hover rounded-md text-secondary-foreground" title="Copy Code">
                      <Copy className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              )}
            </main>
          </div>
        </div>
      )}
    </>
  );
};

export default App;
