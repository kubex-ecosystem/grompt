import path from 'path';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, '.', '');
  return {
    define: {
      /*

      Environment variables for the Application
      They are all prefixed with "process.env." to be compatible with libraries that expect this format.

      !!!!! DON'T LET THEM LEAK INTO THE FRONTEND CODE !!!!
      !!!!! THEY SHOULD ONLY BE USED IN SERVER-SIDE CODE !!!!!

      */
      'process.env.API_KEY': JSON.stringify(env.GRT_API_KEY || ''),
      'process.env.PORT': JSON.stringify(env.PORT || '3000'),
      'process.env.MAX_RETRIES': JSON.stringify(env.MAX_RETRIES || '3'),
      'process.env.RETRY_DELAY_MS': JSON.stringify(env.RETRY_DELAY_MS || '1000'),
      'process.env.DEBUG': JSON.stringify(env.DEBUG || 'false'),
      'process.env.TEST_MODE': JSON.stringify(env.TEST_MODE || 'false'),
      'process.env.LOG_LEVEL': JSON.stringify(env.LOG_LEVEL || 'info'),
      'process.env.SECURITY_ENABLED': JSON.stringify(env.SECURITY_ENABLED || 'true'),
      'process.env.RATE_LIMIT': JSON.stringify(env.RATE_LIMIT || '100'),
      'process.env.RATE_LIMIT_WINDOW': JSON.stringify(env.RATE_LIMIT_WINDOW || '1m'),

      // 'process.env.CORS_ORIGIN': JSON.stringify(env.CORS_ORIGIN || '*'),
      // 'process.env.CORS_CREDENTIALS': JSON.stringify(env.CORS_CREDENTIALS || 'true'),

      'process.env.OPENAI_API_KEY': JSON.stringify(env.OPENAI_API_KEY || ''),
      'process.env.CHATGPT_API_KEY': JSON.stringify(env.CHATGPT_API_KEY || ''),
      'process.env.ANTHROPIC_API_KEY': JSON.stringify(env.ANTHROPIC_API_KEY || ''),
      'process.env.DEEPSEEK_API_KEY': JSON.stringify(env.DEEPSEEK_API_KEY || ''),
      'process.env.OLLAMA_API_KEY': JSON.stringify(env.OLLAMA_API_KEY || ''),
      'process.env.OLLAMA_API_URL': JSON.stringify(env.OLLAMA_API_URL || ''),
      'process.env.GEMINI_API_KEY': JSON.stringify(env.GEMINI_API_KEY || ''),

      // 'process.env.GOOGLE_AI_STUDIO_API_KEY': JSON.stringify(env.GOOGLE_AI_STUDIO_API_KEY || ''),
      // 'process.env.GOOGLE_ANALYTICS_API_KEY': JSON.stringify(env.GOOGLE_ANALYTICS_API_KEY || ''),

      // 'process.env.PINECONE_API_KEY': JSON.stringify(env.PINECONE_API_KEY || ''),
      // 'process.env.GITHUB_PAT_TOKEN': JSON.stringify(env.GITHUB_PAT_TOKEN || ''),

      // 'process.env.VITE_SUPABASE_URL': JSON.stringify(env.VITE_SUPABASE_URL || ''),
      // 'process.env.VITE_SUPABASE_ANON_KEY': JSON.stringify(env.VITE_SUPABASE_ANON_KEY || ''),
      // 'process.env.VITE_SUPABASE_SERVICE_ROLE_KEY': JSON.stringify(env.VITE_SUPABASE_SERVICE_ROLE_KEY || ''),
      // 'process.env.VITE_SUPABASE_BEARER_TOKEN': JSON.stringify(env.VITE_SUPABASE_BEARER_TOKEN || ''),
      // 'process.env.VITE_SUPABASE_FALLBACK_ENABLE': JSON.stringify(env.VITE_SUPABASE_FALLBACK_ENABLE || 'false'),
      // 'process.env.VITE_SUPABASE_FALLBACK_URL': JSON.stringify(env.VITE_SUPABASE_FALLBACK_URL || ''),
      // 'process.env.VITE_SUPABASE_FALLBACK_ANON_KEY': JSON.stringify(env.VITE_SUPABASE_FALLBACK_ANON_KEY || ''),

      // 'process.env.POSTGRES_HOST': JSON.stringify(env.POSTGRES_HOST || ''),
      // 'process.env.POSTGRES_PORT': JSON.stringify(env.POSTGRES_PORT || ''),
      // 'process.env.POSTGRES_USER': JSON.stringify(env.POSTGRES_USER || ''),
      // 'process.env.POSTGRES_PASSWORD': JSON.stringify(env.POSTGRES_PASSWORD || ''),
      // 'process.env.POSTGRES_DB': JSON.stringify(env.POSTGRES_DB || ''),
      // 'process.env.POSTGRES_SSL': JSON.stringify(env.POSTGRES_SSL || ''),
      // 'process.env.POSTGRES_MAX_CLIENTS': JSON.stringify(env.POSTGRES_MAX_CLIENTS || ''),

      // 'process.env.REDIS_HOST': JSON.stringify(env.REDIS_HOST || ''),
      // 'process.env.REDIS_PORT': JSON.stringify(env.REDIS_PORT || ''),
      // 'process.env.REDIS_PASSWORD': JSON.stringify(env.REDIS_PASSWORD || ''),
      // 'process.env.REDIS_SENTINEL_ENABLED': JSON.stringify(env.REDIS_SENTINEL_ENABLED || ''),
      // 'process.env.REDIS_SENTINEL_HOST': JSON.stringify(env.REDIS_SENTINEL_HOST || ''),
      // 'process.env.REDIS_SENTINEL_PORT': JSON.stringify(env.REDIS_SENTINEL_PORT || ''),
      // 'process.env.REDIS_SENTINEL_MASTER_NAME': JSON.stringify(env.REDIS_SENTINEL_MASTER_NAME || ''),

      // 'process.env.RABBITMQ_HOST': JSON.stringify(env.RABBITMQ_HOST || ''),
      // 'process.env.RABBITMQ_PORT': JSON.stringify(env.RABBITMQ_PORT || ''),
      // 'process.env.RABBITMQ_USER': JSON.stringify(env.RABBITMQ_USER || ''),
      // 'process.env.RABBITMQ_PASSWORD': JSON.stringify(env.RABBITMQ_PASSWORD || ''),
      // 'process.env.RABBITMQ_VHOST': JSON.stringify(env.RABBITMQ_VHOST || ''),
      // 'process.env.RABBITMQ_QUEUE_NAME': JSON.stringify(env.RABBITMQ_QUEUE_NAME || ''),
      // 'process.env.RABBITMQ_QUEUE_DURABLE': JSON.stringify(env.RABBITMQ_QUEUE_DURABLE || ''),
      // 'process.env.RABBITMQ_QUEUE_AUTO_DELETE': JSON.stringify(env.RABBITMQ_QUEUE_AUTO_DELETE || ''),

      // 'process.env.MONGODB_URI': JSON.stringify(env.MONGODB_URI || ''),
      // 'process.env.MONGODB_DB': JSON.stringify(env.MONGODB_DB || ''),
      // 'process.env.MONGODB_USER': JSON.stringify(env.MONGODB_USER || ''),
      // 'process.env.MONGODB_PASSWORD': JSON.stringify(env.MONGODB_PASSWORD || ''),
      // 'process.env.MONGODB_AUTH_SOURCE': JSON.stringify(env.MONGODB_AUTH_SOURCE || ''),
      // 'process.env.MONGODB_SSL': JSON.stringify(env.MONGODB_SSL || ''),
      // 'process.env.MONGODB_REPLICA_SET': JSON.stringify(env.MONGODB_REPLICA_SET || ''),
      // 'process.env.MONGODB_MAX_POOL_SIZE': JSON.stringify(env.MONGODB_MAX_POOL_SIZE || ''),
      // 'process.env.MONGODB_MIN_POOL_SIZE': JSON.stringify(env.MONGODB_MIN_POOL_SIZE || ''),
      // 'process.env.MONGODB_CONNECT_TIMEOUT_MS': JSON.stringify(env.MONGODB_CONNECT_TIMEOUT_MS || ''),
      // 'process.env.MONGODB_SOCKET_TIMEOUT_MS': JSON.stringify(env.MONGODB_SOCKET_TIMEOUT_MS || ''),
      // 'process.env.MONGODB_SERVER_SELECTION_TIMEOUT_MS': JSON.stringify(env.MONGODB_SERVER_SELECTION_TIMEOUT_MS || ''),
      // 'process.env.MONGODB_HEARTBEAT_FREQUENCY_MS': JSON.stringify(env.MONGODB_HEARTBEAT_FREQUENCY_MS || ''),

      // 'process.env.BACKUP_ENABLED': JSON.stringify(env.BACKUP_ENABLED || ''),
      // 'process.env.BACKUP_SCHEDULE': JSON.stringify(env.BACKUP_SCHEDULE || ''),
      // 'process.env.BACKUP_RETENTION': JSON.stringify(env.BACKUP_RETENTION || ''),
      // 'process.env.BACKUP_PATH': JSON.stringify(env.BACKUP_PATH || ''),

      // 'process.env.SMTP_HOST': JSON.stringify(env.SMTP_HOST || ''),
      // 'process.env.SMTP_PORT': JSON.stringify(env.SMTP_PORT || ''),
      // 'process.env.SMTP_USER': JSON.stringify(env.SMTP_USER || ''),
      // 'process.env.SMTP_PASSWORD': JSON.stringify(env.SMTP_PASSWORD || ''),
      // 'process.env.SMTP_FROM': JSON.stringify(env.SMTP_FROM || ''),
      // 'process.env.SMTP_FROM_NAME': JSON.stringify(env.SMTP_FROM_NAME || '')
    },
    resolve: {
      alias: {
        '@': path.resolve(__dirname, '.'),
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
        },
      },
      chunkSizeWarningLimit: 1512,
    },
  };
});
