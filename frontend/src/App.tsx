//(trecho de integração do botão “Evolve”)
import { buildLookatniBlob, requestUnifiedDiff } from "./services/evolver";

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
