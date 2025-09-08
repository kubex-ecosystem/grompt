# Environment Variables Documentation

Grompt v2.0 supports multiple environment variables for flexible configuration across different deployment scenarios.

## Core Application Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3000` | Server port for development |
| `NODE_ENV` | `development` | Application environment |
| `DEBUG` | `false` | Enable debug logging |
| `TEST_MODE` | `false` | Enable test mode features |
| `LOG_LEVEL` | `info` | Logging level (error, warn, info, debug) |

## AI Model Providers

### Google Gemini (Primary)

| Variable | Description |
|----------|-------------|
| `GEMINI_API_KEY` | Google Gemini API key for prompt generation |
| `GRT_API_KEY` | Alternative Gemini API key variable |

### OpenAI

| Variable | Description |
|----------|-------------|
| `OPENAI_API_KEY` | OpenAI GPT API key |
| `CHATGPT_API_KEY` | Alternative OpenAI API key variable |

### Other AI Providers

| Variable | Description |
|----------|-------------|
| `ANTHROPIC_API_KEY` | Anthropic Claude API key |
| `DEEPSEEK_API_KEY` | DeepSeek AI API key |
| `OLLAMA_API_KEY` | Ollama local AI API key |
| `OLLAMA_API_URL` | Ollama server URL |

## Security & Rate Limiting

| Variable | Default | Description |
|----------|---------|-------------|
| `SECURITY_ENABLED` | `true` | Enable security features |
| `RATE_LIMIT` | `100` | Rate limit per window |
| `RATE_LIMIT_WINDOW` | `1m` | Rate limit window duration |
| `MAX_RETRIES` | `3` | Maximum API retry attempts |
| `RETRY_DELAY_MS` | `1000` | Delay between retries |

## Future Database Support (Planned)

### PostgreSQL

| Variable | Description |
|----------|-------------|
| `POSTGRES_HOST` | PostgreSQL server host |
| `POSTGRES_PORT` | PostgreSQL server port |
| `POSTGRES_USER` | Database username |
| `POSTGRES_PASSWORD` | Database password |
| `POSTGRES_DB` | Database name |
| `POSTGRES_SSL` | Enable SSL connection |
| `POSTGRES_MAX_CLIENTS` | Maximum connection pool size |

### Redis

| Variable | Description |
|----------|-------------|
| `REDIS_HOST` | Redis server host |
| `REDIS_PORT` | Redis server port |
| `REDIS_PASSWORD` | Redis authentication password |
| `REDIS_SENTINEL_ENABLED` | Enable Redis Sentinel |
| `REDIS_SENTINEL_HOST` | Sentinel server host |
| `REDIS_SENTINEL_PORT` | Sentinel server port |
| `REDIS_SENTINEL_MASTER_NAME` | Sentinel master name |

### MongoDB

| Variable | Description |
|----------|-------------|
| `MONGODB_URI` | Complete MongoDB connection string |
| `MONGODB_DB` | MongoDB database name |
| `MONGODB_USER` | MongoDB username |
| `MONGODB_PASSWORD` | MongoDB password |
| `MONGODB_AUTH_SOURCE` | Authentication database |
| `MONGODB_SSL` | Enable SSL connection |
| `MONGODB_REPLICA_SET` | Replica set name |
| `MONGODB_MAX_POOL_SIZE` | Maximum connection pool size |
| `MONGODB_MIN_POOL_SIZE` | Minimum connection pool size |
| `MONGODB_CONNECT_TIMEOUT_MS` | Connection timeout |

### RabbitMQ

| Variable | Description |
|----------|-------------|
| `RABBITMQ_HOST` | RabbitMQ server host |
| `RABBITMQ_PORT` | RabbitMQ server port |
| `RABBITMQ_USER` | RabbitMQ username |
| `RABBITMQ_PASSWORD` | RabbitMQ password |
| `RABBITMQ_VHOST` | Virtual host |
| `RABBITMQ_QUEUE_NAME` | Default queue name |
| `RABBITMQ_QUEUE_DURABLE` | Queue durability |
| `RABBITMQ_QUEUE_AUTO_DELETE` | Auto-delete queue |

## Usage Examples

### Development

```bash
# Minimal setup for development
GEMINI_API_KEY=your_key_here
DEBUG=true
LOG_LEVEL=debug
```

### Production

```bash
# Production setup with rate limiting
GEMINI_API_KEY=your_production_key
NODE_ENV=production
SECURITY_ENABLED=true
RATE_LIMIT=50
RATE_LIMIT_WINDOW=1m
LOG_LEVEL=warn
```

### Multi-Provider Setup

```bash
# Support multiple AI providers
GEMINI_API_KEY=gemini_key
OPENAI_API_KEY=openai_key
ANTHROPIC_API_KEY=claude_key
OLLAMA_API_URL=http://localhost:11434
```

## Security Notes

⚠️ **Never commit API keys to version control**

- Use `.env.local` for local development
- Use secure environment variable management in production
- API keys are only used server-side, never exposed to frontend
- User-provided API keys are stored only in localStorage (client-side)

## Kubex Principles Applied

- **Radical Simplicity**: One variable = one function
- **Modularity**: Variables grouped by logical domains
- **No Cages**: Support for multiple providers, no vendor lock-in
