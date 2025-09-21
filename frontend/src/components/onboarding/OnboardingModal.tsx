import React from 'react';
import { X } from 'lucide-react';
import onboardingSteps from '../../constants/onboardingSteps';
import { Theme } from '../../constants/themes';

interface OnboardingModalProps {
  showOnboarding: boolean;
  currentStep: number;
  currentTheme: Theme;
  nextOnboardingStep: () => void;
}

const OnboardingModal: React.FC<OnboardingModalProps> = ({
  showOnboarding,
  currentStep,
  currentTheme,
  nextOnboardingStep
}) => {
  if (!showOnboarding || currentStep >= onboardingSteps.length) return null;

  const step = onboardingSteps[currentStep];

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="bg-gray-800 border border-gray-700 rounded-xl w-full max-w-md p-6 shadow-2xl">
        <div className="flex justify-between items-start mb-4">
          <h3 className="text-lg font-semibold text-white">{step.title}</h3>
          <button
            onClick={nextOnboardingStep}
            className="text-gray-400 hover:text-white"
          >
            <X size={20} />
          </button>
        </div>

        <p className="text-gray-300 mb-6">{step.content}</p>

        <div className="flex justify-between items-center">
          <span className="text-xs text-gray-500">
            {currentStep + 1} de {onboardingSteps.length}
          </span>
          <button
            onClick={nextOnboardingStep}
            className="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors"
          >
            {currentStep === onboardingSteps.length - 1 ? 'Finalizar' : 'Próximo'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default OnboardingModal;