/**
 * Basic API Integration Tests for Grompt Frontend
 * Tests the main API functionality and React hooks
 */

import { describe, test, expect, beforeEach, jest } from '@jest/globals';
import { renderHook, act, waitFor } from '@testing-library/react';
import { GromptAPI, APIError, isAPIError } from '../services/api';
import { useGromptAPI, useGeneratePrompt, useProviders } from '../hooks/useGromptAPI';

// Mock fetch globally
global.fetch = jest.fn();
global.EventSource = jest.fn();

describe('GromptAPI Service', () => {
  let api: GromptAPI;

  beforeEach(() => {
    api = new GromptAPI('http://localhost:8080');
    (fetch as jest.Mock).mockClear();
  });

  describe('generatePrompt', () => {
    test('should make correct API call for prompt generation', async () => {
      const mockResponse = {
        id: 'gen_123',
        object: 'prompt.generation',
        createdAt: Date.now(),
        provider: 'claude',
        model: 'claude-sonnet-4',
        prompt: 'Generated prompt content',
        ideas: ['test idea'],
        purpose: 'code',
        usage: { tokens: 100, costUSD: 0.001 }
      };

      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse
      });

      const request = {
        provider: 'claude',
        ideas: ['test idea'],
        purpose: 'code' as const
      };

      const result = await api.generatePrompt(request);

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/v1/generate',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Content-Type': 'application/json'
          }),
          body: JSON.stringify(request)
        })
      );

      expect(result).toEqual(mockResponse);
    });

    test('should handle API errors correctly', async () => {
      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 400,
        statusText: 'Bad Request',
        json: async () => ({ error: 'INVALID_REQUEST', message: 'Provider is required' })
      });

      const request = {
        provider: '',
        ideas: ['test idea'],
        purpose: 'code' as const
      };

      await expect(api.generatePrompt(request)).rejects.toThrow(APIError);
    });

    test('should handle network errors', async () => {
      (fetch as jest.Mock).mockRejectedValueOnce(new Error('Network error'));

      const request = {
        provider: 'claude',
        ideas: ['test idea'],
        purpose: 'code' as const
      };

      await expect(api.generatePrompt(request)).rejects.toThrow(APIError);
    });
  });

  describe('listProviders', () => {
    test('should fetch providers list correctly', async () => {
      const mockProviders = {
        object: 'list',
        data: [
          { name: 'claude', available: true, type: 'anthropic', defaultModel: 'claude-sonnet-4' },
          { name: 'openai', available: false, type: 'openai', error: 'API key not configured' }
        ],
        hasMore: false,
        timestamp: Date.now()
      };

      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockProviders
      });

      const result = await api.listProviders();

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/v1/providers',
        expect.objectContaining({
          method: undefined // GET is default
        })
      );

      expect(result).toEqual(mockProviders);
    });
  });

  describe('healthCheck', () => {
    test('should check API health', async () => {
      const mockHealth = {
        status: 'healthy' as const,
        service: 'grompt-v1',
        timestamp: Date.now(),
        version: '1.0.0',
        dependencies: {
          providers: {
            claude: { status: 'healthy' as const },
            openai: { status: 'unhealthy' as const, error: 'API key invalid' }
          },
          gobe_proxy: { status: 'not_configured' as const }
        }
      };

      (fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockHealth
      });

      const result = await api.healthCheck();
      expect(result).toEqual(mockHealth);
    });
  });

  describe('rate limiting', () => {
    test('should enforce rate limits', async () => {
      const api = new GromptAPI('http://localhost:8080', {
        rateLimitRequests: 1,
        rateLimitWindowMs: 1000
      });

      (fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: async () => ({ providers: [] })
      });

      // First request should succeed
      await api.listProviders();

      // Second request should be rate limited
      await expect(api.listProviders()).rejects.toThrow('Rate limit exceeded');
    });
  });
});

describe('useGeneratePrompt Hook', () => {
  beforeEach(() => {
    (fetch as jest.Mock).mockClear();
  });

  test('should handle successful prompt generation', async () => {
    const mockResponse = {
      id: 'gen_123',
      object: 'prompt.generation',
      createdAt: Date.now(),
      provider: 'claude',
      model: 'claude-sonnet-4',
      prompt: 'Generated prompt content',
      ideas: ['test idea'],
      purpose: 'code',
      usage: { tokens: 100, costUSD: 0.001 }
    };

    (fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => mockResponse
    });

    const { result } = renderHook(() => useGeneratePrompt());

    expect(result.current.loading).toBe(false);
    expect(result.current.data).toBe(null);
    expect(result.current.error).toBe(null);

    await act(async () => {
      await result.current.generateSync({
        provider: 'claude',
        ideas: ['test idea'],
        purpose: 'code'
      });
    });

    expect(result.current.loading).toBe(false);
    expect(result.current.data).toEqual(mockResponse);
    expect(result.current.error).toBe(null);
  });

  test('should handle generation errors', async () => {
    (fetch as jest.Mock).mockResolvedValueOnce({
      ok: false,
      status: 500,
      statusText: 'Internal Server Error',
      json: async () => ({ error: 'GENERATION_FAILED', message: 'Provider unavailable' })
    });

    const { result } = renderHook(() => useGeneratePrompt());

    await act(async () => {
      try {
        await result.current.generateSync({
          provider: 'claude',
          ideas: ['test idea'],
          purpose: 'code'
        });
      } catch (error) {
        // Expected to throw
      }
    });

    expect(result.current.loading).toBe(false);
    expect(result.current.data).toBe(null);
    expect(result.current.error).toBeInstanceOf(APIError);
  });

  test('should allow cancellation of requests', async () => {
    const { result } = renderHook(() => useGeneratePrompt());

    // Start a request
    act(() => {
      result.current.generateSync({
        provider: 'claude',
        ideas: ['test idea'],
        purpose: 'code'
      });
    });

    expect(result.current.loading).toBe(true);

    // Cancel the request
    act(() => {
      result.current.cancel();
    });

    expect(result.current.loading).toBe(false);
  });

  test('should reset state correctly', async () => {
    const { result } = renderHook(() => useGeneratePrompt());

    // Set some state
    act(() => {
      result.current.reset();
    });

    expect(result.current.data).toBe(null);
    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBe(null);
    expect(result.current.progress.isStreaming).toBe(false);
    expect(result.current.progress.content).toBe('');
  });
});

