# Kubernetes Deployment Guide for Restaurant Microservices

This directory contains Kubernetes manifests for deploying the Restaurant Microservices application to a Kubernetes cluster.

## Prerequisites

- Kubernetes cluster (local or cloud-based)
- `kubectl` command-line tool configured to connect to your cluster
- Docker registry with your container images pushed

## Deployment Steps

### 1. Update Image References

Before deploying, update the image references in the `06-restaurant-service.yaml` and `07-feedback-service.yaml` files to point to your Docker registry:

```yaml
image: your-registry/restaurant-service:latest
```

```yaml
image: your-registry/feedback-service:latest
```

### 2. Apply All Resources Using Kustomize

```bash
kubectl apply -k ./
```

Or apply individual files in sequence:

```bash
kubectl apply -f 00-namespace.yaml
kubectl apply -f 01-zookeeper.yaml
kubectl apply -f 02-kafka.yaml
# ... and so on
```

### 3. Verify Deployment

```bash
kubectl get pods -n restaurant-app
kubectl get services -n restaurant-app
kubectl get ingress -n restaurant-app
```

### 4. Access the Services

Once deployed, you can access:
- API Gateway: http://<cluster-ip>
- Restaurant Service: http://<cluster-ip>/api/restaurant
- Feedback Service: http://<cluster-ip>/api/feedback
- Traefik Dashboard: http://<cluster-ip>/dashboard

## Resource Explanations

- **00-namespace.yaml**: Creates a dedicated namespace for the application
- **01-zookeeper.yaml**: Zookeeper deployment for Kafka
- **02-kafka.yaml**: Kafka message broker
- **03-restaurant-db.yaml**: PostgreSQL database for the Restaurant service
- **04-feedback-db.yaml**: PostgreSQL database for the Feedback service
- **05-secrets.yaml**: Secret for database credentials and JWT
- **06-restaurant-service.yaml**: Restaurant ordering microservice
- **07-feedback-service.yaml**: User feedback microservice
- **08-traefik.yaml**: Traefik API Gateway
- **09-ingress.yaml**: Ingress routing rules

## Migrating from Docker Compose

This Kubernetes setup is a cloud-native version of the Docker Compose configuration. Key differences:

1. **Persistence**: Uses PersistentVolumeClaims instead of local volumes
2. **Networking**: Uses Kubernetes Services instead of a Docker bridge network
3. **Configuration**: Uses ConfigMaps and Secrets instead of environment files
4. **Ingress**: Uses Kubernetes Ingress with Traefik instead of port mapping
5. **Scaling**: Supports independent scaling of services

## Future Enhancements

- Add horizontal pod autoscaling (HPA)
- Implement network policies for security
- Add Prometheus and Grafana for monitoring
- Set up GitOps workflow for automatic deployment
