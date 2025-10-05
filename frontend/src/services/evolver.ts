import {
  AIModel,
  AIProvider,
  AnthropicModels,
  DeepSeekModels,
  GeminiModels,
  GenerateContentParams,
  MultiAIWrapper as llm,
  MultiAIConfig,
  OpenAIModels
} from "../core/llm";
import { extractDiffFenced, pack } from "../utils/lookatni";

/** Seleciona provider a partir do env sem travar UI */
export function selectProvider(providers: string[], defaultProvider: string | null, isDemoMode: boolean): string | null {
  if (isDemoMode || providers.length === 0) return null;
  if (defaultProvider && providers.includes(defaultProvider)) return defaultProvider;
  if (providers.includes('gpt-4o')) return 'gpt-4o';
  if (providers.includes('gpt-4o-mini')) return 'gpt-4o-mini';
  if (providers.includes('gpt-4')) return 'gpt-4';
  if (providers.includes('claude-3')) return 'claude-3';
  if (providers.includes('claude-instant-100k')) return 'claude-instant-100k';
  if (providers.includes('gpt-3.5-turbo')) return 'gpt-3.5-turbo';
  return providers[0];
}

/** Monta o payload LookAtni a partir do VFS/FS do app. Troca por fetch('/snapshot') se preferir. */
export async function buildLookatniBlob(grab: (p: string) => string, include: string[]): Promise<string> {
  const files = include.map(p => ({ path: p, content: grab(p) }));
  return pack(files);
}

