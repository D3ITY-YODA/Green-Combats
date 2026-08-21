#!/usr/bin/env bash
set -euo pipefail

# Test the data connectors against public APIs.
# Usage: ./scripts/test-connectors.sh
# Requires: go

echo "=== Testing Open-Meteo Connector ==="
echo ""

# Test Open-Meteo weather data fetch (Kampala coordinates)
echo "Fetching weather data for Kampala (0.3476°N, 32.5825°E)..."
OPEN_METEO_RESPONSE=$(curl -s "https://api.open-meteo.com/v1/forecast?latitude=0.3476&longitude=32.5825&hourly=temperature_2m,precipitation,soil_moisture_0_to_7cm,wind_speed_10m&past_days=7&forecast_days=0")

if echo "$OPEN_METEO_RESPONSE" | grep -q "temperature_2m"; then
    echo "✓ Open-Meteo: received hourly data"
    TEMPERATURE_COUNT=$(echo "$OPEN_METEO_RESPONSE" | python3 -c "import sys,json; d=json.load(sys.stdin); print(len(d.get('hourly',{}).get('temperature_2m',[])))" 2>/dev/null || echo "?")
    echo "  Temperature data points: $TEMPERATURE_COUNT"
else
    echo "✗ Open-Meteo: unexpected response"
    echo "  Response: ${OPEN_METEO_RESPONSE:0:200}"
fi

echo ""
echo "=== Testing NASA POWER Connector ==="
echo ""

# Test NASA POWER data fetch
echo "Fetching NASA POWER data for Kampala..."
NASA_RESPONSE=$(curl -s "https://power.larc.nasa.gov/api/temporal/daily/point?parameters=T2M,PRECTOTCORR,WS2M,ALLSKY_SFC_SW_DWN&community=AG&longitude=32.5825&latitude=0.3476&start=20250101&end=20250107&format=JSON")

if echo "$NASA_RESPONSE" | grep -q "T2M"; then
    echo "✓ NASA POWER: received daily data"
else
    echo "✗ NASA POWER: unexpected response"
    echo "  Response: ${NASA_RESPONSE:0:200}"
fi

echo ""
echo "=== Running Connector Unit Tests ==="
echo ""

cd "$(dirname "$0")/.."
go test -race -count=1 -v ./internal/connectors/... 2>&1 | tail -20

echo ""
echo "=== All connector tests complete ==="
