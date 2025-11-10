import React from 'react';
import { useTranslations } from '../../i18n/useTranslations';

const Footer: React.FC = () => {
  const { t } = useTranslations();

  return (
    <footer className="space-y-3 py-8 text-center text-xs text-[#64748b] dark:text-[#94a3b8]">
      <div className="flex justify-center items-center gap-6 text-xs">
        {/* <a
          href="https://kubex.world"
          target="_blank"
          rel="noopener noreferrer"
          className="transition-colors duration-200 hover:text-[#06b6d4] dark:hover:text-[#38cde4]"
        >
          Kubex Ecosystem
        </a> */}
        <a
          href="https://github.com/kubex-ecosystem/grompt"
          target="_blank"
          rel="noopener noreferrer"
          className="transition-colors duration-200 hover:text-[#06b6d4] dark:hover:text-[#38cde4]"
        >
          GitHub
        </a>
        {/* <a
          href="/humans.txt"
          target="_blank"
          rel="noopener noreferrer"
          className="transition-colors duration-200 hover:text-[#06b6d4] dark:hover:text-[#38cde4]"
        >
          Humans.txt
        </a> */}
        {/* <a
          href="/.well-known/security.txt"
          target="_blank"
          rel="noopener noreferrer"
          className="transition-colors duration-200 hover:text-[#06b6d4] dark:hover:text-[#38cde4]"
        >
          Security
        </a> */}
      </div>
      <div>
        <p>{t('poweredBy')}</p>
        <p className="mt-1 font-orbitron tracking-wider">{t('motto')}</p>
      </div>
      <div className="text-[10px] opacity-70">
        <p>Grompt v2.0 • Made with ❤️ in Brasil • Open Source</p>
      </div>
    </footer>
  );
};

export default Footer;
