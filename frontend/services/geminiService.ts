import { GoogleGenAI, GenerateContentResponse } from "@google/genai";

const API_KEY = process.env.API_KEY;

if (!API_KEY) {
  throw new Error("API_KEY environment variable not set");
}

const ai = new GoogleGenAI({ apiKey: API_KEY });

const BASE_INSTRUCTION = `You are an expert prompt engineer. Your task is to synthesize a user's raw ideas and their stated purpose into a single, comprehensive, and professional prompt for a large language model. The prompt must be structured, clear, and detailed to elicit the best possible response from the AI.

Follow these steps:
1.  **Role and Goal**: Start by defining the role the AI should take and the primary goal of the task.
2.  **Context**: Incorporate all user ideas as context, background, or specific requirements.
3.  **Instructions**: Provide clear, step-by-step instructions.
4.  **Constraints & Formatting**: Define constraints, such as output format (e.g., JSON, markdown), tone, and negative constraints (what to avoid).
5.  **Refinement**: Ensure the language is precise and unambiguous. Use Markdown for structure.

Produce ONLY the final prompt text. Do not include any conversational filler, greetings, or explanations about your process.`;


const SYSTEM_INSTRUCTIONS: Record<string, string> = {
    'Code Generation': `${BASE_INSTRUCTION} \nSpecialize the prompt for generating code. Ask for specific languages, libraries, function signatures, and error handling. The AI should act as a senior software engineer.`,
    'Creative Writing': `${BASE_INSTRUCTION} \nSpecialize the prompt for creative writing. Focus on tone, style, character details, plot points, and setting. The AI should act as a master storyteller or author.`,
    'Data Analysis': `${BASE_INSTRUCTION} \nSpecialize the prompt for data analysis. Request the AI to act as a data scientist. The prompt should specify the data format, the desired analysis (e.g., statistical summary, trend identification), and the output format (e.g., tables, charts description).`,
    'Technical Documentation': `${BASE_INSTRUCTION} \nSpecialize the prompt for technical documentation. The AI should act as a technical writer. The prompt must define the audience, the component or feature to be documented, and the required sections (e.g., overview, API reference, examples).`,
    'Marketing Copy': `${BASE_INSTRUCTION} \nSpecialize the prompt for marketing copy. The AI should act as a senior copywriter. The prompt should detail the target audience, product/service, key benefits, call to action, and desired tone (e.g., witty, professional, urgent).`,
    'General Summarization': `${BASE_INSTRUCTION} \nSpecialize the prompt for summarization. The AI should act as an expert analyst. The prompt should specify the desired length and format of the summary (e.g., bullet points, short paragraph) and the key aspects to focus on.`,
};


interface CraftPromptParams {
  ideas: string[];
  purpose: string;
}

interface RefactorCodeParams {
    systemPrompt: string;
    code: string;
}

export const geminiService = {
  craftPrompt: async ({ ideas, purpose }: CraftPromptParams): Promise<string> => {
    try {
      const userContent = `
        Here are the components for the prompt I need:

        **Purpose:**
        ${purpose}

        **Raw Ideas & Requirements (combine these into the prompt):**
        ${ideas.map(idea => `- ${idea}`).join('\n')}
      `;
      
      const systemInstruction = SYSTEM_INSTRUCTIONS[purpose] || BASE_INSTRUCTION;

      const response: GenerateContentResponse = await ai.models.generateContent({
        model: 'gemini-2.5-flash',
        contents: userContent,
        config: {
            systemInstruction
        }
      });
      return response.text;
    } catch (error) {
      console.error("Error crafting prompt:", error);
      throw new Error("Failed to craft prompt using Gemini API.");
    }
  },

  refactorCode: async ({ systemPrompt, code }: RefactorCodeParams): Promise<string> => {
    try {
        const response: GenerateContentResponse = await ai.models.generateContent({
            model: 'gemini-2.5-flash',
            contents: code,
            config: {
                systemInstruction: systemPrompt
            }
        });
        return response.text;
    } catch (error) {
        console.error("Error refactoring code:", error);
        throw new Error("Failed to refactor code using Gemini API.");
    }
  }
};