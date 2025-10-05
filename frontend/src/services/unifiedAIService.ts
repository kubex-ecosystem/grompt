// Unified AI Service for Grompt frontend
// Communicates with backend's unified API endpoints

import { Idea } from '../../types';
import { configService, type ProviderInfo } from './configService';

export interface UnifiedRequest {
  lang?: string;
  purpose?: string;
  purpose_type?: string;
  ideas?: string[];
  prompt?: string;
  max_tokens?: number;
  model?: string;
  provider?: string;
}

export interface UnifiedResponse {
  response: string;
  provider: string;
  model: string;
  usage?: {
    prompt_tokens?: number;
    completion_tokens?: number;
    total_tokens?: number;
    estimated_cost?: number;
  };
}

export interface GenerationResult {
  prompt: string;
  provider: string;
  model: string;
  usageMetadata?: {
    promptTokenCount?: number;
    candidatesTokenCount?: number;
    totalTokenCount?: number;
  };
}

class UnifiedAIService {
  private baseUrl = '';

  /**
   * Generate a structured prompt using backend's unified API
   */
  async generateStructuredPrompt(
    ideas: Idea[],
    purpose: string,
    provider?: string,
    model?: string
  ): Promise<GenerationResult> {
    try {
      // Get config to determine the best provider to use
      const config = await configService.getConfig();

      // Use provided provider or get default available one
      const targetProvider = provider || config.default_provider;

      if (!targetProvider) {
        throw new Error('No AI providers are configured. Please configure an API key.');
      }

      // Check if provider is available
      const providerInfo = config.providers[targetProvider];
      if (!providerInfo || !providerInfo.available) {
        throw new Error(`Provider ${targetProvider} is not available or configured.`);
      }

      // Convert ideas to strings
      const ideaTexts = ideas.map(idea => idea.text);

      // Prepare request
      const request: UnifiedRequest = {
        ideas: ideaTexts,
        purpose: purpose,
        provider: targetProvider,
        model: model || providerInfo.models[0], // Use first available model if none specified
        lang: 'português',
        max_tokens: 5000,
      };

      // Call unified API
      const response = await fetch('/api/unified', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`API Error ${response.status}: ${errorText}`);
      }

      const data: UnifiedResponse = await response.json();

      // Transform response to match expected format
      return {
        prompt: data.response,
        provider: data.provider,
        model: data.model,
        usageMetadata: {
          promptTokenCount: data.usage?.prompt_tokens,
          candidatesTokenCount: data.usage?.completion_tokens,
          totalTokenCount: data.usage?.total_tokens,
        },
      };
    } catch (error) {
      console.error('Failed to generate structured prompt:', error);

      // Fallback to demo mode
      return this.generateDemoPrompt(ideas, purpose);
    }
  }

  /**
   * Generate a direct prompt (skip prompt engineering)
   */
  async generateDirectPrompt(
    prompt: string,
    provider?: string,
    model?: string,
    maxTokens = 1000
  ): Promise<UnifiedResponse> {
    try {
      const config = await configService.getConfig();
      const targetProvider = provider || config.default_provider;

      if (!targetProvider) {
        throw new Error('No AI providers are configured. Please configure an API key.');
      }

      const request: UnifiedRequest = {
        prompt: prompt,
        provider: targetProvider,
        model: model,
        max_tokens: maxTokens,
      };

      const response = await fetch('/api/unified', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`API Error ${response.status}: ${errorText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Failed to generate direct prompt:', error);
      throw error;
    }
  }

  /**
   * Get available providers and their status
   */
  async getAvailableProviders(): Promise<ProviderInfo[]> {
    return configService.getAvailableProviders();
  }

  /**
   * Check if running in demo mode
   */
  async isDemoMode(): Promise<boolean> {
    return configService.isDemoMode();
  }

  /**
   * Test a specific provider
   */
  async testProvider(provider: string): Promise<{ available: boolean; message: string }> {
    try {
      const response = await fetch(`/api/test?provider=${provider}`, {
        method: 'GET',
      });

      if (!response.ok) {
        return {
          available: false,
          message: `HTTP ${response.status}: ${response.statusText}`,
        };
      }

      const data = await response.json();
      return {
        available: data.available || false,
        message: data.message || 'Unknown status',
      };
    } catch (error) {
      return {
        available: false,
        message: `Connection error: ${error}`,
      };
    }
  }

  /**
   * Fallback demo prompt generation
   */
  private generateDemoPrompt(ideas: Idea[], purpose: string): GenerationResult {
    const ideasText = ideas.map((idea, index) => `- ${idea.text}`).join('\n');

    const demoPrompt = `# ${purpose} Expert Assistant

## Primary Objective
Transform the provided ideas into actionable ${purpose.toLowerCase()} solutions following Kubex principles of radical simplicity and modularity.

## User Requirements
${ideasText}

## Task Instructions
You are an expert ${purpose.toLowerCase()} specialist. Based on the requirements above, provide a comprehensive solution that:

### Key Requirements:
- Follow KUBEX principles: Radical Simplicity, Modularity, No Cages
- Use clear, anti-jargon language
- Provide modular, reusable components
- Ensure outputs are platform-agnostic

### Expected Output Format:
- Use Markdown for clear structure
- Include code examples when applicable
- Provide step-by-step instructions
- Add relevant comments and documentation

### Constraints:
- Avoid vendor lock-in solutions
- Keep complexity minimal
- Focus on practical, implementable solutions
- Use open standards and formats

## Context
This prompt was generated using Grompt, part of the Kubex Ecosystem, following principles of radical simplicity and avoiding technological cages.

---
*Generated in demo mode - Connect your AI provider API key for enhanced AI-powered prompts*`;

    // Simulate token usage for demo
    const estimatedTokens = Math.floor(demoPrompt.length / 4); // Rough estimation

    return {
      prompt: demoPrompt,
      provider: 'demo',
      model: 'demo-model',
      usageMetadata: {
        promptTokenCount: Math.floor(estimatedTokens * 0.3),
        candidatesTokenCount: Math.floor(estimatedTokens * 0.7),
        totalTokenCount: estimatedTokens,
      },
    };
  }
}

// Export singleton instance
export const unifiedAIService = new UnifiedAIService();

// Backwards compatibility - export the main function with same signature as geminiService
export const generateStructuredPrompt = unifiedAIService.generateStructuredPrompt.bind(unifiedAIService);
