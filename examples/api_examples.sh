#!/bin/bash

# API Examples for Packing Service with PostgreSQL
# Make sure the service is running on localhost:8080

BASE_URL="http://localhost:8080/api/v1"

echo "=== Packing Service API Examples ==="
echo

# Test health check
echo "1. Health Check:"
curl -s "$BASE_URL/../health"
echo -e "\n"

# Get current pack sizes
echo "2. Get Current Pack Sizes:"
curl -s "$BASE_URL/config" | jq .
echo -e "\n"

# List all pack sizes (with database info)
echo "3. List All Pack Sizes (from database):"
curl -s "$BASE_URL/pack-sizes" | jq .
echo -e "\n"

# Calculate packing for 501 items
echo "4. Calculate Packing for 501 items:"
curl -s -X POST "$BASE_URL/calculate" \
  -H "Content-Type: application/json" \
  -d '{"items": 501}' | jq .
echo -e "\n"

# Create a new pack size
echo "5. Create New Pack Size (750):"
curl -s -X POST "$BASE_URL/pack-sizes" \
  -H "Content-Type: application/json" \
  -d '{"size": 750, "is_active": true}' | jq .
echo -e "\n"

# List pack sizes again to see the new one
echo "6. List Pack Sizes (after adding 750):"
curl -s "$BASE_URL/pack-sizes" | jq .
echo -e "\n"

# Calculate packing again with the new pack size
echo "7. Calculate Packing for 501 items (with new pack size):"
curl -s -X POST "$BASE_URL/calculate" \
  -H "Content-Type: application/json" \
  -d '{"items": 501}' | jq .
echo -e "\n"

# Update a pack size (deactivate 250)
echo "8. Deactivate pack size 250:"
PACK_ID=$(curl -s "$BASE_URL/pack-sizes" | jq -r '.pack_sizes[] | select(.size == 250) | .id')
curl -s -X PUT "$BASE_URL/pack-sizes/$PACK_ID" \
  -H "Content-Type: application/json" \
  -d '{"size": 250, "is_active": false}' | jq .
echo -e "\n"

# Calculate packing without 250 pack size
echo "9. Calculate Packing for 251 items (without 250 pack size):"
curl -s -X POST "$BASE_URL/calculate" \
  -H "Content-Type: application/json" \
  -d '{"items": 251}' | jq .
echo -e "\n"

echo "=== Examples Complete ==="
