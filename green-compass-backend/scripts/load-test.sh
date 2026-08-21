#!/usr/bin/env bash
set -euo pipefail

# Load test the GreenCompass API.
# Usage: ./scripts/load-test.sh [base_url] [duration]
# Defaults: http://localhost:8080, 30s
# Requires: curl, optionally vegeta (go install github.com/tsenart/vegeta@latest)

BASE_URL="${1:-http://localhost:8080}"
DURATION="${2:-30}"
RATE="${3:-10}"

echo "Load testing $BASE_URL for ${DURATION}s at ${RATE} req/s"
echo ""

# Check if vegeta is available
if command -v vegeta &> /dev/null; then
    echo "Using vegeta for load testing..."
    
    echo "GET $BASE_URL/health" | vegeta attack \
        -duration="${DURATION}s" \
        -rate="${RATE}" \
        -output=results.bin
    
    vegeta report results.bin
    rm -f results.bin
else
    echo "vegeta not found, using curl-based testing..."
    echo "(Install vegeta: go install github.com/tsenart/vegeta@latest)"
    echo ""
    
    TOTAL=0
    SUCCESS=0
    FAILED=0
    START=$(date +%s)
    
    while true; do
        NOW=$(date +%s)
        ELAPSED=$((NOW - START))
        if [ "$ELAPSED" -ge "$DURATION" ]; then
            break
        fi
        
        STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health")
        TOTAL=$((TOTAL + 1))
        
        if [ "$STATUS" = "200" ]; then
            SUCCESS=$((SUCCESS + 1))
        else
            FAILED=$((FAILED + 1))
        fi
        
        sleep 0.$(shuf -i 1-9 -n 1) 2>/dev/null || sleep 0.1
    done
    
    echo ""
    echo "=== Results ==="
    echo "Total requests: $TOTAL"
    echo "Successful:     $SUCCESS"
    echo "Failed:         $FAILED"
    echo "Duration:       ${DURATION}s"
    if [ "$TOTAL" -gt 0 ]; then
        RPS=$(echo "scale=2; $TOTAL / $DURATION" | bc 2>/dev/null || echo "?")
        echo "Requests/sec:   $RPS"
    fi
fi
