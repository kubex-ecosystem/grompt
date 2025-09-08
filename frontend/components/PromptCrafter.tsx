import React, { useState, useCallback, useEffect, useRef } from 'react';
import { Idea, HistoryItem } from '../types';
import { generateStructuredPrompt } from '../services/geminiService';
import { Plus, Trash2, Wand2, Loader, Clipboard, ClipboardCheck, AlertTriangle, X, Lightbulb, History, Eye, XCircle } from 'lucide-react';

// --- IndexedDB Helpers for Autosave ---
const DB_NAME = 'PromptCrafterDB';
const STORE_NAME = 'drafts';
const DB_VERSION = 1;

interface Draft {
  ideas: Idea[];
  purpose: string;
}

const openDB = (): Promise<IDBDatabase> => {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, DB_VERSION);
    request.onerror = () => reject(new Error("Error opening IndexedDB. Your browser might be in private mode."));
    request.onsuccess = () => resolve(request.result);
    request.onupgradeneeded = (event) => {
      const db = (event.target as IDBOpenDBRequest).result;
      if (!db.objectStoreNames.contains(STORE_NAME)) {
        db.createObjectStore(STORE_NAME);
      }
    };
  });
};

// FIX: Added a trailing comma inside the generic type parameter <T,> to disambiguate from JSX syntax in a .tsx file.
const getDraft = async <T,>(key: IDBValidKey): Promise<T | undefined> => {
  const db = await openDB();
  return new Promise((resolve, reject) => {
    const transaction = db.transaction(STORE_NAME, 'readonly');
    const store = transaction.objectStore(STORE_NAME);
    const request = store.get(key);
    request.onerror = () => reject(new Error("Error getting draft from IndexedDB"));
    request.onsuccess = () => resolve(request.result as T | undefined);
  });
};

// FIX: Added a trailing comma inside the generic type parameter <T,> to disambiguate from JSX syntax in a .tsx file.
const setDraft = async <T,>(key: IDBValidKey, value: T): Promise<void> => {
  const db = await openDB();
  return new Promise((resolve, reject) => {
    const transaction = db.transaction(STORE_NAME, 'readwrite');
    const store = transaction.objectStore(STORE_NAME);
    const request = store.put(value, key);
    request.onerror = () => reject(new Error("Error setting draft in IndexedDB"));
    request.onsuccess = () => resolve();
  });
};


interface Example {
  purpose: string;
  ideas: string[];
}

const examples: Example[] = [
  {
    purpose: "Code Generation",
    ideas: [
      "Create a React hook for fetching data from an API.",
      "It should handle loading, error, and data states.",
      "Use the native `fetch` API.",
      "The hook should be written in TypeScript and be well-documented.",
    ],
  },
  {
    purpose: "Creative Writing",
    ideas: [
      "Write a short story opening.",
      "The setting is a neon-lit cyberpunk city in 2077.",
      "The main character is a grizzled detective who is part-cyborg.",
      "It's perpetually raining and the streets are reflective.",
    ],
  },
  {
    purpose: "Marketing Copy",
    ideas: [
      "Draft an email campaign for a new productivity app.",
      "The target audience is busy professionals and university students.",
      "Highlight features like AI-powered task scheduling, calendar sync, and focus mode.",
      "The tone should be encouraging, professional, and slightly urgent.",
    ],
  },
  {
    purpose: "Technical Documentation",
    ideas: [
        "Write the 'Getting Started' section for a new JavaScript library.",
        "The library is called 'ChronoWarp' and it simplifies date manipulation.",
        "Include a simple installation guide using npm.",
        "Provide a clear, concise code example for its primary use case."
    ]
  }
];

