import { Plus } from 'lucide-react';
import * as React from 'react';

interface Theme {
  [key: string]: string;
}

interface IdeasInputProps {
  currentInput: { id: string; text: string; timestamp: Date };
  setCurrentInput: (value: { id: string; text: string; timestamp: Date }) => void;
  addIdea: (ideas: { id: string; text: string; timestamp: Date }) => void;
  currentTheme: Theme;
}

const IdeasInput: React.FC<IdeasInputProps> = ({
  currentInput,
  setCurrentInput,
  addIdea,
  currentTheme
}) => {
  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && e.ctrlKey) {
      addIdea(currentInput);
    }
  };

  return (
    <div id="ideas-input">
      <h2 className="text-xl font-semibold mb-4 text-white">📝 Adicionar Ideias</h2>
      <div className="space-y-4">
        <textarea
          key={currentInput.id}
          value={currentInput.text}
          onChange={(e) => setCurrentInput({ ...currentInput, text: e.target.value })}
          placeholder="Cole suas notas, ideias brutas ou pensamentos desorganizados aqui..."
          className="w-full h-32 px-4 py-3 rounded-lg border border-gray-600 bg-gray-700/80 text-white placeholder-gray-400 focus:ring-2 focus:ring-purple-500 focus:border-purple-500 resize-none transition-colors"
          onKeyDown={handleKeyDown}
        />
        <button
          onClick={() => addIdea(currentInput)}
          disabled={!currentInput.text.trim()}
          className="w-full flex items-center justify-center gap-2 px-4 py-3 rounded-lg bg-gradient-to-r from-purple-600 to-blue-500 text-white font-semibold disabled:opacity-50 disabled:cursor-not-allowed transition-all hover:scale-105 hover:from-purple-700 hover:to-blue-600"
        >
          <Plus size={20} />
          Incluir (Ctrl + Enter)
        </button>
      </div>
    </div>
  );
};

export default IdeasInput;
