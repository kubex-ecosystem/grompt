// import React from 'react';
import React from 'react';
import { purposeKeys } from '../../constants/prompts';
import { useTranslations } from '../../i18n/useTranslations';

interface PurposeSelectorProps {
  purpose: string;
  setPurpose: (value: string) => void;
}
const PurposeSelector: React.FC<PurposeSelectorProps> = React.memo(({ purpose, setPurpose }) => {
  const { t } = useTranslations();
  const purposes = Object.keys(purposeKeys);

  return (
    <div className="mt-6">
      <label htmlFor="purpose-input" className="block text-lg font-medium text-sky-500 dark:text-[#00f0ff] mb-3 light-shadow-sky dark:neon-glow-cyan">{t('purposeLabel')}</label>
      <div className="grid grid-cols-2 sm:grid-cols-3 gap-2 mb-4">
        {purposes.map(p => (
          <button
            key={p}
            onClick={() => setPurpose(p)}
            className={`p-2 rounded-md text-sm font-medium border-2 transition-all duration-200 ${purpose === p
              ? 'bg-emerald-500/20 dark:bg-[#00e676]/20 border-emerald-500 dark:border-[#00e676] text-emerald-700 dark:text-white scale-105'
              : 'bg-slate-100 dark:bg-[#10151b]/50 border-transparent hover:border-emerald-500/50 dark:hover:border-[#00e676]/50 text-slate-600 dark:text-slate-300'
              }`}
          >
            {t(purposeKeys[p])}
          </button>
        ))}
      </div>
      <input
        id="purpose-input"
        type="text"
        value={purpose}
        onChange={(e) => setPurpose(e.target.value)}
        list="purposes-list"
        placeholder={t('customPurposePlaceholder')}
        className="w-full bg-white dark:bg-[#10151b] border-2 border-emerald-500/50 dark:border-[#00e676]/30 rounded-md p-3 focus:outline-none focus:ring-2 focus:ring-emerald-500 dark:focus:ring-[#00e676] focus:border-emerald-500 dark:focus:border-[#00e676] transition-all duration-300 placeholder:text-slate-400 dark:placeholder:text-[#90a4ae]/50 text-emerald-600 dark:text-[#00e676] font-semibold"
      />
      <datalist id="purposes-list">
        {purposes.map(p => <option key={p} value={p} />)}
      </datalist>
    </div>
  );
});

export default PurposeSelector;
