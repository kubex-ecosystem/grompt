import React from 'react';
import { useTranslations } from '../../i18n/useTranslations';
import { Idea } from '../../types';
import { Trash2 } from 'lucide-react';

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

export default IdeasList;
