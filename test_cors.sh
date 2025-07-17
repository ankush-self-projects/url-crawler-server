#!/bin/bash

# Test CORS headers
echo "Testing CORS headers..."

# Test preflight request (OPTIONS)
echo "Testing OPTIONS request..."
curl -X OPTIONS http://localhost:8080/api/urls \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type,Authorization" \
  -v

echo -e "\n\nTesting actual request..."
# Test actual request
curl -X GET http://localhost:8080/api/urls \
  -H "Origin: http://localhost:3000" \
  -v 