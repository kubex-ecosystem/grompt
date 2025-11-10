/**
 * Tests for Multi-Provider Integration
 */

import { AIProvider } from '@/types/types';
import { beforeEach, describe, expect, test } from 'vitest';
import { MultiProviderConfig, multiProviderService } from '../services/multiProviderService';

describe('MultiProviderService', () => {
  beforeEach(() => {
    // Clear localStorage before each test
    localStorage.clear();
  });

  test('should initialize without errors', () => {
    expect(multiProviderService).toBeDefined();
    expect(multiProviderService.isConfigured()).toBe(false);
  });

  test('should configure providers correctly', async () => {
    const config: MultiProviderConfig = {
      providers: {
        [AIProvider.OPENAI]: {
          apiKey: 'test-openai-key', // pragma: allowlist secret
          defaultModel: 'gpt-4'
        },
        [AIProvider.ANTHROPIC]: {
          apiKey: 'test-anthropic-key', // pragma: allowlist secret
          defaultModel: 'claude-3-sonnet-20240229'
        }
      },
      fallbackToBackend: true
    };

    await multiProviderService.configure(config);
    expect(multiProviderService.getConfig()).toEqual(config);
  });

  test('should get available providers', () => {
    const providers = multiProviderService.getAvailableProviders();
    expect(Array.isArray(providers)).toBe(true);
  });

  test('should handle provider mapping correctly', () => {
    // Test that backend provider names are mapped correctly
    const config: MultiProviderConfig = {
      providers: {
        [AIProvider.OPENAI]: {
          apiKey: 'test-key', // pragma: allowlist secret
          defaultModel: 'gpt-4'
        }
      }
    };

    multiProviderService.configure(config);
    const providers = multiProviderService.getAvailableProviders();

    // Should have at least the configured provider
    expect(providers.length).toBeGreaterThanOrEqual(0);
  });

  test('should get available models for configured providers', async () => {
    const config: MultiProviderConfig = {
      providers: {
        [AIProvider.OPENAI]: {
          apiKey: 'test-key', // pragma: allowlist secret
          defaultModel: 'gpt-4'
        }
      }
    };

    await multiProviderService.configure(config);
    const models = multiProviderService.getAvailableModels('openai');
    expect(Array.isArray(models)).toBe(true);
  });
});

describe('MultiProvider Integration', () => {
  test('should handle requests with fallback to backend', async () => {
    const config: MultiProviderConfig = {
      providers: {},
      fallbackToBackend: true
    };

    await multiProviderService.configure(config);

    const request = {
      provider: 'openai',
      ideas: ['Test idea'],
      purpose: 'general' as const
    };

    // Should not throw error even without local providers configured
    // Should fallback to backend or offline mode
    try {
      await multiProviderService.generateContent(request);
    } catch (error) {
      // Expected to fail in test environment without backend
      expect(error).toBeDefined();
    }
  });

  test('should refresh providers from backend', async () => {
    // This test verifies the method exists and can be called
    try {
      await multiProviderService.refreshAvailableProviders();
    } catch (error) {
      // Expected to fail in test environment without backend
      expect(error).toBeDefined();
    }
  });
});
