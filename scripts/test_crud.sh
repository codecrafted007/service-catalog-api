#!/bin/bash

API_URL="http://localhost:8080"
API_KEY=$1

echo "==> [1] Creating 1000 services..."
for i in $(seq 1 1000); do
    curl -s -X POST $API_URL/services \
        -H "Content-Type: application/json" \
        -H "X-API-Key: $API_KEY" \
        -d '{
      "name": "Service '$i'",
      "description": "Description for service '$i'",
      "version": "v1.0.0",
      "changelog": "Initial release for service '$i'"
    }' >/dev/null
done

echo "==> [2] Listing services without any query parameters (should return default 20)..."
curl -s -X GET "$API_URL/services" \
    -H "X-API-Key: $API_KEY" | jq .

echo "==> [3] Filtering services by name with sort=name..."
curl -s -X GET "$API_URL/services?filter=Service%202&sort=name" \
    -H "X-API-Key: $API_KEY" | jq .

echo "==> [3] Filtering services by description with sort=createdAt..."
curl -s -X GET "$API_URL/services?filter=Description%203&sort=createdAt" \
    -H "X-API-Key: $API_KEY" | jq .

echo "==> [4] Get Service by ID (serviceId=10)..."
curl -s -X GET "$API_URL/services/10" \
    -H "X-API-Key: $API_KEY" | jq .

echo "==> [5] Creating a patch version (v1.0.1) for serviceId=10..."
VERSION_ID=$(curl -s -X POST "$API_URL/services/10/versions" \
    -H "Content-Type: application/json" \
    -H "X-API-Key: $API_KEY" \
    -d '{
    "version": "v1.0.1",
    "changelog": "Patch release"
  }' | jq -r '.data.id')

echo "==> New version created with ID: $VERSION_ID"

echo "==> [6] Get Service by ID (serviceId=10) after patch version..."
curl -s -X GET "$API_URL/services/10" \
    -H "X-API-Key: $API_KEY" | jq .

echo "==> [7] Deleting version ID $VERSION_ID..."
curl -s -X DELETE "$API_URL/versions/$VERSION_ID" \
    -H "X-API-Key: $API_KEY"

echo "==> [8] Get Service by ID (serviceId=10) after deletion..."
curl -s -X GET "$API_URL/services/10" \
    -H "X-API-Key: $API_KEY" | jq .
