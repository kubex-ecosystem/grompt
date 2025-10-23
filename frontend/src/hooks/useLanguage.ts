import { useState, useEffect } from 'react';
import { Language } from '../types';
import { translations } from '../i18n/translations';

export const useLanguage = () => {
  const [language, setLanguage] = useState<Language>('en');

  useEffect(() => {
    try {
      const savedLang = localStorage.getItem('language') as Language | null;
      const browserLang = navigator.language.split('-')[0] as Language;
      const initialLang = savedLang || (['en', 'es', 'zh', 'pt'].includes(browserLang) ? browserLang : 'en');
      setLanguage(initialLang);
    } catch (error) {
        console.warn('Could not access localStorage to get language. Using default.', error);
        // Fallback to browser language if localStorage is not available
        const browserLang = navigator.language.split('-')[0] as Language;
        setLanguage((['en', 'es', 'zh', 'pt'].includes(browserLang) ? browserLang : 'en'));
    }
  }, []);

  useEffect(() => {
    try {
      localStorage.setItem('language', language);
    } catch (error) {
      console.warn('Could not access localStorage to save language.', error);
    }
  }, [language]);

  const t = (key: string, params?: Record<string, string>): string => {
    let translation = translations[language][key] || translations.en[key] || key;
    if (params) {
      Object.keys(params).forEach(paramKey => {
        translation = translation.replace(`{${paramKey}}`, params[paramKey]);
      });
    }
    return translation;
  };

  return { language, setLanguage, t };
};
