To build a Go (Golang) CRUD application with MongoDB, you’ll need the official **mongo-go-driver**. Go is excellent for this because it is extremely fast and compiles into a single binary, which makes your Kubernetes Docker images very small.

### 🛠️ Step 1: Initialize the Project
Run these commands in your terminal to set up the Go module and install the driver.

```bash
mkdir go-mongo-crud
cd go-mongo-crud
go mod init go-mongo-crud
go get go.mongodb.org/mongo-driver/mongo
```

---

### 📝 Step 2: The CRUD Implementation
Create a file named `main.go`. This code connects to MongoDB and performs **Create, Read, Update, and Delete** operations.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Task represents our data model
type Task struct {
	ID     string `bson:"_id,omitempty"`
	Title  string `bson:"title"`
	Status string `bson:"status"`
}

func main() {
	// 1. Connect to MongoDB (Use your K8s Service name or Multipass IP)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
    
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("taskdb").Collection("tasks")

	// --- CREATE ---
	newTask := Task{Title: "Learn Kubernetes Networking", Status: "Pending"}
	insertResult, _ := collection.InsertOne(ctx, newTask)
	fmt.Println("Inserted ID:", insertResult.InsertedID)

	// --- READ ---
	var result Task
	filter := bson.D{{"title", "Learn Kubernetes Networking"}}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println("Record not found")
	}
	fmt.Printf("Found Task: %+v\n", result)

	// --- UPDATE ---
	update := bson.D{{"$set", bson.D{{"status", "Completed"}}}}
	_, _ = collection.UpdateOne(ctx, filter, update)
	fmt.Println("Task updated to Completed")

	// --- DELETE ---
	// _, _ = collection.DeleteOne(ctx, filter)
	// fmt.Println("Task deleted")
}
```



---

### 🚀 Step 3: Production Concepts (VPC & K8s)

When you move this code into your **Multipass Cluster**, you need to change two things for it to be "Production Grade":

#### 1. The Connection String
In production, don't use `localhost`. Use the **Kubernetes Service DNS** name we talked about earlier:
```go
// Production URL using K8s DNS
uri := "mongodb://mongodb-service.task-tracker.svc.cluster.local:27017"
```

#### 2. Environment Variables
Remember the **ConfigMaps and Secrets** lesson? Instead of hardcoding the URI, pull it from the environment:
```go
dbUri := os.Getenv("MONGO_URL")
```



---

### 📦 How to Run
1. Make sure your MongoDB container is running (either in K8s or on your VM).
2. Run the Go app:
   ```bash
   go run main.go
   ```

### ✅ Next Step
Since you now have a **Vite (Frontend)** and a **Go (Backend)**, you have a complete stack! 

**Would you like me to show you how to write a Multi-stage Dockerfile for this Go app?** (This is crucial for production because it reduces your image size from 800MB to about 15MB).


To turn the Go code into a real API server, we use the standard `net/http` package. In production, this allows your **Ingress Controller** or **Service** to send traffic to your Pod on port **8080**.

### 🛠️ Step 1: Install a Router
While Go has a built-in router, **Gin** is the production standard because it is fast, handles JSON perfectly, and has great middleware for logging.

Run this in your project folder:
```bash
go get github.com/gin-gonic/gin
```

---

### 📝 Step 2: The API Server Code
Replace your `main.go` with this. This code creates a web server that listens on **8080** and handles JSON requests.

```go
package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Task struct {
	Title  string `json:"title" bson:"title"`
	Status string `json:"status" bson:"status"`
}

var collection *mongo.Collection

func main() {
	// 1. Setup MongoDB Connection
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	collection = client.Database("taskdb").Collection("tasks")

	// 2. Initialize Gin Router
	r := gin.Default()

	// 3. Define Routes (CRUD)
	r.POST("/tasks", createTask)
	r.GET("/tasks", getTasks)

	// 4. Start Server on Port 8080
	// In K8s, this allows the Service to find your app
	r.Run(":8080") 
}

func createTask(c *gin.Context) {
	var newTask Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, _ = collection.InsertOne(context.Background(), newTask)
	c.JSON(http.StatusCreated, newTask)
}

func getTasks(c *gin.Context) {
	cursor, _ := collection.Find(context.Background(), bson.M{})
	var tasks []Task
	cursor.All(context.Background(), &tasks)
	c.JSON(http.StatusOK, tasks)
}
```

---

### 🚀 Step 3: Test Your API
Once you run `go run main.go`, you can test it from your terminal (or using Postman):

**1. Create a Task (POST):**
```bash
curl -X POST http://localhost:8080/tasks \
   -H "Content-Type: application/json" \
   -d '{"title": "Setup Multipass Cluster", "status": "In Progress"}'
```

**2. Get all Tasks (GET):**
```bash
curl http://localhost:8080/tasks
```



---

### 🏗️ Production Bridge: Why Port 8080?
In Kubernetes production, port **8080** is the "sweet spot" because:
1.  **Non-Root:** It doesn't require "Root" privileges (ports below 1024 do), which makes your container more secure.
2.  **Service Mapping:** Your Kubernetes Service will map its own port 80 to your Pod's port 8080.



---

### 📦 The "Multi-Stage" Dockerfile (Crucial for Go)
Since you are running this on **Multipass**, you want your images to be tiny. Use this Dockerfile to turn your 800MB Go environment into a **15MB executable**.

```dockerfile
# Stage 1: Build the binary
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Stage 2: Run the binary
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

**Would you like me to show you the Kubernetes `Deployment` and `Service` YAML to get this Go API running in your task-tracker namespace?**

go mod tidy



Build Image

```
docker build -t my-backend .


```
Run Container

```
docker run -p 8080:8080 my-backend


```

docker run -d \
--name my-backend \
--network my-net \
-p 8081:8080 \
my-backend

```

docker run -d --name my-backend --network my-net -p 8080:8080 -e MONGO_URI="mongodb://admin:admin123@my-mongo:27017/mydatabase?authSource=admin" my-backend

docker network connect my-net my-backend
 docker inspect network my-net