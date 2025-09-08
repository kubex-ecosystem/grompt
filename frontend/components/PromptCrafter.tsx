import { AlertTriangle, BrainCircuit, Clipboard, ClipboardCheck, Code, Eye, History, Lightbulb, Loader, Plus, Share2, Trash2, Wand2, X, XCircle } from 'lucide-react';
import React, { useCallback, useContext, useEffect, useRef, useState } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { LanguageContext } from '../App';
import { generateStructuredPrompt } from '../services/geminiService';
import { HistoryItem, Idea } from '../types';

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
    purpose: "Data Analysis",
    ideas: [
      "Analyze a dataset of customer sales from the last quarter.",
      "The dataset includes columns: 'Date', 'CustomerID', 'ProductCategory', 'Revenue', 'UnitsSold'.",
      "Identify the top 3 product categories by total revenue.",
      "Calculate the average revenue per customer.",
      "Look for any weekly sales trends or seasonality.",
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

const purposeKeys: Record<string, string> = {
  "Code Generation": "purposeCodeGeneration",
  "Creative Writing": "purposeCreativeWriting",
  "Data Analysis": "purposeDataAnalysis",
  "Technical Documentation": "purposeTechnicalDocumentation",
  "Marketing Copy": "purposeMarketingCopy",
  "General Summarization": "purposeGeneralSummarization",
};

const i18n: Record<string, Record<string, string>> = {
  en: {
    "inputIdeasTitle": "1. INPUT IDEAS",
    "loadExample": "Load Example",
    "ideaPlaceholder": "Enter a raw idea, concept, or requirement...",
    "addIdea": "Add Idea",
    "removeIdea": "Remove idea: {idea}",
    "purposeLabel": "Purpose",
    "purposeCodeGeneration": "Code Generation",
    "purposeCreativeWriting": "Creative Writing",
    "purposeDataAnalysis": "Data Analysis",
    "purposeTechnicalDocumentation": "Technical Documentation",
    "purposeMarketingCopy": "Marketing Copy",
    "purposeGeneralSummarization": "General Summarization",
    "customPurposePlaceholder": "Or type a custom purpose...",
    "generatedPromptTitle": "2. GENERATED PROMPT",
    "generatingMessage": "Generating with Gemini...",
    "generationFailedTitle": "Generation Failed",
    "close": "Close",
    "toggleView": "Toggle view mode",
    "copyLink": "Copy shareable link",
    "copyPrompt": "Copy prompt text",
    "promptPlaceholder": "Your professional prompt will appear here.",
    "generateButton": "GENERATE PROMPT",
    "generatingButton": "GENERATING...",
    "historyTitle": "3. PROMPT HISTORY",
    "clearAll": "Clear All",
    "historyPlaceholder": "Your generated prompts will be saved here.",
    "loadPrompt": "Load this prompt",
    "deletePrompt": "Delete this prompt",
    "errorAddIdea": "Please add at least one idea before generating.",
    "errorSpecifyPurpose": "Please specify a purpose for the prompt before generating.",
    "tokens": "tokens",
    "input": "Input",
    "output": "Output",
    "total": "Total",
  },
  es: {
    "inputIdeasTitle": "1. INGRESAR IDEAS",
    "loadExample": "Cargar Ejemplo",
    "ideaPlaceholder": "Introduce una idea, concepto o requisito...",
    "addIdea": "Añadir Idea",
    "removeIdea": "Quitar idea: {idea}",
    "purposeLabel": "Propósito",
    "purposeCodeGeneration": "Generación de Código",
    "purposeCreativeWriting": "Escritura Creativa",
    "purposeDataAnalysis": "Análisis de Datos",
    "purposeTechnicalDocumentation": "Documentación Técnica",
    "purposeMarketingCopy": "Copy de Marketing",
    "purposeGeneralSummarization": "Resumen General",
    "customPurposePlaceholder": "O escribe un propósito personalizado...",
    "generatedPromptTitle": "2. PROMPT GENERADO",
    "generatingMessage": "Generando con Gemini...",
    "generationFailedTitle": "Falló la Generación",
    "close": "Cerrar",
    "toggleView": "Cambiar vista",
    "copyLink": "Copiar enlace para compartir",
    "copyPrompt": "Copiar texto del prompt",
    "promptPlaceholder": "Tu prompt profesional aparecerá aquí.",
    "generateButton": "GENERAR PROMPT",
    "generatingButton": "GENERANDO...",
    "historyTitle": "3. HISTORIAL DE PROMPTS",
    "clearAll": "Limpiar Todo",
    "historyPlaceholder": "Tus prompts generados se guardarán aquí.",
    "loadPrompt": "Cargar este prompt",
    "deletePrompt": "Eliminar este prompt",
    "errorAddIdea": "Por favor, añade al menos una idea antes de generar.",
    "errorSpecifyPurpose": "Por favor, especifica un propósito para el prompt antes de generar.",
    "tokens": "tokens",
    "input": "Entrada",
    "output": "Salida",
    "total": "Total",
  },
  zh: {
    "inputIdeasTitle": "1. 输入想法",
    "loadExample": "加载示例",
    "ideaPlaceholder": "输入一个原始想法、概念或要求...",
    "addIdea": "添加想法",
    "removeIdea": "删除想法: {idea}",
    "purposeLabel": "目的",
    "purposeCodeGeneration": "代码生成",
    "purposeCreativeWriting": "创意写作",
    "purposeDataAnalysis": "数据分析",
    "purposeTechnicalDocumentation": "技术文档",
    "purposeMarketingCopy": "营销文案",
    "purposeGeneralSummarization": "通用总结",
    "customPurposePlaceholder": "或输入自定义目的...",
    "generatedPromptTitle": "2. 生成的提示",
    "generatingMessage": "正在通过 Gemini 生成...",
    "generationFailedTitle": "生成失败",
    "close": "关闭",
    "toggleView": "切换视图",
    "copyLink": "复制分享链接",
    "copyPrompt": "复制提示文本",
    "promptPlaceholder": "您的专业提示将出现在这里。",
    "generateButton": "生成提示",
    "generatingButton": "正在生成...",
    "historyTitle": "3. 提示历史",
    "clearAll": "全部清除",
    "historyPlaceholder": "您生成的提示将保存在这里。",
    "loadPrompt": "加载此提示",
    "deletePrompt": "删除此提示",
    "errorAddIdea": "生成前请至少添加一个想法。",
    "errorSpecifyPurpose": "生成前请为提示指定一个目的。",
    "tokens": "个 token",
    "input": "输入",
    "output": "输出",
    "total": "总计",
  },
};

const useTranslations = () => {
  const { language } = useContext(LanguageContext);
  const t = (key: string, params?: Record<string, string>): string => {
    let translation = i18n[language][key] || i18n['en'][key] || key;
    if (params) {
      Object.keys(params).forEach(paramKey => {
        translation = translation.replace(`{${paramKey}}`, params[paramKey]);
      });
    }
    return translation;
  };
  return { t, language };
};


const formatRelativeTime = (timestamp: number, locale: string): string => {
  try {
    const rtf = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' });
    const seconds = Math.floor((Date.now() - timestamp) / 1000);

    if (seconds < 60) return rtf.format(Math.floor(-seconds), 'second');
    if (seconds < 3600) return rtf.format(-Math.floor(seconds / 60), 'minute');
    if (seconds < 86400) return rtf.format(-Math.floor(seconds / 3600), 'hour');
    if (seconds < 2592000) return rtf.format(-Math.floor(seconds / 86400), 'day');
    if (seconds < 31536000) return rtf.format(-Math.floor(seconds / 2592000), 'month');
    return rtf.format(-Math.floor(seconds / 31536000), 'year');
  } catch (e) {
    console.error("Error formatting relative time", e);
    return new Date(timestamp).toLocaleDateString();
  }
};

// --- Sub-components (Memoized) ---

interface IdeaInputProps {
  currentIdea: string;
  setCurrentIdea: (value: string) => void;
  onAddIdea: () => void;
}
const IdeaInput: React.FC<IdeaInputProps> = React.memo(({ currentIdea, setCurrentIdea, onAddIdea }) => {
  const { t } = useTranslations();
  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      onAddIdea();
    }
  };

  return (
    <div className="flex gap-2">
      <input
        type="text"
        value={currentIdea}
        onChange={(e) => setCurrentIdea(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder={t('ideaPlaceholder')}
        className="flex-grow bg-white dark:bg-[#10151b] border-2 border-slate-300 dark:border-[#7c4dff]/50 rounded-md p-3 focus:outline-none focus:ring-2 focus:ring-indigo-500 dark:focus:ring-[#7c4dff] focus:border-indigo-500 dark:focus:border-[#7c4dff] transition-all duration-300 placeholder:text-slate-400 dark:placeholder:text-[#90a4ae]/50 text-slate-800 dark:text-white"
      />
      <button
        onClick={onAddIdea}
        disabled={!currentIdea.trim()}
        className="bg-indigo-500 dark:bg-[#7c4dff] text-white p-3 rounded-md flex items-center justify-center hover:bg-indigo-600 dark:hover:bg-[#8e24aa] disabled:bg-slate-400 dark:disabled:bg-gray-600 disabled:cursor-not-allowed transition-all duration-300 shadow-lg shadow-indigo-500/30 dark:shadow-[0_0_10px_rgba(124,77,255,0.5)] hover:shadow-xl hover:shadow-indigo-500/40 dark:hover:shadow-[0_0_15px_rgba(142,36,170,0.7)]"
        aria-label={t('addIdea')}
      >
        <Plus size={24} />
      </button>
    </div>
  );
});

interface IdeasListProps {
  ideas: Idea[];
  onRemoveIdea: (id: string) => void;
}
const IdeasList: React.FC<IdeasListProps> = React.memo(({ ideas, onRemoveIdea }) => {
  const { t } = useTranslations();
  return (
    <div className="space-y-3 mt-4 pr-2 max-h-60 overflow-y-auto">
      {ideas.map((idea, index) => (
        <div key={idea.id} className="bg-slate-100/50 dark:bg-[#10151b]/50 p-3 rounded-md flex justify-between items-center border border-transparent hover:border-sky-400/50 dark:hover:border-[#00f0ff]/30 transition-colors duration-300 animate-fade-in" style={{ animationDelay: `${index * 50}ms` }}>
          <span className="text-slate-700 dark:text-[#e0f7fa]">{idea.text}</span>
          <button
            onClick={() => onRemoveIdea(idea.id)}
            className="text-red-500 dark:text-red-400 hover:text-red-600 dark:hover:text-red-300 p-1 rounded-full hover:bg-red-500/10 dark:hover:bg-red-500/20 transition-all duration-200"
            aria-label={t('removeIdea', { idea: idea.text })}
          >
            <Trash2 size={16} />
          </button>
        </div>
      ))}
    </div>
  );
});

interface PurposeSelectorProps {
  purpose: string;
  setPurpose: (value: string) => void;
}
const PurposeSelector: React.FC<PurposeSelectorProps> = React.memo(({ purpose, setPurpose }) => {
  const { t } = useTranslations();
  const purposes = Object.keys(purposeKeys);

  return (
    <div className="mt-6">
      <label htmlFor="purpose-input" className="block text-lg font-medium text-sky-500 dark:text-[#00f0ff] mb-3 light-shadow-sky dark:neon-glow-cyan">{t('purposeLabel')}</label>
      <div className="grid grid-cols-2 sm:grid-cols-3 gap-2 mb-4">
        {purposes.map(p => (
          <button
            key={p}
            onClick={() => setPurpose(p)}
            className={`p-2 rounded-md text-sm font-medium border-2 transition-all duration-200 ${purpose === p
                ? 'bg-emerald-500/20 dark:bg-[#00e676]/20 border-emerald-500 dark:border-[#00e676] text-emerald-700 dark:text-white scale-105'
                : 'bg-slate-100 dark:bg-[#10151b]/50 border-transparent hover:border-emerald-500/50 dark:hover:border-[#00e676]/50 text-slate-600 dark:text-slate-300'
              }`}
          >
            {t(purposeKeys[p])}
          </button>
        ))}
      </div>
      <input
        id="purpose-input"
        type="text"
        value={purpose}
        onChange={(e) => setPurpose(e.target.value)}
        list="purposes-list"
        placeholder={t('customPurposePlaceholder')}
        className="w-full bg-white dark:bg-[#10151b] border-2 border-emerald-500/50 dark:border-[#00e676]/30 rounded-md p-3 focus:outline-none focus:ring-2 focus:ring-emerald-500 dark:focus:ring-[#00e676] focus:border-emerald-500 dark:focus:border-[#00e676] transition-all duration-300 placeholder:text-slate-400 dark:placeholder:text-[#90a4ae]/50 text-emerald-600 dark:text-[#00e676] font-semibold"
      />
      <datalist id="purposes-list">
        {purposes.map(p => <option key={p} value={p} />)}
      </datalist>
    </div>
  );
});

interface PromptHistoryProps {
  history: HistoryItem[];
  onLoad: (item: HistoryItem) => void;
  onDelete: (id: string) => void;
  onClear: () => void;
}
const PromptHistoryDisplay: React.FC<PromptHistoryProps> = React.memo(({ history, onLoad, onDelete, onClear }) => {
  const { t, language } = useTranslations();
  return (
    <div className="lg:col-span-2 bg-white/60 dark:bg-[#10151b]/30 p-6 rounded-lg border-2 border-slate-200 dark:border-sky-500/30 backdrop-blur-sm shadow-2xl shadow-slate-500/10">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-2xl font-bold font-orbitron text-sky-600 dark:text-sky-400 tracking-wider">{t('historyTitle')}</h3>
        {history.length > 0 && (
          <button
            onClick={onClear}
            className="flex items-center gap-2 px-3 py-1 text-sm bg-red-500/10 dark:bg-red-500/20 text-red-600 dark:text-red-300 rounded-full hover:bg-red-500/20 dark:hover:bg-red-500/30 transition-colors duration-200"
            aria-label={t('clearAll')}
          >
            <XCircle size={16} />
            <span>{t('clearAll')}</span>
          </button>
        )}
      </div>
      <div className="max-h-96 overflow-y-auto pr-2 space-y-3">
        {history.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-40 text-slate-400 dark:text-[#90a4ae]/50">
            <History size={48} className="mb-4" />
            <p className="font-semibold text-center">{t('historyPlaceholder')}</p>
          </div>
        ) : (
          history.map(item => (
            <div key={item.id} className="bg-slate-100/50 dark:bg-[#10151b]/50 p-3 rounded-md flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3 border border-transparent hover:border-sky-400/50 dark:hover:border-sky-500/50 transition-colors duration-300">
              <div className="flex-grow overflow-hidden">
                <div className="flex items-center gap-2 mb-1 flex-wrap">
                  <span className="px-2 py-0.5 text-xs font-bold text-white dark:text-black rounded-full bg-sky-500 dark:bg-sky-400 flex-shrink-0">{item.purpose}</span>
                  <span className="text-xs text-slate-500 dark:text-[#90a4ae] truncate">{formatRelativeTime(item.timestamp, language)}</span>
                </div>
                <p className="text-sm text-slate-700 dark:text-[#e0f7fa] line-clamp-2">
                  {item.prompt}
                </p>
              </div>
              <div className="flex items-center gap-2 self-end sm:self-center flex-shrink-0">
                <button
                  onClick={() => onLoad(item)}
                  className="p-2 rounded-md bg-indigo-500/10 dark:bg-indigo-400/20 text-indigo-600 dark:text-indigo-300 hover:bg-indigo-500/20 dark:hover:bg-indigo-400/30 transition-colors duration-200"
                  aria-label={t('loadPrompt')}
                >
                  <Eye size={18} />
                </button>
                <button
                  onClick={() => onDelete(item.id)}
                  className="p-2 rounded-md bg-red-500/10 dark:bg-red-400/20 text-red-600 dark:text-red-300 hover:bg-red-500/20 dark:hover:bg-red-400/30 transition-colors duration-200"
                  aria-label={t('deletePrompt')}
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
});

// --- Main Component ---

const PromptCrafter: React.FC = () => {
  const { t } = useTranslations();
  const [ideas, setIdeas] = useState<Idea[]>([]);
  const [currentIdea, setCurrentIdea] = useState('');
  const [purpose, setPurpose] = useState('Code Generation');
  const [isLoading, setIsLoading] = useState(false);
  const [generatedPrompt, setGeneratedPrompt] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [isCopied, setIsCopied] = useState(false);
  const [isLinkCopied, setIsLinkCopied] = useState(false);
  const [promptHistory, setPromptHistory] = useState<HistoryItem[]>([]);
  const [viewMode, setViewMode] = useState<'preview' | 'raw'>('preview');
  const [tokenUsage, setTokenUsage] = useState<{ input: number; output: number; } | null>(null);
  const isInitialLoad = useRef(true);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

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
    if (promptHistory.length === 0 && localStorage.getItem('promptHistory') === null) {
      return;
    }
    try {
      localStorage.setItem('promptHistory', JSON.stringify(promptHistory));
    } catch (error) {
      console.error("Failed to save history to localStorage", error);
    }
  }, [promptHistory]);

  // Load shared prompt from URL hash on initial load
  useEffect(() => {
    const loadSharedPrompt = () => {
      try {
        const hash = window.location.hash;
        if (hash.startsWith('#prompt=')) {
          const encodedData = hash.substring('#prompt='.length);
          const decodedJson = atob(encodedData);
          const data = JSON.parse(decodedJson) as { ideas: Idea[], purpose: string, prompt: string };

          if (data.ideas && data.purpose && data.prompt) {
            setIdeas(data.ideas);
            setPurpose(data.purpose);
            setGeneratedPrompt(data.prompt);
            setTokenUsage(null); // Tokens are not shared in link
            setError(null);
            window.history.replaceState(null, document.title, window.location.pathname + window.location.search);
          }
        }
      } catch (e) {
        console.error("Failed to parse shared prompt from URL", e);
        setError("The shared link appears to be invalid or corrupted.");
        window.history.replaceState(null, document.title, window.location.pathname + window.location.search);
      }
    };
    loadSharedPrompt();
  }, []);

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
      return;
    }
    const handler = setTimeout(() => {
      try {
        setDraft('currentDraft', { ideas, purpose });
      } catch (e) {
        console.error("Failed to save draft to IndexedDB", e);
      }
    }, 500);
    return () => clearTimeout(handler);
  }, [ideas, purpose]);

  // Auto-resize textarea
  useEffect(() => {
    if (viewMode === 'raw' && textareaRef.current) {
      textareaRef.current.style.height = 'auto';
      const scrollHeight = textareaRef.current.scrollHeight;
      textareaRef.current.style.height = `${scrollHeight}px`;
    }
  }, [generatedPrompt, viewMode]);


  const handleAddIdea = useCallback(() => {
    if (currentIdea.trim() !== '') {
      setIdeas(prev => [...prev, { id: Date.now().toString(), text: currentIdea.trim() }]);
      setCurrentIdea('');
    }
  }, [currentIdea]);

  const handleRemoveIdea = useCallback((id: string) => {
    setIdeas(prev => prev.filter(idea => idea.id !== id));
  }, []);

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
      setTokenUsage({ input: inputTokens, output: outputTokens });

      const newItem: HistoryItem = {
        id: Date.now().toString(),
        prompt: result.prompt,
        purpose,
        ideas: [...ideas],
        timestamp: Date.now(),
        inputTokens,
        outputTokens,
      };
      setPromptHistory(prev => [newItem, ...prev].slice(0, 20));
    } catch (err: any) {
      setError(err.message || 'An unknown error occurred.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCopy = () => {
    if (generatedPrompt) {
      navigator.clipboard.writeText(generatedPrompt);
      setIsCopied(true);
      setTimeout(() => setIsCopied(false), 2000);
    }
  };

  const handleShare = () => {
    if (!generatedPrompt) return;
    try {
      const dataToShare = {
        ideas,
        purpose,
        prompt: generatedPrompt,
      };
      const jsonString = JSON.stringify(dataToShare);
      const encodedData = btoa(jsonString);
      const shareUrl = `${window.location.origin}${window.location.pathname}#prompt=${encodedData}`;

      navigator.clipboard.writeText(shareUrl);
      setIsLinkCopied(true);
      setTimeout(() => setIsLinkCopied(false), 2500);

    } catch (e) {
      console.error("Failed to create share link", e);
      setError("An unexpected error occurred while creating the shareable link.");
    }
  };

  const handleLoadExample = useCallback(() => {
    setError(null);
    setGeneratedPrompt('');
    setTokenUsage(null);
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
    setTokenUsage(item.inputTokens && item.outputTokens ? { input: item.inputTokens, output: item.outputTokens } : null);
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

  return (
    <div className="grid lg:grid-cols-2 gap-8">
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
        <IdeaInput currentIdea={currentIdea} setCurrentIdea={setCurrentIdea} onAddIdea={handleAddIdea} />
        <IdeasList ideas={ideas} onRemoveIdea={handleRemoveIdea} />
        <PurposeSelector purpose={purpose} setPurpose={setPurpose} />
      </div>

      {/* Output Section */}
      <div className="flex flex-col">
        <div className="bg-white/60 dark:bg-[#10151b]/30 p-6 rounded-lg border-2 border-slate-200 dark:border-[#00e676]/30 backdrop-blur-sm flex-grow flex flex-col shadow-2xl shadow-slate-500/10">
          <h3 className="text-2xl font-bold font-orbitron text-emerald-600 dark:text-[#00e676] tracking-wider mb-4 light-shadow-emerald dark:neon-glow-green">{t('generatedPromptTitle')}</h3>
          <div className="flex-grow bg-slate-100 dark:bg-[#0a0f14] rounded-md min-h-[300px] flex flex-col border border-slate-200 dark:border-transparent">
            <div className="relative flex-grow p-4 overflow-y-auto">
              {isLoading && (
                <div className="absolute inset-0 flex flex-col items-center justify-center bg-slate-100/80 dark:bg-[#0a0f14]/80 backdrop-blur-sm z-20">
                  <Loader className="animate-spin text-sky-500 dark:text-[#00f0ff]" size={48} />
                  <p className="mt-4 text-sky-500 dark:text-[#00f0ff] font-semibold">{t('generatingMessage')}</p>
                </div>
              )}
              {error && (
                <div className="absolute inset-0 flex flex-col items-center justify-center bg-slate-100/80 dark:bg-[#0a0f14]/80 backdrop-blur-sm p-4 z-20">
                  <div className="bg-red-100 dark:bg-red-900/50 border border-red-300 dark:border-red-500 text-red-700 dark:text-red-300 p-4 rounded-lg flex flex-col items-center gap-2 text-center shadow-lg">
                    <AlertTriangle size={32} />
                    <h4 className="font-bold">{t('generationFailedTitle')}</h4>
                    <p className="text-sm">{error}</p>
                    <button onClick={() => setError(null)} className="mt-2 bg-red-500/20 dark:bg-red-500/50 hover:bg-red-500/30 dark:hover:bg-red-500/80 text-red-800 dark:text-white px-3 py-1 text-sm rounded-md flex items-center gap-1">
                      <X size={14} /> {t('close')}
                    </button>
                  </div>
                </div>
              )}
              {generatedPrompt && !isLoading && !error && (
                <>
                  <div className="absolute top-2 right-2 flex gap-2 z-10">
                    <button onClick={() => setViewMode(viewMode === 'raw' ? 'preview' : 'raw')} className="bg-white dark:bg-[#10151b] p-2 rounded-md text-sky-500 dark:text-[#00f0ff] hover:bg-sky-100 dark:hover:bg-[#00f0ff] hover:text-sky-600 dark:hover:text-black transition-all duration-200" aria-label={t('toggleView')}>
                      {viewMode === 'raw' ? <Eye size={20} /> : <Code size={20} />}
                    </button>
                    <button onClick={handleShare} className="bg-white dark:bg-[#10151b] p-2 rounded-md text-sky-500 dark:text-[#00f0ff] hover:bg-sky-100 dark:hover:bg-[#00f0ff] hover:text-sky-600 dark:hover:text-black transition-all duration-200" aria-label={t('copyLink')}>
                      {isLinkCopied ? <ClipboardCheck size={20} /> : <Share2 size={20} />}
                    </button>
                    <button onClick={handleCopy} className="bg-white dark:bg-[#10151b] p-2 rounded-md text-sky-500 dark:text-[#00f0ff] hover:bg-sky-100 dark:hover:bg-[#00f0ff] hover:text-sky-600 dark:hover:text-black transition-all duration-200" aria-label={t('copyPrompt')}>
                      {isCopied ? <ClipboardCheck size={20} /> : <Clipboard size={20} />}
                    </button>
                  </div>
                  {viewMode === 'preview' ? (
                    <div className="prose prose-sm max-w-none text-slate-800 dark:text-[#e0f7fa] dark:prose-invert prose-headings:font-orbitron prose-code:font-plex-mono prose-code:bg-slate-200 dark:prose-code:bg-[#10151b] prose-code:px-1.5 prose-code:py-0.5 prose-code:rounded-md prose-pre:bg-slate-200 dark:prose-pre:bg-[#10151b]">
                      <ReactMarkdown remarkPlugins={[remarkGfm]}>{generatedPrompt}</ReactMarkdown>
                    </div>
                  ) : (
                    <textarea
                      title={`(t('${generatedPrompt}'))`}
                      ref={textareaRef}
                      readOnly
                      value={generatedPrompt}
                      className="w-full h-auto bg-transparent resize-none border-none focus:ring-0 p-0 m-0 font-plex-mono text-sm leading-relaxed text-slate-800 dark:text-[#e0f7fa] overflow-hidden"
                      rows={1}
                    />
                  )}
                </>
              )}
              {!generatedPrompt && !isLoading && !error && (
                <div className="flex flex-col items-center justify-center h-full text-slate-400 dark:text-[#90a4ae]/50">
                  <Wand2 size={48} className="mb-4" />
                  <p className="font-semibold text-center">{t('promptPlaceholder')}</p>
                </div>
              )}
            </div>
            {tokenUsage && (
              <div className="flex-shrink-0 bg-slate-200/50 dark:bg-[#10151b]/50 p-2 border-t border-slate-200 dark:border-slate-700/50">
                <div className="flex items-center justify-center gap-4 text-xs text-slate-600 dark:text-slate-400 font-semibold">
                  <BrainCircuit size={16} className="text-sky-500 dark:text-sky-400" />
                  <span>{t('input')}: <span className="font-bold text-sky-600 dark:text-sky-300">{tokenUsage.input}</span> {t('tokens')}</span>
                  <span className="text-slate-300 dark:text-slate-600">|</span>
                  <span>{t('output')}: <span className="font-bold text-emerald-600 dark:text-emerald-300">{tokenUsage.output}</span> {t('tokens')}</span>
                  <span className="text-slate-300 dark:text-slate-600">|</span>
                  <span>{t('total')}: <span className="font-bold text-indigo-600 dark:text-indigo-300">{tokenUsage.input + tokenUsage.output}</span> {t('tokens')}</span>
                </div>
              </div>
            )}
          </div>
          <button
            onClick={handleGenerate}
            disabled={isLoading || ideas.length === 0 || !purpose.trim()}
            className="w-full mt-6 bg-gradient-to-r from-emerald-500 to-sky-500 dark:from-[#00e676] dark:to-[#00f0ff] text-white dark:text-black font-bold font-orbitron text-lg p-4 rounded-lg flex items-center justify-center gap-3 hover:scale-105 disabled:scale-100 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300 shadow-lg shadow-sky-500/40 dark:shadow-[0_0_15px_rgba(0,230,118,0.5)] hover:shadow-xl hover:shadow-sky-500/50 dark:hover:shadow-[0_0_25px_rgba(0,240,255,0.7)]"
          >
            {isLoading ? <Loader className="animate-spin" size={28} /> : <Wand2 size={28} />}
            {isLoading ? t('generatingButton') : t('generateButton')}
          </button>
        </div>
      </div>

      {/* Prompt History Section */}
      <PromptHistoryDisplay
        history={promptHistory}
        onLoad={handleLoadFromHistory}
        onDelete={handleDeleteFromHistory}
        onClear={handleClearHistory}
      />
    </div>
  );
};

export default PromptCrafter;
