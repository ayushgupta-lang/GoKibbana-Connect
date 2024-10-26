package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ayush/mongo-kibana/config"
	"github.com/ayush/mongo-kibana/routes"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var esClient *elasticsearch.Client

// LogMessage defines the structure of the log message
type LogMessage struct {
	Timestamp time.Time `json:"@timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Service   string    `json:"service"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Status    int       `json:"status"`
	Latency   float64   `json:"latency"`
}

// Initialize Elasticsearch client
func initElasticsearch() *elasticsearch.Client {
	return config.NewElasticsearchClient()
}

func ConnectDB() *mongo.Client {
	// MongoDB connection logic...
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Ping to ensure connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	return client
}

func main() {
	client = ConnectDB()
	esClient = initElasticsearch()

	r := gin.Default()

	// Apply the logging middleware
	r.Use(RequestLogger)

	// Register routes
	routes.RegisterUserRoutes(r, client)

	r.Run(":8080")
}

// RequestLogger middleware to log requests
func RequestLogger(c *gin.Context) {
	startTime := time.Now()

	// Process request
	c.Next()

	// Calculate the time taken to process the request
	latency := time.Since(startTime)
	statusCode := c.Writer.Status()

	// Create log message
	logMsg := LogMessage{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "Request processed",
		Service:   "mongo-crud-service",
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Status:    statusCode,
		Latency:   float64(latency.Milliseconds()),
	}

	// Send log to Elasticsearch
	go sendLogToElasticsearch(logMsg)
}

// Send logs to Elasticsearch
func sendLogToElasticsearch(logMsg LogMessage) {
	jsonLog, err := json.Marshal(logMsg)
	if err != nil {
		log.Printf("Error marshaling log message: %s", err)
		return
	}

	res, err := esClient.Index(
		"mongo-crud-logs", // Index name
		bytes.NewReader(jsonLog),
		esClient.Index.WithContext(context.Background()),
	)
	if err != nil {
		log.Printf("Error indexing log to Elasticsearch: %s", err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error response from Elasticsearch: %s", res.String())
	}
}
