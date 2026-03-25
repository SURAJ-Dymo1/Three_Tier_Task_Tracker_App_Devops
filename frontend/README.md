# Project Set Up
To create a React application using **Vite**, you can use the following command. Vite is significantly faster than the old `create-react-app` because it uses native ES modules.

### 🚀 The Quick Start Command
Run this in your terminal:

```bash
npm create vite@latest my-react-app -- --template react
```

---

### 🛠️ Step-by-Step Installation
After running the command above, follow these three steps to get your development server running:

1. **Enter the project folder:**
   ```bash
   cd my-react-app
   ```
2. **Install the dependencies:**
   ```bash
   npm install
   ```
3. **Start the development server:**
   ```bash
   npm run dev
   ```

---

### 💡 Why use Vite instead of Create-React-App?
In a production environment and in your **Multipass** cluster, Vite is preferred for several reasons:

* **Speed:** Vite starts the server almost instantly, regardless of project size.
* **Hot Module Replacement (HMR):** When you save a file, the browser updates immediately without a full page reload.
* **Optimized Build:** It uses **Rollup** for production builds, resulting in smaller, faster chunks for your Kubernetes pods to serve.



---

### 🌐 Running it in your Multipass VM
If you are running this inside your **k8s-worker** or **k8s-master** and want to access it from your laptop's browser, you need to tell Vite to listen on all network interfaces.

**Change your `package.json` dev script to:**
```json
"dev": "vite --host 0.0.0.0"
```
Now, you can access the app at `http://<VM-IP>:5173`.

**Would you like me to show you the Dockerfile to containerize this Vite app so you can deploy it to your Kubernetes cluster?**

To deploy your Vite-React app to your Multipass Kubernetes cluster, we need a **two-stage Dockerfile**. 

In production, we don't use the Vite development server (`npm run dev`) because it's slow and heavy. Instead, we **build** the app into static HTML/JS/CSS files and serve them using **Nginx**, which is high-performance and lightweight.

---

### 🐳 The Production Dockerfile
Create a file named `Dockerfile` in your React project root:

```dockerfile
# Stage 1: Build the React app
FROM node:20-alpine AS build
WORKDIR /app

# Install dependencies first (better caching)
COPY package*.json ./
RUN npm install

# Copy the rest of the code and build
COPY . .
RUN npm run build

# Stage 2: Serve with Nginx
FROM nginx:stable-alpine

# Copy the build output from Stage 1 to Nginx's public folder
COPY --from=build /app/dist /usr/share/nginx/html

# Expose port 80 (standard for HTTP)
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
```



---

### 🛠️ Step-by-Step Production Setup

#### 1. Add a `.dockerignore` file
This prevents your local `node_modules` and logs from being copied into the image, making the build much faster.
```text
node_modules
dist
.git
npm-debug.log
```

#### 2. Build the Image
If you are using a local registry or pushing to Docker Hub:
```bash
docker build -t your-username/react-task-frontend:v1 .
```

---

### 🏗️ Why this is "Actual Production" standard:

1.  **Small Image Size:** By using the `nginx:alpine` image in the second stage, your final image will be around **20MB**, compared to **800MB+** if you kept the Node.js environment.
2.  **Security:** The final image doesn't contain your source code or `npm` tools—only the compiled, minified production files.
3.  **Performance:** Nginx is specifically designed to serve static files (like your Vite build) much faster than a Node.js server could.



---

### 🚀 Integrating with your Go Backend
When you run this in Kubernetes, your React app will try to talk to the Go API. Since Nginx is serving the frontend on port 80, you have two choices for networking:

* **Simple:** Hardcode the Go Service IP (not recommended).
* **Pro (Ingress):** Use a Kubernetes **Ingress Controller** so that both the Frontend and Backend are on the same domain (e.g., `my-app.local`).

---

### ✅ Next Step: The Kubernetes Manifest
Now that you have Dockerfiles for both **Go (Backend)** and **React (Frontend)**, we need to tell Kubernetes how to run them together.

**Would you like me to provide the `deployment.yaml` that sets up both the Frontend and Backend in your Multipass cluster?**