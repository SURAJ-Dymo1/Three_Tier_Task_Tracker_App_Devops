kubectl run debug --rm -it --image=busybox -- sh

In production you should **never hard-code secrets** like this in code:

```go
"mongodb://admin:admin123@my-mongo:27017/mydatabase?authSource=admin"
```

Real systems use **environment variables + secret managers**.
This is the standard approach in containers, **Docker**, and **Kubernetes**.

I'll show you the **proper industry workflow**.

---

# 1️⃣ Modify Your Go Code (Use Environment Variable)

Replace the hardcoded URL with an environment variable.

```go
import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	mongoURI := os.Getenv("MONGO_URI")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
}
```

Now your **secret is not in code**.

---

# 2️⃣ Pass Secret in Docker

Run container like this:

```bash
docker run -d \
--name my-backend \
--network my-net \
-p 8081:8080 \
-e MONGO_URI="mongodb://admin:admin123@mongodb:27017/mydatabase?authSource=admin" \
my-backend
```

Your Go app reads it from:

```bash
os.Getenv("MONGO_URI")
```

---

# 3️⃣ Even Better (Use `.env` file)

Create file:

```
.env
```

```
MONGO_URI=mongodb://admin:admin123@mongodb:27017/mydatabase?authSource=admin
```

Run:

```bash
docker run --env-file .env my-backend
```

⚠️ Add `.env` to `.gitignore`.

---

# 4️⃣ Production Method in Docker Compose

Example `docker-compose.yml`

```yaml
version: "3.9"

services:

  mongodb:
    image: mongo
    container_name: mongodb
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: admin123
    networks:
      - my-net

  backend:
    build: .
    ports:
      - "8081:8080"
    environment:
      MONGO_URI: mongodb://admin:admin123@mongodb:27017/mydatabase?authSource=admin
    depends_on:
      - mongodb
    networks:
      - my-net

networks:
  my-net:
```

Run everything:

```bash
docker compose up -d
```

---

# 5️⃣ Production Method in Kubernetes (Best Practice)

Use **Secrets**.

Create secret:

```bash
kubectl create secret generic mongo-secret \
--from-literal=MONGO_URI="mongodb://admin:admin123@mongodb:27017/mydatabase?authSource=admin"
```

Then use it in Deployment:

```yaml
env:
- name: MONGO_URI
  valueFrom:
    secretKeyRef:
      name: mongo-secret
      key: MONGO_URI
```

Now your app still reads:

```go
os.Getenv("MONGO_URI")
```

But the value is stored securely in **Kubernetes Secrets.

---

# 6️⃣ Enterprise Level Secret Storage

In real companies secrets are stored in tools like:

* HashiCorp Vault
* AWS Secrets Manager
* Azure Key Vault

Kubernetes pulls secrets dynamically from these systems.

---

# 🚀 Recommended DevOps Structure

```
project
│
├── backend
│   ├── main.go
│   ├── Dockerfile
│
├── docker-compose.yml
│
├── .env
│
└── .gitignore
```

`.gitignore`

```
.env
```

---

# ⭐ DevOps Tip

Never store secrets in:

* GitHub
* Dockerfile
* Source code
* Images

Always use:

```
Environment Variables
Secrets
Secret Managers
```

---

✅ If you want, I can also show you **how companies structure a real 3-tier production project (React + Go + Mongo + Docker + Kubernetes + CI/CD)** so your DevOps project becomes **portfolio-level strong**.




Yes, in a **Kubernetes cluster** you should use **ConfigMaps and Secrets** to avoid putting configuration or credentials in your code. Both are standard production practices.

---

# 1️⃣ Difference Between ConfigMap and Secret

| Type                     | Use Case                    | Example                           |
| ------------------------ | --------------------------- | --------------------------------- |
| **Kubernetes ConfigMap** | Non-sensitive configuration | App settings, URLs, feature flags |
| **Kubernetes Secret**    | Sensitive data              | Passwords, tokens, DB credentials |

Example for your app:

* Mongo host → **ConfigMap**
* Mongo username/password → **Secret**

---

# 2️⃣ Create Kubernetes Secret (for MongoDB credentials)

```bash
kubectl create secret generic mongo-secret \
--from-literal=MONGO_USER=admin \
--from-literal=MONGO_PASSWORD=admin123
```

Verify:

```bash
kubectl get secrets
```

---

# 3️⃣ Create ConfigMap (for non-sensitive config)

Example `configmap.yaml`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: backend-config
data:
  MONGO_HOST: mongodb
  MONGO_DB: mydatabase
