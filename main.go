package main

import (
	"InternBorobitApp/Handlers"
	"InternBorobitApp/Repos"
	"InternBorobitApp/Services"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Get database connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	log.Printf("DB_USER: %s", dbUser)
	log.Printf("DB_PASSWORD: %s", dbPassword)
	log.Printf("DB_HOST: %s", dbHost)
	log.Printf("DB_NAME: %s", dbName)

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbName == "" {
		log.Fatalf("Error: Required environment variables are not set")
	}

	// Load the CA certificate
	caCert, err := ioutil.ReadFile("/app/global-bundle.pem")
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}

	// Create a CA certificate pool and add the CA certificate to it
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Configure TLS options
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	// Create a MongoDB client with the TLS options
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s/%s?tls=true&replicaSet=rs0&readPreference=secondaryPreferred&retryWrites=false", dbUser, dbPassword, dbHost, dbName)
	clientOptions := options.Client().ApplyURI(mongoURI).SetTLSConfig(tlsConfig)

	// MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatalf("Failed to disconnect MongoDB: %v", err)
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
