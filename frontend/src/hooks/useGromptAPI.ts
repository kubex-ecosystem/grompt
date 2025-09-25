/**
 * React hooks for Grompt API integration
 * Provides React-friendly interface to the Grompt backend API
 */

import { useCallback, useEffect, useRef, useState } from 'react'
import {
  APIError,
  GenerateRequest,
  GenerateResponse,
  HealthResponse,
  isAPIError,
  Provider
} from '../services/api'
import { enhancedAPI } from '../services/enhancedAPI'

interface UseGeneratePromptState {
  generateStream: (request: GenerateRequest) => Promise<void>
  generateSync: (request: GenerateRequest) => Promise<GenerateResponse>
  cancel: () => void
  reset: () => void
  data: GenerateResponse | null
  loading: boolean
  error: APIError | null
  progress: {
    isStreaming: boolean
    content: string
    usage?: { tokens: number; costUSD: number }
  }
}

interface UseProvidersState {
  providers: Provider[]
  loading: boolean
  error: APIError | null
  lastFetched: number | null
}

interface UseHealthState {
  isHealthy: boolean
  health: HealthResponse | null
  loading: boolean
  error: APIError | null
  lastChecked: number | null
}

/**
 * Hook for generating prompts with both sync and streaming support
 */
function useGeneratePrompt() {
  const [state, setState] = useState<UseGeneratePromptState>({
    data: null,
    loading: false,
    error: null,
    progress: {
      isStreaming: false,
      content: '',
      usage: undefined
    },
    generateStream: async () => { },
    generateSync: async () => { return Promise.resolve({} as GenerateResponse) },
    cancel: () => { },
    reset: () => { }
  })

  const abortControllerRef = useRef<AbortController | null>(null)

  const generateSync = useCallback(async (request: GenerateRequest) => {
    // Cancel any ongoing request
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
    }

    setState(prev => ({
      ...prev,
      loading: true,
      error: null,
      data: null,
      progress: { isStreaming: false, content: '', usage: undefined }
    }))

    try {
      const response = await enhancedAPI.generatePrompt(request)
      setState(prev => ({
        ...prev,
        loading: false,
        data: response,
        progress: { isStreaming: false, content: response.prompt, usage: response.usage }
      }))
      return response
    } catch (error) {
      const apiError = isAPIError(error) ? error : new APIError(
        error instanceof Error ? error.message : 'Unknown error',
        0,
        'UNKNOWN_ERROR'
      )
      setState(prev => ({
        ...prev,
        loading: false,
        error: apiError
      }))
      throw apiError
    }
  }, [])

  const generateStream = useCallback(async (request: GenerateRequest) => {
    // Cancel any ongoing request
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
    }

    abortControllerRef.current = new AbortController()

    setState(prev => ({
      ...prev,
      loading: true,
      error: null,
      data: null,
      progress: { isStreaming: true, content: '', usage: undefined }
    }))

    try {
      await enhancedAPI.generatePromptStream(
        request,
        (content) => {
          setState(prev => ({
            ...prev,
            progress: {
              ...prev.progress,
              content: prev.progress.content + content
            }
          }))
        },
        (usage) => {
          setState(prev => {
            const finalData: GenerateResponse = {
              id: `stream_${Date.now()}`,
              object: 'prompt.generation',
              createdAt: Date.now() / 1000,
              provider: request.provider,
              model: request.model || '',
              prompt: prev.progress.content,
              ideas: request.ideas,
              purpose: request.purpose || 'general',
              usage,
              metadata: {
                temperature: request.temperature,
                maxTokens: request.maxTokens
              }
            }

            return {
              ...prev,
              data: finalData,
              progress: {
                isStreaming: false,
                content: prev.progress.content,
                usage
              }
            }
          })
        },
        (error) => {
          setState(prev => ({
            ...prev,
            loading: false,
            error: new APIError(error, 0, 'STREAM_ERROR'),
            progress: {
              isStreaming: false,
              content: prev.progress.content,
              usage: undefined
            }
          }))
        }
      )
    } catch (error) {
      const apiError = isAPIError(error) ? error : new APIError(
        error instanceof Error ? error.message : 'Streaming error',
        0,
        'STREAM_ERROR'
      )
      setState(prev => ({
        ...prev,
        loading: false,
        error: apiError,
        progress: {
          isStreaming: false,
          content: prev.progress.content,
          usage: undefined
        }
      }))
      throw apiError
    }
  }, [])

  const cancel = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
      abortControllerRef.current = null
    }
    setState(prev => ({
      ...prev,
      loading: false,
      progress: {
        ...prev.progress,
        isStreaming: false
      }
    }))
  }, [])

  const reset = useCallback(() => {
    cancel()
    setState({
      data: null,
      loading: false,
      error: null,
      progress: {
        isStreaming: false,
        content: '',
        usage: undefined
      },
      generateStream,
      generateSync,
      cancel,
      reset() {
        return;
      },
    })
  }, [cancel])

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort()
      }
    }
  }, [])

  return {
    ...state,
    generateSync,
    generateStream,
    cancel,
    reset
  }
}

