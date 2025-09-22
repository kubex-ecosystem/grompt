import {
  AIModel,
  AIProvider,
  AIResponse,
  BaseProvider,
  GenerateContentParams,
  MultiAIConfig,
  OpenAIModels
} from "../../../types/types";
import OpenAI from "openai";


export class OpenAIProvider extends BaseProvider {
  private openai: OpenAI;
  private config?: MultiAIConfig['providers'][AIProvider.OPENAI];

  constructor(apiKey: string, defaultModel: OpenAIModels, config?: MultiAIConfig['providers'][AIProvider.OPENAI]) {
    super(defaultModel);
    this.config = config;
    this.openai = new OpenAI({
      apiKey,
      baseURL: config?.options?.baseURL,
      organization: config?.options?.organization,
      project: config?.options?.project,
    });
  }

  async generateContent(params: {
    prompt: string;
    systemInstruction?: string;
    model?: AIModel;
    options?: GenerateContentParams['options'];
  }): Promise<AIResponse> {
    const model = (params.model as OpenAIModels) || (this.defaultModel as OpenAIModels);

    const role: React.AriaRole = 'user';
    let messages: OpenAI.Chat.Completions.ChatCompletionMessageParam[];
    if (params.systemInstruction) {
      messages = [
        { role: 'system', content: params.systemInstruction },
        { role: 'user', content: params.prompt }
      ];
    } else {
      messages = [{ role: 'user', content: params.prompt }];
    }

    try {
      const completion = await this.openai.chat.completions.create({
        model,
        messages,
        temperature: params.options?.temperature ?? this.config?.options?.defaultQuery?.temperature,
        max_tokens: params.options?.maxTokens ?? this.config?.options?.defaultQuery?.max_tokens ?? 4000,
        top_p: params.options?.topP ?? this.config?.options?.defaultQuery?.top_p,
        frequency_penalty: this.config?.options?.defaultQuery?.frequency_penalty,
        presence_penalty: this.config?.options?.defaultQuery?.presence_penalty,
        stop: params.options?.stopSequences,
        stream: false,
      });

      const choice = completion.choices[0];

      return {
        text: choice?.message?.content || '',
        provider: AIProvider.OPENAI,
        model,
        usage: {
          promptTokens: completion.usage?.prompt_tokens,
          completionTokens: completion.usage?.completion_tokens,
          totalTokens: completion.usage?.total_tokens,
        },
        finishReason: choice?.finish_reason || undefined,
        cached: false,
      };
    } catch (error) {
      console.error('OpenAI API error:', error);
      throw new Error(`OpenAI generation failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  }

  async *streamContent(params: {
    prompt: string;
    systemInstruction?: string;
    model?: AIModel;
    options?: GenerateContentParams['options'];
  }): AsyncIterable<string> {
    const model = (params.model as OpenAIModels) || (this.defaultModel as OpenAIModels);

    let messages: OpenAI.Chat.Completions.ChatCompletionMessageParam[];
    if (params.systemInstruction) {
      messages = [
        { role: 'system', content: params.systemInstruction },
        { role: 'user', content: params.prompt }
      ];
    } else {
      messages = [{ role: 'user', content: params.prompt }];
    }

    try {
      const stream = await this.openai.chat.completions.create({
        model,
        messages,
        temperature: params.options?.temperature,
        max_tokens: params.options?.maxTokens ?? 4000,
        top_p: params.options?.topP,
        stop: params.options?.stopSequences,
        stream: true,
      });

      for await (const chunk of stream) {
        const content = chunk.choices?.[0]?.delta?.content;
        if (content) yield content;
      }
    } catch (error) {
      console.error('OpenAI streaming error:', error);
      throw new Error(`OpenAI streaming failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  }
}
