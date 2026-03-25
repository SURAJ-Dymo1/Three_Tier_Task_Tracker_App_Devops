package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
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
	mongoURI := os.Getenv("MONGO_URI")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic("MongoDB connection failed")
	}

	collection = client.Database("mydatabase").Collection("sample_collection")

	r := gin.Default()

	// CORS Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Routes
	r.POST("/tasks", createTask)
	r.GET("/tasks", getTasks)
	r.DELETE("/tasks/:title", deleteTask)

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var tasks []Task

	if err = cursor.All(ctx, &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
func deleteTask(c *gin.Context) {
	// 1. Get the title from the URL parameter (e.g., /tasks/Learn-K8s)
	title := c.Param("title")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 2. Define the filter
	filter := bson.M{"title": title}

	// 3. Execute Delete
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	// 4. Check if anything was actually deleted
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No task found with that title"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
