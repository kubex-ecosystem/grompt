import { GoogleGenAI } from "@google/genai";
import {
  AIProvider,
  BaseProvider,
  GeminiModels,
  type AIModel,
  type AIResponse,
  type GenerateContentParams,
  type MultiAIConfig,
} from "../types";

function readText(resp: any): string {
  // cobre variantes do SDK (prop direta, método, parts…)
  return (
    resp?.text ??
    resp?.response?.text?.() ??
    (Array.isArray(resp?.response?.candidates)
      ? resp.response.candidates
        .flatMap((c: any) => c?.content?.parts ?? [])
        .map((p: any) => p?.text ?? "")
        .join("")
      : "") ??
    ""
  );
}

function readUsage(resp: any) {
  const u = resp?.usage ?? resp?.response?.usageMetadata;
  return u
    ? {
      promptTokens: u.promptTokenCount,
      completionTokens: u.candidatesTokenCount,
      totalTokens: u.totalTokenCount,
    }
    : undefined;
}

export class GeminiProvider extends BaseProvider {
  private ai: GoogleGenAI;
  private config?: MultiAIConfig["providers"][AIProvider.GEMINI];

  constructor(
    apiKey: string,
    defaultModel: GeminiModels,
    config?: MultiAIConfig["providers"][AIProvider.GEMINI]
  ) {
    super(defaultModel);
    this.ai = new GoogleGenAI({ apiKey });
    this.config = config;
  }

  /** monta a sessão de chat com as configs unificadas */
  private createChat(
    model: GeminiModels,
    systemInstruction?: string,
    options?: GenerateContentParams["options"]
  ) {
    return this.ai.chats.create({
      model,
      config: {
        // geração
        temperature:
          options?.temperature ??
          this.config?.options?.generationConfig?.temperature,
        topP: options?.topP ?? this.config?.options?.generationConfig?.topP,
        maxOutputTokens:
          options?.maxTokens ??
          this.config?.options?.generationConfig?.maxOutputTokens ??
          4000,
        stopSequences: options?.stopSequences,

        // safety & system
        safetySettings: this.config?.options?.safetySettings,
        systemInstruction, // string direto funciona nesse layer
      },
    });
  }

  async generateContent(params: {
    prompt: string;
    systemInstruction?: string;
    model?: AIModel;
    options?: GenerateContentParams["options"];
  }): Promise<AIResponse> {
    const model =
      (params.model as GeminiModels) || (this.defaultModel as GeminiModels);

    try {
      // chat efêmero por call (sem histórico compartilhado)
      const chat = this.createChat(
        model,
        params.systemInstruction,
        params.options
      );
      const resp = await chat.sendMessage({ message: params.prompt });

      return {
        text: readText(resp),
        provider: AIProvider.GEMINI,
        model,
        usage: readUsage(resp),
        finishReason: [resp.promptFeedback?.blockReason || "Unknown", resp.promptFeedback?.blockReasonMessage || "N/A"].join(" - "),
        cached: false,
      };
    } catch (error) {
      console.error("Gemini API error:", error);
      throw new Error(
        `Gemini generation failed: ${error instanceof Error ? error.message : "Unknown error"
        }`
      );
    }
  }

  async *streamContent(params: {
    prompt: string;
    systemInstruction?: string;
    model?: AIModel;
    options?: GenerateContentParams["options"];
  }): AsyncIterable<string> {
    const model =
      (params.model as GeminiModels) || (this.defaultModel as GeminiModels);

    try {
      const chat = this.createChat(
        model,
        params.systemInstruction,
        params.options
      );
      const stream = await chat.sendMessageStream({ message: params.prompt });

      // algumas versões expõem `.stream`, outras já são AsyncIterable
      const it: AsyncIterable<any> = (stream as any).stream ?? (stream as any);

      for await (const chunk of it) {
        // chunk pode ter `text()` ou `text`
        const t =
          (chunk && typeof chunk.text === "function" && chunk.text()) ||
          (typeof chunk?.text === "string" ? chunk.text : "");
        if (t) yield t;
      }
    } catch (error) {
      console.error("Gemini streaming error:", error);
      throw new Error(
        `Gemini streaming failed: ${error instanceof Error ? error.message : "Unknown error"
        }`
      );
    }
  }
}
