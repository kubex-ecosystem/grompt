import { GenerateContentResponse, GoogleGenAI } from "@google/genai";
import { Idea } from '../types';

const API_KEY = process.env.API_KEY;

if (!API_KEY) {
  console.warn("Gemini API key not found. Using mock data. Please set the API_KEY environment variable.");
}

const ai = new GoogleGenAI({ apiKey: API_KEY! });
// FIX: Updated model to 'gemini-2.5-flash' to comply with guidelines.
const model = 'gemini-2.5-flash';

export interface PromptGenerationResult {
  prompt: string;
  usageMetadata?: {
    promptTokenCount?: number;
    candidatesTokenCount?: number;
    totalTokenCount?: number;
  };
}

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

const MAX_RETRIES = 3;
const RETRY_DELAY_MS = 1000;

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
  if (!API_KEY) {
    // Mock response for environments without an API key
    await new Promise(resolve => setTimeout(resolve, 1500));
    const ideasText = ideas.map((idea, index) => `- ${idea.text}`).join('\n');
    const mockPrompt = `
**Persona:** You are a mock AI assistant.

**Primary Objective:** Demonstrate a structured prompt based on user input.

**Context:** This is a mock response because the Gemini API key is not configured.

**User's Ideas Provided:**
${ideasText}

**Stated Purpose:**
${purpose}

**Expected Output (Mock):**
- A well-formatted response that addresses each of the user's ideas.
- Clear structure using Markdown.
- Adherence to the requested purpose.
`;
    // Mock token counts
    const inputTokens = Math.floor(mockPrompt.length / 4);
    const outputTokens = Math.floor(ideasText.length / 4);

    return {
      prompt: mockPrompt,
      usageMetadata: {
        promptTokenCount: inputTokens,
        candidatesTokenCount: outputTokens,
        totalTokenCount: inputTokens + outputTokens
      }
    };
  }

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
      return {
        // FIX: Using response.text directly as per Gemini API guidelines.
        prompt: response.text || '',
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
