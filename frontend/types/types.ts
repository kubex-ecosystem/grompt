export interface PromptHistoryItem {
  id: number;
  prompt: string;
  timestamp: number;
  inputs: {
    ideas: string[];
    purpose: string;
  };
}
