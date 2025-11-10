/**
 * Multi-Provider Configuration Component
 * Allows users to configure API keys and settings for all available providers
 */

import { AlertCircle, Check, Eye, EyeOff, Key, Settings, TestTube, X } from 'lucide-react'
import * as React from 'react'
import { useEffect, useState } from 'react'
import { multiProviderService, type MultiProviderConfig } from '../../services/multiProviderService'
import { AIProvider } from '@/types/types'

export interface MultiProviderConfigProps {
  isOpen: boolean
  onClose: () => void
  onConfigUpdate?: (config: MultiProviderConfig) => void
}

interface ProviderState {
  [AIProvider.OPENAI]: {
    apiKey: string
    defaultModel: string
    baseURL?: string
    testing: boolean
    testResult?: boolean
    showKey: boolean
  }
  [AIProvider.ANTHROPIC]: {
    apiKey: string
    defaultModel: string
    baseURL?: string
    testing: boolean
    testResult?: boolean
    showKey: boolean
  }
  [AIProvider.GEMINI]: {
    apiKey: string
    defaultModel: string
    baseURL?: string
    testing: boolean
    testResult?: boolean
    showKey: boolean
  }
}

const DEFAULT_MODELS = {
  [AIProvider.OPENAI]: 'gpt-4',
  [AIProvider.ANTHROPIC]: 'claude-3-sonnet-20240229',
  [AIProvider.GEMINI]: 'gemini-pro'
}

const PROVIDER_LABELS = {
  [AIProvider.OPENAI]: 'OpenAI',
  [AIProvider.ANTHROPIC]: 'Anthropic (Claude)',
  [AIProvider.GEMINI]: 'Google Gemini'
}

