import { useContext } from 'react';
import { LanguageContext, LanguageContextType } from '../context/LanguageContext';

export const useTranslations = (): Omit<LanguageContextType, 'setLanguage'> => {
  const { language, t } = useContext(LanguageContext);
  return { language, t };
};
