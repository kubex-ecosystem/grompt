import React from 'react';
import { useTranslations } from '../../i18n/useTranslations';
import { Plus } from 'lucide-react';

interface IdeaInputProps {
  currentIdea: string;
  setCurrentIdea: (value: string) => void;
  onAddIdea: () => void;
  disabled?: boolean;
}

const IdeaInput: React.FC<IdeaInputProps> = React.memo(({ currentIdea, setCurrentIdea, onAddIdea, disabled }) => {
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
        disabled={!currentIdea.trim() || disabled}
        className="bg-indigo-500 dark:bg-[#7c4dff] text-white p-3 rounded-md flex items-center justify-center hover:bg-indigo-600 dark:hover:bg-[#8e24aa] disabled:bg-slate-400 dark:disabled:bg-gray-600 disabled:cursor-not-allowed transition-all duration-300 shadow-lg shadow-indigo-500/30 dark:shadow-[0_0_10px_rgba(124,77,255,0.5)] hover:shadow-xl hover:shadow-indigo-500/40 dark:hover:shadow-[0_0_15px_rgba(142,36,170,0.7)]"
        aria-label={t('addIdea')}
      >
        <Plus size={24} />
      </button>
    </div>
  );
});

export default IdeaInput;