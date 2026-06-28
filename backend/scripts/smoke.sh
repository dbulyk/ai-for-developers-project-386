#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"

echo "=== Calendar Booking API smoke test ==="
echo "Building server..."
go build -o bin/smoke-server ./cmd/server

cleanup() {
  echo "Cleaning up..."
  if [[ -n "${SERVER_PID:-}" ]]; then
    kill "${SERVER_PID}" 2>/dev/null || true
    wait "${SERVER_PID}" 2>/dev/null || true
  fi
}
trap cleanup EXIT

echo "Starting server on ${BASE_URL}..."
./bin/smoke-server &
SERVER_PID=$!

for i in {1..30}; do
  if curl -sf "${BASE_URL}/health" >/dev/null 2>&1; then
    break
  fi
  sleep 0.1
done

curl -sf "${BASE_URL}/health" >/dev/null || { echo "Server did not start"; exit 1; }
echo "Server is ready"

echo "Creating event type..."
ET_RESPONSE=$(curl -sf -X POST "${BASE_URL}/admin/event-types" \
  -H "Content-Type: application/json" \
  -d '{"name":"Consultation","description":"30 min call","durationMinutes":30}')
ET_ID=$(echo "${ET_RESPONSE}" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")
echo "Event type id: ${ET_ID}"

echo "Listing public event types..."
curl -sf "${BASE_URL}/public/event-types" | python3 -m json.tool

echo "Fetching slots..."
curl -sf "${BASE_URL}/public/event-types/${ET_ID}/slots" | python3 -m json.tool

echo "Creating booking..."
# Pick a slot tomorrow at 09:00 UTC; adjust if it falls on a weekend.
START_TIME=$(date -u -d "tomorrow 09:00" +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u -v+1d -v09H -v00M -v00S +%Y-%m-%dT%H:%M:%SZ)
BK_RESPONSE=$(curl -sf -X POST "${BASE_URL}/public/bookings" \
  -H "Content-Type: application/json" \
  -d "{\"eventTypeId\":\"${ET_ID}\",\"guestName\":\"Alice\",\"startTime\":\"${START_TIME}\"}")
BK_ID=$(echo "${BK_RESPONSE}" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")
echo "Booking id: ${BK_ID}"

echo "Listing admin bookings..."
curl -sf "${BASE_URL}/admin/bookings" | python3 -m json.tool

echo "Deleting booking..."
curl -sf -X DELETE "${BASE_URL}/admin/bookings/${BK_ID}" -o /dev/null

echo "Deleting event type..."
curl -sf -X DELETE "${BASE_URL}/admin/event-types/${ET_ID}" -o /dev/null

echo "=== Smoke test passed ==="
