import OnboardingSteps from '../constants/onboardingSteps.js';

const OnboardingModal = ({
  showOnboarding,
  currentStep,
  currentTheme,
  nextOnboardingStep
}) => {
  if (!showOnboarding) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
      <div className={`${currentTheme.cardBg} rounded-xl p-6 max-w-md border ${currentTheme.border} shadow-xl`}>
        <h3 className="text-xl font-bold mb-4">
          {OnboardingSteps[currentStep].title}
        </h3>
        <p className={`${currentTheme.textSecondary} mb-6`}>
          {OnboardingSteps[currentStep].content}
        </p>
        <div className="flex justify-between">
          <span className="text-sm text-gray-500">
            {currentStep + 1} de {OnboardingSteps.length}
          </span>
          <button
            onClick={nextOnboardingStep}
            className={`px-4 py-2 rounded-lg ${currentTheme.button}`}
          >
            {currentStep < OnboardingSteps.length - 1 ? 'Próximo' : 'Finalizar'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default OnboardingModal;
