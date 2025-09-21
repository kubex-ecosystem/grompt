#!/bin/bash

# Test script for the analyzer gateway
# This script demonstrates the gateway functionality

echo "🚀 Testing Analyzer Gateway"
echo "=========================="

# Set OpenAI API key (replace with your key)
export OPENAI_API_KEY="your-openai-api-key-here"

# Start the gateway in background
echo "Starting gateway..."
./dist/analyzer-gw &
GATEWAY_PID=$!

# Wait for gateway to start
sleep 2

# Test health endpoint
echo "Testing health endpoint..."
curl -s http://localhost:8080/healthz | jq

# Test providers endpoint
echo "Testing providers endpoint..."
curl -s http://localhost:8080/v1/providers | jq

# Test chat endpoint with SSE
echo "Testing chat endpoint..."
curl -X POST http://localhost:8080/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "openai",
    "model": "gpt-4o-mini",
    "messages": [
      {"role": "user", "content": "Hello! Say hi in one sentence."}
    ],
    "temperature": 0.7
  }' \
  --no-buffer

# Cleanup
echo "Stopping gateway..."
kill $GATEWAY_PID

echo "✅ Test completed!"
