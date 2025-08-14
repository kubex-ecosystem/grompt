import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

// Importar as traduções
import enUS from './locales/en-US.json';
import ptBR from './locales/pt-BR.json';

const resources = {
  'pt-BR': {
    translation: ptBR
  },
  'en-US': {
    translation: enUS
  }
};

// Verificar se está no browser
const isClient = typeof window !== 'undefined';

const i18nConfig = {
  resources,
  lng: 'pt-BR', // idioma padrão
  fallbackLng: 'pt-BR',

  interpolation: {
    escapeValue: false, // React já faz escape
  },

  // Só usar detecção de idioma no cliente
  ...(isClient && {
    detection: {
      order: ['localStorage', 'navigator', 'htmlTag'],
      lookupLocalStorage: 'i18nextLng',
      caches: ['localStorage'],
    }
  })
};

if (!i18n.isInitialized) {
  if (isClient) {
    // No cliente, usar o detector de idioma
    import('i18next-browser-languagedetector').then((LanguageDetector) => {
      i18n
        .use(LanguageDetector.default)
        .use(initReactI18next)
        .init(i18nConfig);
    });
  } else {
    // No servidor, inicializar sem detector
    i18n
      .use(initReactI18next)
      .init(i18nConfig);
  }
}

export default i18n;
