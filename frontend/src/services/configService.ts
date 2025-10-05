// Configuration service for Grompt frontend
// Handles communication with backend /api/config endpoint

export interface ProviderInfo {
  name: string;
  display_name: string;
  available: boolean;
  configured: boolean;
  models: string[];
  endpoint?: string;
  status: 'ready' | 'needs_api_key' | 'offline';
}

export interface ServerConfig {
  server: {
    name: string;
    version: string;
    port: string;
    status: string;
  };
  providers: Record<string, ProviderInfo>;
  available_providers: string[];
  default_provider: string;
  environment: {
    demo_mode: boolean;
  };
  // Backwards compatibility
  openai_available: boolean;
  deepseek_available: boolean;
  ollama_available: boolean;
  claude_available: boolean;
  gemini_available: boolean;
  chatgpt_available: boolean;
}

class ConfigService {
  private config: ServerConfig | null = null;
  private cache: Map<string, any> = new Map();
  private cacheTimeout = 5 * 60 * 1000; // 5 minutes

  /**
   * Get server configuration from backend
   */
  async getConfig(forceRefresh = false): Promise<ServerConfig> {
    const cacheKey = 'server_config';
    const cached = this.cache.get(cacheKey);

    if (!forceRefresh && cached && (Date.now() - cached.timestamp) < this.cacheTimeout) {
      return cached.data;
    }

    try {
      const response = await fetch('/api/config', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`Config fetch failed: ${response.statusText}`);
      }

      const config: ServerConfig = await response.json();

      // Cache the result
      this.cache.set(cacheKey, {
        data: config,
        timestamp: Date.now(),
      });

      this.config = config;
      return config;
    } catch (error) {
      console.error('Failed to fetch config from backend:', error);

      // Return demo configuration as fallback
      return this.getDemoConfig();
    }
  }

  /**
   * Get available providers
   */
  async getAvailableProviders(): Promise<ProviderInfo[]> {
    const config = await this.getConfig();
    return config.available_providers.map(name => config.providers[name]);
  }

  /**
   * Get specific provider info
   */
  async getProvider(name: string): Promise<ProviderInfo | null> {
    const config = await this.getConfig();
    return config.providers[name] || null;
  }

  /**
   * Get default provider name
   */
  async getDefaultProvider(): Promise<string> {
    const config = await this.getConfig();
    return config.default_provider;
  }

  /**
   * Check if running in demo mode
   */
  async isDemoMode(): Promise<boolean> {
    const config = await this.getConfig();
    return config.environment.demo_mode;
  }

  /**
   * Update provider configuration (for future use)
   */
  async updateProviderConfig(provider: string, apiKey: string): Promise<boolean> {
    try {
      const response = await fetch('/api/config', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          [`${provider}_api_key`]: apiKey,
        }),
      });

      if (response.ok) {
        // Clear cache to force refresh
        this.cache.clear();
        return true;
      }
      return false;
    } catch (error) {
      console.error('Failed to update provider config:', error);
      return false;
    }
  }

  /**
   * Get server status
   */
  async getServerStatus(): Promise<{
    name: string;
    version: string;
    status: string;
    port: string;
  }> {
    const config = await this.getConfig();
    return config.server;
  }

  /**
   * Clear service cache
   */
  clearCache(): void {
    this.cache.clear();
  }

  /**
   * Fallback demo configuration
   */
  private getDemoConfig(): ServerConfig {
    return {
      server: {
        name: 'Grompt Server (Demo)',
        version: '1.0.0',
        port: '8080',
        status: 'demo',
      },
      providers: {
        gemini: {
          name: 'gemini',
          display_name: 'Google Gemini',
          available: false,
          configured: false,
          models: ['gemini-2.0-flash', 'gemini-2.0-flash-exp'],
          status: 'needs_api_key',
        },
      },
      available_providers: [],
      default_provider: 'gemini',
      environment: {
        demo_mode: true,
      },
      openai_available: false,
      deepseek_available: false,
      ollama_available: false,
      claude_available: false,
      gemini_available: false,
      chatgpt_available: false,
    };
  }
}

// Export singleton instance
export const configService = new ConfigService();
