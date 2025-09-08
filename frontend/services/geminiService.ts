import { GenerateContentResponse, GoogleGenAI } from "@google/genai";
import { Idea } from '../types';

const API_KEY = process.env.API_KEY;

if (!API_KEY) {
  console.warn("Gemini API key not found. Running in demo mode with simulated responses.");
}

// Function to get user's API key from localStorage
const getUserApiKey = (): string | null => {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('userApiKey');
  }
  return null;
};

// Function to get active API key (user's key takes precedence)
const getActiveApiKey = (): string | null => {
  return getUserApiKey() || API_KEY || null;
};

const createAIInstance = (apiKey: string) => {
  return new GoogleGenAI({ apiKey });
};

// FIX: Updated model to 'gemini-2.5-flash' to comply with guidelines.
const model = 'gemini-2.5-flash';

export interface PromptGenerationResult {
  prompt: string;
  // FIX: Made token count properties optional to match the Gemini API response type.
  usageMetadata?: {
    promptTokenCount?: number;
    candidatesTokenCount?: number;
    totalTokenCount?: number;
  };
}

// Demo mode simulation
const generateDemoPrompt = (ideas: Idea[], purpose: string): PromptGenerationResult => {
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
  const estimatedTokens = Math.floor(demoPrompt.length / 4); // Rough estimation: 1 token ≈ 4 characters

  return {
    prompt: demoPrompt,
    usageMetadata: {
      promptTokenCount: Math.floor(estimatedTokens * 0.3),
      candidatesTokenCount: Math.floor(estimatedTokens * 0.7),
      totalTokenCount: estimatedTokens
    }
  };
};

const createMetaPrompt = (ideas: Idea[], purpose: string): string => {
  const ideasText = ideas.map((idea, index) => `- ${idea.text}`).join('\n');

  return `
You are a world-class prompt engineering expert, a key copilot in the Kubex Ecosystem. Your mission is to transform a user's raw, unstructured ideas into a clean, effective, and professional prompt for large language models, adhering strictly to Kubex principles.

**KUBEX PRINCIPLES:**
- **Simplicidade Radical:** The prompt must be direct, pragmatic, and anti-jargon. "One command = one result."
- **Modularidade:** The output should be well-structured and easily usable.
- **Sem Jaulas (No Cages):** The prompt should use open, clear formats (like Markdown) and not lock the user into a specific model's quirks.

**User's Raw Ideas:**
${ideasText}

**Desired Purpose of the Final Prompt:**
${purpose}

**Your Task:**
Based on the ideas and purpose provided, generate a single, comprehensive, and well-structured prompt in Markdown format. The generated prompt must be ready for immediate use.

**Key Directives:**
1.  **Define a Persona:** Start the prompt by defining a clear, expert role for the AI (e.g., "You are an expert software architect specializing in cloud-native solutions...").
2.  **State the Objective:** Clearly articulate the main goal of the prompt in a "Primary Objective" section.
3.  **Provide Structure:** Use Markdown (headings, lists, bold text) to create a logical and readable structure. Use sections like "Requirements", "Expected Output", "Constraints", and "Context".
4.  **List Key Requirements:** Systematically break down the user's ideas into specific, actionable requirements in a bulleted list.
5.  **Specify Output Format:** If applicable, describe the desired output format (e.g., JSON schema, code block with language specifier, table).
6.  **Include Constraints:** Add any negative constraints or things to avoid (e.g., "Do not use libraries outside the standard library.").
7.  **Add Context:** If implied by the ideas, add a "Context" section to provide background that will help the AI generate a better response.

**OUTPUT CONTRACT:**
Return ONLY the generated prompt in Markdown. Do not include any introductory text, explanations, or concluding remarks like "Here is the prompt:". Your response must start directly with the Markdown content of the prompt (e.g., starting with a '# Title' or '**Persona:**').
`;
};

const MAX_RETRIES = parseInt(process.env.MAX_RETRIES || "3", 10);
const RETRY_DELAY_MS = parseInt(process.env.RETRY_DELAY_MS || "1000", 10);

/**
 * Checks if an error is likely a transient network issue and thus retryable.
 * @param error The error object.
 * @returns True if the error is retryable, false otherwise.
 */
const isRetryableError = (error: any): boolean => {
  const errorMessage = (error?.message || '').toLowerCase();
  // Simple check for network-related errors.
  return errorMessage.includes('fetch') || errorMessage.includes('network');
};

/**
 * Formats an API error into a user-friendly string.
 * @param error The error object from the API call.
 * @returns A user-friendly error message string.
 */
const formatApiError = (error: any): string => {
  // Keep a detailed log for developers
  console.error("Gemini API Error:", error);

  if (error instanceof Error) {
    // Provide a more specific message for common, user-actionable errors.
    if (error.message.includes('API key not valid')) {
      return 'API Configuration Error: The API key is not valid. Please contact the administrator.';
    }
    // For other errors, provide a message that includes the API's feedback.
    return `An error occurred: ${error.message}. Please check your input or try again. If the issue persists, it may be a network problem.`;
  }

  return "An unknown error occurred while communicating with the Gemini API. Please check your connection and try again.";
};

export const generateStructuredPrompt = async (ideas: Idea[], purpose: string): Promise<PromptGenerationResult> => {
  const activeApiKey = getActiveApiKey();

  // If no API key available, use demo mode
  if (!activeApiKey) {
    return generateDemoPrompt(ideas, purpose);
  }

  const ai = createAIInstance(activeApiKey);

  const isRetryableError = (error: any): boolean => {
    if (error.message?.includes('rate limit') || error.message?.includes('quota')) {
      return true;
    }
    if (error.status === 429 || error.status === 500 || error.status === 502 || error.status === 503) {
      return true;
    }
    return false;
  };

  const formatApiError = (error: any): string => {
    if (error.message?.includes('API key')) {
      return "API key issue. Please check your Gemini API key configuration.";
    }
    if (error.message?.includes('rate limit') || error.message?.includes('quota')) {
      return "Rate limit exceeded. Please try again later.";
    }
    return `Gemini API error: ${error.message || 'Unknown error occurred'}`;
  };

  let lastError: any = null;

  for (let attempt = 1; attempt <= MAX_RETRIES; attempt++) {
    try {
      const metaPrompt = createMetaPrompt(ideas, purpose);
      const response: GenerateContentResponse = await ai.models.generateContent({
        model: model,
        contents: metaPrompt,
        config: {
          temperature: 0.3,
          topP: 0.9,
          topK: 40,
        }
      });

      const generatedText = response.text;
      if (!generatedText) {
        throw new Error("Empty response received from Gemini API");
      }

      return {
        prompt: generatedText,
        usageMetadata: response.usageMetadata
      };
    } catch (error) {
      lastError = error;
      if (isRetryableError(error) && attempt < MAX_RETRIES) {
        console.warn(`API call failed (attempt ${attempt}/${MAX_RETRIES}), retrying after delay...`);
        await new Promise(resolve => setTimeout(resolve, RETRY_DELAY_MS * attempt)); // Simple exponential backoff
      } else {
        // Not a retryable error or max retries reached
        throw new Error(formatApiError(error));
      }
    }
  }
  // This fallback should ideally not be reached, but it's here for safety.
  throw new Error(formatApiError(lastError || new Error("Failed to generate prompt after multiple retries.")));
};
