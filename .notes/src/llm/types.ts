

// Core enums e tipos compartilhados

export enum AIProvider {
  GEMINI = 'gemini',
  OPENAI = 'openai',
  ANTHROPIC = 'anthropic',
}

export enum GeminiModels {
  GEMINI_FLASH = 'gemini-2.5-flash',
  GEMINI_PRO = 'gemini-1.5-pro-latest',
}

export enum OpenAIModels {
  GPT_4 = 'gpt-4',
  GPT_4_TURBO = 'gpt-4-turbo',
  GPT_3_5_TURBO = 'gpt-3.5-turbo',
  GPT_4O = 'gpt-4o',
  GPT_4O_MINI = 'gpt-4o-mini',
}

export enum AnthropicModels {
  CLAUDE_3_OPUS = 'claude-3-opus-20240229',
  CLAUDE_3_SONNET = 'claude-3-sonnet-20240229',
  CLAUDE_3_HAIKU = 'claude-3-haiku-20240307',
  // Se sua org tiver acesso ao build novo, troque aqui.
  CLAUDE_3_5_SONNET = 'claude-3-5-sonnet-20240620',
}

export type AIModel = GeminiModels | OpenAIModels | AnthropicModels;

export interface AIResponse {
  text: string;
  provider: AIProvider;
  model: AIModel;
  cached?: boolean;
  usage?: {
    promptTokens?: number;
    completionTokens?: number;
    totalTokens?: number;
  };
  finishReason?: string;
}

export interface MultiAIConfig {
  providers: {
    [AIProvider.GEMINI]?: {
      apiKey: string;
      defaultModel: GeminiModels;
      options?: {
        safetySettings?: any[];
        generationConfig?: {
          temperature?: number;
          topP?: number;
          topK?: number;
          maxOutputTokens?: number;
        };
      };
    };
    [AIProvider.OPENAI]?: {
      apiKey: string;
      defaultModel: OpenAIModels;
      options?: {
        baseURL?: string;
        organization?: string;
        project?: string;
        defaultQuery?: {
          temperature?: number;
          max_tokens?: number;
          top_p?: number;
          frequency_penalty?: number;
          presence_penalty?: number;
        };
      };
    };
    [AIProvider.ANTHROPIC]?: {
      apiKey: string;
      defaultModel: AnthropicModels;
      options?: {
        baseURL?: string;
        defaultHeaders?: Record<string, string>;
        defaultQuery?: {
          temperature?: number;
          max_tokens?: number;
          top_p?: number;
          top_k?: number;
        };
      };
    };
  };
  defaultProvider: AIProvider;
  enableCache?: boolean;
}

export interface GenerateContentParams {
  prompt: string;
  systemInstruction?: string;
  provider?: AIProvider;
  model?: AIModel;
  options?: {
    temperature?: number;
    maxTokens?: number;
    topP?: number;
    stream?: boolean;
    stopSequences?: string[];
  };
}

export interface CraftPromptParams {
  ideas: string[];
  purpose: string;
  provider?: AIProvider;
  model?: AIModel;
  options?: GenerateContentParams['options'];
}

export interface RefactorCodeParams {
  systemPrompt: string;
  code: string;
  provider?: AIProvider;
  model?: AIModel;
  options?: GenerateContentParams['options'];
}

export abstract class BaseProvider {
  protected defaultModel: AIModel;

  constructor(defaultModel: AIModel) {
    this.defaultModel = defaultModel;
  }

  abstract generateContent(params: {
    prompt: string;
    systemInstruction?: string;
    model?: AIModel;
    options?: GenerateContentParams['options'];
  }): Promise<AIResponse>;

  // Opcional: streaming
  streamContent?(params: {
    prompt: string;
    systemInstruction?: string;
    model?: AIModel;
    options?: GenerateContentParams['options'];
  }): AsyncIterable<string>;
}