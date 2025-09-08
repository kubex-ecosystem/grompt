// FIX: Removed self-import of Idea which conflicts with the interface declaration below.
export interface Idea {
  id: string;
  text: string;
}

export interface HistoryItem {
  id:string;
  prompt: string;
  purpose: string;
  ideas: Idea[];
  timestamp: number;
  inputTokens?: number;
  outputTokens?: number;
}