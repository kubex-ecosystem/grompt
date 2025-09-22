/**
 * Multi-Provider Service
 * Integrates the MultiAIWrapper with the backend API and enhancedAPI
 * Provides unified interface for all AI providers (OpenAI, Anthropic, Gemini)
 */

import { MultiAIWrapper } from '../core/llm/wrapper/MultiAIWrapper'
import { AIProvider, AIModel, MultiAIConfig } from '../types/types'
import { enhancedAPI } from './enhancedAPI'
import { GenerateRequest, GenerateResponse, Provider } from './api'

/**
 * Configuration for multi-provider service
 */
export interface MultiProviderConfig {
  providers: {
    [AIProvider.OPENAI]?: {
      apiKey: string
      defaultModel: string
      baseURL?: string
    }
    [AIProvider.ANTHROPIC]?: {
      apiKey: string
      defaultModel: string
      baseURL?: string
    }
    [AIProvider.GEMINI]?: {
      apiKey: string
      defaultModel: string
      baseURL?: string
    }
  }
  fallbackToBackend?: boolean
  cacheResponses?: boolean
}

/**
 * Provider mapping from backend to frontend types
 */
const PROVIDER_MAPPING: Record<string, AIProvider> = {
  'openai': AIProvider.OPENAI,
  'claude': AIProvider.ANTHROPIC,
  'anthropic': AIProvider.ANTHROPIC,
  'gemini': AIProvider.GEMINI,
  'google': AIProvider.GEMINI
}

/**
 * Reverse mapping from frontend to backend types
 */
const BACKEND_PROVIDER_MAPPING: Record<AIProvider, string> = {
  [AIProvider.OPENAI]: 'openai',
  [AIProvider.ANTHROPIC]: 'claude',
  [AIProvider.GEMINI]: 'gemini'
}

/**
 * Multi-Provider Service Class
 * Handles both local MultiAIWrapper and backend API calls
 */
export class MultiProviderService {
  private multiAI: MultiAIWrapper | null = null
  private config: MultiProviderConfig | null = null
  private availableProviders: Provider[] = []

  constructor() {
    this.initialize()
  }

  /**
   * Initialize the service and load available providers
   */
  private async initialize() {
    try {
      // Load available providers from backend
      await this.refreshAvailableProviders()
    } catch (error) {
      console.warn('Failed to load backend providers:', error)
    }
  }

  /**
   * Configure the multi-provider service
   */
  async configure(config: MultiProviderConfig) {
    this.config = config

    // Convert to MultiAIConfig format
    const multiAIConfig: MultiAIConfig = {
      providers: {}
    }

    // Configure OpenAI
    if (config.providers[AIProvider.OPENAI]) {
      const openaiConfig = config.providers[AIProvider.OPENAI]!
      multiAIConfig.providers[AIProvider.OPENAI] = {
        apiKey: openaiConfig.apiKey,
        options: {
          baseURL: openaiConfig.baseURL,
          defaultQuery: {}
        }
      }
    }

    // Configure Anthropic
    if (config.providers[AIProvider.ANTHROPIC]) {
      const anthropicConfig = config.providers[AIProvider.ANTHROPIC]!
      multiAIConfig.providers[AIProvider.ANTHROPIC] = {
        apiKey: anthropicConfig.apiKey,
        options: {
          baseURL: anthropicConfig.baseURL,
          defaultQuery: {}
        }
      }
    }

    // Configure Gemini
    if (config.providers[AIProvider.GEMINI]) {
      const geminiConfig = config.providers[AIProvider.GEMINI]!
      multiAIConfig.providers[AIProvider.GEMINI] = {
        apiKey: geminiConfig.apiKey,
        options: {
          generationConfig: {}
        }
      }
    }

    // Initialize MultiAIWrapper if we have any provider configured
    if (Object.keys(multiAIConfig.providers).length > 0) {
      this.multiAI = new MultiAIWrapper(multiAIConfig)
    }
  }

  /**
   * Refresh available providers from backend
   */
  async refreshAvailableProviders(): Promise<Provider[]> {
    try {
      const response = await enhancedAPI.getProviders()
      this.availableProviders = response
      return response
    } catch (error) {
      console.error('Failed to refresh providers:', error)
      return []
    }
  }

  /**
   * Get all available providers (both backend and local)
   */
  getAvailableProviders(): Provider[] {
    const backendProviders = [...this.availableProviders]
    const localProviders: Provider[] = []

    // Add local providers that are configured
    if (this.config?.providers[AIProvider.OPENAI]) {
      localProviders.push({
        name: 'openai-local',
        available: true,
        type: 'openai',
        defaultModel: this.config.providers[AIProvider.OPENAI]!.defaultModel
      })
    }

    if (this.config?.providers[AIProvider.ANTHROPIC]) {
      localProviders.push({
        name: 'anthropic-local',
        available: true,
        type: 'anthropic',
        defaultModel: this.config.providers[AIProvider.ANTHROPIC]!.defaultModel
      })
    }

    if (this.config?.providers[AIProvider.GEMINI]) {
      localProviders.push({
        name: 'gemini-local',
        available: true,
        type: 'gemini',
        defaultModel: this.config.providers[AIProvider.GEMINI]!.defaultModel
      })
    }

    // Merge and deduplicate
    const allProviders = [...backendProviders, ...localProviders]
    const uniqueProviders = allProviders.filter((provider, index, self) =>
      index === self.findIndex(p => p.name === provider.name)
    )

    return uniqueProviders
  }