/**
 * Hook for fetching and managing providers list
 */
function useProviders(autoFetch: boolean = true) {
  const [state, setState] = useState<UseProvidersState>({
    providers: [],
    loading: false,
    error: null,
    lastFetched: null
  })

  const fetchProviders = useCallback(async (force: boolean = false) => {
    // Don't refetch if recently fetched (unless forced)
    if (!force && state.lastFetched && Date.now() - state.lastFetched < 30000) {
      return state.providers
    }

    setState(prev => ({
      ...prev,
      loading: true,
      error: null
    }))

    try {
      const providers = await enhancedAPI.getProviders()
      setState(prev => ({
        ...prev,
        loading: false,
        providers: providers,
        lastFetched: Date.now()
      }))
      return providers
    } catch (error) {
      const apiError = isAPIError(error) ? error : new APIError(
        error instanceof Error ? error.message : 'Failed to fetch providers',
        0,
        'FETCH_PROVIDERS_ERROR'
      )
      setState(prev => ({
        ...prev,
        loading: false,
        error: apiError
      }))
      throw apiError
    }
  }, [state.lastFetched, state.providers])

  // Auto-fetch on mount if enabled
  useEffect(() => {
    if (autoFetch && state.providers.length === 0 && !state.loading) {
      fetchProviders()
    }
  }, [autoFetch, state.providers.length, state.loading, fetchProviders])

  const getAvailableProviders = useCallback(() => {
    return state.providers.filter(provider => provider.available)
  }, [state.providers])

  const getProviderByName = useCallback((name: string) => {
    return state.providers.find(provider => provider.name === name)
  }, [state.providers])

  return {
    ...state,
    fetchProviders,
    getAvailableProviders,
    getProviderByName,
    refresh: () => fetchProviders(true)
  }
}

/**
 * Hook for health monitoring
 */
function useHealth(autoCheck: boolean = false, intervalMs: number = 60000) {
  const [state, setState] = useState<UseHealthState>({
    isHealthy: false,
    health: null,
    loading: false,
    error: null,
    lastChecked: null
  })

  const intervalRef = useRef<NodeJS.Timeout | null>(null)

  const checkHealth = useCallback(async () => {
    setState(prev => ({
      ...prev,
      loading: true,
      error: null
    }))

    try {
      const health = await enhancedAPI.getHealth()
      setState(prev => ({
        ...prev,
        loading: false,
        health,
        lastChecked: Date.now()
      }))
      return health
    } catch (error) {
      const apiError = isAPIError(error) ? error : new APIError(
        error instanceof Error ? error.message : 'Health check failed',
        0,
        'HEALTH_CHECK_ERROR'
      )
      setState(prev => ({
        ...prev,
        loading: false,
        error: apiError,
        lastChecked: Date.now()
      }))
      throw apiError
    }
  }, [])

  // Auto health checking
  useEffect(() => {
    if (autoCheck) {
      // Initial check
      checkHealth()

      // Set up interval
      intervalRef.current = setInterval(checkHealth, intervalMs)

      return () => {
        if (intervalRef.current) {
          clearInterval(intervalRef.current)
        }
      }
    }
  }, [autoCheck, intervalMs, checkHealth])

  const isHealthy = useCallback(() => {
    return state.health?.status === 'healthy'
  }, [state.health])

  const getUnhealthyProviders = useCallback(() => {
    if (!state.health?.dependencies?.providers) return []

    return Object.entries(state.health.dependencies.providers)
      .filter(([, status]) => status.status !== 'healthy')
      .map(([name, status]) => ({ name, ...status }))
  }, [state.health])

  return {
    ...state,
    checkHealth,
    isHealthy,
    getUnhealthyProviders
  }
}

/**
 * Hook for rate limit monitoring
 */
function useRateLimit() {
  const [rateLimitStatus, setRateLimitStatus] = useState({
    canMakeRequest: true,
    timeUntilReset: 0
  })

  const updateStatus = useCallback(() => {
    // Enhanced API handles rate limiting internally
    setRateLimitStatus({
      canMakeRequest: enhancedAPI.isConnected(),
      timeUntilReset: 0
    })
  }, [])

  // Update status periodically
  useEffect(() => {
    const interval = setInterval(updateStatus, 1000)
    return () => clearInterval(interval)
  }, [updateStatus])

  return {
    ...rateLimitStatus,
    updateStatus
  }
}

/**
 * Combined hook for comprehensive API state management
 */
export function useGromptAPI(options: {
  autoFetchProviders?: boolean
  autoCheckHealth?: boolean
  healthCheckInterval?: number
} = {}) {
  const generatePrompt = useGeneratePrompt()
  const providers = useProviders(options.autoFetchProviders)
  const health = useHealth(options.autoCheckHealth, options.healthCheckInterval)
  const rateLimit = useRateLimit()

  return {
    generatePrompt,
    providers,
    health,
    rateLimit,
    api: enhancedAPI // Direct access to enhanced API instance
  }
}