describe('useProviders Hook', () => {
  beforeEach(() => {
    (fetch as jest.Mock).mockClear();
  });

  test('should auto-fetch providers on mount', async () => {
    const mockProviders = {
      object: 'list',
      data: [
        { name: 'claude', available: true, type: 'anthropic', defaultModel: 'claude-sonnet-4' },
        { name: 'openai', available: false, type: 'openai', error: 'API key not configured' }
      ],
      hasMore: false,
      timestamp: Date.now()
    };

    (fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => mockProviders
    });

    const { result } = renderHook(() => useProviders(true));

    expect(result.current.loading).toBe(true);

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.providers).toEqual(mockProviders.data);
    expect(result.current.error).toBe(null);
  });

  test('should filter available providers', async () => {
    const mockProviders = {
      object: 'list',
      data: [
        { name: 'claude', available: true, type: 'anthropic', defaultModel: 'claude-sonnet-4' },
        { name: 'openai', available: false, type: 'openai', error: 'API key not configured' }
      ],
      hasMore: false,
      timestamp: Date.now()
    };

    (fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => mockProviders
    });

    const { result } = renderHook(() => useProviders(true));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    const availableProviders = result.current.getAvailableProviders();
    expect(availableProviders).toHaveLength(1);
    expect(availableProviders[0].name).toBe('claude');
  });

  test('should find providers by name', async () => {
    const mockProviders = {
      object: 'list',
      data: [
        { name: 'claude', available: true, type: 'anthropic', defaultModel: 'claude-sonnet-4' }
      ],
      hasMore: false,
      timestamp: Date.now()
    };

    (fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => mockProviders
    });

    const { result } = renderHook(() => useProviders(true));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    const claudeProvider = result.current.getProviderByName('claude');
    expect(claudeProvider).toBeDefined();
    expect(claudeProvider?.name).toBe('claude');

    const nonExistentProvider = result.current.getProviderByName('nonexistent');
    expect(nonExistentProvider).toBeUndefined();
  });
});

describe('API Error Utilities', () => {
  test('should identify API errors correctly', () => {
    const apiError = new APIError('Test error', 400, 'TEST_ERROR');
    const regularError = new Error('Regular error');

    expect(isAPIError(apiError)).toBe(true);
    expect(isAPIError(regularError)).toBe(false);
    expect(isAPIError(null)).toBe(false);
    expect(isAPIError({})).toBe(false);
  });

  test('should create API errors with correct properties', () => {
    const apiError = new APIError('Test message', 400, 'TEST_CODE', { extra: 'data' });

    expect(apiError.message).toBe('Test message');
    expect(apiError.status).toBe(400);
    expect(apiError.code).toBe('TEST_CODE');
    expect(apiError.details).toEqual({ extra: 'data' });
    expect(apiError.name).toBe('APIError');
  });
});

describe('Integration Tests', () => {
  test('should handle full prompt generation workflow', async () => {
    const mockProviders = {
      object: 'list',
      data: [
        { name: 'claude', available: true, type: 'anthropic', defaultModel: 'claude-sonnet-4' }
      ],
      hasMore: false,
      timestamp: Date.now()
    };

    const mockGeneration = {
      id: 'gen_123',
      object: 'prompt.generation',
      createdAt: Date.now(),
      provider: 'claude',
      model: 'claude-sonnet-4',
      prompt: 'Generated prompt content',
      ideas: ['test idea'],
      purpose: 'code',
      usage: { tokens: 100, costUSD: 0.001 }
    };

    (fetch as jest.Mock)
      .mockResolvedValueOnce({
        ok: true,
        json: async () => mockProviders
      })
      .mockResolvedValueOnce({
        ok: true,
        json: async () => mockGeneration
      });

    const { result } = renderHook(() => useGromptAPI({
      autoFetchProviders: true
    }));

    // Wait for providers to load
    await waitFor(() => {
      expect(result.current.providers.loading).toBe(false);
    });

    expect(result.current.providers.providers).toHaveLength(1);

    // Generate a prompt
    await act(async () => {
      await result.current.generatePrompt.generateSync({
        provider: 'claude',
        ideas: ['test idea'],
        purpose: 'code'
      });
    });

    expect(result.current.generatePrompt.data).toEqual(mockGeneration);
    expect(result.current.generatePrompt.error).toBe(null);
  });
});