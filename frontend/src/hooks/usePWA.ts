/**
 * PWA Hook for Service Worker and Installation Management
 */

import * as React from 'react';

export interface PWAState {
  isInstallable: boolean;
  isInstalled: boolean;
  isOffline: boolean;
  isUpdateAvailable: boolean;
  isLoading: boolean;
  queuedRequestsCount: number;
}

export interface PWAActions {
  installApp: () => Promise<void>;
  updateApp: () => Promise<void>;
  syncOfflineData: () => Promise<void>;
  clearOfflineData: () => Promise<void>;
}

interface BeforeInstallPromptEvent extends Event {
  readonly platforms: string[];
  readonly userChoice: Promise<{
    outcome: 'accepted' | 'dismissed';
    platform: string;
  }>;
  prompt(): Promise<void>;
}

declare global {
  interface WindowEventMap {
    beforeinstallprompt: BeforeInstallPromptEvent;
  }
}

export function usePWA(): PWAState & PWAActions {
  const [state, setState] = React.useState<PWAState>({
    isInstallable: false,
    isInstalled: false,
    isOffline: !navigator.onLine,
    isUpdateAvailable: false,
    isLoading: true,
    queuedRequestsCount: 0
  });

  const [deferredPrompt, setDeferredPrompt] = React.useState<BeforeInstallPromptEvent | null>(null);
  const [serviceWorkerRegistration, setServiceWorkerRegistration] = React.useState<ServiceWorkerRegistration | null>(null);

  // Initialize PWA functionality
  React.useEffect(() => {
    initializePWA();
    setupEventListeners();

    return () => {
      cleanupEventListeners();
    };
  }, []);

  // Monitor queued requests count
  React.useEffect(() => {
    const updateQueueCount = async () => {
      try {
        // Import enhanced API dynamically to avoid circular deps
        const { enhancedAPI } = await import('../services/enhancedAPI');
        const count = await enhancedAPI.getQueuedRequestsCount();
        setState(prev => ({ ...prev, queuedRequestsCount: count }));
      } catch (error) {
        console.warn('Failed to get queued requests count:', error);
      }
    };

    updateQueueCount();
    const interval = setInterval(updateQueueCount, 30000); // Check every 30 seconds

    return () => clearInterval(interval);
  }, [state.isOffline]);

  const initializePWA = async () => {
    try {
      // Check if app is already installed
      const isInstalled = window.matchMedia('(display-mode: standalone)').matches ||
                         (window.navigator as any).standalone ||
                         document.referrer.includes('android-app://');

      setState(prev => ({ ...prev, isInstalled }));

      // Register service worker
      if ('serviceWorker' in navigator) {
        const registration = await navigator.serviceWorker.register('/sw.js', {
          scope: '/'
        });

        setServiceWorkerRegistration(registration);
        console.log('✅ Service Worker registered:', registration.scope);

        // Check for updates
        registration.addEventListener('updatefound', () => {
          const newWorker = registration.installing;
          if (newWorker) {
            newWorker.addEventListener('statechange', () => {
              if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
                setState(prev => ({ ...prev, isUpdateAvailable: true }));
              }
            });
          }
        });

        // Listen for messages from service worker
        navigator.serviceWorker.addEventListener('message', handleServiceWorkerMessage);
      }

    } catch (error) {
      console.error('❌ PWA initialization failed:', error);
    }

    setState(prev => ({ ...prev, isLoading: false }));
  };

  const setupEventListeners = () => {
    // Installation prompt
    window.addEventListener('beforeinstallprompt', handleBeforeInstallPrompt);

    // Online/offline status
    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    // App installed
    window.addEventListener('appinstalled', handleAppInstalled);
  };

  const cleanupEventListeners = () => {
    window.removeEventListener('beforeinstallprompt', handleBeforeInstallPrompt);
    window.removeEventListener('online', handleOnline);
    window.removeEventListener('offline', handleOffline);
    window.removeEventListener('appinstalled', handleAppInstalled);

    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.removeEventListener('message', handleServiceWorkerMessage);
    }
  };

  const handleBeforeInstallPrompt = (event: BeforeInstallPromptEvent) => {
    event.preventDefault();
    setDeferredPrompt(event);
    setState(prev => ({ ...prev, isInstallable: true }));
  };

  const handleOnline = () => {
    setState(prev => ({ ...prev, isOffline: false }));
    // Trigger sync when back online
    syncOfflineData();
  };

  const handleOffline = () => {
    setState(prev => ({ ...prev, isOffline: true }));
  };

  const handleAppInstalled = () => {
    setState(prev => ({
      ...prev,
      isInstalled: true,
      isInstallable: false
    }));
    setDeferredPrompt(null);
  };

  const handleServiceWorkerMessage = (event: MessageEvent) => {
    if (event.data?.type === 'SYNC_OFFLINE_QUEUE') {
      syncOfflineData();
    }
  };

  // Actions
  const installApp = async (): Promise<void> => {
    if (!deferredPrompt) {
      throw new Error('App não está pronto para instalação');
    }

    try {
      await deferredPrompt.prompt();
      const choiceResult = await deferredPrompt.userChoice;

      if (choiceResult.outcome === 'accepted') {
        console.log('✅ User accepted PWA installation');
      } else {
        console.log('❌ User dismissed PWA installation');
      }

      setDeferredPrompt(null);
      setState(prev => ({ ...prev, isInstallable: false }));
    } catch (error) {
      console.error('❌ PWA installation failed:', error);
      throw error;
    }
  };

  const updateApp = async (): Promise<void> => {
    if (!serviceWorkerRegistration) {
      throw new Error('Service Worker não está registrado');
    }

    try {
      const newWorker = serviceWorkerRegistration.waiting;
      if (newWorker) {
        // Send message to skip waiting
        newWorker.postMessage({ type: 'SKIP_WAITING' });

        // Wait for controller change
        await new Promise<void>((resolve) => {
          navigator.serviceWorker.addEventListener('controllerchange', () => {
            resolve();
          }, { once: true });
        });

        // Reload the page
        window.location.reload();
      }
    } catch (error) {
      console.error('❌ App update failed:', error);
      throw error;
    }
  };

  const syncOfflineData = async (): Promise<void> => {
    try {
      const { enhancedAPI } = await import('../services/enhancedAPI');

      if (enhancedAPI.isConnected()) {
        // Trigger sync - the enhanced API handles this internally
        console.log('🔄 Syncing offline data...');

        // Update queue count after sync
        setTimeout(async () => {
          const count = await enhancedAPI.getQueuedRequestsCount();
          setState(prev => ({ ...prev, queuedRequestsCount: count }));
        }, 1000);
      }
    } catch (error) {
      console.error('❌ Offline data sync failed:', error);
      throw error;
    }
  };

  const clearOfflineData = async (): Promise<void> => {
    try {
      const { enhancedAPI } = await import('../services/enhancedAPI');
      await enhancedAPI.clearCache();

      setState(prev => ({ ...prev, queuedRequestsCount: 0 }));
      console.log('✅ Offline data cleared');
    } catch (error) {
      console.error('❌ Failed to clear offline data:', error);
      throw error;
    }
  };

  return {
    ...state,
    installApp,
    updateApp,
    syncOfflineData,
    clearOfflineData
  };
}

