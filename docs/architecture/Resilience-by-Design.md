# Resilience by Design — Kortex/Kubex (Outline)

> Descrever a arquitetura resiliente e o “offline-first com melhoria progressiva”.
> Documento curto e objetivo, focado em produto.

## 1. Princípios
- *Offline-first*, reconexão automática, *event broadcasting* interno
- UI com indicadores de conectividade e frescor
- Fallback seguro sem MCP, upgrades com WebSockets quando presente

## 2. Padrão de Estados
- `AppContext` global (Next.js/TS strict)
- Filas de eventos, handlers idempotentes
- Retentativas com *backoff* e circuit breakers

## 3. Integração MCP
- Conexão WebSocket com reconexão
- Logs e métricas em tempo real
- Estratégia de degradação graciosa

## 4. Build/Deploy
- Static export (`output: 'export'`) para GitHub Pages
- Imagens `unoptimized: true`, trailingSlash
- CI com lint/test/build e checagens de qualidade

## 5. Aceite
- Build estável
- Sem *runtime errors* offline
- Realtime upgrades funcionando quando MCP disponível