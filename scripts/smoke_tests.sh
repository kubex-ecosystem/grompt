#!/bin/bash

# Smoke tests for Analyzer Gateway - Day 1 validation
# Tests the golden path endpoints with robust health checks

set -e

BASE_URL="http://localhost:8666"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
log_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((TESTS_PASSED++))
}

log_error() {
    echo -e "${RED}[FAIL]${NC} $1"
    ((TESTS_FAILED++))
}

# Test helper function
test_endpoint() {
    local name="$1"
    local method="$2"
    local endpoint="$3"
    local expected_status="$4"
    local data="$5"

    log_info "Testing $name..."

    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "HTTPSTATUS:%{http_code}" -H "Accept: application/json" "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X "$method" -H "Content-Type: application/json" -H "Accept: application/json" -d "$data" "$BASE_URL$endpoint")
    fi

    http_code=$(echo "$response" | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo "$response" | sed -e 's/HTTPSTATUS\:.*//g')

    if [ "$http_code" -eq "$expected_status" ]; then
        log_success "$name - HTTP $http_code"
        echo "Response: $body" | jq . 2>/dev/null || echo "Response: $body"
        echo
        return 0
    else
        log_error "$name - Expected HTTP $expected_status, got $http_code"
        echo "Response: $body"
        echo
        return 1
    fi
}

# Check if server is running
check_server() {
    log_info "Checking if server is running at $BASE_URL..."
    if curl -s "$BASE_URL/healthz" > /dev/null; then
        log_success "Server is running"
        return 0
    else
        log_error "Server is not running. Start with: ./dist/analyzer_linux_amd64 gateway serve -p 8666"
        exit 1
    fi
}

# Main test suite
main() {
    echo "🧪 Analyzer Gateway Smoke Tests - Day 1 Golden Path"
    echo "=================================================="
    echo

    # Check server availability
    check_server

    # Test 1: Health check
    test_endpoint "Health Check" "GET" "/healthz" 200

    # Test 2: List providers
    test_endpoint "List Providers" "GET" "/v1/providers" 200

    # Test 3: Advise - exec mode
    exec_payload='{
        "mode": "exec",
        "context": {
            "repository": "test-repo",
            "scorecard": {"chi": 70}
        },
        "options": {
            "timeout_sec": 30
        }
    }'
    test_endpoint "Advise - Exec Mode" "POST" "/v1/advise" 200 "$exec_payload"

    # Test 4: Advise - code mode
    code_payload='{
        "mode": "code",
        "context": {
            "repository": "test-repo",
            "hotspots": ["src/main.go", "internal/api.go"]
        }
    }'
    test_endpoint "Advise - Code Mode" "POST" "/v1/advise" 200 "$code_payload"

    # Test 5: Production status
    test_endpoint "Production Status" "GET" "/v1/status" 200

    # Test 6: Invalid advise mode (should fail)
    invalid_payload='{
        "mode": "invalid",
        "context": {}
    }'
    test_endpoint "Invalid Advise Mode" "POST" "/v1/advise" 400 "$invalid_payload"

    # Test Summary
    echo
    echo "🎯 Test Summary"
    echo "==============="
    echo -e "${GREEN}Tests Passed: $TESTS_PASSED${NC}"
    echo -e "${RED}Tests Failed: $TESTS_FAILED${NC}"

    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "${GREEN}🎉 All tests passed! Golden path is working.${NC}"
        exit 0
    else
        echo -e "${RED}❌ Some tests failed. Check the output above.${NC}"
        exit 1
    fi
}

# BF1_MODE test
test_bf1_mode() {
    echo
    log_info "Testing BF1_MODE functionality..."

    # Set BF1_MODE and restart server (this would be manual in real scenario)
    export BF1_MODE=true

    # Test that BF1 headers are present
    log_info "Checking BF1 mode headers..."
    response=$(curl -s -I -X POST -H "Content-Type: application/json" -d '{"mode":"exec"}' "$BASE_URL/v1/advise")

    if echo "$response" | grep -q "X-BF1-Mode: true"; then
        log_success "BF1 mode headers present"
    else
        log_error "BF1 mode headers missing"
    fi
}

function test_webhook_endpoints() {
    log_info "🔄 Testing Meta-Recursive Webhook Endpoints..."

    # Test webhook health
    test_endpoint "Webhook Health" "GET" "/v1/webhooks/health" 200

    # Test webhook processing with mock GitHub push event
    local github_payload='{
        "repository": {
            "full_name": "test-owner/test-repo",
            "name": "test-repo",
            "owner": {
                "login": "test-owner"
            }
        },
        "commits": [
            {
                "id": "abc123",
                "message": "Add new feature",
                "author": {
                    "name": "Test User"
                }
            }
        ]
    }'

    # GitHub webhook with proper headers
    log_info "Testing GitHub webhook processing..."
    response=$(curl -s -w "%{http_code}" -X POST \
        -H "Content-Type: application/json" \
        -H "X-GitHub-Event: push" \
        -d "$github_payload" \
        "$BASE_URL/v1/webhooks")

    status_code="${response: -3}"
    if [ "$status_code" = "202" ] || [ "$status_code" = "200" ]; then
        log_success "GitHub webhook processing: Status $status_code"
    else
        log_error "GitHub webhook processing failed: Status $status_code"
        echo "Response: ${response%???}"
    fi

    log_info "🎉 Webhook tests completed!"
}

# Run main test suite
main

# Optionally run BF1 tests
if [ "$1" = "--bf1" ]; then
    test_bf1_mode
fi

# Test webhook endpoints if requested
if [ "$1" = "--webhooks" ] || [ "$2" = "--webhooks" ]; then
    test_webhook_endpoints
fi
