package main

import (
	"context"
	"fmt"
	"log"

	"example/todo-go/controllers"
	"example/todo-go/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server         *gin.Engine
	todoService    services.TodoService
	todoController controllers.TodoController
	ctx            context.Context
	todoCollection *mongo.Collection
	client         *mongo.Client
	err            error
)

func init() {
	ctx = context.TODO()

	connection := options.Client().ApplyURI("mongodb://admin:password@localhost:27020")
	client, err = mongo.Connect(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB on 27020!")

	// controller -> service -> collection
	todoCollection = client.Database("development").Collection("todos")
	todoService = services.NewTodoService(todoCollection, ctx)
	todoController = controllers.New(todoService)
	server = gin.Default()
}

func closeMongoDB() {
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func main() {
	// close   MongoDB connection
	defer closeMongoDB()

	basePath := server.Group("/v1")
	todoController.RegisterTodoRoutes(basePath)

	log.Fatal(server.Run(":9090"))
}
