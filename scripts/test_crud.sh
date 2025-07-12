#!/bin/bash

API_KEY="$1"
BASE_URL="http://localhost:8080"

echo "Creating a new service..."
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/services" \
    -H "Content-Type: application/json" \
    -H "X-API-Key: $API_KEY" \
    -d '{"name":"Test Service","description":"Test Description"}')

echo "$CREATE_RESPONSE"
SERVICE_ID=$(echo "$CREATE_RESPONSE" | jq -r '.data')

echo "Created service with ID: $SERVICE_ID"
echo

echo "Fetching created service..."
curl -s -H "X-API-Key: $API_KEY" "$BASE_URL/services/$SERVICE_ID" | jq
echo

echo "Updating service..."
curl -s -X PUT "$BASE_URL/services/$SERVICE_ID" \
    -H "Content-Type: application/json" \
    -H "X-API-Key: $API_KEY" \
    -d '{"name":"Updated Service","description":"Updated Description"}' | jq
echo

echo "Fetching updated service..."
curl -s -H "X-API-Key: $API_KEY" "$BASE_URL/services/$SERVICE_ID" | jq
echo

echo "Deleting service..."
curl -s -X DELETE "$BASE_URL/services/$SERVICE_ID" \
    -H "X-API-Key: $API_KEY" | jq
echo

echo "Verifying deletion..."
curl -s -H "X-API-Key: $API_KEY" "$BASE_URL/services/$SERVICE_ID" | jq
