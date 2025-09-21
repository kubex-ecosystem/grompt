# Reforma do Backend do Grompt (versão enxuta/gateway)

## Contexto (essencial)

* Projeto: **Grompt** (app leve, <5 MB, PWA-friendly) para criação de prompts profissionais e agentes.
* Filosofia: **backend mínimo** (“BE do front”) apenas como **gateway**: streaming, cadência, proteção de custos e proxy.
* Serviços “pesados” (auth forte, storage, billing, políticas, uploads grandes etc.) = **delegados ao GoBE** (serviço principal).
* UX: rápido, bonito, responsivo; **SSE** como primeiro cidadão; mesmo padrão de camadas já usado no Analyzer/PWA.

## Objetivo desta rodada

Entregar um **BE enxuto** para o Grompt com:

1. **SSE streaming** de modelos (multiprovider via SDKs oficiais)
2. **Gateway/Proxy** para o GoBE (auth/flash-drive/storage por lá)
3. **Controle de cadência e custo** (rate-limit, idempotency, timeouts)
4. **Compatível com PWA** e “4 camadas” de entrega (edge-friendly)

## Regras de ouro

* Domínio/framework-agnóstico: use **ports/adapters**; nada de acoplamento de domínio ao framework.
* **Não** implementar auth pesada, storage, billing: **proxie para o GoBE** (URLs fornecidas via ENV).
* Foco em **SSE**, **latência baixa**, **robustez mínima** (429/5xx retry curto).
* Usar SDKs oficiais: `@google/genai` (Chat API: `ai.chats.create`), `@anthropic-ai/sdk`, `openai`.
* Nada de tipos React no backend.

## Contratos imutáveis (não alterar)

```ts
// packages/ports/provider.ts
export type ProviderId = "openai"|"anthropic"|"gemini";
export interface GenOptions { temperature?: number; topP?: number; maxTokens?: number; stopSequences?: string[]; }
export interface GenInput { system?: string; user: string; }
export interface GenUsage { prompt: number; completion: number; total: number; }
export interface GenResult { text: string; provider: ProviderId; model: string; usage?: GenUsage; finishReason?: string; }
export interface ProviderPort {
  generate(input: GenInput, options?: GenOptions): Promise<GenResult>;
  stream?(input: GenInput, options?: GenOptions): AsyncIterable<string>;
}
```

```ts
// packages/ports/http.ts
export type HttpMethod = "GET"|"POST"|"PUT"|"PATCH"|"DELETE";
export interface RequestContext {
  method: HttpMethod; path: string; params: Record<string,string>;
  query: Record<string,string|string[]>; headers: Record<string,string>;
  body?: unknown; raw?: { req:any; res:any }; user?: { id:string; org?:string };
}
export interface ResponsePayload { status?: number; headers?: Record<string,string>; body?: any; }
export interface HttpServerPort {
  route(method: HttpMethod, path: string, h:(ctx:RequestContext)=>Promise<ResponsePayload>): void;
  sse?(path: string, h:(ctx:RequestContext)=>Promise<ResponsePayload>): void;
  start(): Promise<void>; stop(): Promise<void>;
}
```

## Endpoints (v1)

* `POST /v1/generate` → inicia geração (com opção de streaming)
* `GET  /v1/generate/stream` → **SSE** de chunks (provider/model/opts via query)
* `POST /v1/proxy/*` → proxy simples para GoBE (pass-through seguro)
* `GET  /v1/providers` → lista providers/modelos disponíveis (estático/ENV)

> **Sem** rotas de storage, auth, billing, upload pesado: **use `/v1/proxy/*` → GoBE**.

## Variáveis de ambiente

```plaintext
OPENAI_API_KEY=...
ANTHROPIC_API_KEY=...
GEMINI_API_KEY=...
GOBE_BASE_URL=https://gobe.example.com
PORT=3000
HTTP_ADAPTER=fastify                # adapter pluggable
DEFAULT_PROVIDER=gemini
DEFAULT_MODEL=gemini-2.5-flash
RATE_LIMIT_RPS=3
REQUEST_TIMEOUT_MS=30000
```

## Tarefas (em ordem) + critérios de aceite

### 1) Ports + HTTP Adapter mínimo

