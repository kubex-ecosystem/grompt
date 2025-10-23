<<<<<<< HEAD:frontend/types/index.ts
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

export * from './agents';
=======
export * from "@/core/llm/providers/anthropic";
export * from "@/core/llm/providers/gemini";
export * from "@/core/llm/providers/openai";
export * from "@/core/llm/wrapper/MultiAIWrapper";
export * from "@/types/types";

>>>>>>> origin/main:frontend/src/types/index.ts