/** Pede um UNIFIED DIFF pequeno. */
export async function requestUnifiedDiff(blob: string, task: string): Promise<string> {
  const system = [
    "You are a rigorous code refactorer.",
    "Return ONLY one unified diff fenced with diff.",
    "Small, reviewable changes. Preserve behavior and paths (a/ b/ headers)."
  ].join("\n");

  const user = [
    "Context (LookAtni):",
    blob,
    "",
    "Task:",
    task,
    "",
    "Constraints:",
    "- Keep exports stable; no broad rewrites.",
    "- Prefer extracting helpers, removing duplication, adding small tests.",
    "- If you add files, include proper diff headers."
  ].join("\n");

  const config: MultiAIConfig = {
    apiKey: process.env.API_KEY || undefined,
    azureApiKey: process.env.AZURE_API_KEY || undefined,
    azureEndpoint: process.env.AZURE_ENDPOINT || undefined,
    azureDeployment: process.env.AZURE_DEPLOYMENT || undefined,
    anthropicApiKey: process.env.ANTHROPIC_API_KEY || undefined,
    cohereApiKey: process.env.COHERE_API_KEY || undefined,
    googleApiKey: process.env.GOOGLE_API_KEY || undefined,
    googleModel: process.env.GOOGLE_MODEL || undefined,
    huggingfaceApiKey: process.env.HUGGINGFACE_API_KEY || undefined,
    deepseekApiKey: process.env.DEEPSEEK_API_KEY || undefined,

    defaultProvider: AIProvider.GEMINI,
    defaultModel: GeminiModels.GEMINI_FLASH,
    enableCache: true,

    providers: {
      [AIProvider.GEMINI]: {
        apiKey: process.env.GOOGLE_API_KEY || "",
        defaultModel: process.env.GEMINI_MODEL?.search(GeminiModels.GEMINI_PRO) !== -1 ? GeminiModels.GEMINI_PRO : GeminiModels.GEMINI_FLASH,
        baseURL: process.env.GOOGLE_BASE_URL || "",
        project: process.env.GOOGLE_PROJECT || "",
        location: process.env.GOOGLE_LOCATION || "",
        options: {
          safetySettings: [],
          generationConfig: {
            maxOutputTokens: process.env.GOOGLE_MAX_TOKENS ? parseInt(process.env.GOOGLE_MAX_TOKENS) : 1024,
            topP: process.env.GOOGLE_TOP_P ? parseFloat(process.env.GOOGLE_TOP_P) : 1,
            topK: process.env.GOOGLE_TOP_K ? parseInt(process.env.GOOGLE_TOP_K) : 40,
            temperature: process.env.GOOGLE_TEMPERATURE ? parseFloat(process.env.GOOGLE_TEMPERATURE) : 0.7,
          },
        },
      },
      [AIProvider.OPENAI]: {
        apiKey: process.env.OPENAI_API_KEY || "",
        defaultModel: process.env.OPENAI_MODEL?.search('gpt-4o') !== -1 ? OpenAIModels.GPT_4O :
          process.env.OPENAI_MODEL?.search('gpt-4o-mini') !== -1 ? OpenAIModels.GPT_4O_MINI :
            process.env.OPENAI_MODEL?.search('gpt-4') !== -1 ? OpenAIModels.GPT_4 : OpenAIModels.GPT_3_5_TURBO,
        options: {
          baseURL: process.env.OPENAI_BASE_URL || "",
          defaultQuery: {
            temperature: process.env.OPENAI_TEMPERATURE ? parseFloat(process.env.OPENAI_TEMPERATURE) : 0.7,
            max_tokens: process.env.OPENAI_MAX_TOKENS ? parseInt(process.env.OPENAI_MAX_TOKENS) : 1024,
            top_p: process.env.OPENAI_TOP_P ? parseFloat(process.env.OPENAI_TOP_P) : 1,
            top_k: process.env.OPENAI_TOP_K ? parseInt(process.env.OPENAI_TOP_K) : 40,
          },
        },
      },
      [AIProvider.ANTHROPIC]: {
        apiKey: process.env.ANTHROPIC_API_KEY || "",
        defaultModel: process.env.ANTHROPIC_MODEL?.search('claude-3') !== -1 ? AnthropicModels.CLAUDE_3_SONNET :
          process.env.ANTHROPIC_MODEL?.search('claude-instant-100k') !== -1 ? AnthropicModels.CLAUDE_3_OPUS :
            AnthropicModels.CLAUDE_3_5_SONNET,
        options: {
          baseURL: process.env.ANTHROPIC_BASE_URL || "",
          defaultQuery: {
            temperature: process.env.ANTHROPIC_TEMPERATURE ? parseFloat(process.env.ANTHROPIC_TEMPERATURE) : 0.7,
            max_tokens: process.env.ANTHROPIC_MAX_TOKENS ? parseInt(process.env.ANTHROPIC_MAX_TOKENS) : 1024,
            top_p: process.env.ANTHROPIC_TOP_P ? parseFloat(process.env.ANTHROPIC_TOP_P) : 1,
            top_k: process.env.ANTHROPIC_TOP_K ? parseInt(process.env.ANTHROPIC_TOP_K) : 40,
          },
        },
      },
      [AIProvider.DEEPSEEK]: {
        apiKey: process.env.DEEPSEEK_API_KEY || "",
        defaultModel: process.env.DEEPSEEK_MODEL?.search(DeepSeekModels.DEEPSEEK_V1) !== -1 ? DeepSeekModels.DEEPSEEK_V1 : DeepSeekModels.DEEPSEEK_V2,
        options: {
          baseURL: process.env.DEEPSEEK_BASE_URL || "",
          defaultQuery: {
            temperature: process.env.DEEPSEEK_TEMPERATURE ? parseFloat(process.env.DEEPSEEK_TEMPERATURE) : 0.7,
            max_tokens: process.env.DEEPSEEK_MAX_TOKENS ? parseInt(process.env.DEEPSEEK_MAX_TOKENS) : 1024,
            top_p: process.env.DEEPSEEK_TOP_P ? parseFloat(process.env.DEEPSEEK_TOP_P) : 1,
            top_k: process.env.DEEPSEEK_TOP_K ? parseInt(process.env.DEEPSEEK_TOP_K) : 40,
          },
        },
      },
    }
  };

  const selectedModel = config.providers[AIProvider.GEMINI]?.defaultModel
    || config.providers[AIProvider.OPENAI]?.defaultModel
    || config.providers[AIProvider.ANTHROPIC]?.defaultModel
    || config.providers[AIProvider.DEEPSEEK]?.defaultModel
    || GeminiModels.GEMINI_FLASH;

  const providers = new llm(config)
  console.log("Available LLM providers:", providers);

  const selectedProvider = selectProvider(providers.getAvailableProviders(), process.env.DEFAULT_LLM_PROVIDER || null, !process.env.API_KEY);
  if (!selectedProvider) throw new Error("No LLM provider selected");

  const providerWorks = await providers.testProvider(selectedProvider as AIProvider);
  if (!providerWorks) throw new Error(`Failed to initialize provider ${selectedProvider}`);
  console.log("Selected LLM provider:", selectedProvider);

  const generateContentParams: GenerateContentParams = {
    prompt: blob,
    provider: selectedProvider as AIProvider,
    model: selectedModel as AIModel,
    systemInstruction: system,
    options: {
      temperature: 0.7,
      maxTokens: 1024,
      stopSequences: ["```", "### End"],
      stream: false,
      topP: 1,
      topK: 40,
    }
  };

  const response = await providers.generateContent(generateContentParams)
  if (!response || !response.text) throw new Error("No response from LLM");
  console.log("LLM response received");

  const out = providers.getUsageInfo(response)
  if (!out) throw new Error("Failed to get usage info from LLM response");

  return extractDiffFenced(out);
}