  /**
   * Generate content using the best available method
   */
  async generateContent(request: GenerateRequest): Promise<GenerateResponse> {
    const provider = this.mapProviderName(request.provider)

    // Try local MultiAIWrapper first if configured and provider is available locally
    if (this.multiAI && this.isProviderAvailableLocally(provider)) {
      try {
        return await this.generateWithLocalProvider(request, provider)
      } catch (error) {
        console.warn('Local provider failed, falling back to backend:', error)

        // If fallback is disabled, throw the error
        if (!this.config?.fallbackToBackend) {
          throw error
        }
      }
    }

    // Fallback to backend API
    return await enhancedAPI.generatePrompt(request)
  }

  /**
   * Generate content with streaming using the best available method
   */
  async generateContentStream(
    request: GenerateRequest,
    onChunk: (content: string) => void,
    onComplete: (usage?: { tokens: number; costUSD: number }) => void,
    onError: (error: string) => void
  ): Promise<void> {
    const provider = this.mapProviderName(request.provider)

    // Try local MultiAIWrapper first if configured and provider is available locally
    if (this.multiAI && this.isProviderAvailableLocally(provider)) {
      try {
        return await this.generateStreamWithLocalProvider(request, provider, onChunk, onComplete, onError)
      } catch (error) {
        console.warn('Local provider streaming failed, falling back to backend:', error)

        // If fallback is disabled, throw the error
        if (!this.config?.fallbackToBackend) {
          throw error
        }
      }
    }

    // Fallback to backend API
    return await enhancedAPI.generatePromptStream(request, onChunk, onComplete, onError)
  }

  /**
   * Check if provider is available locally
   */
  private isProviderAvailableLocally(provider: AIProvider): boolean {
    return this.config?.providers[provider] !== undefined
  }

  /**
   * Map provider name from backend format to frontend enum
   */
  private mapProviderName(providerName: string): AIProvider {
    return PROVIDER_MAPPING[providerName.toLowerCase()] || AIProvider.ANTHROPIC
  }

  /**
   * Generate content using local MultiAIWrapper
   */
  private async generateWithLocalProvider(request: GenerateRequest, provider: AIProvider): Promise<GenerateResponse> {
    if (!this.multiAI) {
      throw new Error('MultiAI wrapper not initialized')
    }

    const ideas = request.ideas.join('\n')
    const systemInstruction = request.context?.systemInstruction ||
      `Generate a ${request.purpose || 'general'} prompt based on the provided ideas.`

    const response = await this.multiAI.generateContent({
      prompt: ideas,
      systemInstruction,
      provider,
      model: request.model as AIModel,
      options: {
        temperature: request.temperature,
        maxTokens: request.maxTokens
      }
    })

    // Convert to backend response format
    return {
      id: `local_${Date.now()}`,
      object: 'prompt.generation',
      createdAt: Date.now() / 1000,
      provider: BACKEND_PROVIDER_MAPPING[provider],
      model: response.model,
      prompt: response.text,
      ideas: request.ideas,
      purpose: request.purpose || 'general',
      usage: response.usage ? {
        tokens: response.usage.totalTokens || 0,
        costUSD: 0 // We don't track cost locally
      } : undefined,
      metadata: {
        temperature: request.temperature,
        maxTokens: request.maxTokens,
        localGeneration: true
      }
    }
  }

  /**
   * Generate content with streaming using local MultiAIWrapper
   */
  private async generateStreamWithLocalProvider(
    request: GenerateRequest,
    provider: AIProvider,
    onChunk: (content: string) => void,
    onComplete: (usage?: { tokens: number; costUSD: number }) => void,
    onError: (error: string) => void
  ): Promise<void> {
    if (!this.multiAI) {
      throw new Error('MultiAI wrapper not initialized')
    }

    const ideas = request.ideas.join('\n')
    const systemInstruction = request.context?.systemInstruction ||
      `Generate a ${request.purpose || 'general'} prompt based on the provided ideas.`

    try {
      const stream = this.multiAI.streamContent({
        prompt: ideas,
        systemInstruction,
        provider,
        model: request.model as AIModel,
        options: {
          temperature: request.temperature,
          maxTokens: request.maxTokens
        }
      })

      for await (const chunk of stream) {
        onChunk(chunk)
      }

      onComplete({ tokens: 0, costUSD: 0 })
    } catch (error) {
      onError(error instanceof Error ? error.message : 'Unknown streaming error')
      throw error
    }
  }

  /**
   * Test connection to a provider
   */
  async testProvider(providerName: string): Promise<boolean> {
    const provider = this.mapProviderName(providerName)

    if (this.multiAI && this.isProviderAvailableLocally(provider)) {
      try {
        await this.multiAI.generateContent({
          prompt: 'Test',
          provider,
          options: { maxTokens: 10 }
        })
        return true
      } catch (error) {
        console.error(`Local provider ${providerName} test failed:`, error)
        return false
      }
    }

    // Test via backend
    try {
      await enhancedAPI.generatePrompt({
        provider: providerName,
        ideas: ['Test'],
        maxTokens: 10
      })
      return true
    } catch (error) {
      console.error(`Backend provider ${providerName} test failed:`, error)
      return false
    }
  }

  /**
   * Get available models for a provider
   */
  getAvailableModels(providerName: string): string[] {
    const provider = this.mapProviderName(providerName)

    if (this.multiAI && this.isProviderAvailableLocally(provider)) {
      return this.multiAI.getAvailableModels(provider)
    }

    // Fallback to backend provider info
    const backendProvider = this.availableProviders.find(p => p.name === providerName)
    return backendProvider?.defaultModel ? [backendProvider.defaultModel] : []
  }

  /**
   * Check if service is properly configured
   */
  isConfigured(): boolean {
    return this.multiAI !== null || this.availableProviders.length > 0
  }

  /**
   * Get current configuration
   */
  getConfig(): MultiProviderConfig | null {
    return this.config
  }
}

/**
 * Default instance
 */
export const multiProviderService = new MultiProviderService()