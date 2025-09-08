export type Theme = 'light' | 'dark';
export type Language = 'en' | 'es' | 'zh' | 'pt';

export interface Idea {
  id: string;
  text: string;
}

export interface HistoryItem {
  id: string;
  prompt: string;
  purpose: string;
  ideas: Idea[];
  timestamp: number;
  inputTokens?: number;
  outputTokens?: number;
  totalTokens?: number;
}

export interface Draft {
  ideas: Idea[];
  purpose: string;
}
