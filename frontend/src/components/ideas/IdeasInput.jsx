import { Plus } from 'lucide-react';

const IdeasInput = ({
  currentInput,
  setCurrentInput,
  addIdea,
  currentTheme
}) => {
  return (
    <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg`} id="ideas-input">
      <h2 className="text-xl font-semibold mb-4">📝 Adicionar Ideias</h2>
      <div className="space-y-4">
        <textarea
          value={currentInput}
          onChange={(e) => setCurrentInput(e.target.value)}
          placeholder="Cole suas notas, ideias brutas ou pensamentos desorganizados aqui..."
          className={`w-full h-32 px-4 py-3 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-blue-500 resize-none`}
          onKeyDown={(e) => {
            if (e.key === 'Enter' && e.ctrlKey) {
              addIdea();
            }
          }}
        />
        <button
          onClick={addIdea}
          disabled={!currentInput.trim()}
          className={`w-full flex items-center justify-center gap-2 px-4 py-3 rounded-lg ${currentTheme.button} disabled:opacity-50 disabled:cursor-not-allowed transition-all`}
        >
          <Plus size={20} />
          Incluir (Ctrl + Enter)
        </button>
      </div>
    </div>
  );
};

export default IdeasInput;
