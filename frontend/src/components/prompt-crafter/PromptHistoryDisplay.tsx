import React from 'react';
import { useTranslations } from '../../i18n/useTranslations';
import { HistoryItem } from '../../types';
import { formatRelativeTime } from '../../lib/time';
import { History, XCircle, Eye, Trash2 } from 'lucide-react';

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

export default PromptHistoryDisplay;