```

Apply it:

```bash
kubectl apply -f configmap.yaml
```

---

# 4️⃣ Use ConfigMap + Secret in Deployment

Example `deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-backend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: go-backend
  template:
    metadata:
      labels:
        app: go-backend
    spec:
      containers:
      - name: backend
        image: my-backend:latest
        ports:
        - containerPort: 8080

        env:
        - name: MONGO_HOST
          valueFrom:
            configMapKeyRef:
              name: backend-config
              key: MONGO_HOST

        - name: MONGO_DB
          valueFrom:
            configMapKeyRef:
              name: backend-config
              key: MONGO_DB

        - name: MONGO_USER
          valueFrom:
            secretKeyRef:
              name: mongo-secret
              key: MONGO_USER

        - name: MONGO_PASSWORD
          valueFrom:
            secretKeyRef:
              name: mongo-secret
              key: MONGO_PASSWORD
```

---

# 5️⃣ Update Go Code to Read Variables

Instead of hardcoding:

```go
mongodb://admin:admin123@mongodb:27017/mydatabase
```

Use environment variables:

```go
mongoUser := os.Getenv("MONGO_USER")
mongoPass := os.Getenv("MONGO_PASSWORD")
mongoHost := os.Getenv("MONGO_HOST")
mongoDB := os.Getenv("MONGO_DB")

uri := fmt.Sprintf(
 "mongodb://%s:%s@%s:27017/%s?authSource=admin",
 mongoUser, mongoPass, mongoHost, mongoDB,
)
```

---

# 6️⃣ Why This Is the Production Standard

Benefits:

* Secrets not stored in source code
* Credentials can be rotated without rebuilding images
* Secure configuration management
* Works well with CI/CD pipelines

Large companies combine this with tools like:

* HashiCorp Vault
* AWS Secrets Manager

---

✅ **Best practice architecture for your project**

```
React (Frontend)
      ↓
Go Backend (Deployment)
      ↓
MongoDB (StatefulSet)
```

with:

* ConfigMap → configuration
* Secret → credentials
* Service → networking

---

💡 If you're building a **DevOps portfolio project (React + Go + Mongo + Kubernetes)**, I can also show you a **complete production-ready Kubernetes setup (Deployments, Services, Ingress, ConfigMaps, Secrets)** that looks like a real company architecture.






For a **React app**, environment variables are handled differently than backend apps because React is **compiled at build time**. The common approach is:

* Use **environment variables during build**
* Use **ConfigMap in Kubernetes**
* Use **ENV in Docker**

I'll show the standard workflow used in production.

---

# 1️⃣ Use Environment Variable in React

In React (especially with **Vite** or **Create React App**) you must prefix variables.

Since your frontend earlier used port `5174`, you are likely using **Vite**.

For **Vite**, use:

```javascript
const API_URL = import.meta.env.VITE_API_URL
```

Example:

```javascript
const API_URL = import.meta.env.VITE_API_URL + "/tasks"
```

---

# 2️⃣ Local Development (.env)

Create a `.env` file in your React project:

```env
VITE_API_URL=http://localhost:8080
```

Then use:

```javascript
const API_URL = `${import.meta.env.VITE_API_URL}/tasks`
```

Restart the dev server:

```bash
npm run dev
```

---

# 3️⃣ Docker Setup

React environment variables must be passed **during build**.

Example **Dockerfile**:

```dockerfile
FROM node:20 AS build

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

ARG Backend_API_URL
ENV Backend_API_URL=$Backend_API_URL

RUN npm run build

FROM nginx:alpine

COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80
```

Build image:

```bash
docker build \
--build-arg Backend_API_URL=http://my-backend:8080 \
-t my-frontend .
```

---

# 4️⃣ Kubernetes ConfigMap (Production)

Create ConfigMap.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: frontend-config
data:
  VITE_API_URL: "http://backend-service:8080"
```

Apply:

```bash
kubectl apply -f configmap.yaml
```

---

# 5️⃣ Use ConfigMap in Deployment

```yaml
env:
- name: VITE_API_URL
  valueFrom:
    configMapKeyRef:
      name: frontend-config
      key: VITE_API_URL
```

But remember: **React needs the value during build**, so usually we:

* Inject config during container startup
* Or rebuild image with env values.

---

# 6️⃣ Real Production Pattern (Best)

Large companies usually do this:

```text
React → Nginx → Backend Service → MongoDB
```

and frontend calls:

```javascript
const API_URL = "/api/tasks"
```

Then **Nginx** reverse proxy routes `/api` to backend.

Example:

```nginx
location /api {
    proxy_pass http://backend-service:8080;
}
```

This avoids hardcoding URLs entirely.

---

# ⭐ Recommended Structure for Your Project