const PromptCrafter: React.FC = () => {
  const [ideas, setIdeas] = useState<Idea[]>([]);
  const [currentIdea, setCurrentIdea] = useState('');
  const [purpose, setPurpose] = useState('Code Generation');
  const [isLoading, setIsLoading] = useState(false);
  const [generatedPrompt, setGeneratedPrompt] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [isCopied, setIsCopied] = useState(false);
  const [promptHistory, setPromptHistory] = useState<HistoryItem[]>([]);
  const isInitialLoad = useRef(true);

  // Load history from localStorage
  useEffect(() => {
    try {
        const storedHistory = localStorage.getItem('promptHistory');
        if (storedHistory) {
            setPromptHistory(JSON.parse(storedHistory));
        }
    } catch (error) {
        console.error("Failed to load history from localStorage", error);
        setPromptHistory([]);
    }
  }, []);

  // Save history to localStorage
  useEffect(() => {
    // Avoid clearing history on initial empty state
    if (promptHistory.length === 0 && localStorage.getItem('promptHistory') === null) {
        return;
    }
    try {
        localStorage.setItem('promptHistory', JSON.stringify(promptHistory));
    } catch (error) {
        console.error("Failed to save history to localStorage", error);
    }
  }, [promptHistory]);
  
  // Autosave: Load draft from IndexedDB on mount
  useEffect(() => {
    const loadDraft = async () => {
      try {
        const savedDraft = await getDraft<Draft>('currentDraft');
        if (savedDraft) {
          setIdeas(savedDraft.ideas);
          setPurpose(savedDraft.purpose);
        }
      } catch (e) {
        console.error("Failed to load draft from IndexedDB", e);
      } finally {
        isInitialLoad.current = false;
      }
    };
    loadDraft();
  }, []);

  // Autosave: Save draft to IndexedDB on changes
  useEffect(() => {
    if (isInitialLoad.current) {
      return; // Don't save during initial load cycle
    }

    const handler = setTimeout(() => {
      try {
        setDraft('currentDraft', { ideas, purpose });
      } catch (e) {
        console.error("Failed to save draft to IndexedDB", e);
      }
    }, 500); // Debounce saving

    return () => {
      clearTimeout(handler);
    };
  }, [ideas, purpose]);


  const handleAddIdea = useCallback(() => {
    if (currentIdea.trim() !== '') {
      setIdeas(prev => [...prev, { id: Date.now().toString(), text: currentIdea.trim() }]);
      setCurrentIdea('');
    }
  }, [currentIdea]);

  const handleRemoveIdea = useCallback((id: string) => {
    setIdeas(prev => prev.filter(idea => idea.id !== id));
  }, []);
  
  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAddIdea();
    }
  };

  const handleGenerate = async () => {
    if (ideas.length === 0 || !purpose.trim()) {
      setError("Please add at least one idea and specify a purpose before generating.");
      return;
    }
    setError(null);
    setIsLoading(true);
    setGeneratedPrompt('');
    try {
      const prompt = await generateStructuredPrompt(ideas, purpose);
      setGeneratedPrompt(prompt);
      const newItem: HistoryItem = {
        id: Date.now().toString(),
        prompt,
        purpose,
        ideas: [...ideas],
        timestamp: Date.now(),
      };
      setPromptHistory(prev => [newItem, ...prev].slice(0, 20)); // Keep last 20 prompts
    } catch (err: any) {
      setError(err.message || 'An unknown error occurred.');
    } finally {
      setIsLoading(false);
    }
  };
  
  const handleCopy = () => {
    if(generatedPrompt) {
        navigator.clipboard.writeText(generatedPrompt);
        setIsCopied(true);
        setTimeout(() => setIsCopied(false), 2000);
    }
  };
  
  const handleLoadExample = useCallback(() => {
    setError(null);
    setGeneratedPrompt('');
    const randomIndex = Math.floor(Math.random() * examples.length);
    const example = examples[randomIndex];
    setPurpose(example.purpose);
    const exampleIdeas: Idea[] = example.ideas.map((text, index) => ({
      id: `example-${Date.now()}-${index}`,
      text,
    }));
    setIdeas(exampleIdeas);
  }, []);

  const handleLoadFromHistory = useCallback((item: HistoryItem) => {
    setIdeas(item.ideas);
    setPurpose(item.purpose);
    setGeneratedPrompt(item.prompt);
    setError(null);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }, []);

  const handleDeleteFromHistory = useCallback((id: string) => {
    setPromptHistory(prev => prev.filter(item => item.id !== id));
  }, []);

  const handleClearHistory = useCallback(() => {
    setPromptHistory([]);
    localStorage.removeItem('promptHistory');
  }, []);

  const formatRelativeTime = (timestamp: number): string => {
    const now = new Date();
    const seconds = Math.floor((now.getTime() - timestamp) / 1000);
    if (seconds < 60) return "just now";
    let interval = seconds / 31536000;
    if (interval > 1) return Math.floor(interval) + " years ago";
    interval = seconds / 2592000;
    if (interval > 1) return Math.floor(interval) + " months ago";
    interval = seconds / 86400;
    if (interval > 1) return Math.floor(interval) + " days ago";
    interval = seconds / 3600;
    if (interval > 1) return Math.floor(interval) + " hours ago";
    interval = seconds / 60;
    if (interval > 1) return Math.floor(interval) + " minutes ago";
    return Math.floor(seconds) + " seconds ago";
  };

  const purposes = [
    "Code Generation",
    "Creative Writing",
    "Data Analysis",
    "Technical Documentation",
    "Marketing Copy",
    "General Summarization",
  ];

  const IdeaInput = () => (
    <div className="flex gap-2">
      <input
        type="text"
        value={currentIdea}
        onChange={(e) => setCurrentIdea(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="Enter a raw idea, concept, or requirement..."
        className="flex-grow bg-white dark:bg-[#10151b] border-2 border-slate-300 dark:border-[#7c4dff]/50 rounded-md p-3 focus:outline-none focus:ring-2 focus:ring-indigo-500 dark:focus:ring-[#7c4dff] focus:border-indigo-500 dark:focus:border-[#7c4dff] transition-all duration-300 placeholder:text-slate-400 dark:placeholder:text-[#90a4ae]/50 text-slate-800 dark:text-white"
      />
      <button
        onClick={handleAddIdea}
        disabled={!currentIdea.trim()}
        className="bg-indigo-500 dark:bg-[#7c4dff] text-white p-3 rounded-md flex items-center justify-center hover:bg-indigo-600 dark:hover:bg-[#8e24aa] disabled:bg-slate-400 dark:disabled:bg-gray-600 disabled:cursor-not-allowed transition-all duration-300 shadow-lg shadow-indigo-500/30 dark:shadow-[0_0_10px_rgba(124,77,255,0.5)] hover:shadow-xl hover:shadow-indigo-500/40 dark:hover:shadow-[0_0_15px_rgba(142,36,170,0.7)]"
        aria-label="Add Idea"
      >
        <Plus size={24} />
      </button>
    </div>
  );

  const IdeasList = () => (
    <div className="space-y-3 mt-4 pr-2 max-h-60 overflow-y-auto">
      {ideas.map((idea, index) => (
        <div key={idea.id} className="bg-slate-100/50 dark:bg-[#10151b]/50 p-3 rounded-md flex justify-between items-center border border-transparent hover:border-sky-400/50 dark:hover:border-[#00f0ff]/30 transition-colors duration-300 animate-fade-in" style={{ animationDelay: `${index * 50}ms` }}>
          <span className="text-slate-700 dark:text-[#e0f7fa]">{idea.text}</span>
          <button
            onClick={() => handleRemoveIdea(idea.id)}
            className="text-red-500 dark:text-red-400 hover:text-red-600 dark:hover:text-red-300 p-1 rounded-full hover:bg-red-500/10 dark:hover:bg-red-500/20 transition-all duration-200"
            aria-label={`Remove idea: ${idea.text}`}
          >
            <Trash2 size={16} />
          </button>
        </div>
      ))}
    </div>
  );

  const PurposeSelector = () => (
    <div className="mt-6">
      <label htmlFor="purpose-input" className="block text-lg font-medium text-sky-500 dark:text-[#00f0ff] mb-2 light-shadow-sky dark:neon-glow-cyan">Purpose</label>
      <input
        id="purpose-input"
        type="text"
        value={purpose}
        onChange={(e) => setPurpose(e.target.value)}
        list="purposes-list"
        placeholder="e.g., Code Generation, Creative Writing..."
        className="w-full bg-white dark:bg-[#10151b] border-2 border-emerald-500/50 dark:border-[#00e676]/30 rounded-md p-3 focus:outline-none focus:ring-2 focus:ring-emerald-500 dark:focus:ring-[#00e676] focus:border-emerald-500 dark:focus:border-[#00e676] transition-all duration-300 placeholder:text-slate-400 dark:placeholder:text-[#90a4ae]/50 text-emerald-600 dark:text-[#00e676] font-semibold"
      />
      <datalist id="purposes-list">
        {purposes.map(p => <option key={p} value={p} />)}
      </datalist>
    </div>
  );

  const PromptHistory = () => (
    <div className="lg:col-span-2 bg-white/60 dark:bg-[#10151b]/30 p-6 rounded-lg border-2 border-slate-200 dark:border-sky-500/30 backdrop-blur-sm shadow-2xl shadow-slate-500/10">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-2xl font-bold font-orbitron text-sky-600 dark:text-sky-400 tracking-wider">3. PROMPT HISTORY</h3>
        {promptHistory.length > 0 && (
          <button
            onClick={handleClearHistory}
            className="flex items-center gap-2 px-3 py-1 text-sm bg-red-500/10 dark:bg-red-500/20 text-red-600 dark:text-red-300 rounded-full hover:bg-red-500/20 dark:hover:bg-red-500/30 transition-colors duration-200"
            aria-label="Clear all history"
          >
            <XCircle size={16} />
            <span>Clear All</span>
          </button>
        )}
      </div>
      <div className="max-h-96 overflow-y-auto pr-2 space-y-3">
        {promptHistory.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-40 text-slate-400 dark:text-[#90a4ae]/50">
            <History size={48} className="mb-4" />
            <p className="font-semibold text-center">Your generated prompts will be saved here.</p>
          </div>
        ) : (
          promptHistory.map(item => (
            <div key={item.id} className="bg-slate-100/50 dark:bg-[#10151b]/50 p-3 rounded-md flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3 border border-transparent hover:border-sky-400/50 dark:hover:border-sky-500/50 transition-colors duration-300">
              <div className="flex-grow overflow-hidden">
                <div className="flex items-center gap-2 mb-1">
                  <span className="px-2 py-0.5 text-xs font-bold text-white dark:text-black rounded-full bg-sky-500 dark:bg-sky-400 flex-shrink-0">{item.purpose}</span>
                  <span className="text-xs text-slate-500 dark:text-[#90a4ae] truncate">{formatRelativeTime(item.timestamp)}</span>
                </div>
                <p className="text-sm text-slate-700 dark:text-[#e0f7fa] line-clamp-2">
                  {item.prompt}
                </p>
              </div>
              <div className="flex items-center gap-2 self-end sm:self-center flex-shrink-0">
                <button 
                  onClick={() => handleLoadFromHistory(item)}
                  className="p-2 rounded-md bg-indigo-500/10 dark:bg-indigo-400/20 text-indigo-600 dark:text-indigo-300 hover:bg-indigo-500/20 dark:hover:bg-indigo-400/30 transition-colors duration-200"
                  aria-label="Load this prompt"
                >
                  <Eye size={18} />
                </button>
                <button 
                  onClick={() => handleDeleteFromHistory(item.id)}
                  className="p-2 rounded-md bg-red-500/10 dark:bg-red-400/20 text-red-600 dark:text-red-300 hover:bg-red-500/20 dark:hover:bg-red-400/30 transition-colors duration-200"
                  aria-label="Delete this prompt"
                >
                  <Trash2 size={18} />
                </button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );

  return (
    <div className="grid lg:grid-cols-2 gap-8">
      {/* Input Section */}
      <div className="bg-white/60 dark:bg-[#10151b]/30 p-6 rounded-lg border-2 border-slate-200 dark:border-[#7c4dff]/30 backdrop-blur-sm shadow-2xl shadow-slate-500/10">
        <div className="flex justify-between items-center mb-4">
            <h3 className="text-2xl font-bold font-orbitron text-indigo-600 dark:text-[#7c4dff] tracking-wider">1. INPUT IDEAS</h3>
            <button
                onClick={handleLoadExample}
                className="flex items-center gap-2 px-3 py-1 text-sm bg-slate-200/50 dark:bg-[#10151b] text-indigo-600 dark:text-[#7c4dff] rounded-full hover:bg-slate-300/50 dark:hover:bg-slate-800 transition-colors duration-200"
                aria-label="Load an example"
            >
                <Lightbulb size={16} />
                <span>Load Example</span>
            </button>
        </div>
        <IdeaInput />
        <IdeasList />
        <PurposeSelector />
      </div>

      {/* Output Section */}
      <div className="flex flex-col">
        <div className="bg-white/60 dark:bg-[#10151b]/30 p-6 rounded-lg border-2 border-slate-200 dark:border-[#00e676]/30 backdrop-blur-sm flex-grow flex flex-col shadow-2xl shadow-slate-500/10">
          <h3 className="text-2xl font-bold font-orbitron text-emerald-600 dark:text-[#00e676] tracking-wider mb-4 light-shadow-emerald dark:neon-glow-green">2. GENERATED PROMPT</h3>
          <div className="flex-grow bg-slate-100 dark:bg-[#0a0f14] rounded-md p-4 min-h-[300px] relative overflow-y-auto border border-slate-200 dark:border-transparent">
            {isLoading && (
              <div className="absolute inset-0 flex flex-col items-center justify-center bg-slate-100/80 dark:bg-[#0a0f14]/80 backdrop-blur-sm">
                 <Loader className="animate-spin text-sky-500 dark:text-[#00f0ff]" size={48} />
                 <p className="mt-4 text-sky-500 dark:text-[#00f0ff] font-semibold">Generating with Gemini...</p>
              </div>
            )}
            {error && (
              <div className="absolute inset-0 flex flex-col items-center justify-center bg-slate-100/80 dark:bg-[#0a0f14]/80 backdrop-blur-sm p-4">
                 <div className="bg-red-100 dark:bg-red-900/50 border border-red-300 dark:border-red-500 text-red-700 dark:text-red-300 p-4 rounded-lg flex flex-col items-center gap-2 text-center shadow-lg">
                    <AlertTriangle size={32} />
                    <h4 className="font-bold">Generation Failed</h4>
                    <p className="text-sm">{error}</p>
                    <button onClick={() => setError(null)} className="mt-2 bg-red-500/20 dark:bg-red-500/50 hover:bg-red-500/30 dark:hover:bg-red-500/80 text-red-800 dark:text-white px-3 py-1 text-sm rounded-md flex items-center gap-1">
                        <X size={14}/> Close
                    </button>
                 </div>
              </div>
            )}
            {generatedPrompt && !isLoading && !error && (
               <>
                <button onClick={handleCopy} className="absolute top-2 right-2 bg-white dark:bg-[#10151b] p-2 rounded-md text-sky-500 dark:text-[#00f0ff] hover:bg-sky-100 dark:hover:bg-[#00f0ff] hover:text-sky-600 dark:hover:text-black transition-all duration-200">
                    {isCopied ? <ClipboardCheck size={20} /> : <Clipboard size={20} />}
                </button>
                <pre className="whitespace-pre-wrap font-plex-mono text-sm leading-relaxed text-slate-800 dark:text-[#e0f7fa]">
                    {generatedPrompt}
                </pre>
               </>
            )}
            {!generatedPrompt && !isLoading && !error && (
                <div className="flex flex-col items-center justify-center h-full text-slate-400 dark:text-[#90a4ae]/50">
                    <Wand2 size={48} className="mb-4"/>
                    <p className="font-semibold text-center">Your professional prompt will appear here.</p>
                </div>
            )}
          </div>
          <button 
            onClick={handleGenerate} 
            disabled={isLoading || ideas.length === 0 || !purpose.trim()}
            className="w-full mt-6 bg-gradient-to-r from-emerald-500 to-sky-500 dark:from-[#00e676] dark:to-[#00f0ff] text-white dark:text-black font-bold font-orbitron text-lg p-4 rounded-lg flex items-center justify-center gap-3 hover:scale-105 disabled:scale-100 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300 shadow-lg shadow-sky-500/40 dark:shadow-[0_0_15px_rgba(0,230,118,0.5)] hover:shadow-xl hover:shadow-sky-500/50 dark:hover:shadow-[0_0_25px_rgba(0,240,255,0.7)]"
          >
            {isLoading ? <Loader className="animate-spin" size={28} /> : <Wand2 size={28} />}
            {isLoading ? 'GENERATING...' : 'GENERATE PROMPT'}
          </button>
        </div>
      </div>
      
      {/* Prompt History Section */}
      <PromptHistory />
    </div>
  );
};

export default PromptCrafter;