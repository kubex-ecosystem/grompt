import { useEffect } from 'react';

interface AnalyticsEvent {
  action: string;
  category: string;
  label?: string;
  value?: number;
}

// Simple, privacy-first analytics
class KubexAnalytics {
  private static instance: KubexAnalytics;
  private isEnabled: boolean = false;

  private constructor() {
    // Only enable in production
    this.isEnabled = window.location.hostname === 'kubex.world';
  }

  public static getInstance(): KubexAnalytics {
    if (!KubexAnalytics.instance) {
      KubexAnalytics.instance = new KubexAnalytics();
    }
    return KubexAnalytics.instance;
  }

  public track(event: AnalyticsEvent): void {
    if (!this.isEnabled) return;

    // Simple privacy-first tracking - no personal data
    const data = {
      ...event,
      timestamp: new Date().toISOString(),
      userAgent: navigator.userAgent.substring(0, 100), // Truncated for privacy
      referrer: document.referrer || 'direct',
      url: window.location.pathname,
    };

    // Send to your own analytics endpoint or console log for now
    console.log('📊 Analytics:', data);

    // You can later integrate with privacy-focused analytics like:
    // - Plausible Analytics
    // - Fathom Analytics
    // - Your own endpoint
  }

  public trackPageView(): void {
    this.track({
      action: 'page_view',
      category: 'navigation',
      label: window.location.pathname,
    });
  }

  public trackPromptGeneration(mode: 'demo' | 'api'): void {
    this.track({
      action: 'prompt_generated',
      category: 'user_interaction',
      label: mode,
    });
  }

  public trackApiKeyConnected(): void {
    this.track({
      action: 'api_key_connected',
      category: 'user_interaction',
    });
  }

  public trackLanguageChange(language: string): void {
    this.track({
      action: 'language_changed',
      category: 'user_preference',
      label: language,
    });
  }

  public trackThemeToggle(theme: string): void {
    this.track({
      action: 'theme_toggled',
      category: 'user_preference',
      label: theme,
    });
  }
}

export const analytics = KubexAnalytics.getInstance();

// Hook for easy usage in components
export const useAnalytics = () => {
  useEffect(() => {
    analytics.trackPageView();
  }, []);

  return analytics;
};

export default analytics;
