
import { createProvider } from "../../../.notes/bkp";
import { extractDiffFenced, pack } from "../utils/lookatni";

/** Seleciona provider a partir do env sem travar UI */
const llm = createProvider();

/** Monta o payload LookAtni a partir do VFS/FS do app. Troca por fetch('/snapshot') se preferir. */
export async function buildLookatniBlob(grab: (p: string) => string, include: string[]): Promise<string> {
  const files = include.map(p => ({ path: p, content: grab(p) }));
  return pack(files);
}

/** Pede um UNIFIED DIFF pequeno. */
export async function requestUnifiedDiff(blob: string, task: string): Promise<string> {
  const system = [
    "You are a rigorous code refactorer.",
    "Return ONLY one unified diff fenced with diff.",
    "Small, reviewable changes. Preserve behavior and paths (a/ b/ headers)."
  ].join("\n");

  const user = [
    "Context (LookAtni):",
    blob,
    "",
    "Task:",
    task,
    "",
    "Constraints:",
    "- Keep exports stable; no broad rewrites.",
    "- Prefer extracting helpers, removing duplication, adding small tests.",
    "- If you add files, include proper diff headers."
  ].join("\n");

  const out = await llm.generate({ system, user, temperature: 0.15, maxTokens: 1800 });
  return extractDiffFenced(out.text);
}

