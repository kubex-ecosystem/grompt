import React from 'react';
import { X } from 'lucide-react';
import { DemoMode } from '../../config/demoMode';
import { Theme } from '../../constants/themes';

interface EducationalModalProps {
  showEducational: boolean;
  educationalTopic: string | null;
  currentTheme: Theme;
  setShowEducational: (value: boolean) => void;
}

const EducationalModal: React.FC<EducationalModalProps> = ({
  showEducational,
  educationalTopic,
  currentTheme,
  setShowEducational
}) => {
  if (!showEducational || !educationalTopic) return null;

  const content = DemoMode.education[educationalTopic];
  if (!content) return null;

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="bg-gray-800 border border-gray-700 rounded-xl w-full max-w-2xl p-6 shadow-2xl">
        <div className="flex justify-between items-start mb-4">
          <h3 className="text-xl font-semibold text-white">{content.title}</h3>
          <button
            onClick={() => setShowEducational(false)}
            className="text-gray-400 hover:text-white"
          >
            <X size={20} />
          </button>
        </div>

        <p className="text-gray-300 mb-6">{content.description}</p>

        <div className="space-y-2 mb-6">
          <h4 className="text-lg font-medium text-white">Benefícios:</h4>
          <ul className="space-y-2">
            {content.benefits.map((benefit, index) => (
              <li key={index} className="text-gray-300 flex items-start gap-2">
                <span className="text-purple-400 mt-1">•</span>
                {benefit}
              </li>
            ))}
          </ul>
        </div>

        <div className="flex justify-end">
          <button
            onClick={() => setShowEducational(false)}
            className="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors"
          >
            Entendi
          </button>
        </div>
      </div>
    </div>
  );
};

export default EducationalModal;