import { AlertTriangle, BrainCircuit, Check, Clipboard, ClipboardCheck, Code, Copy, Eye, Loader, Share2, Wand2, X } from 'lucide-react';
import React, { useEffect, useRef, useState } from 'react';
import ReactMarkdown from 'react-markdown';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { oneDark, oneLight } from 'react-syntax-highlighter/dist/esm/styles/prism';
import remarkGfm from 'remark-gfm';
import { useTranslations } from '../../i18n/useTranslations';
import { Idea, Theme } from '../../types';


interface GeneratedPromptProps {
  isLoading: boolean;
  error: string | null;
  setError: (error: string | null) => void;
  generatedPrompt: string;
  ideas: Idea[];
  purpose: string;
  tokenUsage: { input: number; output: number; total: number; } | null;
  theme: Theme;
}

const GeneratedPrompt: React.FC<GeneratedPromptProps> = ({ isLoading, error, setError, generatedPrompt, ideas, purpose, tokenUsage, theme }) => {
  const { t } = useTranslations();
  const [isCopied, setIsCopied] = useState(false);
  const [isLinkCopied, setIsLinkCopied] = useState(false);
  const [viewMode, setViewMode] = useState<'preview' | 'raw'>('preview');
  const [copiedCode, setCopiedCode] = useState<string | null>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  // Auto-resize textarea
  useEffect(() => {
    if (viewMode === 'raw' && textareaRef.current) {
      textareaRef.current.style.height = 'auto';
      const scrollHeight = textareaRef.current.scrollHeight;
      textareaRef.current.style.height = `${scrollHeight}px`;
    }
  }, [generatedPrompt, viewMode]);

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

  const handleCodeCopy = async (code: string) => {
    try {
      await navigator.clipboard.writeText(code);
      setCopiedCode(code);
      setTimeout(() => setCopiedCode(null), 2000);
    } catch (err) {
      console.error('Failed to copy code to clipboard:', err);
    }
  };

  const markdownComponents = {
    code({ node, inline, className, children, ...props }: any) {
      const match = /language-(\w+)/.exec(className || '');
      const codeString = String(children).replace(/\n$/, '');
      return !inline && match ? (
        <div className="relative group my-4">
          <SyntaxHighlighter
            style={theme === 'dark' ? oneDark : oneLight}
            language={match[1]}
            PreTag="div"
            className="rounded-md"
            customStyle={{ margin: 0, padding: '1rem', backgroundColor: theme === 'dark' ? '#10151b' : '#f1f5f9' }}
            {...props}
          >
            {codeString}
          </SyntaxHighlighter>
          <button
            onClick={() => handleCodeCopy(codeString)}
            className={`absolute top-2 right-2 p-1.5 rounded-md text-white opacity-0 group-hover:opacity-100 transition-all duration-200 ${copiedCode === codeString
                ? 'bg-emerald-600'
                : 'bg-slate-700/50 hover:bg-slate-600/70'
              }`}
            aria-label={copiedCode === codeString ? t('copied') : t('copyCode')}
            title={copiedCode === codeString ? t('copied') : t('copyCode')}
          >
            {copiedCode === codeString ? <Check size={16} /> : <Copy size={16} />}
          </button>
        </div>
      ) : (
        <code className={className} {...props}>
          {children}
        </code>
      )
    }
  };


  return (
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
                <div className="prose prose-sm max-w-none text-slate-800 dark:text-[#e0f7fa] dark:prose-invert prose-headings:font-orbitron prose-code:font-plex-mono prose-code:bg-slate-200 dark:prose-code:bg-[#10151b] prose-code:px-1.5 prose-code:py-0.5 prose-code:rounded-md">
                  <ReactMarkdown remarkPlugins={[remarkGfm]} components={markdownComponents}>{generatedPrompt}</ReactMarkdown>
                </div>
              ) : (
                <textarea
                  title='Generated Prompt'
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
              <span>{t('total')}: <span className="font-bold text-indigo-600 dark:text-indigo-300">{tokenUsage.total}</span> {t('tokens')}</span>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default GeneratedPrompt;
