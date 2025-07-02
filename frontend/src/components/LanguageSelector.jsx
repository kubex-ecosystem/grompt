import React from 'react';
import { useTranslation } from 'react-i18next';
import { Globe } from 'lucide-react';

const LanguageSelector = ({ currentTheme }) => {
  const { i18n } = useTranslation();

  const languages = [
    { code: 'pt-BR', name: 'Português', flag: '🇧🇷' },
    { code: 'en-US', name: 'English', flag: '🇺🇸' }
  ];

  const handleLanguageChange = (langCode) => {
    i18n.changeLanguage(langCode);
  };

  return (
    <div className="relative group">
      <button
        className={`flex items-center gap-2 px-3 py-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
      >
        <Globe size={16} />
        <span className="text-sm">
          {languages.find(lang => lang.code === i18n.language)?.flag || '🌍'}
        </span>
      </button>
      
      <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 z-50">
        {languages.map((language) => (
          <button
            key={language.code}
            onClick={() => handleLanguageChange(language.code)}
            className={`w-full flex items-center gap-3 px-4 py-3 text-left hover:bg-gray-100 dark:hover:bg-gray-700 first:rounded-t-lg last:rounded-b-lg transition-colors ${
              i18n.language === language.code ? 'bg-blue-50 dark:bg-blue-900 text-blue-600 dark:text-blue-400' : 'text-gray-700 dark:text-gray-300'
            }`}
          >
            <span className="text-lg">{language.flag}</span>
            <span className="text-sm font-medium">{language.name}</span>
          </button>
        ))}
      </div>
    </div>
  );
};

export default LanguageSelector;
