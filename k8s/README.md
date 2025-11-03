# Kubernetes Deployment Guide

This directory contains Kubernetes manifests for deploying the weather service.

## Components

1. **ConfigMap** (`configmap.yaml`)
   - Contains application configuration
   - Environment variables for timeouts and addresses
   - Non-sensitive configuration

2. **Secrets** (`secrets.yaml`)
   - Database connection string
   - Redis password
   - Sensitive configuration

3. **Deployment** (`deployment.yaml`)
   - Runs the weather service
   - 2 replicas by default
   - Resource limits and requests
   - Health checks configured
   - Prometheus annotations

4. **Service** (`service.yaml`)
   - ClusterIP type
   - Exposes port 8080
   - Internal access within cluster

5. **HorizontalPodAutoscaler** (`hpa.yaml`)
   - Scales based on CPU and Memory
   - 2-10 replicas
   - 70% CPU threshold
   - 80% Memory threshold

6. **ServiceMonitor** (`servicemonitor.yaml`)
   - Prometheus integration
   - 15s scrape interval
   - Monitors /metrics endpoint

## Deployment Steps

1. Create namespace (optional):
   ```bash
   kubectl create namespace weather
   ```

2. Apply ConfigMap and Secrets:
   ```bash
   kubectl apply -f configmap.yaml
   kubectl apply -f secrets.yaml
   ```

3. Deploy application:
   ```bash
   kubectl apply -f deployment.yaml
   kubectl apply -f service.yaml
   ```

4. Apply HPA and monitoring:
   ```bash
   kubectl apply -f hpa.yaml
   kubectl apply -f servicemonitor.yaml
   ```

## Prerequisites

- Kubernetes cluster with:
  - Metrics Server (for HPA)
  - Prometheus Operator (for ServiceMonitor)
  - Redis and PostgreSQL installed

## Verification

1. Check deployment:
   ```bash
   kubectl get deployments
   kubectl get pods
   ```

2. Verify service:
   ```bash
   kubectl get svc
   ```

3. Check HPA:
   ```bash
   kubectl get hpa
   ```

4. Monitor metrics:
   ```bash
   kubectl port-forward svc/weather-api 8080:8080
   curl localhost:8080/metrics
   ```

## Notes

- The ServiceMonitor assumes Prometheus Operator is installed
- Adjust resource limits based on actual usage
- Update DB_URL in secrets.yaml with actual database credentials
- Consider using sealed-secrets or external secret management