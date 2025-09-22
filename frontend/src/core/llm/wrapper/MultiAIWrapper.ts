import { AnthropicProvider } from "../providers/anthropic";
import { GeminiProvider } from "../providers/gemini";
import { OpenAIProvider } from "../providers/openai";
import {
  AIModel,
  AIProvider,
  AIResponse,
  CraftPromptParams,
  GenerateContentParams,
  MultiAIConfig,
  RefactorCodeParams
} from "../../../types/types";

export class MultiAIWrapper {
  private providers: Map<AIProvider, { generateContent: any; streamContent?: any }>;
  private config: MultiAIConfig;
  private cache: Map<string, AIResponse>;

  private readonly SYSTEM_INSTRUCTIONS: Record<string, string> = {
    'Code Generation': 'You are an expert software developer. Generate clean, efficient, and well-documented code based on the requirements provided.',
    'Code Refactoring': 'You are an expert code reviewer. Refactor the provided code to improve readability, performance, and maintainability while preserving functionality.',
    'General Summarization': 'You are a helpful assistant that creates clear and concise summaries.',
    'Documentation': 'You are a technical writer. Create comprehensive and clear documentation.',
    'Bug Fixing': 'You are a debugging expert. Analyze the code and fix any bugs or issues found.',
    'Code Review': 'You are an experienced code reviewer. Provide constructive feedback on code quality, best practices, and potential improvements.',
    'API Design': 'You are an API design expert. Create well-structured, RESTful APIs with proper documentation.',
    'Database Design': 'You are a database architect. Design efficient, normalized database schemas.',
    'Testing': 'You are a testing expert. Write comprehensive unit tests and integration tests.',
    'Security Audit': 'You are a security expert. Identify potential security vulnerabilities and recommend fixes.',
  };

  constructor(config: MultiAIConfig) {
    this.config = config;
    this.providers = new Map();
    this.cache = new Map();
    this.initializeProviders();
  }

  private initializeProviders(): void {
    const { providers } = this.config;

    if (providers[AIProvider.GEMINI]) {
      const gem = providers[AIProvider.GEMINI]!;
      this.providers.set(
        AIProvider.GEMINI,
        new GeminiProvider(gem.apiKey, gem.defaultModel, gem)
      );
    }

    if (providers[AIProvider.OPENAI]) {
      const oai = providers[AIProvider.OPENAI]!;
      this.providers.set(
        AIProvider.OPENAI,
        new OpenAIProvider(oai.apiKey, oai.defaultModel, oai)
      );
    }

    if (providers[AIProvider.ANTHROPIC]) {
      const claude = providers[AIProvider.ANTHROPIC]!;
      this.providers.set(
        AIProvider.ANTHROPIC,
        new AnthropicProvider(claude.apiKey, claude.defaultModel, claude)
      );
    }
  }

  private generateCacheKey(params: any): string {
    return JSON.stringify(params);
  }

  private getProvider(providerType?: AIProvider) {
    const target = providerType || this.config.defaultProvider;
    const p = this.providers.get(target);
    if (!p) throw new Error(`Provider ${target} not configured or initialized`);
    return { instance: p, type: target };
  }

  public async generateContent(params: GenerateContentParams): Promise<AIResponse> {
    const { instance, type } = this.getProvider(params.provider);
    const cacheKey = this.generateCacheKey({ ...params, provider: type });

    if (this.config.enableCache && this.cache.has(cacheKey)) {
      return { ...this.cache.get(cacheKey)!, cached: true };
    }

    const response = await instance.generateContent({
      prompt: params.prompt,
      systemInstruction: params.systemInstruction,
      model: params.model as AIModel,
      options: params.options,
    });

    if (this.config.enableCache) this.cache.set(cacheKey, response);
    return response;
  }

  public async *streamContent(params: GenerateContentParams): AsyncIterable<string> {
    const { instance, type } = this.getProvider(params.provider);
    if (!instance.streamContent) throw new Error(`Provider ${type} does not support streaming`);

    yield* instance.streamContent({
      prompt: params.prompt,
      systemInstruction: params.systemInstruction,
      model: params.model as AIModel,
      options: params.options,
    });
  }

  public async craftPrompt(params: CraftPromptParams): Promise<AIResponse> {
    const systemInstruction =
      this.SYSTEM_INSTRUCTIONS[params.purpose] ?? this.SYSTEM_INSTRUCTIONS['General Summarization'];

    const userContent = `**Purpose:** ${params.purpose}
**Raw Ideas & Requirements:**
${params.ideas.map(i => `- ${i}`).join('\n')}`;

    return this.generateContent({
      prompt: userContent,
      systemInstruction,
      provider: params.provider,
      model: params.model,
      options: params.options,
    });
  }

  public async refactorCode(params: RefactorCodeParams): Promise<AIResponse> {
    return this.generateContent({
      prompt: params.code,
      systemInstruction: params.systemPrompt,
      provider: params.provider,
      model: params.model,
      options: params.options,
    });
  }

  public getAvailableProviders(): AIProvider[] {
    return Array.from(this.providers.keys());
  }

  public clearCache(): void {
    this.cache.clear();
  }

  public getCacheStats(): { size: number; keys: string[] } {
    return { size: this.cache.size, keys: Array.from(this.cache.keys()) };
  }

  public async testProvider(provider: AIProvider): Promise<boolean> {
    try {
      const response = await this.generateContent({
        prompt: "Return exactly: Test successful.",
        provider,
        options: { maxTokens: 30, temperature: 0, stopSequences: ["\n"] }
      });
      return /Test successful/.test(response.text);
    } catch {
      return false;
    }
  }

  public getUsageInfo(response: AIResponse): string {
    if (!response.usage) return 'Usage information not available';
    const { promptTokens, completionTokens, totalTokens } = response.usage;
    return `Tokens - Prompt: ${promptTokens ?? 'N/A'}, Completion: ${completionTokens ?? 'N/A'}, Total: ${totalTokens ?? 'N/A'}`;
  }
}
