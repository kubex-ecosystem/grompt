
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
}
