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
