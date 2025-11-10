import type { Agent, AgentsGenerationResult, StoredAgent } from '@/types';

const BASE_PATH = '/api/v1/agents';

async function parseJSON<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || `Request failed with status ${response.status}`);
  }
  return response.json() as Promise<T>;
}

async function parseText(response: Response): Promise<string> {
  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || `Request failed with status ${response.status}`);
  }
  return response.text();
}

export const agentsService = {
  async list(): Promise<StoredAgent[]> {
    const response = await fetch(BASE_PATH, {
      method: 'GET',
      headers: {
        'Accept': 'application/json',
      },
      credentials: 'same-origin',
    });
    return parseJSON<StoredAgent[]>(response);
  },

  async create(agent: Agent): Promise<StoredAgent> {
    const response = await fetch(BASE_PATH, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
      body: JSON.stringify(agent),
      credentials: 'same-origin',
    });
    return parseJSON<StoredAgent>(response);
  },

  async update(id: number, agent: Agent): Promise<StoredAgent> {
    const response = await fetch(`${BASE_PATH}/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
      body: JSON.stringify(agent),
      credentials: 'same-origin',
    });
    return parseJSON<StoredAgent>(response);
  },

  async remove(id: number): Promise<void> {
    const response = await fetch(`${BASE_PATH}/${id}`, {
      method: 'DELETE',
      credentials: 'same-origin',
    });
    if (!response.ok && response.status !== 204) {
      const message = await response.text();
      throw new Error(message || `Request failed with status ${response.status}`);
    }
  },

  async generate(requirements: string): Promise<AgentsGenerationResult> {
    const response = await fetch(`${BASE_PATH}/generate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
      body: JSON.stringify({ requirements }),
      credentials: 'same-origin',
    });
    return parseJSON<AgentsGenerationResult>(response);
  },

  async exportMarkdown(): Promise<string> {
    const response = await fetch(`${BASE_PATH}.md`, {
      method: 'GET',
      headers: {
        'Accept': 'text/markdown',
      },
      credentials: 'same-origin',
    });
    return parseText(response);
  },
};

export type AgentsService = typeof agentsService;
