# Arquitetura do Backend do Grompt

Deixar o BE do Grompt **idempotente, multiprovider plugável, stream-first, observável e barato** — sem travar seu FE e sem refém de um único vendor.

---

# 1) Arquitetura em camadas (hex + plugins)

**Camadas (do centro pra fora):**

1. **Domain (core)**

   * `Job`, `Session`, `Prompt`, `Diff`, `Snapshot`, `Provider`, `Task`.
   * Regras de negócio (FSM) e contratos estáveis.
2. **Ports (interfaces)**

   * `ProviderPort` (LLM), `SnapshotPort` (LookAtni), `StorePort` (persistência), `QueuePort`, `CachePort`, `MeteringPort`.
3. **Adapters (drivers)**

   * Providers: `openai`, `anthropic`, `gemini` (SDKs oficiais).
   * Store: `postgres` (prisma/knex) ou `gorm` (se Go).
   * Cache: `redis`.
   * Queue: `bullmq` / `asynq` (Go).
   * Metering: `opentelemetry`, `prometheus`, `sentry`.
4. **API (edge)**

   * HTTP (Fastify/Express) + **SSE/WebSocket** (stream).
   * Auth (JWT/Bearer), rate-limit, idempotency.

**Padrão de extensão:**

* Plugin registry: `providers.register("anthropic", new AnthropicAdapter(cfg))`.
* **Feature flags**: `FF_STREAM_PATCH`, `FF_JSON_MODE`, `FF_EVAL_LITE`.

---

# 2) Contratos que NÃO podem quebrar (invariantes)

```ts
// core/ports/provider.ts
export type ProviderId = "openai"|"anthropic"|"gemini";
export interface GenOptions {
  temperature?: number; topP?: number; maxTokens?: number;
  stopSequences?: string[];
}
export interface GenInput {
  system?: string;
  user: string;   // prompt final ou payload mergeado
}
export interface GenUsage { prompt: number; completion: number; total: number; }

export interface GenResult {
  text: string;
  provider: ProviderId;
  model: string;
  usage?: GenUsage;
  finishReason?: string;
}

export interface ProviderPort {
  generate(input: GenInput, options?: GenOptions): Promise<GenResult>;
  stream?(input: GenInput, options?: GenOptions): AsyncIterable<string>;
}
```

```ts
// core/domain/job.ts (FSM enxuta)
export type JobState = "queued"|"snapshotting"|"prompting"|"diffing"|"done"|"error";
export interface Job {
  id: string;
  sessionId: string;
  state: JobState;
  provider: ProviderId;
  model: string;
  promptPurpose: "craft"|"refactor"|"summarize"|"custom";
  createdAt: string; updatedAt: string;
}
```

---

# 3) Fluxos críticos (E2E)

## 3.1 Evolve (refactor com diff)

1. **POST** `/jobs/evolve` → cria `Job(queued)` com `idempotency-key`.
2. **Queue** consome, muda pra `snapshotting`, usa **LookAtni** (`SnapshotPort`) p/ zip/manifest (inclui/exclude).
3. `prompting`: monta prompt com **PromptRegistry** (ver §5).
4. `diffing`: chama `ProviderPort.generate` pedindo **diff unificado**; valida patch.
5. `done`: salva `diff + usage + costs`, emite evento SSE p/ FE.
6. **Download** do `.patch` e `.latx` (se quiser).

## 3.2 Craft Prompt

* Idêntico, sem snapshot. Só passa por `PromptRegistry` + `ProviderPort`.

## 3.3 Streaming

* **SSE** em `/stream/:jobId`: cada chunk em \~1–50ms, bufferizado para o FE com backpressure leve.
* **Fallback** WebSocket só se precisar bidirecional.

---

# 4) APIs (limpas e versionadas)

```
POST   /v1/jobs/evolve           # cria job
GET    /v1/jobs/:id              # status + resumo
GET    /v1/jobs/:id/stream       # SSE live
GET    /v1/jobs/:id/artifacts    # links patch/zip
POST   /v1/providers/test        # sanity check
GET    /v1/providers             # list providers/models
```

### Exemplos (Fastify)

```ts
fastify.post("/v1/jobs/evolve", { preHandler: [auth, idem, rl] }, createEvolveJob);
fastify.get("/v1/jobs/:id/stream", sseJobStream);
```

* `auth`: Bearer/JWT.
* `idem`: cabeçalho `Idempotency-Key` → evita duplicar custo.
* `rl`: rate-limit por IP/org/user.

---

# 5) Prompt Registry (fonte única de verdade)

Pasta versionada `prompts/` (YAML/MDX/JSON5) + loader. Cada prompt tem **purpose**, **inputs**, **guards**.

````yaml
# prompts/refactor-small.yml
id: refactor-small
purpose: small-reviewable-refactor
system: |
  You are a senior engineer...
  Return ONLY a unified diff fenced with ```diff.
inputs:
  - codebase_context
  - style_guide
guards:
  max_lines: 4000
  must_include: ["```diff"]