* Criar `packages/ports` (↑ contratos) e `adapters/http/FastifyAdapter` implementando `HttpServerPort`.
* **Aceite**: `GET /v1/health` → 200; adapter sob controle (`start/stop`).

### 2) Providers + MultiAI (SDKs oficiais)

* `packages/llm`: `OpenAIProvider`, `AnthropicProvider`, `GeminiProvider`.

  * **Gemini**: `GoogleGenAI` + `ai.chats.create({ model, config })` → `sendMessage`/`sendMessageStream`, `systemInstruction`/geração em `config`.
  * OpenAI: `chat.completions.create` (generate/stream).
  * Anthropic: `messages.create` (generate/stream).
* `MultiAIWrapper` simples: escolhe provider/model; retorna `GenResult`; retry curto 429/5xx + timeouts.
* **Aceite**: smoke test retorna texto; streaming emite chunks; `usage` quando houver.

### 3) Endpoints /v1 (generate + SSE + providers)

* `POST /v1/generate` → recebe `{ provider?, model?, system?, user, options? }` → dispara `generate` (sincrono) **ou** sugere uso de `/stream`.
* `GET /v1/generate/stream` (SSE) → lê query/body, chama `stream()` e escreve `data: {...}` por chunk.
* `GET /v1/providers` → devolve lista a partir do wrapper/ENV.
* **Aceite**: latência baixa; SSE estável; erros mapeados (4xx/5xx) com mensagem curta.

### 4) Proxy para GoBE

* `POST /v1/proxy/*` → encaminha métodos/headers/body para `${GOBE_BASE_URL}/*` (whitelist básica de caminhos).
* Copiar `Authorization` e `X-Request-Id` quando existirem.
* **Aceite**: round-trip ok; **sem** logar payload sensível; timeout/controlado.

### 5) Guard-rails de custo/latência

* **Rate-limit** simples (RPS via ENV) por IP/org (se header presente).
* **Idempotency-Key** opcional no `/v1/generate` (cache em memória simples com TTL curto).
* **Timeout** por requisição (AbortController).
* **Aceite**: excedeu → 429 JSON curto; timeout → 504 JSON curto; repetir mesma key → reaproveita.

### 6) Observabilidade mínima

* **Logs estruturados** (pino) com `req_id`, `provider`, `model`, `lat_ms`.
* **/metrics** (Prometheus) com:

  * `grompt_gateway_tokens_total{provider,model,type="prompt|completion"}`
  * `grompt_gateway_job_latency_seconds` (histograma light)
* **Aceite**: `/metrics` expõe counters/histogram; logs mostram request→provider.

## Não escopo (explícito)

* Nada de DB pesado, storage de artifacts, snapshots, auth complexa: **tudo via GoBE**.
* Nada de queue/worker nesta rodada (só se estritamente necessário p/ SSE não bloquear).

## Boas práticas obrigatórias

* Sem dependência de React/DOM no BE.
* Tipagem completa (TS estrito).
* Erros com mensagens curtas e status correto.
* Sem segredos em log; truncar conteúdo longo.

## Smoke tests (dev têm que passar)

* `POST /v1/generate` com provider default → 200 com texto.
* `GET /v1/generate/stream` → recebe chunks em SSE.
* `GET /v1/providers` → lista modelos.
* `POST /v1/proxy/echo` (apontado ao GoBE de teste) → round-trip ok.
* Rate-limit funciona; timeout funciona; idempotency reutiliza.

## Notas sobre Gemini (importante)

* Usar `import { GoogleGenAI } from "@google/genai"`.
* Chat API:

  ```ts
  const chat = ai.chats.create({ model, config: { systemInstruction, temperature, topP, maxOutputTokens, stopSequences } });
  const resp = await chat.sendMessage({ message: prompt });      // non-stream
  const s = await chat.sendMessageStream({ message: prompt });   // stream
  ```

* Ler texto de forma defensiva: `resp.text` **ou** `resp.response.text()`.

## Entrega

* Código e README curto (2–5 bullets de decisões).
* Pequenos testes de fumaça (providers e endpoints).
* PR único desta rodada, bem focado, sem invadir escopo do GoBE.

---

> Resultado esperado: um **gateway rápido e leve** para o Grompt, com SSE impecável e controle de cadência/custos, delegando o resto ao **GoBE**.
