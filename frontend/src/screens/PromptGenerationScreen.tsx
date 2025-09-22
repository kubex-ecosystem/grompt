import * as React from 'react';
import IdeasInput from '../components/ideas/IdeasInput';
import IdeasList from '../components/ideas/IdeasList';
import OutputPanel from '../components/settings/OutputPanel';
import { Purpose } from '../hooks/usePromptCrafter';
import { UseGromptAPIState } from '../hooks/useGromptAPI';

interface Theme {
  [key: string]: string;
}

interface PromptGenerationScreenProps {
  // Ideas state
  currentInput: string;
  setCurrentInput: (value: string) => void;
  ideas: Array<{ id: string; text: string; timestamp: Date }>;
  editingId: string | null;
  editingText: string;
  setEditingText: (value: string) => void;

  // Prompt configuration
  purpose: Purpose;
  setPurpose: (value: Purpose) => void;
  customPurpose: string;
  setCustomPurpose: (value: string) => void;
  maxLength: number;
  setMaxLength: (value: number) => void;

  // Generation state
  generatedPrompt: string;
  isGenerating: boolean;
  copied: boolean;

  // Actions
  addIdea: (text: string) => void;
  removeIdea: (id: string) => void;
  startEditing: (id: string, text: string) => void;
  saveEdit: () => void;
  cancelEdit: () => void;
  generatePrompt: () => void;
  copyToClipboard: () => void;

  // Theme and API
  currentTheme: Theme;
  apiGenerateState: UseGromptAPIState['generatePrompt'];
}

const PromptGenerationScreen: React.FC<PromptGenerationScreenProps> = ({
  currentInput,
  setCurrentInput,
  ideas,
  editingId,
  editingText,
  setEditingText,
  purpose,
  setPurpose,
  customPurpose,
  setCustomPurpose,
  maxLength,
  setMaxLength,
  generatedPrompt,
  isGenerating,
  copied,
  addIdea,
  removeIdea,
  startEditing,
  saveEdit,
  cancelEdit,
  generatePrompt,
  copyToClipboard,
  currentTheme,
  apiGenerateState
}) => {
  const purposeOptions: Purpose[] = ['Código', 'Imagem', 'Análise', 'Escrita', 'Outros'];

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="text-center">
        <h1 className="text-3xl font-bold mb-2 text-white">📝 Geração de Prompts</h1>
        <p className="text-gray-400">Transforme suas ideias em prompts eficazes para IA</p>
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">

        {/* Left Column: Input and Configuration */}
        <div className="space-y-6">
          {/* Ideas Input Card */}
          <div className="bg-gray-800/50 border border-gray-700/80 rounded-xl p-6 backdrop-blur-sm transition-all duration-300 hover:border-purple-500/50">
            <IdeasInput
              currentInput={currentInput}
              setCurrentInput={setCurrentInput}
              addIdea={addIdea}
              currentTheme={currentTheme}
            />
          </div>

          {/* Prompt Configuration */}
          <div className="bg-gray-800/50 border border-gray-700/80 rounded-xl p-6 backdrop-blur-sm transition-all duration-300 hover:border-purple-500/50">
            <h2 className="text-xl font-semibold mb-4 text-white">⚙️ Configurações do Prompt</h2>

            <div className="space-y-4">
              {/* Purpose Selection */}
              <div>
                <label className="text-sm font-medium mb-2 text-white flex items-center">
                  Propósito do Prompt
                </label>
                <div className="space-y-2">
                  <div className="flex gap-2 flex-wrap">
                    {purposeOptions.map((option) => (
                      <button
                        key={option}
                        onClick={() => setPurpose(option)}
                        className={`px-3 py-2 rounded-lg text-sm border transition-colors ${
                          purpose === option
                            ? 'bg-purple-600 text-white border-purple-600'
                            : 'bg-gray-700/80 border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white'
                        }`}
                      >
                        {option}
                      </button>
                    ))}
                  </div>
                  {purpose === 'Outros' && (
                    <input
                      type="text"
                      value={customPurpose}
                      onChange={(e) => setCustomPurpose(e.target.value)}
                      placeholder="Descreva o objetivo do prompt..."
                      className="w-full px-3 py-2 rounded-lg border border-gray-600 bg-gray-700/80 text-white placeholder-gray-400 focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
                    />
                  )}
                </div>
              </div>

              {/* Max Length */}
              <div>
                <label className="text-sm font-medium mb-2 text-white flex items-center" htmlFor="max-length-slider">
                  Tamanho Máximo: {maxLength.toLocaleString()} caracteres
                </label>
                <input
                  id='max-length-slider'
                  type="range"
                  min="500"
                  max="130000"
                  step="500"
                  value={maxLength}
                  onChange={(e) => setMaxLength(parseInt(e.target.value))}
                  className="w-full"
                />
              </div>
            </div>
          </div>
        </div>

        {/* Center Column: Ideas List */}
        <div className="bg-gray-800/50 border border-gray-700/80 rounded-xl p-6 backdrop-blur-sm transition-all duration-300 hover:border-purple-500/50">
          <IdeasList
            ideas={ideas}
            editingId={editingId}
            editingText={editingText}
            setEditingText={setEditingText}
            startEditing={startEditing}
            saveEdit={saveEdit}
            cancelEdit={cancelEdit}
            removeIdea={removeIdea}
            generatePrompt={generatePrompt}
            isGenerating={isGenerating}
            outputType="prompt"
            currentTheme={currentTheme}
            apiGenerateState={apiGenerateState}
          />
        </div>

        {/* Right Column: Output Panel */}
        <div className="bg-gray-800/50 border border-gray-700/80 rounded-xl p-6 backdrop-blur-sm transition-all duration-300 hover:border-purple-500/50">
          <OutputPanel
            generatedPrompt={generatedPrompt}
            copyToClipboard={copyToClipboard}
            copied={copied}
            outputType="prompt"
            agentFramework="crewai"
            agentProvider={agentProvider}
            maxLength={maxLength}
            mcpServers={[]}
            currentTheme={currentTheme}
            apiGenerateState={apiGenerateState}
          />
        </div>
      </div>
    </div>
  );
};

export default PromptGenerationScreen;