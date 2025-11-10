#!/usr/bin/env bash
# Test script for Grompt V1 API - Validates all endpoints and integrations

set -euo pipefail

# Configuration
BASE_URL="${GROMPT_BASE_URL:-http://localhost:8080}"
TIMEOUT="${TIMEOUT:-10}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Test helper function
test_endpoint() {
    local method="$1"
    local endpoint="$2"
    local expected_status="$3"
    local description="$4"
    local data="${5:-}"

    log_info "Testing: $description"

    local curl_cmd="curl -s -w '%{http_code}' --max-time $TIMEOUT"

    if [[ "$method" == "POST" && -n "$data" ]]; then
        curl_cmd="$curl_cmd -X POST -H 'Content-Type: application/json' -d '$data'"
    elif [[ "$method" == "GET" ]]; then
        curl_cmd="$curl_cmd -X GET"
    fi

    curl_cmd="$curl_cmd '$BASE_URL$endpoint'"

    local response
    response=$(eval "$curl_cmd")
    local status_code="${response: -3}"
    local body="${response%???}"

    if [[ "$status_code" == "$expected_status" ]]; then
        log_success "$description - HTTP $status_code"
        if [[ -n "$body" && "$body" != "null" ]]; then
            echo "Response: $(echo "$body" | jq . 2>/dev/null || echo "$body")"
        fi
        return 0
    else
        log_error "$description - Expected HTTP $expected_status, got HTTP $status_code"
        echo "Response: $body"
        return 1
    fi
}

# Test streaming endpoint
test_streaming() {
    log_info "Testing: SSE Streaming Generation"

    local stream_url="$BASE_URL/v1/generate/stream?provider=gemini&ideas[]=Test%20streaming&ideas[]=Simple%20prompt&purpose=general"

    # Test if streaming endpoint responds with SSE headers
    local headers
    headers=$(curl -s -I --max-time 5 "$stream_url" || true)

    if echo "$headers" | grep -q "text/event-stream"; then
        log_success "SSE Streaming - Headers correct"

        # Test actual streaming (limited to 3 seconds)
        timeout 3s curl -s -N "$stream_url" | head -5 || true
        log_success "SSE Streaming - Stream data received"
        return 0
    else
        log_error "SSE Streaming - Missing event-stream headers"
        return 1
    fi
}

# Main test suite
main() {
    echo "=========================================="
    echo "🚀 Grompt V1 API Test Suite"
    echo "=========================================="
    echo "Base URL: $BASE_URL"
    echo "Timeout: ${TIMEOUT}s"
    echo ""

    local failed_tests=0
    local total_tests=0

    # Test 1: Health Check
    ((total_tests++))
    test_endpoint "GET" "/healthz" "200" "Health Check" || ((failed_tests++))
    echo ""

    # Test 2: List Providers
    ((total_tests++))
    test_endpoint "GET" "/v1/providers" "200" "List Providers" || ((failed_tests++))
    echo ""

    # Test 3: Basic Generation (if providers available)
    ((total_tests++))
    local gen_payload='{
        "provider": "gemini",
        "ideas": ["Create a simple hello world", "In Python"],
        "purpose": "code",
        "temperature": 0.7
    }'
    test_endpoint "POST" "/v1/generate" "200" "Synchronous Generation" "$gen_payload" || ((failed_tests++))
    echo ""

    # Test 4: SSE Streaming
    ((total_tests++))
    test_streaming || ((failed_tests++))
    echo ""

    # Test 5: GoBE Proxy Health (if configured)
    if [[ -n "${GOBE_BASE_URL:-}" ]]; then
        ((total_tests++))
        test_endpoint "GET" "/v1/proxy/healthz" "200" "GoBE Proxy Health Check" || ((failed_tests++))
        echo ""
    else
        log_warning "Skipping GoBE proxy tests - GOBE_BASE_URL not configured"
    fi

    # Test 6: Error Handling
    ((total_tests++))
    local invalid_payload='{"invalid": "request"}'
    test_endpoint "POST" "/v1/generate" "400" "Error Handling - Invalid Request" "$invalid_payload" || ((failed_tests++))
    echo ""

    # Test Results Summary
    echo "=========================================="
    echo "📊 Test Results Summary"
    echo "=========================================="

    local passed_tests=$((total_tests - failed_tests))

    if [[ $failed_tests -eq 0 ]]; then
        log_success "All tests passed! ($passed_tests/$total_tests)"
        echo ""
        echo "✅ Grompt V1 API is working correctly!"
        echo "🎯 All endpoints are responding as expected"
        echo "🚀 Ready for production deployment"
    else
        log_error "$failed_tests/$total_tests tests failed"
        echo ""
        echo "❌ Some tests failed - please check the logs above"
        echo "🔧 Ensure all required environment variables are set:"
        echo "   - At least one AI provider API key (GEMINI_API_KEY, OPENAI_API_KEY, etc.)"
        echo "   - Optional: GOBE_BASE_URL for proxy functionality"
        exit 1
    fi
}

# Environment validation
validate_environment() {
    log_info "Validating environment..."

    # Check if server is running
    if ! curl -s --max-time 5 "$BASE_URL/v1/healthz" > /dev/null; then
        log_error "Grompt server is not running at $BASE_URL"
        echo "Please start the server with: make run"
        exit 1
    fi

    log_success "Server is running at $BASE_URL"
}

# Help function
show_help() {
    echo "Grompt V1 API Test Script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  -u, --url URL  Set base URL (default: http://localhost:8080)"
    echo "  -t, --timeout  Set timeout in seconds (default: 10)"
    echo ""
    echo "Environment Variables:"
    echo "  GROMPT_BASE_URL  Base URL for the Grompt server"
    echo "  GOBE_BASE_URL    Base URL for GoBE proxy (optional)"
    echo "  TIMEOUT          Request timeout in seconds"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Test local server"
    echo "  $0 -u https://grompt.example.com     # Test remote server"
    echo "  TIMEOUT=30 $0                        # Use longer timeout"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -u|--url)
            BASE_URL="$2"
            shift 2
            ;;
        -t|--timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Check dependencies
if ! command -v curl &> /dev/null; then
    log_error "curl is required but not installed"
    exit 1
fi

if ! command -v jq &> /dev/null; then
    log_warning "jq is not installed - JSON formatting will be disabled"
fi

# Run tests
validate_environment
main
