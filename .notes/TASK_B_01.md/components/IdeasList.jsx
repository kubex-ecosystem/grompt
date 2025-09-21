import { Edit3, Trash2, Wand2 } from 'lucide-react';

const IdeasList = ({
  ideas,
  editingId,
  editingText,
  setEditingText,
  startEditing,
  saveEdit,
  cancelEdit,
  removeIdea,
  generatePrompt,
  isGenerating,
  outputType,
  currentTheme
}) => {
  return (
    <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg`}>
      <h2 className="text-xl font-semibold mb-4">💡 Suas Ideias ({ideas.length})</h2>
      <div className="space-y-3 max-h-96 overflow-y-auto">
        {ideas.length === 0 ? (
          <p className={`${currentTheme.textSecondary} text-center py-8`}>
            Adicione suas primeiras ideias ao lado ←
          </p>
        ) : (
          ideas.map((idea) => (
            <div key={idea.id} className={`p-3 rounded-lg border ${currentTheme.border} bg-opacity-50`}>
              {editingId === idea.id ? (
                <div className="space-y-2">
                  <textarea
                    value={editingText}
                    onChange={(e) => setEditingText(e.target.value)}
                    className={`w-full px-2 py-1 rounded border ${currentTheme.input} text-sm`}
                    rows="2"
                  />
                  <div className="flex gap-1">
                    <button
                      onClick={saveEdit}
                      className="px-2 py-1 bg-green-600 text-white rounded text-xs hover:bg-green-700"
                    >
                      Salvar
                    </button>
                    <button
                      onClick={cancelEdit}
                      className={`px-2 py-1 rounded text-xs ${currentTheme.buttonSecondary}`}
                    >
                      Cancelar
                    </button>
                  </div>
                </div>
              ) : (
                <>
                  <p className="text-sm mb-2">{idea.text}</p>
                  <div className="flex justify-end gap-1">
                    <button
                      onClick={() => startEditing(idea.id, idea.text)}
                      className={`p-1 rounded ${currentTheme.buttonSecondary} hover:bg-opacity-80`}
                    >
                      <Edit3 size={14} />
                    </button>
                    <button
                      onClick={() => removeIdea(idea.id)}
                      className="p-1 rounded bg-red-600 text-white hover:bg-red-700"
                    >
                      <Trash2 size={14} />
                    </button>
                  </div>
                </>
              )}
            </div>
          ))
        )}
      </div>

      {ideas.length > 0 && (
        <button
          onClick={generatePrompt}
          disabled={isGenerating}
          className={`w-full mt-4 flex items-center justify-center gap-2 px-4 py-3 rounded-lg bg-gradient-to-r ${outputType === 'prompt'
              ? 'from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700'
              : 'from-green-600 to-blue-600 hover:from-green-700 hover:to-blue-700'
            } text-white disabled:opacity-50 disabled:cursor-not-allowed transition-all transform hover:scale-105`}
        >
          <Wand2 size={20} className={isGenerating ? 'animate-spin' : ''} />
          {isGenerating
            ? `Gerando ${outputType === 'prompt' ? 'prompt' : 'agent'}...`
            : `Criar ${outputType === 'prompt' ? 'Prompt' : 'Agent'} 🚀`
          }
        </button>
      )}
    </div>
  );
};

export default IdeasList;
