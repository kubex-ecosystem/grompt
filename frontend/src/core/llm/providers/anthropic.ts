import { default as Anthropic } from "@anthropic-ai/sdk";

const {
  Messages,
} = Anthropic;

import { MessagesPage } from "openai/resources/beta/threads/messages";
import {
  AIProvider,
  AnthropicModels,
  BaseProvider,
  type AIModel,
  type AIResponse,
  type GenerateContentParams,
  type MultiAIConfig
} from "../types";

export class AnthropicProvider extends BaseProvider {
  private anthropic: Anthropic;
  private config?: MultiAIConfig['providers'][AIProvider.ANTHROPIC];

  constructor(apiKey: string, defaultModel: AnthropicModels, config?: MultiAIConfig['providers'][AIProvider.ANTHROPIC]) {
    super(defaultModel);
    this.config = config;
    this.anthropic = new Anthropic({
      apiKey,
      baseURL: config?.options?.baseURL,
      defaultHeaders: config?.options?.defaultHeaders,
    });
  }

  async generateContent(params: {
    prompt: string;
    systemInstruction?: string;
    model?: AIModel;
    options?: GenerateContentParams['options'];
  }): Promise<AIResponse> {
    const model = (params.model as AnthropicModels) || (this.defaultModel as AnthropicModels);

    try {
      const msg = await this.anthropic.messages.create({
        model,
        max_tokens: params.options?.maxTokens ?? this.config?.options?.defaultQuery?.max_tokens ?? 4000,
        temperature: params.options?.temperature ?? this.config?.options?.defaultQuery?.temperature,
        top_p: params.options?.topP ?? this.config?.options?.defaultQuery?.top_p,
        top_k: this.config?.options?.defaultQuery?.top_k,
        system: params.systemInstruction,
        stop_sequences: params.options?.stopSequences,
        messages: [{ role: "user", content: params.prompt }],
        stream: false,
      });

      const first = msg.content?.[0];
      const text = (first && first.type === "text") ? first.text : "";

      return {
        text,
        provider: AIProvider.ANTHROPIC,
        model,
        usage: {
          promptTokens: msg.usage?.input_tokens,
          completionTokens: msg.usage?.output_tokens,
          totalTokens: (msg.usage?.input_tokens ?? 0) + (msg.usage?.output_tokens ?? 0),
        },
        finishReason: msg.stop_reason ?? undefined,
        cached: false,
      };
    } catch (error) {
      console.error("Anthropic API error:", error);
      throw new Error(`Anthropic generation failed: ${error instanceof Error ? error.message : "Unknown error"}`);
    }
  }

  async *streamContent(params: {
    prompt: string;
    systemInstruction?: string;
    model?: AIModel;
    options?: GenerateContentParams['options'];
  }): AsyncIterable<string> {
    const model = (params.model as AnthropicModels) || (this.defaultModel as AnthropicModels);

    try {
      const stream = await this.anthropic.messages.create({
        model,
        max_tokens: params.options?.maxTokens ?? 4000,
        temperature: params.options?.temperature,
        top_p: params.options?.topP,
        system: params.systemInstruction,
        stop_sequences: params.options?.stopSequences,
        messages: [{ role: "user", content: params.prompt }],
        stream: true,
      });

      for await (const ev of stream as AsyncIterable<MessagesPage>) {
        const iterator = ev.getNextPage();
        for await (const message of (await iterator).iterPages()) {
          if (typeof message === "string") {
            yield message;
          }
        }
      }
    } catch (error) {
      console.error("Anthropic streaming error:", error);
      throw new Error(`Anthropic streaming failed: ${error instanceof Error ? error.message : "Unknown error"}`);
    }
  }
}
