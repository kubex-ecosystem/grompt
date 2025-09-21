import ConfigurationPanel from './components/ConfigurationPanel.jsx';
import DemoStatusFooter from './components/DemoStatusFooter.jsx';
import EducationalModal from './components/EducationalModal.jsx';
import Header from './components/Header.jsx';
import IdeasInput from './components/IdeasInput.jsx';
import IdeasList from './components/IdeasList.jsx';
import OnboardingModal from './components/OnboardingModal.jsx';
import OutputPanel from './components/OutputPanel.jsx';
import { themes } from './constants/themes.js';
import usePromptCrafter from './hooks/usePromptCrafter.js';

const PromptCrafter = () => {
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
  } = usePromptCrafter();

  const currentTheme = darkMode ? themes.dark : themes.light;

  return (
    <div className={`min-h-screen ${currentTheme.bg} ${currentTheme.text} p-4 transition-colors duration-300`}>
      <div className="max-w-7xl mx-auto">

        {/* Header */}
        <Header
          darkMode={darkMode}
          setDarkMode={setDarkMode}
          currentTheme={currentTheme}
          startOnboarding={startOnboarding}
          showEducation={showEducation}
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

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">

          {/* Input Section */}
          <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg`}>
            <IdeasInput
              currentInput={currentInput}
              setCurrentInput={setCurrentInput}
              addIdea={addIdea}
              currentTheme={currentTheme}
            />

            {/* Configuration */}
            <div className="mt-6">
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
              />
            </div>
          </div>

          {/* Ideas List */}
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
          />

          {/* Output Panel */}
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
          />
        </div>

        {/* Demo Status Footer */}
        <DemoStatusFooter />
      </div>
    </div>
  );
};

export default PromptCrafter;
