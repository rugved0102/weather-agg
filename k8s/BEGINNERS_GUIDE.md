mini# Kubernetes Guide for Beginners

## What We Created

Think of Kubernetes (K8s) as a manager for your containers. We created 6 configuration files that tell Kubernetes how to run your weather service:

1. **ConfigMap** - Like a settings file
   - Stores non-secret settings
   - Example: what port to use, timeout values

2. **Secret** - Like a password file
   - Stores sensitive information
   - Example: database passwords

3. **Deployment** - The main application setup
   - Tells K8s how to run your weather service
   - How many copies to run (replicas)
   - What resources it needs (CPU/memory)

4. **Service** - Like a network router
   - Makes your application accessible
   - Routes traffic to your pods

5. **HPA** (Horizontal Pod Autoscaler)
   - Automatically adds or removes copies of your app
   - Based on CPU/memory usage

6. **ServiceMonitor**
   - Helps Prometheus collect metrics
   - Monitors your application's health

## Step-by-Step Testing Guide

### 1. First, Install Required Tools

```powershell
# Install kubectl (Kubernetes command-line tool)
# For Windows, using chocolatey:
choco install kubernetes-cli

# Install minikube (Local Kubernetes cluster)
choco install minikube
```

### 2. Start Local Cluster

```powershell
# Start minikube
minikube start

# Verify it's running
kubectl get nodes
```

### 3. Build and Load Your Image

```powershell
# Build your image
docker build -t weather-agg-weather-api:latest .

# Load image into minikube
minikube image load weather-agg-weather-api:latest
```

### 4. Deploy Your Application

```powershell
# Apply configurations one by one
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/hpa.yaml
kubectl apply -f k8s/servicemonitor.yaml
```

### 5. Check Everything is Running

```powershell
# Check pods are running
kubectl get pods

# Check service is created
kubectl get services

# Check deployments
kubectl get deployments

# Check HPA
kubectl get hpa
```

### 6. Access Your Application

```powershell
# Create a tunnel to your service
minikube service weather-api --url

# In another terminal, test the service
curl http://[URL_FROM_ABOVE]/weather?city=London
```

### 7. Monitor Your Application

```powershell
# View pod logs
kubectl logs -l app=weather-api

# Watch pod scaling
kubectl get pods -w

# Check metrics
kubectl top pods
```

## Understanding What's Happening

1. **When you deploy**:
   - Kubernetes creates pods (containers)
   - Starts your weather service
   - Sets up networking
   - Configures monitoring

2. **When requests come in**:
   - Service routes traffic to pods
   - HPA watches resource usage
   - Adds/removes pods as needed

3. **For monitoring**:
   - Prometheus collects metrics
   - You can see request counts, response times
   - Helps identify issues

## Common Commands for Debugging

```powershell
# See pod details
kubectl describe pod [pod-name]

# Check pod logs
kubectl logs [pod-name]

# Get into a pod (like docker exec)
kubectl exec -it [pod-name] -- /bin/sh

# Delete everything and start over
kubectl delete -f k8s/
```

## What to Watch For

1. **Pod Status**
   - Should show "Running"
   - If "Error" or "CrashLoopBackOff", check logs

2. **Resource Usage**
   - Watch CPU and Memory
   - HPA will scale if too high

3. **Logs**
   - Check for errors
   - Monitor request flow

## Prerequisites Check

✅ Minikube installed
✅ kubectl installed
✅ Docker running
✅ Weather service image built
✅ Kubernetes configs ready

## Troubleshooting Tips

1. **Pods not starting?**
   ```powershell
   kubectl describe pod [pod-name]
   ```

2. **Can't access service?**
   ```powershell
   kubectl get svc weather-api
   minikube service list
   ```

3. **Application errors?**
   ```powershell
   kubectl logs -l app=weather-api
   ```

4. **Resource issues?**
   ```powershell
   kubectl top pods
   kubectl get events
   ```

## Next Steps

1. Learn about pod logs and debugging
2. Understand Kubernetes resource limits
3. Explore Kubernetes dashboard
   ```powershell
   minikube dashboard
   ```

## Need Help?

If you see errors:
1. Check pod status and logs
2. Verify configurations
3. Ensure prerequisites are met
4. Ask for specific error messages