export function MultiProviderConfig({ isOpen, onClose, onConfigUpdate }: MultiProviderConfigProps) {
  const [providers, setProviders] = useState<ProviderState>({
    [AIProvider.OPENAI]: {
      apiKey: '',
      defaultModel: DEFAULT_MODELS[AIProvider.OPENAI],
      testing: false,
      showKey: false
    },
    [AIProvider.ANTHROPIC]: {
      apiKey: '',
      defaultModel: DEFAULT_MODELS[AIProvider.ANTHROPIC],
      testing: false,
      showKey: false
    },
    [AIProvider.GEMINI]: {
      apiKey: '',
      defaultModel: DEFAULT_MODELS[AIProvider.GEMINI],
      testing: false,
      showKey: false
    }
  })

  const [saving, setSaving] = useState(false)

  // Load existing configuration on mount
  useEffect(() => {
    const existingConfig = multiProviderService.getConfig()
    if (existingConfig) {
      setProviders(prev => {
        const updated = { ...prev }

        Object.keys(existingConfig.providers).forEach(providerKey => {
          const provider = providerKey as AIProvider
          const providerConfig = existingConfig.providers[provider]
          if (providerConfig) {
            updated[provider] = {
              ...updated[provider],
              apiKey: providerConfig.apiKey || '',
              defaultModel: providerConfig.defaultModel || DEFAULT_MODELS[provider],
              baseURL: providerConfig.baseURL
            }
          }
        })

        return updated
      })
    }
  }, [isOpen])

  const updateProvider = (provider: AIProvider, field: keyof ProviderState[AIProvider], value: any) => {
    setProviders(prev => ({
      ...prev,
      [provider]: {
        ...prev[provider],
        [field]: value
      }
    }))
  }

  const testProvider = async (provider: AIProvider) => {
    updateProvider(provider, 'testing', true)
    updateProvider(provider, 'testResult', undefined)

    try {
      // Create a temporary config for testing
      const testConfig: MultiProviderConfig = {
        providers: {
          [provider]: {
            apiKey: providers[provider].apiKey,
            defaultModel: providers[provider].defaultModel,
            baseURL: providers[provider].baseURL
          }
        },
        fallbackToBackend: false
      }

      // Configure the service temporarily
      await multiProviderService.configure(testConfig)

      // Test the provider
      const backendProviderName = provider === AIProvider.ANTHROPIC ? 'anthropic' :
        provider === AIProvider.GEMINI ? 'gemini' : 'openai'
      const result = await multiProviderService.testProvider(backendProviderName)

      updateProvider(provider, 'testResult', result)
    } catch (error) {
      console.error(`Failed to test ${provider}:`, error)
      updateProvider(provider, 'testResult', false)
    } finally {
      updateProvider(provider, 'testing', false)
    }
  }

  const saveConfiguration = async () => {
    setSaving(true)

    try {
      const config: MultiProviderConfig = {
        providers: {},
        fallbackToBackend: true,
        cacheResponses: true
      }

      // Only include providers with API keys
      Object.entries(providers).forEach(([providerKey, providerState]) => {
        if (providerState.apiKey.trim()) {
          const provider = providerKey as AIProvider
          config.providers[provider] = {
            apiKey: providerState.apiKey.trim(),
            defaultModel: providerState.defaultModel,
            baseURL: providerState.baseURL
          }
        }
      })

      // Configure the service
      await multiProviderService.configure(config)

      // Refresh providers
      await multiProviderService.refreshAvailableProviders()

      // Notify parent component
      onConfigUpdate?.(config)

      // Store in localStorage for persistence
      localStorage.setItem('multiProviderConfig', JSON.stringify(config))

      onClose()
    } catch (error) {
      console.error('Failed to save configuration:', error)
      alert('Falha ao salvar configuração. Verifique as chaves de API.')
    } finally {
      setSaving(false)
    }
  }

  const toggleKeyVisibility = (provider: AIProvider) => {
    updateProvider(provider, 'showKey', !providers[provider].showKey)
  }

  const getTestIcon = (provider: AIProvider) => {
    const { testing, testResult } = providers[provider]

    if (testing) {
      return <div className="animate-spin w-4 h-4 border-2 border-blue-400 border-t-transparent rounded-full" />
    }

    if (testResult === true) {
      return <Check size={16} className="text-green-400" />
    }

    if (testResult === false) {
      return <X size={16} className="text-red-400" />
    }

    return <TestTube size={16} className="text-gray-400" />
  }

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-gray-800 rounded-xl border border-gray-700 max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center gap-3">
              <Settings className="text-purple-400" size={24} />
              <h2 className="text-xl font-bold text-white">Configuração dos Providers</h2>
            </div>
            <button
              title='Fechar'
              onClick={onClose}
              className="text-gray-400 hover:text-white p-2"
            >
              <X size={20} />
            </button>
          </div>

          <div className="space-y-6">
            {Object.entries(providers).map(([providerKey, providerState]) => {
              const provider = providerKey as AIProvider
              return (
                <div key={provider} className="bg-gray-900/50 rounded-lg border border-gray-700 p-4">
                  <div className="flex items-center gap-3 mb-4">
                    <Key className="text-blue-400" size={20} />
                    <h3 className="text-lg font-semibold text-white">
                      {PROVIDER_LABELS[provider]}
                    </h3>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {/* API Key */}
                    <div>
                      <label className="block text-sm font-medium text-gray-300 mb-2">
                        Chave da API
                      </label>
                      <div className="relative">
                        <input
                          type={providerState.showKey ? 'text' : 'password'}
                          value={providerState.apiKey}
                          onChange={(e) => updateProvider(provider, 'apiKey', e.target.value)}
                          placeholder={`Digite sua chave da API ${PROVIDER_LABELS[provider]}`}
                          className="w-full px-3 py-2 pr-10 rounded-lg border border-gray-600 bg-gray-700 text-white focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
                        />
                        <button
                          type="button"
                          onClick={() => toggleKeyVisibility(provider)}
                          className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-white"
                        >
                          {providerState.showKey ? <EyeOff size={16} /> : <Eye size={16} />}
                        </button>
                      </div>
                    </div>

                    {/* Default Model */}
                    <div>
                      <label className="block text-sm font-medium text-gray-300 mb-2">
                        Modelo Padrão
                      </label>
                      <input
                        title='Digite o modelo padrão'
                        type="text"
                        value={providerState.defaultModel}
                        onChange={(e) => updateProvider(provider, 'defaultModel', e.target.value)}
                        className="w-full px-3 py-2 rounded-lg border border-gray-600 bg-gray-700 text-white focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
                      />
                    </div>

                    {/* Base URL (optional) */}
                    <div>
                      <label className="block text-sm font-medium text-gray-300 mb-2">
                        URL Base (opcional)
                      </label>
                      <input
                        type="url"
                        value={providerState.baseURL || ''}
                        onChange={(e) => updateProvider(provider, 'baseURL', e.target.value)}
                        placeholder="https://api.example.com"
                        className="w-full px-3 py-2 rounded-lg border border-gray-600 bg-gray-700 text-white focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
                      />
                    </div>

                    {/* Test Button */}
                    <div className="flex items-end">
                      <button
                        onClick={() => testProvider(provider)}
                        disabled={!providerState.apiKey || providerState.testing}
                        className="w-full flex items-center justify-center gap-2 px-4 py-2 rounded-lg bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white"
                      >
                        {getTestIcon(provider)}
                        Testar Conexão
                      </button>
                    </div>
                  </div>

                  {/* Test Result */}
                  {providerState.testResult !== undefined && (
                    <div className={`mt-3 p-3 rounded-lg flex items-center gap-2 ${providerState.testResult
                      ? 'bg-green-900/50 text-green-400 border border-green-700'
                      : 'bg-red-900/50 text-red-400 border border-red-700'
                      }`}>
                      {providerState.testResult ? (
                        <>
                          <Check size={16} />
                          Conexão bem-sucedida! Provider configurado corretamente.
                        </>
                      ) : (
                        <>
                          <AlertCircle size={16} />
                          Falha na conexão. Verifique a chave da API e tente novamente.
                        </>
                      )}
                    </div>
                  )}
                </div>
              )
            })}
          </div>

          {/* Action Buttons */}
          <div className="flex justify-end gap-3 mt-6 pt-6 border-t border-gray-700">
            <button
              onClick={onClose}
              className="px-4 py-2 rounded-lg border border-gray-600 text-gray-300 hover:bg-gray-700"
            >
              Cancelar
            </button>
            <button
              onClick={saveConfiguration}
              disabled={saving}
              className="px-4 py-2 rounded-lg bg-purple-600 hover:bg-purple-700 disabled:bg-gray-600 text-white flex items-center gap-2"
            >
              {saving ? (
                <div className="animate-spin w-4 h-4 border-2 border-white border-t-transparent rounded-full" />
              ) : (
                <Settings size={16} />
              )}
              {saving ? 'Salvando...' : 'Salvar Configuração'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
