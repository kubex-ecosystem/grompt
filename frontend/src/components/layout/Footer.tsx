import React from 'react';
import { useTranslations } from '../../i18n/useTranslations';

const Footer: React.FC = () => {
  const { t } = useTranslations();

  return (
    <footer className="text-center mt-12 text-slate-500 dark:text-[#90a4ae] text-xs space-y-3">
      <div className="flex justify-center items-center gap-6 text-xs">
        <a
          href="https://kubex.world"
          target="_blank"
          rel="noopener noreferrer"
          className="hover:text-sky-500 dark:hover:text-[#00f0ff] transition-colors duration-200"
        >
          Kubex Ecosystem
        </a>
        <a
          href="https://github.com/kubex-ecosystem/gemx/grompt"
          target="_blank"
          rel="noopener noreferrer"
          className="hover:text-sky-500 dark:hover:text-[#00f0ff] transition-colors duration-200"
        >
          GitHub
        </a>
        <a
          href="/humans.txt"
          target="_blank"
          rel="noopener noreferrer"
          className="hover:text-sky-500 dark:hover:text-[#00f0ff] transition-colors duration-200"
        >
          Humans.txt
        </a>
        <a
          href="/.well-known/security.txt"
          target="_blank"
          rel="noopener noreferrer"
          className="hover:text-sky-500 dark:hover:text-[#00f0ff] transition-colors duration-200"
        >
          Security
        </a>
      </div>
      <div>
        <p>{t('poweredBy')}</p>
        <p className="mt-1 font-orbitron tracking-wider">{t('motto')}</p>
      </div>
      <div className="text-[10px] opacity-70">
        <p>Grompt v2.0 • Made with ❤️ in Brasil • Open Source • No Lock-in</p>
      </div>
    </footer>
  );
};

export default Footer;
