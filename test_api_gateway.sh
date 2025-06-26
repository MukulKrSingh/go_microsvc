#!/bin/bash
# =============================================
# API Gateway Test Script
# =============================================
# This script tests the API Gateway endpoints
# to ensure that the gateway is properly routing
# requests to the appropriate microservices.
# =============================================

# Set text colors for better readability
GREEN='\033[0;32m'  # Success messages
RED='\033[0;31m'    # Error messages
YELLOW='\033[0;33m' # Warning/info messages
BLUE='\033[0;34m'   # Section headers
NC='\033[0m'        # No Color (reset)

# API Gateway endpoint
GATEWAY="http://localhost"

# Function to print colored messages
print_message() {
  local color=$1
  local message=$2
  echo -e "${color}${message}${NC}"
}

# Function to test an endpoint
test_endpoint() {
  local endpoint=$1
  local method=${2:-GET}
  local data=$3
  local description=$4

  print_message $BLUE "\n=== Testing $description ==="
  print_message $YELLOW "Endpoint: $method $endpoint"
  
  # Execute the request
  if [ "$method" == "GET" ]; then
    response=$(curl -s -w "\n%{http_code}" "$endpoint")
  else
    response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" "$endpoint")
  fi
  
  # Extract status code (the last line)
  status_code=$(echo "$response" | tail -n1)
  
  # Extract the response body (everything except the last line)
  body=$(echo "$response" | sed '$d')
  
  # Check if the request was successful
  if [ "$status_code" -ge 200 ] && [ "$status_code" -lt 300 ]; then
    print_message $GREEN "✅ Success (Status: $status_code)"
    # Pretty print JSON if possible
    echo "$body" | python3 -m json.tool 2>/dev/null || echo "$body"
  else
    print_message $RED "❌ Failed (Status: $status_code)"
    echo "$body"
  fi
}

print_message $BLUE "========================================"
print_message $BLUE "    API GATEWAY TEST SCRIPT"
print_message $BLUE "========================================"

# Test API Gateway health
test_endpoint "$GATEWAY" "GET" "" "API Gateway Health"

# Test Restaurant Service via API Gateway
test_endpoint "$GATEWAY/api/restaurant/food-items" "GET" "" "Restaurant Service - Food Items"

# Test authentication via API Gateway
test_endpoint "$GATEWAY/api/restaurant/auth" "POST" '{"username":"testuser", "password":"password123"}' "Restaurant Service - Authentication"

# Test Feedback Service via API Gateway
test_endpoint "$GATEWAY/api/feedback/health" "GET" "" "Feedback Service - Health Check"

print_message $BLUE "\n========================================"
print_message $BLUE "    TEST COMPLETE"
print_message $BLUE "========================================"