````

Runtime:

```ts
const { system, template } = registry.load("refactor-small");
const user = template.render({ codebase_context, style_guide });
provider.generate({ system, user }, options);
```

**Vantagem**: trocar vendor sem tocar prompt, e dá pra rodar **evals** por prompt.

---

# 6) Persistência e infra

* **DB**: Postgres (ou SQLite no dev). Tabelas:

  * `jobs(id, state, provider, model, usage_json, costs_cents, created_at, updated_at)`
  * `snapshots(id, job_id, storage_url, bytes, sha256)`
  * `artifacts(id, job_id, kind('patch'|'zip'), storage_url, sha256)`
  * `events(id, job_id, type, payload_json, ts)`
* **Storage**: S3/MinIO pra `.patch`, `.zip/.latx`, manifests.
* **Cache**: Redis (response cache com chave `hash(prompt+opts+model)` + TTL/SLRU).
* **Queue**: BullMQ (Node) ou Asynq (Go).
* **Idempotency**: tabela `idempotency_keys(key, job_id, ttl)`.

---

# 7) Observabilidade & custos (sem isso, a fatura te come)

* **OpenTelemetry** (HTTP, queue, provider-call).
* Métricas Prometheus:

  * `grompt_job_latency_seconds{purpose,provider,model,state}` (histograma)
  * `grompt_tokens_total{provider,model,type="prompt|completion"}`
  * `grompt_cost_cents_total{provider,model}`
  * `grompt_stream_backpressure_events_total`
* **Logs estruturados** (pino/winston) com `job_id`, `session_id`, `provider_call_id`.
* **Sentry** pra exceções + breadcrumb da FSM (state transitions).

---

# 8) Segurança e robustez

* **Rate limit** por org/user (tokens/min e jobs/min).
* **Cotas** (limite diário em custo).
* **Retry/backoff** p/ 429/5xx dos providers (com jitter, no máximo 3 tentativas).
* **Timeouts**: 20–60s, cancelável (AbortController).
* **Sanitização** do prompt (no logs): nunca logar código completo se não tiver flag.
* **Conteúdo sensível**: redigir PII nos logs.

---

# 9) Mercado (perspectiva pragmática)

* **Vendor lock**: abstrair provider agora vale grana quando mudar preço/modelo.
* **Custos**: cache + idempotência reduzem 10–40% do gasto em prompts repetidos.
* **SLA**: fila com prioridade (e.g., “craft-prompt” > “refactor”) melhora percepção do usuário mesmo sem hardware extra.
* **Auditoria**: métricas e artifacts guardados viram **compliance**/histórico (ajuda em B2B).

---

# 10) Plano de migração (sem dor)

1. **Branch “core-ports”**: criar `ports/` + `providers/` (OpenAI/Claude/Gemini) com testes de fumaça.
2. **Branch “queue-stream”**: colocar SSE + fila e idempotency.
3. **Branch “prompt-registry”**: externalizar prompts, congelar contratos.
4. **Cutover**: feature flags — FE continua chamando `POST /evolve`; o handler novo liga no fluxo FSM.
5. **Cleanup**: apagar chamadas direct-to-SDK no handler antigo; manter fallback 1 sprint.

---

# 11) Snippets de referência

## 11.1 Provider adapter (Claude) com retry

```ts
const RETRIES = 3;
for (let i=0;i<RETRIES;i++) {
  try {
    return await anthropic.generate(input, opts);
  } catch (e:any) {
    const s = e?.status ?? e?.response?.status;
    if ((s===429 || (s>=500&&s<600)) && i<RETRIES-1) {
      await sleep(150*(i+1) + Math.random()*100);
      continue;
    }
    throw e;
  }
}
```

## 11.2 SSE stream (Fastify)

```ts
reply.raw.writeHead(200, {
  "Content-Type": "text/event-stream",
  "Cache-Control": "no-cache",
  Connection: "keep-alive",
});
const send = (d:any)=>reply.raw.write(`data: ${JSON.stringify(d)}\n\n`);

queue.onJobUpdate(jobId,(evt)=>send(evt));
request.raw.on("close", ()=>unsubscribe(jobId));
```

## 11.3 Idempotency

```ts
const key = req.headers["idempotency-key"];
if (key) {
  const existing = await idemRepo.get(key);
  if (existing) return reply.send(existing);
}
// cria job, persiste mapping key->job
```

---

# 12) Testes que pegam bug caro

* **Contract tests** por provider (inputs iguais → shape igual).
* **Snapshot determinístico** (mesma seleção de arquivos => mesmo hash/bytes).
* **Diff validator** (aplica patch e roda `tsc eslint vitest - silent`).
* **Chaos**: simular 429 em 30% das chamadas do provider (flag).
