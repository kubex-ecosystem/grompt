import path from 'path';
import { fileURLToPath } from 'url';
import { defineConfig, loadEnv } from 'vite';

const getEnvResilient = (mode: string, envDir: string) => {
  // Estratégia de fallback para diferentes locais de .env
  const envPaths = [
    envDir,                    // diretório raiz
    path.join(envDir, 'config'), // ./config/
    path.join(envDir, '..'),   // diretório pai
  ];

  let finalEnv = {};

  for (const envPath of envPaths) {
    try {
      const env = loadEnv(mode, envPath, ['GEMINI_', 'GITHUB_', 'JIRA_', 'API_']);
      finalEnv = { ...finalEnv, ...env };
    } catch (error) {
      // Silently continue to next path
      continue;
    }
  }

  return finalEnv;
};

export default defineConfig(({ mode }: { mode: string }) => {
  const env: Record<string, string> = getEnvResilient(mode, process.cwd());
  return {
    root: '.',
    base: './',
    publicDir: 'public',
    cacheDir: 'node_modules/.vite',
    mode: mode,
    define: {
      'process.env.API_KEY': JSON.stringify(env.VITE_GEMINI_API_KEY || env.GEMINI_API_KEY || ""),
      'process.env.GEMINI_API_KEY': JSON.stringify(env.VITE_GEMINI_API_KEY || env.GEMINI_API_KEY || ""),
      'process.env.GITHUB_PAT': JSON.stringify(env.VITE_GITHUB_PAT || env.GITHUB_PAT || ""),
      'process.env.JIRA_API_TOKEN': JSON.stringify(env.VITE_JIRA_API_TOKEN || env.JIRA_API_TOKEN || ""),
      'process.env.JIRA_INSTANCE_URL': JSON.stringify(env.VITE_JIRA_INSTANCE_URL || env.JIRA_INSTANCE_URL || ""),
      'process.env.JIRA_USER_EMAIL': JSON.stringify(env.VITE_JIRA_USER_EMAIL || env.JIRA_USER_EMAIL || "")
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('.', import.meta.url)),
      }
    },
    build: {
      rollupOptions: {
        output: {
          manualChunks: {
            vendor: [
              'react',
              'react-dom'
            ],
          },
          plugins: [
            {
              name: 'watch-external',
              handleHotUpdate({ file, server }: { file: string; server: any }) {
                if (file.endsWith('shared/config.json')) {
                  server.ws.send({ type: 'full-reload' });
                }
              }
            },
            {
              name: 'configure-response-headers',
              configureServer(server: any) {
                server.middlewares.use((req: any, res: any, next: any) => {
                  res.setHeader('Access-Control-Allow-Origin', '*');
                  res.setHeader('Access-Control-Allow-Methods', 'GET,POST,PUT,DELETE,OPTIONS');
                  res.setHeader('Access-Control-Allow-Credentials', 'true');

                  // Set CORS headers (for preflight requests)
                  if (req.method === 'OPTIONS') {
                    res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization');
                    res.setHeader('Access-Control-Max-Age', '86400'); // 24 hours
                    res.setHeader('Permissions-Policy', 'geolocation=(), microphone=(), camera=(), payment=()');

                    res.writeHead(204);
                    res.end();
                    return;
                  }

                  // Expose headers
                  res.setHeader('Access-Control-Max-Age', '86400'); // 24 hours
                  res.setHeader('Permissions-Policy', 'geolocation=(), microphone=(), camera=(), payment=()');
                  res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization, X-Requested-With');
                  res.setHeader('Access-Control-Expose-Headers', 'Content-Length, X-Total-Count, X-Page, X-Per-Page, X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset');

                  // Security headers
                  res.setHeader('Strict-Transport-Security', 'max-age=31536000; includeSubDomains; preload');
                  res.setHeader('Content-Security-Policy', "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self' https://analyzer.kubex.world; frame-ancestors 'none'; form-action 'self'; base-uri 'self';");
                  res.setHeader('X-Content-Type-Options', 'nosniff');
                  res.setHeader('X-DNS-Prefetch-Control', 'off');
                  res.setHeader('X-Frame-Options', 'DENY');
                  res.setHeader('X-XSS-Protection', '1; mode=block');
                  res.setHeader('X-Download-Options', 'noopen');
                  res.setHeader('X-Permitted-Cross-Domain-Policies', 'none');

                  next();
                });
              },
              onwarn: (warning: any, warn: any) => {
                if (
                  warning.code === 'UNUSED_EXTERNAL_IMPORT' &&
                  warning.message.includes('was ignored.')
                ) {
                  return;
                }
                warn(warning);
              },
            }
          ]

        },
      },
      outDir: 'dist',
      sourcemap: true,
      chunkSizeWarningLimit: 900,
    },
    css: {
      preprocessorOptions: {
        scss: {
          additionalData: `@import "@/styles/variables.scss";`
        }
      }
    },
    optimizeDeps: {
      include: ['react', 'react-dom', 'framer-motion', 'uuid',],
    },
    esbuild: {
      drop: ['console', 'debugger'],
    },
  };
});

