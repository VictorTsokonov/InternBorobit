package main

import (
	"InternBorobitApp/Handlers"
	"InternBorobitApp/Repos"
	"InternBorobitApp/Services"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Get database connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s/%s", dbUser, dbPassword, dbHost, dbName)

	// MongoDB connection
	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {

		}
	}(client, ctx)

	// Initialize repository and service
	gameRepo := Repos.NewMongoGameRepository(client, "InternBorobit", "Games")
	gameService := Services.NewGameService(gameRepo)
	gameHandler := Handlers.NewGameHandler(gameService)

	// Set up router and routes
	router := mux.NewRouter()
	router.HandleFunc("/games", gameHandler.CreateGame).Methods("POST")
	router.HandleFunc("/games/{id}", gameHandler.GetGameByID).Methods("GET")
	router.HandleFunc("/games/{id}", gameHandler.UpdateGame).Methods("PUT")
	router.HandleFunc("/games/{id}", gameHandler.DeleteGame).Methods("DELETE")
	router.HandleFunc("/games", gameHandler.ListGames).Methods("GET")

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
