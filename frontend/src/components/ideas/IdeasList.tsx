import { AlertCircle, Edit3, Loader2, Trash2, Wand2 } from 'lucide-react';
import * as React from 'react';
import { useGromptAPI } from '../../hooks/useGromptAPI';
import type { Idea } from '@/types';
import { OutputType } from '../../hooks/usePromptCrafter';

interface Theme {
  [key: string]: string;
}

const {
  generatePrompt,
} = useGromptAPI({});


interface IdeasListProps {
  ideas: Idea[];
  editingId: string | null;
  editingText: string;
  setEditingText: (value: string) => void;
  startEditing: (id: string, text: string) => void;
  saveEdit: () => void;
  cancelEdit: () => void;
  removeIdea: (id: string) => void;
  generatePrompt: () => void;
  isGenerating: boolean;
  outputType: OutputType;
  currentTheme: Theme;
  apiGenerateState?: typeof generatePrompt;
}

const IdeasList: React.FC<IdeasListProps> = ({
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
  currentTheme,
  apiGenerateState
}) => {
  // Show API state information
  const isAPIGenerating = apiGenerateState?.loading || apiGenerateState?.progress?.isStreaming;
  const apiError = apiGenerateState?.error;

  return (
    <div>
      <h2 className="text-xl font-semibold mb-4 text-white">💡 Suas Ideias ({ideas.length})</h2>

      {/* API Status Indicator */}
      {apiError && (
        <div className="mb-4 p-3 bg-red-900/50 border border-red-700 rounded-lg flex items-center gap-2">
          <AlertCircle size={16} className="text-red-400" />
          <span className="text-red-400 text-sm">
            Erro na API: {apiError.message}
          </span>
        </div>
      )}

      {/* Streaming Progress */}
      {apiGenerateState?.progress?.isStreaming && (
        <div className="mb-4 p-3 bg-purple-900/50 border border-purple-700 rounded-lg">
          <div className="flex items-center gap-2 mb-2">
            <Loader2 size={16} className="text-purple-400 animate-spin" />
            <span className="text-purple-400 text-sm">Gerando prompt...</span>
          </div>
          {apiGenerateState.progress.content && (
            <div className="bg-gray-800/50 p-2 rounded text-xs font-mono text-gray-300 max-h-20 overflow-y-auto">
              {apiGenerateState.progress.content}
            </div>
          )}
        </div>
      )}

      <div className="space-y-3 max-h-96 overflow-y-auto">
        {ideas.length === 0 ? (
          <p className="text-gray-400 text-center py-8">
            Adicione suas primeiras ideias ao lado ←
          </p>
        ) : (
          ideas.map((idea) => (
            <div key={idea.id} className="p-3 rounded-lg border border-gray-600 bg-gray-700/50 backdrop-blur-sm">
              {editingId === idea.id ? (
                <div className="space-y-2">
                  <textarea
                    title='Edite sua ideia aqui'
                    value={editingText}
                    onChange={(e) => setEditingText(e.target.value)}
                    className="w-full px-2 py-1 rounded border border-gray-600 bg-gray-700/80 text-white text-sm resize-none"
                    rows={2}
                  />
                  <div className="flex gap-1">
                    <button
                      onClick={saveEdit}
                      className="px-2 py-1 bg-green-600 text-white rounded text-xs hover:bg-green-700 transition-colors"
                    >
                      Salvar
                    </button>
                    <button
                      onClick={cancelEdit}
                      className="px-2 py-1 rounded text-xs bg-gray-700/80 border border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white transition-colors"
                    >
                      Cancelar
                    </button>
                  </div>
                </div>
              ) : (
                <>
                  <p className="text-sm mb-2 text-gray-300">{idea.text}</p>
                  <div className="flex justify-end gap-1">
                    <button
                      title='Edite sua ideia aqui'
                      onClick={() => startEditing(idea.id, idea.text)}
                      className="p-1 rounded bg-gray-700/80 border border-gray-600 text-gray-400 hover:bg-gray-700 hover:text-white transition-colors"
                    >
                      <Edit3 size={14} />
                    </button>
                    <button
                      title='Remova sua ideia aqui'
                      onClick={() => removeIdea(idea.id)}
                      className="p-1 rounded bg-red-600 text-white hover:bg-red-700 transition-colors"
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
          title='Generate Prompt or Agent'
          onClick={() => generatePrompt}
          disabled={isGenerating || isAPIGenerating}
          className={`w-full mt-4 flex items-center justify-center gap-2 px-4 py-3 rounded-lg bg-gradient-to-r ${outputType === 'prompt'
            ? 'from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700'
            : 'from-green-600 to-blue-600 hover:from-green-700 hover:to-blue-700'
            } text-white font-semibold disabled:opacity-50 disabled:cursor-not-allowed transition-all transform hover:scale-105 shadow-lg`}
        >
          {isGenerating || isAPIGenerating ? (
            <>
              <Wand2 size={20} className="animate-spin" />
              {apiGenerateState?.progress?.isStreaming
                ? 'Gerando (Streaming)...'
                : `Gerando ${outputType === 'prompt' ? 'prompt' : 'agent'}...`
              }
            </>
          ) : (
            <>
              <Wand2 size={20} />
              {`Criar ${outputType === 'prompt' ? 'Prompt' : 'Agent'} 🚀`}
            </>
          )}
        </button>
      )}
    </div>
  );
};

export default IdeasList;
