//(trecho de integração do botão “Evolve”)
import { AlertTriangle } from 'lucide-react';
import * as React from "react";
import { useEffect, useState } from 'react';
import PromptCrafter from './components/features/PromptCrafter';
import Footer from './components/layout/Footer';
import Header from './components/layout/Header';
import { LanguageContext } from "./context/LanguageContext";
import { useLanguage } from './hooks/useLanguage';
import { useTheme } from './hooks/useTheme';
import { useAnalytics } from './services/analytics';
import { buildLookatniBlob, requestUnifiedDiff } from "./services/evolver";
import { initStorage } from './services/storageService';

// defina a lista mínima do teu app:
const EVOLVE_INCLUDE = [
  "index.html",
  "src/main.tsx",
  "src/App.tsx",
  "src/services/geminiService.ts",   // se existir
  "src/styles.css",
  "vite.config.ts",
  "package.json",
  "tsconfig.json"
];

// Função que “pega” arquivos do runtime do teu ambiente.
// Se não tiver __VFS__, substitui por fetch('/snapshot/<path>')
function grabFile(path: string): string {
  // @ts-expect-error ambiente do AI Studio/preview pode expor isso; se não, troca por fetch()
  const b = window.__VFS__?.[path];
  if (typeof b === "string") return b;
  throw new Error(`VFS missing: ${path} (troque grabFile para fetch('/snapshot/${path}'))`);
}

async function handleVirtuousEvolve() {
  setEvolveStep("generating_prompt");
  const selfPrompt = [
    "Perform a tiny but valuable refactor improving cohesion/clarity.",
    "Prefer lazy-chunk UI pieces and remove any dead code.",
    "Return ONLY a unified diff fenced with ```diff."
  ].join("\n");

  const blob = await buildLookatniBlob(grabFile, EVOLVE_INCLUDE);
  setEvolveStep("refactoring");
  const diff = await requestUnifiedDiff(blob, selfPrompt);
  setEvolveDiff(diff);
  setEvolveStep("done");
}


const App: React.FC = () => {
  const [evolveStep, setEvolveStep] = useState<"idle" | "generating_prompt" | "refactoring" | "done">("idle");
  const [evolveDiff, setEvolveDiff] = useState<string | null>(null);
  const [theme, toggleTheme] = useTheme();
  const { language, setLanguage, t } = useLanguage();
  const [isApiKeyMissing, setIsApiKeyMissing] = useState(false);
  const languageContextValue = { language, setLanguage, t };

  // Initialize analytics
  useAnalytics();

  // Initialize the storage service on app load
  useEffect(() => {
    initStorage();
    // Check for the API key availability on mount
    if (!process.env.API_KEY) {
      console.log("Running in demo mode - API key not configured");
      setIsApiKeyMissing(true);
    }
  }, []);

  return (
    <LanguageContext value={languageContextValue}>
      <div className="min-h-screen text-slate-800 dark:text-[#e0f7fa] font-plex-mono p-4 sm:p-6 lg:p-8">
        <div className="max-w-7xl mx-auto">
          <Header theme={theme} toggleTheme={toggleTheme} />
          {isApiKeyMissing && (
            <div className="bg-blue-100 border-l-4 border-blue-500 text-blue-700 p-4 rounded-md mb-6 dark:bg-blue-900/30 dark:text-blue-300 dark:border-blue-600 flex items-start gap-3" role="alert">
              <AlertTriangle className="h-6 w-6 text-blue-500 dark:text-blue-400 flex-shrink-0 mt-0.5" />
              <div>
                <p className="font-bold">{t('apiKeyMissingTitle')}</p>
                <p>{t('apiKeyMissingText')}</p>
              </div>
            </div>
          )}
          <main>
            <PromptCrafter theme={theme} isApiKeyMissing={isApiKeyMissing} />
          </main>
          <Footer />
        </div>
      </div>
    </LanguageContext>
  );
};

export default App;
function setEvolveStep(arg0: string) {
  throw new Error("Function not implemented.");
}

function setEvolveDiff(diff: string) {
  throw new Error("Function not implemented.");
}

