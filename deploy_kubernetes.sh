#!/bin/bash
# =============================================================================
# Kubernetes Deployment Script for Restaurant Microservices
# =============================================================================
# This script deploys the Restaurant Microservices to a Kubernetes cluster.
# It performs the following steps:
# 1. Updates image references in deployment files
# 2. Creates the necessary Kubernetes resources
# 3. Verifies the deployment
# =============================================================================

set -e

# Set text colors for better readability
GREEN='\033[0;32m'  # Success messages
RED='\033[0;31m'    # Error messages
YELLOW='\033[0;33m' # Warning/info messages
BLUE='\033[0;34m'   # Section headers
NC='\033[0m'        # No Color (reset)

# Default registry if not provided
DEFAULT_REGISTRY="localhost:5000"
REGISTRY=${1:-$DEFAULT_REGISTRY}

# Log function - prints a message with timestamp and color
log() {
  local message=$1
  local color=${2:-$BLUE}
  echo -e "${color}[$(date +"%T")] $message${NC}"
}

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    log "kubectl is not installed. Please install kubectl before proceeding." $RED
    exit 1
fi

# Check if connected to a Kubernetes cluster
if ! kubectl cluster-info &> /dev/null; then
    log "Not connected to a Kubernetes cluster. Please check your kubeconfig." $RED
    exit 1
fi

# Convert Docker Compose images to Kubernetes compatible images
log "Setting up Kubernetes deployment for Restaurant Microservices..." $BLUE

# Update image paths in Kubernetes manifests
log "Updating image references in deployment files..." $YELLOW
sed -i.bak "s|\${DOCKER_REGISTRY}/restaurant-service:latest|${REGISTRY}/restaurant-service:latest|g" kubernetes/06-restaurant-service.yaml
sed -i.bak "s|\${DOCKER_REGISTRY}/feedback-service:latest|${REGISTRY}/feedback-service:latest|g" kubernetes/07-feedback-service.yaml
log "Image references updated." $GREEN

# Ask for confirmation before deploying
read -p "Do you want to deploy to Kubernetes now? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    log "Deployment canceled." $YELLOW
    exit 0
fi

# Apply Kubernetes manifests
log "Applying Kubernetes manifests..." $BLUE
kubectl apply -k kubernetes/

# Wait for deployment to complete
log "Waiting for deployment to complete..." $YELLOW
kubectl wait --namespace restaurant-app --for=condition=ready pod --all --timeout=300s

# Show deployed resources
log "Kubernetes resources created:" $GREEN
echo
log "Pods:" $BLUE
kubectl get pods -n restaurant-app
echo
log "Services:" $BLUE
kubectl get services -n restaurant-app
echo
log "Ingresses:" $BLUE
kubectl get ingress -n restaurant-app
echo

# Show access information
log "Access Information:" $GREEN
INGRESS_IP=$(kubectl get service -n restaurant-app traefik -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
if [ -z "$INGRESS_IP" ]; then
    INGRESS_IP=$(kubectl get service -n restaurant-app traefik -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
fi

if [ -z "$INGRESS_IP" ]; then
    log "No external IP assigned yet. For local clusters, you may need to use port-forwarding:" $YELLOW
    echo "kubectl port-forward -n restaurant-app svc/traefik 8080:80"
    echo "Then access the services at http://localhost:8080"
else
    log "API Gateway: http://${INGRESS_IP}" $GREEN
    log "Restaurant Service: http://${INGRESS_IP}/api/restaurant" $GREEN
    log "Feedback Service: http://${INGRESS_IP}/api/feedback" $GREEN
    log "Traefik Dashboard: http://${INGRESS_IP}:8082/dashboard/" $GREEN
fi

log "Kubernetes deployment completed successfully!" $GREEN