```
frontend/
 ├── .env
 ├── Dockerfile
 ├── src/
 │   └── api.js

backend/
 ├── main.go
 ├── Dockerfile

k8s/
 ├── frontend-deployment.yaml
 ├── backend-deployment.yaml
 ├── configmap.yaml
```

---

# 🚀 DevOps Tip

When deploying React apps in **Kubernetes**, the most scalable pattern is:

* React served via **Nginx**
* Backend exposed through **Service**
* Frontend calls `/api` instead of hardcoded URL.

---

✅ If you want, I can also show you a **full working example of your 3-tier project (React + Go + Mongo) running in Docker Compose and Kubernetes** with environment variables and networking done correctly.


















Great 👍 Since you're building a **3-tier DevOps project (React + Go + MongoDB)**, I'll show you a **production-style architecture** used in real companies with **Docker**, **Docker Compose**, and later deployable to **Kubernetes**.

This will make your project **very strong for DevOps / MLOps portfolios**.

---

# 🏗️ Final Architecture

```text
Browser
   │
   ▼
React Frontend (Nginx)
   │
   ▼
Go Backend API (Gin)
   │
   ▼
MongoDB Database
```

Inside Docker network:

```text
frontend  →  backend  →  mongodb
```

---

# 📁 Project Structure (Recommended)

```
three-tier-task-tracker
│
├── frontend
│   ├── Dockerfile
│   ├── nginx.conf
│   └── React App
│
├── backend
│   ├── Dockerfile
│   └── main.go
│
├── database
│
├── docker-compose.yml
│
└── k8s
    ├── frontend-deployment.yaml
    ├── backend-deployment.yaml
    ├── mongodb-deployment.yaml
```

---

# 🐳 1. Backend Dockerfile (Go)

```dockerfile
# Build stage
FROM golang:1.22 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o main .

# Run stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
```

---

# 🐳 2. Frontend Dockerfile (React)

```dockerfile
# Build React
FROM node:20 AS builder

WORKDIR /app
COPY . .
RUN npm install
RUN npm run build

# Serve with nginx
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
```

---

# 🌐 3. Nginx Reverse Proxy

`frontend/nginx.conf`

```nginx
events {}

http {

  server {
    listen 80;

    location / {
      root /usr/share/nginx/html;
      index index.html;
      try_files $uri /index.html;
    }

    location /api/ {
      proxy_pass http://backend:8080/;
    }

  }
}
```

Now React calls:

```javascript
fetch("/api/tasks")
```

Nginx forwards request to backend.

---

# 🐳 4. Docker Compose (Best for DevOps Projects)

`docker-compose.yml`

```yaml
version: "3.9"

services:

  mongodb:
    image: mongo
    container_name: mongodb
    restart: always
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: admin123
    ports:
      - "27017:27017"

  backend:
    build: ./backend
    container_name: backend
    depends_on:
      - mongodb
    environment:
      MONGO_URI: mongodb://admin:admin123@mongodb:27017/mydatabase?authSource=admin
    ports:
      - "8080:8080"

  frontend:
    build: ./frontend
    container_name: frontend
    depends_on:
      - backend
    ports:
      - "3000:80"
```

Run everything:

```bash
docker compose up -d
```

Open:

```
http://localhost:3000
```

---

# 🔐 Backend Code (Use ENV Variable)

Instead of hardcoding Mongo URL:

```go
mongoURI := os.Getenv("MONGO_URI")
client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
```

---

# ☸️ Kubernetes Version

Later convert to **Kubernetes**.

Services communicate like this:

```text
frontend → backend-service → mongodb-service
```

Example backend service URL:

```
mongodb://admin:admin123@mongodb-service:27017
```

---

# 🔐 Kubernetes Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: mongodb-secret
type: Opaque
stringData:
  MONGO_URI: mongodb://admin:admin123@mongodb-service:27017/mydatabase?authSource=admin
```

---

# 🚀 DevOps Tools You Can Add (Highly Recommended)

To make this a **real DevOps portfolio project**:

| Tool           | Purpose                     |
| -------------- | --------------------------- |
| Docker         | Containerization            |
| Docker Compose | Local multi-container setup |
| Kubernetes     | Container orchestration     |
| Helm           | Kubernetes package manager  |
| GitHub Actions | CI/CD pipeline              |
| Prometheus     | Monitoring                  |
| Grafana        | Visualization               |

---

# 💼 What This Project Shows in Resume

You can write:

```
Built a 3-tier cloud-native task tracking system using
React, Go (Gin), and MongoDB.

