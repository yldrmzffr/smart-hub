# Smart Hub Kubernetes Deployment Guide 🚀

## Configuration Files and Settings 🛠️

### deployment.yaml (Smart Hub)
- 🔄 **Replicas**: Set with `spec.replicas` field
- 💻 **Container Resources**:
  ```yaml
  resources:
    limits:
      cpu: "1"      # 👉 Maximum CPU usage
      memory: "512Mi" # 👉 Maximum memory usage
    requests:
      cpu: "500m"    # 👉 Guaranteed CPU
      memory: "256Mi" # 👉 Guaranteed memory
  ```
- 🏥 **Health Check Strategy**:
    - Readiness Probe:
      ```yaml
      readinessProbe:
        tcpSocket:
          port: 50051
        initialDelaySeconds: 5  # 👉 Start checking after 5s
        periodSeconds: 10       # 👉 Check every 10s
      ```
    - Liveness Probe:
      ```yaml
      livenessProbe:
        tcpSocket:
          port: 50051
        initialDelaySeconds: 15 # 👉 Start checking after 15s
        periodSeconds: 20       # 👉 Check every 20s
      ```

### postgres-deployment.yaml 🐘
- 💾 **Storage Setup**:
  ```yaml
  volumeMounts:
    - name: postgres-storage
      mountPath: /var/lib/postgresql/data  # 👉 Data storage path
  ```
- 🏥 **PostgreSQL Health Check**:
  ```yaml
  readinessProbe:
    exec:
      command: ["pg_isready", "-U", "postgres"]  # 👉 Checks if DB is ready
  ```

### service.yaml 🔗
- 🌐 **Service Type**: `ClusterIP` (internal cluster communication)
- 🔌 **Port Setup**:
  ```yaml
  ports:
    - protocol: TCP
      port: 50051        # 👉 Service port
      targetPort: 50051  # 👉 Pod port
  ```

### postgres-pvc.yaml 💾
- 📦 **Storage Request**:
  ```yaml
  resources:
    requests:
      storage: 1Gi  # 👉 Requested storage size
  ```
- 📝 **Access Mode**: `ReadWriteOnce` (single pod write access)

### secret.yaml 🔐
- 🔑 **Credentials**:
  ```yaml
  stringData:
    username: "postgres"     # 👉 DB username
    password: "postgres123"  # 👉 DB password
  ```

## Health Check Strategies 🏥

### Smart Hub Health Checks
- ✅ TCP Socket check: Checks if app listens on port 50051
- ⏰ Readiness: Checks if pod can receive traffic
- 💓 Liveness: Checks if pod is running

### PostgreSQL Health Checks
- ✅ pg_isready: Checks database connection
- ⏰ Readiness: Starts after 5s, checks every 10s
- 💓 Liveness: Starts after 30s, checks every 20s

## Resource Management Tips 💡

### Smart Hub
- 🎯 CPU request/limit: 2:1 (500m:1)
- 🎯 Memory request/limit: 2:1 (256Mi:512Mi)
- 🔍 Optimized for average workload

### PostgreSQL
- 🎯 CPU request/limit: 2:1 (1:2)
- 🎯 Memory request/limit: 2:1 (1Gi:2Gi)
- 🔍 Higher resources for database

## Local Development with Minikube 🛠️

### Prerequisites
- Install Minikube
- Install kubectl
- Docker installed

### Start Local Environment
```bash
# Start Minikube
minikube start

# Build local image
docker build -t smart-hub:latest .

# Load image into Minikube
minikube image load smart-hub:latest

# Apply Kubernetes configs
kubectl apply -f secret.yaml
kubectl apply -f postgres-pvc.yaml
kubectl apply -f postgres-deployment.yaml
kubectl apply -f postgres-service.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

# Check status
kubectl get all
```

### Cleanup
```bash
# Delete all resources
kubectl delete -f .

# Stop Minikube
minikube stop
```