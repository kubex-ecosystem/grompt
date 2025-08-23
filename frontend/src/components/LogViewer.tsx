import { Download, Eye, RotateCcw } from 'lucide-react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { logger } from '../utils/logger';

interface LogViewerProps {
  currentTheme: any;
  isOpen?: boolean;
  onClose?: () => void;
}

const LogViewer = ({ currentTheme, isOpen: externalIsOpen, onClose }: LogViewerProps) => {
  const { t } = useTranslation();
  const [internalIsOpen, setInternalIsOpen] = useState(false);
  const [logs] = useState(() => logger.getLogs());

  // Use external control if provided, otherwise internal state
  const isOpen = externalIsOpen !== undefined ? externalIsOpen : internalIsOpen;
  const setIsOpen = onClose ? (value: boolean) => {
    if (!value) onClose();
  } : setInternalIsOpen;

  const handleExportLogs = () => {
    logger.exportLogs();
  };

  const handleClearLogs = () => {
    logger.clearLogs();
    setIsOpen(false);
  };

  const handleOpenViewer = () => {
    if (externalIsOpen === undefined) {
      setInternalIsOpen(true);
    }
  };

  const getLogIcon = (type: string) => {
    const icons: Record<string, string> = {
      navigation: '🧭',
      request: '📡',
      response: '✅',
      error: '❌'
    };
    return icons[type] || '📝';
  };

  const getLogLevel = (level: string) => {
    const colors: Record<string, string> = {
      info: 'text-blue-500',
      warn: 'text-yellow-500',
      error: 'text-red-500'
    };
    return colors[level] || 'text-gray-500';
  };

  if (!isOpen) {
    return (
      <button
        onClick={handleOpenViewer}
        className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
        title="Ver Logs de Navegação e Requisições"
      >
        <Eye size={16} />
      </button>
    );
  }

  return (
    <div className="fixed inset-0 z-50 bg-black bg-opacity-50 flex items-center justify-center p-4">
      <div className={`${currentTheme.cardBg} rounded-lg shadow-xl max-w-4xl w-full max-h-[80vh] flex flex-col`}>
        {/* Header */}
        <div className={`p-4 border-b ${currentTheme.border} flex justify-between items-center`}>
          <h2 className="text-xl font-bold">Logs de Sistema</h2>
          <div className="flex gap-2">
            <button
              onClick={handleExportLogs}
              className={`p-2 rounded-lg ${currentTheme.button} text-white transition-colors`}
              title="Exportar Logs"
            >
              <Download size={16} />
            </button>
            <button
              onClick={handleClearLogs}
              className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
              title="Limpar Logs"
            >
              <RotateCcw size={16} />
            </button>
            <button
              onClick={() => setIsOpen(false)}
              className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
              title="Fechar"
            >
              ✕
            </button>
          </div>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-auto p-4">
          {logs.length === 0 ? (
            <div className="text-center text-gray-500 py-8">
              <p>Nenhum log registrado ainda.</p>
              <p className="text-sm mt-2">Navegue pela aplicação ou faça requisições para ver os logs aqui.</p>
            </div>
          ) : (
            <div className="space-y-2">
              {logs.slice(-50).reverse().map((log, index) => (
                <div
                  key={index}
                  className={`p-3 rounded-lg border ${currentTheme.border} bg-opacity-50 hover:bg-opacity-75 transition-colors`}
                >
                  <div className="flex items-center gap-2 mb-1">
                    <span className="text-lg">{getLogIcon(log.type)}</span>
                    <span className={`font-medium ${getLogLevel(log.level)}`}>
                      {log.type.toUpperCase()}
                    </span>
                    <span className="text-sm text-gray-500">{log.timestamp}</span>
                  </div>
                  <div className="text-sm">
                    {log.type === 'navigation' && (
                      <div>
                        <span className="text-gray-500">De:</span> {log.data.from} →{' '}
                        <span className="text-gray-500">Para:</span> {log.data.to}
                        {log.data.method && (
                          <span className="ml-2 text-gray-400">({log.data.method})</span>
                        )}
                      </div>
                    )}
                    {log.type === 'request' && (
                      <div>
                        <span className="font-medium">{log.data.method}</span> {log.data.url}
                        {log.data.body && (
                          <div className="text-gray-500 mt-1">Body: {log.data.body}</div>
                        )}
                      </div>
                    )}
                    {log.type === 'response' && (
                      <div>
                        <span className={log.data.status >= 400 ? 'text-red-500' : 'text-green-500'}>
                          {log.data.status} {log.data.statusText}
                        </span>{' '}
                        {log.data.url}
                        {log.data.duration && (
                          <span className="ml-2 text-gray-400">({log.data.duration})</span>
                        )}
                      </div>
                    )}
                    {log.type === 'error' && (
                      <div>
                        <div className="text-red-500 font-medium">{log.data.message}</div>
                        {log.data.context && (
                          <div className="text-gray-500 mt-1">Context: {log.data.context}</div>
                        )}
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className={`p-4 border-t ${currentTheme.border} text-sm text-gray-500`}>
          Total de logs: {logs.length} (mostrando últimos 50)
        </div>
      </div>
    </div>
  );
};

export default LogViewer;
