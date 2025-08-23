// Logger para requisições e navegação
class Logger {
  private isDev = process.env.NODE_ENV === 'development';
  private logs: Array<{
    timestamp: string;
    type: 'navigation' | 'request' | 'response' | 'error';
    data: any;
    level: 'info' | 'warn' | 'error';
  }> = [];

  private formatTimestamp(): string {
    return new Date().toISOString().split('T')[1].slice(0, -1); // HH:mm:ss.sss
  }

  private addLog(type: 'navigation' | 'request' | 'response' | 'error', data: any, level: 'info' | 'warn' | 'error' = 'info') {
    const logEntry = {
      timestamp: this.formatTimestamp(),
      type,
      data,
      level
    };

    this.logs.push(logEntry);

    // Manter apenas os últimos 100 logs
    if (this.logs.length > 100) {
      this.logs.shift();
    }

    if (this.isDev) {
      this.consoleLog(logEntry);
    }
  }

  private consoleLog(entry: {
    timestamp: string;
    type: 'navigation' | 'request' | 'response' | 'error';
    data: any;
    level: 'info' | 'warn' | 'error';
  }) {
    const emoji: Record<string, string> = {
      navigation: '🧭',
      request: '📡',
      response: '✅',
      error: '❌'
    };

    const levelColor: Record<string, string> = {
      info: 'color: #2563eb',
      warn: 'color: #f59e0b',
      error: 'color: #dc2626'
    };

    console.groupCollapsed(
      `%c${emoji[entry.type]} [${entry.timestamp}] ${entry.type.toUpperCase()}`,
      levelColor[entry.level]
    );
    console.log(entry.data);
    console.groupEnd();
  }

  // Navegação
  logNavigation(from: string, to: string, method: 'push' | 'replace' | 'back' = 'push') {
    this.addLog('navigation', {
      from,
      to,
      method,
      userAgent: navigator.userAgent.split(' ')[0] // Apenas o browser
    });
  }

  // Requisições HTTP
  logRequest(url: string, method: string, body?: any, headers?: any) {
    this.addLog('request', {
      url: url.replace(window.location.origin, ''), // URL relativa
      method,
      body: body ? JSON.stringify(body).slice(0, 200) + '...' : undefined,
      headers: headers ? Object.keys(headers) : undefined
    });
  }

  // Respostas HTTP
  logResponse(url: string, status: number, data?: any, duration?: number) {
    const level = status >= 400 ? 'error' : status >= 300 ? 'warn' : 'info';
    this.addLog('response', {
      url: url.replace(window.location.origin, ''),
      status,
      statusText: this.getStatusText(status),
      duration: duration ? `${duration}ms` : undefined,
      data: data ? JSON.stringify(data).slice(0, 100) + '...' : undefined
    }, level);
  }

  // Erros
  logError(error: Error, context?: string, extra?: any) {
    this.addLog('error', {
      message: error.message,
      stack: error.stack?.split('\n').slice(0, 3), // Apenas 3 primeiras linhas do stack
      context,
      extra
    }, 'error');
  }

  // Obter todos os logs
  getLogs() {
    return [...this.logs];
  }

  // Limpar logs
  clearLogs() {
    this.logs = [];
    console.clear();
  }

  // Exportar logs como JSON
  exportLogs() {
    const dataStr = JSON.stringify(this.logs, null, 2);
    const dataBlob = new Blob([dataStr], { type: 'application/json' });
    const url = URL.createObjectURL(dataBlob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `grompt-logs-${new Date().toISOString().split('T')[0]}.json`;
    link.click();
    URL.revokeObjectURL(url);
  }

  private getStatusText(status: number): string {
    const statusTexts: { [key: number]: string } = {
      200: 'OK',
      201: 'Created',
      204: 'No Content',
      400: 'Bad Request',
      401: 'Unauthorized',
      403: 'Forbidden',
      404: 'Not Found',
      500: 'Internal Server Error',
      502: 'Bad Gateway',
      503: 'Service Unavailable'
    };
    return statusTexts[status] || 'Unknown';
  }
}

// Instância global do logger
export const logger = new Logger();

// Hook para interceptar navegação do Next.js
export const useNavigationLogger = () => {
  if (typeof window !== 'undefined') {
    // Interceptar mudanças de URL
    const originalPushState = history.pushState;
    const originalReplaceState = history.replaceState;

    history.pushState = function (...args) {
      logger.logNavigation(window.location.pathname, args[2] as string, 'push');
      return originalPushState.apply(this, args);
    };

    history.replaceState = function (...args) {
      logger.logNavigation(window.location.pathname, args[2] as string, 'replace');
      return originalReplaceState.apply(this, args);
    };

    // Interceptar botão voltar
    window.addEventListener('popstate', () => {
      logger.logNavigation('history', window.location.pathname, 'back');
    });
  }
};

// Hook para interceptar fetch requests
export const useFetchLogger = () => {
  if (typeof window !== 'undefined') {
    const originalFetch = window.fetch;

    window.fetch = async function (input: RequestInfo | URL, init?: RequestInit) {
      const url = typeof input === 'string' ? input : input.toString();
      const method = init?.method || 'GET';
      const startTime = Date.now();

      logger.logRequest(url, method, init?.body, init?.headers);

      try {
        const response = await originalFetch(input, init);
        const duration = Date.now() - startTime;

        // Tentar obter dados da resposta para log (sem consumir o stream)
        const responseClone = response.clone();
        let responseData;
        try {
          responseData = await responseClone.text();
        } catch {
          responseData = '[Binary or stream data]';
        }

        logger.logResponse(url, response.status, responseData, duration);

        return response;
      } catch (error) {
        const duration = Date.now() - startTime;
        logger.logError(error as Error, `Fetch failed for ${url}`, { duration });
        throw error;
      }
    };
  }
};

export default logger;
