export interface PromptHistoryItem {
  id: number;
  prompt: string;
  timestamp: number;
  inputs: {
    ideas: string[];
    purpose: string;
  };
}

export enum AIProvider {
  GEMINI = 'gemini',
  OPENAI = 'openai',
  ANTHROPIC = 'anthropic',
  DEEPSEEK = 'deepseek',
}

export enum GeminiModels {
  GEMINI_FLASH = 'gemini-2.0-flash',
  GEMINI_PRO = 'gemini-2.5-pro',
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

export enum DeepSeekModels {
  DEEPSEEK_V1 = 'deepseek-v1',
  DEEPSEEK_V2 = 'deepseek-v2',
}

export type AIModel = GeminiModels | OpenAIModels | AnthropicModels | DeepSeekModels;

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
  // API keys for different providers
  apiKey?: string;
  azureApiKey?: string;
  azureEndpoint?: string;
  azureDeployment?: string;
  anthropicApiKey?: string;
  cohereApiKey?: string;
  googleApiKey?: string;
  googleModel?: string;
  huggingfaceApiKey?: string;
  deepseekApiKey?: string;

  // Provider-specific configurations
  providers: {
    [AIProvider.GEMINI]?: {
      apiKey: string;
      defaultModel: GeminiModels;
      models?: string[];
      baseURL?: string;
      project?: string;
      location?: string;
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
      models?: string[];
      options?: {
        baseURL?: string;
        organization?: string;
        project?: string;
        defaultQuery?: {
          temperature?: number;
          max_tokens?: number;
          top_p?: number;
          top_k?: number;
          frequency_penalty?: number;
          presence_penalty?: number;
        };
      };
    };
    [AIProvider.ANTHROPIC]?: {
      apiKey: string;
      defaultModel: AnthropicModels;
      models?: string[];
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
    [AIProvider.DEEPSEEK]?: {
      apiKey: string;
      defaultModel: DeepSeekModels;
      models?: string[];
      options?: {
        baseURL?: string;
        defaultQuery?: {
          temperature?: number;
          max_tokens?: number;
          top_p?: number;
          top_k?: number;
        };
      };
    };
  };
  defaultProvider?: AIProvider;
  defaultModel?: AIModel;
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
    topK?: number;
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