• Containerized applications using Docker
• Implemented multi-container orchestration using Docker Compose
• Created Kubernetes manifests for scalable deployment
• Implemented Nginx reverse proxy for API routing
• Managed secrets using Kubernetes Secrets
• Designed microservice-based architecture
```

This is **very attractive for DevOps / MLOps jobs**.

---

# ⭐ If you want, I can also show you the **NEXT LEVEL version of this project**:

I can upgrade it to include:

✅ CI/CD pipeline
✅ Kubernetes autoscaling
✅ Monitoring dashboard
✅ Production-grade architecture

This becomes a **very powerful DevOps portfolio project**.








Good — let’s make it **simple and practical for your 2-node Kubernetes cluster** (1 master + 1 worker in **Kubernetes**).

You have 3 services:

1️⃣ React frontend
2️⃣ Go backend
3️⃣ MongoDB database

All will communicate using **Kubernetes Services (DNS names)**.

---

# 1️⃣ Important Changes Before Building Images

## Backend → MongoDB

Inside Kubernetes, containers **do NOT use container names**.

They use **Service names**.

So change your backend connection string from:

```text
mongodb://admin:admin123@my-mongo:27017/mydatabase?authSource=admin
```

to

```text
mongodb://admin:admin123@mongodb-service:27017/mydatabase?authSource=admin
```

Because Kubernetes DNS will resolve:

```text
mongodb-service → MongoDB Pod
```

---

## React → Backend

Change your `.env` before building React image.

❌ Current

```env
VITE_BACKEND_API_URL=http://localhost:8080
```

✅ Change to

```env
VITE_BACKEND_API_URL=http://backend-service:8080
```

Because frontend will call:

```text
backend-service → backend pod
```

Then **build new React Docker image**.

---

# 2️⃣ MongoDB Deployment

`mongodb-deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mongodb
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mongodb
  template:
    metadata:
      labels:
        app: mongodb
    spec:
      containers:
      - name: mongodb
        image: mongo:latest
        ports:
        - containerPort: 27017
        env:
        - name: MONGO_INITDB_ROOT_USERNAME
          value: "admin"
        - name: MONGO_INITDB_ROOT_PASSWORD
          value: "admin123"
```

---

# MongoDB Service

`mongodb-service.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: mongodb-service
spec:
  selector:
    app: mongodb
  ports:
  - port: 27017
    targetPort: 27017
  type: ClusterIP
```

---

# 3️⃣ Backend Deployment

`backend-deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend
        image: yourdockerhub/backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: MONGO_URI
          value: mongodb://admin:admin123@mongodb-service:27017/mydatabase?authSource=admin
```

---

# Backend Service

`backend-service.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: backend-service
spec:
  selector:
    app: backend
  ports:
  - port: 8080
    targetPort: 8080
  type: ClusterIP
```

---

# 4️⃣ Frontend Deployment

`frontend-deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: frontend
  template:
    metadata:
      labels:
        app: frontend
    spec:
      containers:
      - name: frontend
        image: yourdockerhub/frontend:latest
        ports:
        - containerPort: 80
```

---

# Frontend Service

To access from browser we use **NodePort**.

`frontend-service.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: frontend-service
spec:
  selector:
    app: frontend
  ports:
  - port: 80
    targetPort: 80
    nodePort: 30007
  type: NodePort
```

---

# 5️⃣ Apply Everything

Run in your master node:

```bash
kubectl apply -f mongodb-deployment.yaml
kubectl apply -f mongodb-service.yaml

kubectl apply -f backend-deployment.yaml
kubectl apply -f backend-service.yaml

kubectl apply -f frontend-deployment.yaml
kubectl apply -f frontend-service.yaml
```

---

# 6️⃣ Check Pods

```bash
kubectl get pods
```

Example:

```text
mongodb-xxxx
backend-xxxx
frontend-xxxx
```

---

# 7️⃣ Access Application

Find worker node IP:

```bash
kubectl get nodes -o wide
```

Example:

```
192.168.64.10
```

Open browser:

```text
http://192.168.64.10:30007
```

Your **React app will load**.

---

# 8️⃣ Final Communication Flow

Inside **Kubernetes** cluster:

```
React frontend
      ↓
backend-service:8080
      ↓
mongodb-service:27017
```

---

# 9️⃣ Summary of What You Must Change

### React

```
VITE_BACKEND_API_URL=http://backend-service:8080
```

Then rebuild Docker image.

---

### Backend

Use env variable:

```go
mongoURI := os.Getenv("MONGO_URI")
```

Kubernetes will pass:

```
mongodb://admin:admin123@mongodb-service:27017/mydatabase?authSource=admin
```

---

✅ This setup will **work perfectly in your 2-node Kubernetes cluster**.

---

If you want, I can also show you **3 very common mistakes people make in Kubernetes 3-tier apps** (which cause backend ↔ database connection failures). This will save you hours of debugging.


