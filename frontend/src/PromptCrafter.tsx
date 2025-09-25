import * as React from 'react';
import DemoStatusFooter from './components/demo/DemoStatusFooter';
import EducationalModal from './components/onboarding/EducationalModal';
import Header from './components/layout/Header';
import OnboardingModal from './components/onboarding/OnboardingModal';
import ScreenNavigation, { Screen } from './components/navigation/ScreenNavigation';
import PromptGenerationScreen from './screens/PromptGenerationScreen';
import AgentManagementScreen from './screens/AgentManagementScreen';
import { themes } from './constants/themes';
import usePromptCrafter from './hooks/usePromptCrafter';
import { useGromptAPI } from './hooks/useGromptAPI';

const PromptCrafter: React.FC = () => {
  // Screen navigation state
  const [currentScreen, setCurrentScreen] = React.useState<Screen>('prompts');

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

      {/* Screen Navigation */}
      <ScreenNavigation
        currentScreen={currentScreen}
        setCurrentScreen={setCurrentScreen}
      />

      {/* Main Content - Screen Based */}
      {currentScreen === 'prompts' && (
        <PromptGenerationScreen
          currentInput={currentInput}
          setCurrentInput={setCurrentInput}
          ideas={ideas}
          editingId={editingId}
          editingText={editingText}
          setEditingText={setEditingText}
          purpose={purpose}
          setPurpose={setPurpose}
          customPurpose={customPurpose}
          setCustomPurpose={setCustomPurpose}
          maxLength={maxLength}
          setMaxLength={setMaxLength}
          generatedPrompt={generatedPrompt}
          isGenerating={isGenerating}
          copied={copied}
          addIdea={addIdea}
          removeIdea={removeIdea}
          startEditing={startEditing}
          saveEdit={saveEdit}
          cancelEdit={cancelEdit}
          generatePrompt={generatePrompt}
          copyToClipboard={copyToClipboard}
          currentTheme={currentTheme}
          apiGenerateState={apiGenerate}
        />
      )}

      {currentScreen === 'agents' && (
        <AgentManagementScreen
          currentInput={currentInput}
          setCurrentInput={setCurrentInput}
          ideas={ideas}
          editingId={editingId}
          editingText={editingText}
          setEditingText={setEditingText}
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
          generatedPrompt={generatedPrompt}
          isGenerating={isGenerating}
          copied={copied}
          addIdea={addIdea}
          removeIdea={removeIdea}
          startEditing={startEditing}
          saveEdit={saveEdit}
          cancelEdit={cancelEdit}
          generatePrompt={generatePrompt}
          copyToClipboard={copyToClipboard}
          showEducation={showEducation}
          handleFeatureClick={handleFeatureClick}
          currentTheme={currentTheme}
          apiGenerateState={apiGenerate}
          providers={providers}
        />
      )}

      {/* Demo Status Footer */}
      <DemoStatusFooter />
    </div>
  );
};

export default PromptCrafter;