import DemoMode from '@/config/DemoMode';

const EducationalModal = ({
  showEducational,
  educationalTopic,
  currentTheme,
  setShowEducational
}) => {
  if (!showEducational || !educationalTopic) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
      <div className={`${currentTheme.cardBg} rounded-xl p-6 max-w-lg border ${currentTheme.border} shadow-xl`}>
        <h3 className="text-xl font-bold mb-4">
          {DemoMode.education[educationalTopic].title}
        </h3>
        <p className={`${currentTheme.textSecondary} mb-4`}>
          {DemoMode.education[educationalTopic].description}
        </p>
        <div className="mb-6">
          <h4 className="font-semibold mb-2">Benefícios:</h4>
          <ul className="space-y-1">
            {DemoMode.education[educationalTopic].benefits.map((benefit, index) => (
              <li key={index} className={currentTheme.textSecondary}>
                {benefit}
              </li>
            ))}
          </ul>
        </div>
        <button
          onClick={() => setShowEducational(false)}
          className={`px-4 py-2 rounded-lg ${currentTheme.button} w-full`}
        >
          Entendi!
        </button>
      </div>
    </div>
  );
};

export default EducationalModal;
