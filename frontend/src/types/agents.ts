export interface Agent {
  Title: string;
  Role: string;
  Skills: string[];
  Restrictions: string[];
  PromptExample: string;
}

export interface StoredAgent extends Agent {
  ID: number;
}

export interface AgentsGenerationResult {
  agents: Agent[];
  markdown: string;
}
