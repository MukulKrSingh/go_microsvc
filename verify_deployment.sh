#!/bin/bash
# =============================================================================
# Microservices Deployment Verification Script
# =============================================================================
# This script checks if all the microservices are running correctly and all
# network connections are properly established. It verifies:
# 1. Docker containers are running
# 2. Network connectivity between services
# 3. Kafka topic configuration
# 4. API endpoints accessibility
# =============================================================================

# Set text colors for better readability
GREEN='\033[0;32m'  # Success messages
RED='\033[0;31m'    # Error messages
YELLOW='\033[0;33m' # Warning/info messages
BLUE='\033[0;34m'   # Section headers
NC='\033[0m'        # No Color (reset)

# === Utility Functions ===

# Log function - prints a message with timestamp and color
log() {
  local message=$1
  local color=${2:-$BLUE}
  echo -e "${color}[$(date +"%T")] $message${NC}"
}

# Success function - prints a success message
success() {
  local message=$1
  log "$message" $GREEN
}

# Error function - prints an error message and exits the script
error() {
  local message=$1
  log "$message" $RED
  exit 1
}

# === Prerequisite Checks ===

# Verify Docker is installed and running
if ! docker info > /dev/null 2>&1; then
  error "Docker is not running or not installed. Please start Docker first."
fi

# Check if any services are running via docker-compose
log "Checking if services are running..." $YELLOW
if ! docker-compose ps | grep -q "Up"; then
  error "No services are running. Please start them with 'docker-compose up -d' first."
fi

# === Container Health Checks ===
# Verify each required container is running correctly
log "Checking container status:" $BLUE
# List of all required containers that should be running
containers=("zookeeper" "kafka" "restaurant-db" "restaurant-service" "feedback-db" "feedback-service")
for container in "${containers[@]}"; do
  # Get the container status using Docker inspect
  status=$(docker inspect --format '{{.State.Status}}' "$container" 2>/dev/null)
  
  # Verify container is in "running" state
  if [ "$status" == "running" ]; then
    success "✅ $container is running"
  else
    error "❌ $container is not running properly (status: $status)"
  fi
done

echo

# === Network Connectivity Checks ===
# Verify all services can communicate with each other within the Docker network
log "Checking network connectivity between services:" $BLUE

# Check restaurant service connection to Kafka
# This is critical for publishing order events
if docker exec restaurant-service ping -c 2 kafka > /dev/null 2>&1; then
  success "✅ restaurant-service can reach kafka"
else
  error "❌ restaurant-service cannot reach kafka"
fi

# Check feedback service connection to Kafka
# This is critical for consuming order events
if docker exec feedback-service ping -c 2 kafka > /dev/null 2>&1; then
  success "✅ feedback-service can reach kafka"
else
  error "❌ feedback-service cannot reach kafka"
fi

# Check restaurant service connection to its database
# This is necessary for order management and persistence
if docker exec restaurant-service ping -c 2 restaurant-db > /dev/null 2>&1; then
  success "✅ restaurant-service can reach restaurant-db"
else
  error "❌ restaurant-service cannot reach restaurant-db"
fi

# Check feedback service connection to its database
# This is necessary for feedback storage and analytics
if docker exec feedback-service ping -c 2 feedback-db > /dev/null 2>&1; then
  success "✅ feedback-service can reach feedback-db"
else
  error "❌ feedback-service cannot reach feedback-db"
fi

echo

# === Kafka Configuration Checks ===
# Verify Kafka topics are set up correctly
log "Checking Kafka topics:" $BLUE
topics=$(docker exec kafka kafka-topics --bootstrap-server kafka:9092 --list 2>/dev/null)

# The "orders" topic is used for communication between services
if [[ $topics == *"orders"* ]]; then
  success "✅ 'orders' topic exists in Kafka"
else
  # This is not an error as the topic can be auto-created when first message is sent
  log "Topic 'orders' doesn't exist yet. This is expected if no messages have been sent." $YELLOW
fi

echo

# === API Accessibility Checks ===
# Verify service APIs are accessible from the host machine
log "Checking API accessibility from host machine:" $BLUE

# Check restaurant service API accessibility
# This test hits a public endpoint that doesn't require authentication
if curl -s "http://localhost:8080/food-items" > /dev/null; then
  success "✅ Restaurant service API is accessible on port 8080"
else
  error "❌ Cannot access restaurant service API on port 8080"
fi

# Check feedback service API accessibility
# We test the health endpoint which doesn't require authentication
if curl -s "http://localhost:8081/health" > /dev/null 2>&1 || curl -s "http://localhost:8081" > /dev/null 2>&1; then
  success "✅ Feedback service API is accessible on port 8081"
else
  error "❌ Cannot access feedback service API on port 8081"
fi

echo
# Final success message
log "All containers and network connections are working properly!" $GREEN
log "You can now run the integration tests with ./test_microservices.sh" $GREEN
