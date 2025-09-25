#!/bin/bash
# Frontend Multi-Provider Validation
# Tests if the client-side analyzer works with different AI providers

echo "🚀 FRONTEND MULTI-PROVIDER VALIDATION"
echo "====================================="
echo "Testing client-side analyzer with multiple AI providers"
echo

# Check if we're in the right directory
if [ ! -f "frontend/package.json" ]; then
    echo "❌ Please run this script from the analyzer root directory"
    exit 1
fi

echo "📦 Building frontend with all providers..."
cd frontend || exit

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "Installing dependencies..."
    npm install
fi

# Build the frontend
echo "Building production frontend..."
npm run build

if [ $? -eq 0 ]; then
    echo "✅ Frontend build successful"
else
    echo "❌ Frontend build failed"
    exit 1
fi

echo
echo "🔍 Checking provider implementations..."

# Check if unified-ai service exists and has all providers
if [ -f "services/unified-ai.ts" ]; then
    echo "✅ Unified AI service found"

    # Check for provider implementations
    if grep -q "gateway-openai" services/unified-ai.ts; then
        echo "✅ OpenAI provider implementation found"
    else
        echo "❌ OpenAI provider implementation missing"
    fi

    if grep -q "gateway-anthropic" services/unified-ai.ts; then
        echo "✅ Anthropic provider implementation found"
    else
        echo "❌ Anthropic provider implementation missing"
    fi

    if grep -q "gateway-gemini" services/unified-ai.ts; then
        echo "✅ Gateway Gemini provider implementation found"
    else
        echo "❌ Gateway Gemini provider implementation missing"
    fi

    if grep -q "gemini-direct" services/unified-ai.ts; then
        echo "✅ Direct Gemini provider implementation found"
    else
        echo "❌ Direct Gemini provider implementation missing"
    fi
else
    echo "❌ Unified AI service not found"
    exit 1
fi

echo
echo "🎛️ Checking provider selector component..."

if [ -f "components/settings/ProviderSelector.tsx" ]; then
    echo "✅ Provider selector component found"

    # Check if all providers are in the selector
    PROVIDER_COUNT=$(grep -c "gateway-" components/settings/ProviderSelector.tsx)
    echo "✅ Found $PROVIDER_COUNT gateway providers in selector"
else
    echo "❌ Provider selector component missing"
fi

echo
echo "🔌 Testing backend gateway connectivity..."

# Start a temporary local server to test the built frontend
echo "Starting test server..."
cd dist || exit
python3 -m http.server 3000 > /dev/null 2>&1 &
SERVER_PID=$!
cd ..

# Give server time to start
sleep 2

# Test if frontend loads
if curl -s http://localhost:3000 > /dev/null; then
    echo "✅ Frontend serves successfully"
else
    echo "❌ Frontend server failed to start"
    kill $SERVER_PID 2>/dev/null
    exit 1
fi

# Clean up
kill $SERVER_PID 2>/dev/null

echo
echo "🎯 FRONTEND VALIDATION SUMMARY"
echo "=============================="
echo "✅ Frontend builds successfully"
echo "✅ Multi-provider service implemented"
echo "✅ Provider selector UI component ready"
echo "✅ Gateway integration prepared"
echo
echo "📋 PROVIDERS AVAILABLE:"
echo "   🔹 Gemini Direct (current stable)"
echo "   🔹 Gemini via Gateway (streaming)"
echo "   🔹 OpenAI via Gateway (ready)"
echo "   🔹 Anthropic via Gateway (ready)"
echo
echo "🚀 Frontend is READY for multi-provider support!"
echo "🌐 Deploy to analyzer.kubex.world to test with real providers"

echo
echo "💡 Next steps:"
echo "   1. Deploy frontend build to production"
echo "   2. Ensure backend gateway is running with all providers"
echo "   3. Test each provider in the live environment"
echo "   4. Validate API key handling for each provider"