// Utility hook for PWA features availability
export function usePWAFeatures() {
  const [features, setFeatures] = React.useState({
    serviceWorker: 'serviceWorker' in navigator,
    pushNotifications: 'PushManager' in window,
    backgroundSync: 'serviceWorker' in navigator && 'sync' in window.ServiceWorkerRegistration.prototype,
    installPrompt: 'BeforeInstallPromptEvent' in window || 'onbeforeinstallprompt' in window,
    webShare: 'share' in navigator,
    badgeAPI: 'setAppBadge' in navigator,
    wakeLock: 'wakeLock' in navigator
  });

  React.useEffect(() => {
    // Additional feature detection can be added here
    const isStandalone = window.matchMedia('(display-mode: standalone)').matches;
    setFeatures(prev => ({ ...prev, isStandalone }));
  }, []);

  return features;
}

// Utility functions for PWA management
export const PWAUtils = {
  // Check if app is installed
  isInstalled(): boolean {
    return window.matchMedia('(display-mode: standalone)').matches ||
           (window.navigator as any).standalone ||
           document.referrer.includes('android-app://');
  },

  // Get installation platforms
  async getInstallationPlatforms(): Promise<string[]> {
    if ('getInstalledRelatedApps' in navigator) {
      try {
        const relatedApps = await (navigator as any).getInstalledRelatedApps();
        return relatedApps.map((app: any) => app.platform);
      } catch (error) {
        console.warn('Failed to get installation platforms:', error);
      }
    }
    return [];
  },

  // Share content using Web Share API
  async share(data: { title?: string; text?: string; url?: string }): Promise<boolean> {
    if ('share' in navigator) {
      try {
        await navigator.share(data);
        return true;
      } catch (error) {
        if ((error as Error).name !== 'AbortError') {
          console.error('Web Share API failed:', error);
        }
      }
    }
    return false;
  },

  // Set app badge (if supported)
  async setBadge(count?: number): Promise<void> {
    if ('setAppBadge' in navigator) {
      try {
        await (navigator as any).setAppBadge(count);
      } catch (error) {
        console.warn('Failed to set app badge:', error);
      }
    }
  },

  // Clear app badge
  async clearBadge(): Promise<void> {
    if ('clearAppBadge' in navigator) {
      try {
        await (navigator as any).clearAppBadge();
      } catch (error) {
        console.warn('Failed to clear app badge:', error);
      }
    }
  }
};