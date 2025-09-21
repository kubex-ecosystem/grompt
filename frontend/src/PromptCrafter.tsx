import * as React from 'react';
import ConfigurationPanel from './components/settings/ConfigurationPanel';
import DemoStatusFooter from './components/demo/DemoStatusFooter';
import EducationalModal from './components/onboarding/EducationalModal';
import Header from './components/layout/Header';
import IdeasInput from './components/ideas/IdeasInput';
import IdeasList from './components/ideas/IdeasList';
import OnboardingModal from './components/onboarding/OnboardingModal';
import OutputPanel from './components/settings/OutputPanel';
import { themes } from './constants/themes';
import usePromptCrafter from './hooks/usePromptCrafter';
import { useGromptAPI } from './hooks/useGromptAPI';

const PromptCrafter: React.FC = () => {
  // Initialize API hooks
  const { generatePrompt: apiGenerate, providers, health } = useGromptAPI({
    autoFetchProviders: true,
    autoCheckHealth: true,
    healthCheckInterval: 60000
  });

  const {
    // State
    darkMode,
    currentInput,
    ideas,
    editingId,
    editingText,
    purpose,
    customPurpose,
    maxLength,
    generatedPrompt,
    isGenerating,
    copied,
    outputType,
    agentFramework,
    agentRole,
    agentTools,
    agentProvider,
    mcpServers,
    customMcpServer,
    showOnboarding,
    currentStep,
    showEducational,
    educationalTopic,

    // Setters
    setDarkMode,
    setCurrentInput,
    setEditingText,
    setPurpose,
    setCustomPurpose,
    setMaxLength,
    setOutputType,
    setAgentFramework,
    setAgentRole,
    setAgentTools,
    setAgentProvider,
    setMcpServers,
    setCustomMcpServer,
    setShowEducational,

    // Actions
    addIdea,
    removeIdea,
    startEditing,
    saveEdit,
    cancelEdit,
    generatePrompt,
    copyToClipboard,
    handleFeatureClick,
    startOnboarding,
    nextOnboardingStep,
    showEducation
  } = usePromptCrafter({ apiGenerate });

  // Use dark theme always to match Analyzer design
  const currentTheme = themes.dark;

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <Header
        darkMode={darkMode}
        setDarkMode={setDarkMode}
        currentTheme={currentTheme}
        startOnboarding={startOnboarding}
        showEducation={showEducation}
        providers={providers}
        health={health}
      />

      {/* Onboarding Modal */}
      <OnboardingModal
        showOnboarding={showOnboarding}
        currentStep={currentStep}
        currentTheme={currentTheme}
        nextOnboardingStep={nextOnboardingStep}
      />

      {/* Educational Modal */}
      <EducationalModal
        showEducational={showEducational}
        educationalTopic={educationalTopic}
        currentTheme={currentTheme}
        setShowEducational={setShowEducational}
      />

      {/* Main Content Grid - Following Analyzer layout patterns */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mt-8">

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

          {/* Configuration Panel */}
          <div className="bg-gray-800/50 border border-gray-700/80 rounded-xl p-6 backdrop-blur-sm transition-all duration-300 hover:border-purple-500/50">
            <h2 className="text-xl font-semibold mb-4 text-white">⚙️ Configurações</h2>
            <ConfigurationPanel
              outputType={outputType}
              setOutputType={setOutputType}
              agentFramework={agentFramework}
              setAgentFramework={setAgentFramework}
              agentProvider={agentProvider}
              setAgentProvider={setAgentProvider}
              agentRole={agentRole}
              setAgentRole={setAgentRole}
              agentTools={agentTools}
              setAgentTools={setAgentTools}
              mcpServers={mcpServers}
              setMcpServers={setMcpServers}
              customMcpServer={customMcpServer}
              setCustomMcpServer={setCustomMcpServer}
              purpose={purpose}
              setPurpose={setPurpose}
              customPurpose={customPurpose}
              setCustomPurpose={setCustomPurpose}
              maxLength={maxLength}
              setMaxLength={setMaxLength}
              currentTheme={currentTheme}
              showEducation={showEducation}
              handleFeatureClick={handleFeatureClick}
              providers={providers}
            />
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
            outputType={outputType}
            currentTheme={currentTheme}
            apiGenerateState={apiGenerate}
          />
        </div>

        {/* Right Column: Output Panel */}
        <div className="bg-gray-800/50 border border-gray-700/80 rounded-xl p-6 backdrop-blur-sm transition-all duration-300 hover:border-purple-500/50">
          <OutputPanel
            generatedPrompt={generatedPrompt}
            copyToClipboard={copyToClipboard}
            copied={copied}
            outputType={outputType}
            agentFramework={agentFramework}
            agentProvider={agentProvider}
            maxLength={maxLength}
            mcpServers={mcpServers}
            currentTheme={currentTheme}
            apiGenerateState={apiGenerate}
          />
        </div>
      </div>

      {/* Demo Status Footer */}
      <DemoStatusFooter />
    </div>
  );
};

export default PromptCrafter;