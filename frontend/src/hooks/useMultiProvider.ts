/**
 * Hook for managing multi-provider configuration and state
 */

import { useEffect, useCallback } from 'react'
import { multiProviderService, MultiProviderConfig } from '../services/multiProviderService'

export function useMultiProvider() {
  const loadPersistedConfig = useCallback(async () => {
    try {
      const storedConfig = localStorage.getItem('multiProviderConfig')
      if (storedConfig) {
        const config: MultiProviderConfig = JSON.parse(storedConfig)
        await multiProviderService.configure(config)
        await multiProviderService.refreshAvailableProviders()
        console.log('Multi-provider configuration loaded from localStorage')
      }
    } catch (error) {
      console.error('Failed to load persisted multi-provider config:', error)
    }
  }, [])

  const saveConfig = useCallback(async (config: MultiProviderConfig) => {
    try {
      await multiProviderService.configure(config)
      localStorage.setItem('multiProviderConfig', JSON.stringify(config))
      await multiProviderService.refreshAvailableProviders()
    } catch (error) {
      console.error('Failed to save multi-provider config:', error)
      throw error
    }
  }, [])

  const clearConfig = useCallback(() => {
    localStorage.removeItem('multiProviderConfig')
    // Reset to empty config
    multiProviderService.configure({ providers: {} })
  }, [])

  // Load configuration on mount
  useEffect(() => {
    loadPersistedConfig()
  }, [loadPersistedConfig])

  return {
    loadPersistedConfig,
    saveConfig,
    clearConfig,
    service: multiProviderService
  }
}