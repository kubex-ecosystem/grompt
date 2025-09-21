
export type FileUnit = { path: string; content: string };

export function pack(files: FileUnit[]): string {
  return files.map(f => `/// ${f.path} ///\n${f.content}\n`).join("\n");
}

export function extractDiffFenced(text: string): string {
  const m = text.match(/```diff([\s\S]*?)```/);
  return (m ? m[1] : text).trim();
}
