

import { GoogleGenerativeAI, type GenerateContentResponse } from "@google/genai";
import {
  AIProvider,
  BaseProvider,
  GeminiModels,
  type AIModel,
  type AIResponse,
  type GenerateContentParams,
  type MultiAIConfig
} from "../types";

export class GeminiProvider extends BaseProvider {
  private genAI: GoogleGenerativeAI;
  private config?: MultiAIConfig['providers'][AIProvider.GEMINI];

  constructor(apiKey: string, defaultModel: GeminiModels, config?: MultiAIConfig['providers'][AIProvider.GEMINI]) {
    super(defaultModel);
    this.genAI = new GoogleGenerativeAI(apiKey);
    this.config = config;
  }

  async generateContent(params: {
    prompt: string;
    systemInstruction?: string;
    model?: AIModel;
    options?: GenerateContentParams['options'];
  }): Promise<AIResponse> {
    const model = (params.model as GeminiModels) || (this.defaultModel as GeminiModels);

    const generationConfig = {
      temperature: params.options?.temperature ?? this.config?.options?.generationConfig?.temperature,
      topP: params.options?.topP ?? this.config?.options?.generationConfig?.topP,
      topK: this.config?.options?.generationConfig?.topK,
      maxOutputTokens: params.options?.maxTokens
        ?? this.config?.options?.generationConfig?.maxOutputTokens
        ?? 4000,
      stopSequences: params.options?.stopSequences
    };

    const modelInstance = this.genAI.getGenerativeModel({
      model,
      generationConfig,
      safetySettings: this.config?.options?.safetySettings,
      systemInstruction: params.systemInstruction
        ? { role: "system", parts: [{ text: params.systemInstruction }] }
        : undefined,
    });

    try {
      const result: GenerateContentResponse = await modelInstance.generateContent(params.prompt);
      const response = result.response;

      return {
        text: response?.text() ?? "",
        provider: AIProvider.GEMINI,
        model,
        usage: {
          promptTokens: response?.usageMetadata?.promptTokenCount,
          completionTokens: response?.usageMetadata?.candidatesTokenCount,
          totalTokens: response?.usageMetadata?.totalTokenCount,
        },
        finishReason: response?.candidates?.[0]?.finishReason,
        cached: false,
      };
    } catch (error) {
      console.error('Gemini API error:', error);
      throw new Error(`Gemini generation failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  }

  async *streamContent(params: {
    prompt: string;
    systemInstruction?: string;
    model?: AIModel;
    options?: GenerateContentParams['options'];
  }): AsyncIterable<string> {
    const model = (params.model as GeminiModels) || (this.defaultModel as GeminiModels);

    const modelInstance = this.genAI.getGenerativeModel({
      model,
      systemInstruction: params.systemInstruction
        ? { role: "system", parts: [{ text: params.systemInstruction }] }
        : undefined,
      generationConfig: {
        temperature: params.options?.temperature ?? this.config?.options?.generationConfig?.temperature,
        maxOutputTokens: params.options?.maxTokens
          ?? this.config?.options?.generationConfig?.maxOutputTokens
          ?? 4000,
        stopSequences: params.options?.stopSequences
      }
    });

    try {
      const streamResp = await modelInstance.generateContentStream(params.prompt);
      for await (const chunk of streamResp.stream) {
        const t = chunk.text();
        if (t) yield t;
      }
    } catch (error) {
      console.error('Gemini streaming error:', error);
      throw new Error(`Gemini streaming failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  }
}