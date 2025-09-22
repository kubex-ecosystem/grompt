/**
 * PWA Status Component
 * Shows online/offline status, installation prompt, and sync status
 */

import * as React from 'react';
import {
  Download,
  RefreshCw,
  Wifi,
  WifiOff,
  Cloud,
  CloudOff,
  Smartphone,
  X,
  Check,
  AlertCircle
} from 'lucide-react';
import { usePWA, PWAUtils } from '../../hooks/usePWA';

interface PWAStatusProps {
  className?: string;
}

const PWAStatus: React.FC<PWAStatusProps> = ({ className = '' }) => {
  const {
    isInstallable,
    isInstalled,
    isOffline,
    isUpdateAvailable,
    isLoading,
    queuedRequestsCount,
    installApp,
    updateApp,
    syncOfflineData,
    clearOfflineData
  } = usePWA();

  const [showDetails, setShowDetails] = React.useState(false);
  const [isInstalling, setIsInstalling] = React.useState(false);
  const [isUpdating, setIsUpdating] = React.useState(false);
  const [isSyncing, setIsSyncing] = React.useState(false);
  const [lastSyncTime, setLastSyncTime] = React.useState<Date | null>(null);

  // Auto-hide install prompt after some time
  React.useEffect(() => {
    if (isInstallable) {
      const timer = setTimeout(() => {
        setShowDetails(false);
      }, 10000); // Hide after 10 seconds

      return () => clearTimeout(timer);
    }
  }, [isInstallable]);

  const handleInstall = async () => {
    setIsInstalling(true);
    try {
      await installApp();
      // Show success message or handle success
    } catch (error) {
      console.error('Installation failed:', error);
      // Show error message
    } finally {
      setIsInstalling(false);
    }
  };

  const handleUpdate = async () => {
    setIsUpdating(true);
    try {
      await updateApp();
    } catch (error) {
      console.error('Update failed:', error);
    } finally {
      setIsUpdating(false);
    }
  };

  const handleSync = async () => {
    setIsSyncing(true);
    try {
      await syncOfflineData();
      setLastSyncTime(new Date());
    } catch (error) {
      console.error('Sync failed:', error);
    } finally {
      setIsSyncing(false);
    }
  };

  const handleShare = async () => {
    const shared = await PWAUtils.share({
      title: 'Grompt - AI Prompt Engineering Tool',
      text: 'Transforme suas ideias em prompts eficazes para modelos de IA',
      url: window.location.href
    });

    if (!shared) {
      // Fallback to clipboard
      try {
        await navigator.clipboard.writeText(window.location.href);
        // Show success toast
      } catch (error) {
        console.error('Failed to copy to clipboard:', error);
      }
    }
  };

  if (isLoading) {
    return null;
  }

  return (
    <div className={`pwa-status ${className}`}>
      {/* Main Status Bar */}
      <div className="flex items-center gap-3 p-3 bg-gray-800/50 border border-gray-700/50 rounded-lg backdrop-blur-sm">
        {/* Connection Status */}
        <div className="flex items-center gap-2">
          {isOffline ? (
            <WifiOff size={16} className="text-red-400" />
          ) : (
            <Wifi size={16} className="text-green-400" />
          )}
          <span className="text-sm text-gray-300">
            {isOffline ? 'Offline' : 'Online'}
          </span>
        </div>

        {/* Queued Requests */}
        {queuedRequestsCount > 0 && (
          <div className="flex items-center gap-2 px-2 py-1 bg-orange-600/20 border border-orange-600/30 rounded text-orange-300">
            <Cloud size={14} />
            <span className="text-xs">
              {queuedRequestsCount} pendente{queuedRequestsCount > 1 ? 's' : ''}
            </span>
          </div>
        )}

        {/* Sync Button */}
        {(isOffline || queuedRequestsCount > 0) && (
          <button
            onClick={handleSync}
            disabled={isSyncing}
            className="p-1 rounded hover:bg-gray-700/50 text-gray-400 hover:text-white transition-colors"
            title="Sincronizar dados offline"
          >
            <RefreshCw size={14} className={isSyncing ? 'animate-spin' : ''} />
          </button>
        )}

        {/* Status Toggle */}
        <button
          onClick={() => setShowDetails(!showDetails)}
          className="ml-auto p-1 rounded hover:bg-gray-700/50 text-gray-400 hover:text-white transition-colors"
        >
          <AlertCircle size={16} />
        </button>
      </div>

      {/* Detailed Status Panel */}
      {showDetails && (
        <div className="mt-2 p-4 bg-gray-800/80 border border-gray-700/50 rounded-lg backdrop-blur-sm space-y-3">
          {/* PWA Installation */}
          {isInstallable && !isInstalled && (
            <div className="flex items-center justify-between p-3 bg-purple-600/10 border border-purple-600/30 rounded-lg">
              <div className="flex items-center gap-3">
                <Smartphone size={20} className="text-purple-400" />
                <div>
                  <h4 className="text-sm font-medium text-white">Instalar App</h4>
                  <p className="text-xs text-gray-400">
                    Use o Grompt como aplicativo nativo
                  </p>
                </div>
              </div>
              <button
                onClick={handleInstall}
                disabled={isInstalling}
                className="flex items-center gap-2 px-3 py-1 bg-purple-600 hover:bg-purple-700 text-white text-sm rounded transition-colors disabled:opacity-50"
              >
                {isInstalling ? (
                  <RefreshCw size={14} className="animate-spin" />
                ) : (
                  <Download size={14} />
                )}
                {isInstalling ? 'Instalando...' : 'Instalar'}
              </button>
            </div>
          )}

          {/* App Update */}
          {isUpdateAvailable && (
            <div className="flex items-center justify-between p-3 bg-blue-600/10 border border-blue-600/30 rounded-lg">
              <div className="flex items-center gap-3">
                <RefreshCw size={20} className="text-blue-400" />
                <div>
                  <h4 className="text-sm font-medium text-white">Atualização Disponível</h4>
                  <p className="text-xs text-gray-400">
                    Nova versão com melhorias e correções
                  </p>
                </div>
              </div>
              <button
                onClick={handleUpdate}
                disabled={isUpdating}
                className="flex items-center gap-2 px-3 py-1 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded transition-colors disabled:opacity-50"
              >
                {isUpdating ? (
                  <RefreshCw size={14} className="animate-spin" />
                ) : (
                  <RefreshCw size={14} />
                )}
                {isUpdating ? 'Atualizando...' : 'Atualizar'}
              </button>
            </div>
          )}

          {/* Offline Status */}
          {isOffline && (
            <div className="p-3 bg-yellow-600/10 border border-yellow-600/30 rounded-lg">
              <div className="flex items-center gap-3">
                <CloudOff size={20} className="text-yellow-400" />
                <div>
                  <h4 className="text-sm font-medium text-white">Modo Offline</h4>
                  <p className="text-xs text-gray-400">
                    Funcionalidades limitadas com templates locais
                  </p>
                </div>
              </div>
            </div>
          )}

          {/* Sync Status */}
          {lastSyncTime && (
            <div className="text-xs text-gray-500 text-center">
              Última sincronização: {lastSyncTime.toLocaleTimeString()}
            </div>
          )}

          {/* Actions */}
          <div className="flex gap-2 pt-2 border-t border-gray-700/50">
            <button
              onClick={handleShare}
              className="flex-1 text-xs py-2 px-3 bg-gray-700/50 hover:bg-gray-700 text-gray-300 hover:text-white rounded transition-colors"
            >
              Compartilhar
            </button>
            <button
              onClick={clearOfflineData}
              className="flex-1 text-xs py-2 px-3 bg-gray-700/50 hover:bg-gray-700 text-gray-300 hover:text-white rounded transition-colors"
            >
              Limpar Cache
            </button>
            <button
              onClick={() => setShowDetails(false)}
              className="p-2 hover:bg-gray-700/50 text-gray-400 hover:text-white rounded transition-colors"
            >
              <X size={14} />
            </button>
          </div>
        </div>
      )}

      {/* Install Prompt (Mobile-style) */}
      {isInstallable && !showDetails && (
        <div className="fixed bottom-4 left-4 right-4 max-w-sm mx-auto">
          <div className="bg-purple-600 text-white p-4 rounded-lg shadow-lg">
            <div className="flex items-center gap-3">
              <Smartphone size={24} />
              <div className="flex-1">
                <h4 className="font-medium">Instalar Grompt</h4>
                <p className="text-sm opacity-90">
                  Acesso rápido e funcionalidade offline
                </p>
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => setShowDetails(false)}
                  className="p-2 hover:bg-white/20 rounded transition-colors"
                >
                  <X size={16} />
                </button>
                <button
                  onClick={handleInstall}
                  disabled={isInstalling}
                  className="px-3 py-1 bg-white text-purple-600 text-sm font-medium rounded hover:bg-gray-100 transition-colors disabled:opacity-50"
                >
                  {isInstalling ? 'Instalando...' : 'Instalar'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default PWAStatus;