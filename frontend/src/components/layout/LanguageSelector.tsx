import { ChevronDown } from 'lucide-react';
import React, { useContext, useEffect, useRef, useState } from 'react';
import { LanguageContext } from '../../context/LanguageContext';
import { useTranslations } from '../../i18n/useTranslations';
import { Language } from '../../types';

const languageOptions: { code: Language; flag: React.ReactNode; }[] = [
  { code: 'en', flag: <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 72 48"><path fill="#bd3d44" d="M0 0h640v480H0" /><path d="M0 55.3h640M0 129h640M0 203h640M0 277h640M0 351h640M0 425h640" /><path fill="#192f5d" d="M0 0h364.8v258.5H0" /><path fill="#fff" d="m14 0 9 27L0 10h28L5 27z" /><path fill="none" d="m0 0 16 11h61 61 61 61 60L47 37h61 61 60 61L16 63h61 61 61 61 60L47 89h61 61 60 61L16 115h61 61 61 61 60L47 141h61 61 60 61L16 166h61 61 61 61 60L47 192h61 61 60 61L16 218h61 61 61 61 60z" /></svg> },
  { code: 'es', flag: <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 72 48"><path fill="#ff4b55" d="M0 0h72v48H0z" /><path fill="#ffd500" d="M0 12h72v24H0z" /><path fill="#ff4b55" d="M30 18h2v12h-2zm-2 0h-2v12h2zm-2 1h-1v10h1v-4h2v-2h-2v-4z" /><path fill="#ffd500" d="M30 18h-2v12h2v-5h-1v-2h1v-5z" /><path fill="#ff4b55" d="M31 23h-6v1h6v-1z" /></svg> },
  { code: 'pt', flag: <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 72 48"><path fill="#009b3a" d="M0 0h72v48H0z" /><path fill="#ffdf00" d="m36 6 22 18-22 18-22-18z" /><path fill="#002776" d="M36 33a9 9 0 1 0 0-18 9 9 0 0 0 0 18z" /><path fill="#fff" d="M26 23.5a16 16 0 0 1 20 1v-2a18 18 0 0 0-20-1z" /></svg> },
  { code: 'zh', flag: <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 72 48"><path fill="#ff4b55" d="M0 0h72v48H0z" /><path fill="#ffd500" d="m18 9 2.2 6.8h7.1L21 20l2.2 6.8L18 22.1l-5.3 4.7L15 20l-6.2-4.2h7.1zm12.5 1.5 1 .8-.5 1.3-1.2-.3.1-1.4zm-1 16 .3-1.3-1.2.3.9 1zm6.5-13.5.8-1-1.3.6.5.4zm-4 14-1.3-.5.5-1.3.8 1zm3.5-1.5-.7 1.1.9.8.3-1.4-.5-.5z" /></svg> },
];

const LanguageSelector: React.FC = () => {
  const { language, setLanguage } = useContext(LanguageContext);
  const { t } = useTranslations();
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const currentLanguage = languageOptions.find(lang => lang.code === language) || languageOptions[0];

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleSelectLanguage = (langCode: Language) => {
    setLanguage(langCode);
    setIsOpen(false);
  };

  return (
    <div className="relative" ref={dropdownRef}>
      <button
        type="button"
        title={t('lang_' + currentLanguage.code)}
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center gap-2 bg-slate-200/50 dark:bg-[#10151b] p-2 rounded-full transition-colors duration-200 hover:bg-slate-300/50 dark:hover:bg-slate-800"
        aria-haspopup="true"
        aria-expanded={isOpen}
        aria-label="Select language"
      >
        <div className="w-6 h-6 rounded-full overflow-hidden">{currentLanguage.flag}</div>
        <span className="font-semibold text-sm uppercase text-slate-700 dark:text-slate-300">{currentLanguage.code}</span>
        <ChevronDown size={16} className={`text-slate-600 dark:text-slate-400 transition-transform duration-300 ${isOpen ? 'rotate-180' : ''}`} />
      </button>
      {isOpen && (
        <div
          className="absolute right-0 top-full mt-2 w-48 bg-white/70 dark:bg-[#10151b]/70 backdrop-blur-md rounded-lg shadow-2xl border border-slate-200 dark:border-sky-500/20 py-1 animate-fade-in z-50"
          role="menu"
        >
          {languageOptions.map(lang => (
            <button
              type="button"
              title={t('lang_' + lang.code)}
              key={lang.code}
              onClick={() => handleSelectLanguage(lang.code)}
              className="w-full flex items-center gap-3 px-4 py-2 text-left text-sm text-slate-700 dark:text-slate-200 hover:bg-sky-100 dark:hover:bg-sky-500/20"
              role="menuitem"
              disabled={language === lang.code}
            >
              <div className="w-6 h-6 rounded-full overflow-hidden flex-shrink-0">{lang.flag}</div>
              <span>{t(`lang_${lang.code}`)}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  );
};

export default LanguageSelector;
