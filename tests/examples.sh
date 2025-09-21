#!/bin/bash

echo "🎮 DEMO - Arquitetura Analyzer"
echo "====================================="

# Start gateway se não estiver rodando
if ! curl -s http://localhost:8080/healthz > /dev/null; then
    echo "🚀 Iniciando gateway..."
    cd /srv/apps/LIFE/KUBEX/analyzer || exit
    ./dist/analyzer-gw &
    sleep 2
fi

# 1. Health Check
echo "1️⃣  Health Check:"
curl -s http://localhost:8080/healthz | jq -C . || curl -s http://localhost:8080/healthz
echo ""

# 2. Providers List
echo "2️⃣  Providers Disponíveis:"
curl -s http://localhost:8080/v1/providers | jq -C .providers || curl -s http://localhost:8080/v1/providers
echo ""

# 3. Provider Config
echo "3️⃣  Configuração dos Providers:"
curl -s http://localhost:8080/v1/providers | jq -C .config || curl -s http://localhost:8080/v1/providers
echo ""

# 4. Chat endpoint demo (fake response para demo)
echo "4️⃣  Chat Endpoint (estrutura da API):"
echo "POST /v1/chat"
echo "{"
echo '  "provider": "gemini",'
echo '  "messages": [{"role": "user", "content": "Analyze this project"}],'
echo '  "temperature": 0.7,'
echo '  "meta": {'
echo '    "analysisType": "general",'
echo '    "projectContext": "# My Project..."'
echo '  }'
echo "}"
echo ""

# 5. Estrutura de resposta SSE
echo "5️⃣  Estrutura de Resposta (SSE):"
echo "data: {\"content\": \"Analyzing project...\", \"done\": false}"
echo "data: {\"content\": \"Based on the context...\", \"done\": false}"
echo "data: {\"done\": true, \"usage\": {\"tokens\": 150, \"latency_ms\": 1200, \"cost_usd\": 0.0003}}"
echo ""

# 6. Arquitetura
echo "7️⃣  Nova Arquitetura:"
echo "Frontend (React) ←→ Gateway (Go) ←→ AI Providers"
echo "     │                   │              │"
echo "     │                   │              ├── OpenAI"
echo "     │                   │              ├── Gemini"
echo "     │                   │              ├── Anthropic (soon)"
echo "     │                   │              └── Groq (soon)"
echo "     │                   │"
echo "     │                   ├── Registry (YAML config)"
echo "     │                   ├── SSE Transport"
echo "     │                   ├── Usage Metrics"
echo "     │                   └── BYOK Support"
echo "     │"
echo "     ├── Direct Gemini (mantido)"
echo "     ├── Gateway API (novo)"
echo "     └── Unified Service (híbrido)"
echo ""
