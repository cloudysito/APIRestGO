package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudysito/apirestgo/handlers"
	"github.com/cloudysito/apirestgo/middleware"
	"github.com/cloudysito/apirestgo/repository"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading the .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		fmt.Printf("Error connecting to MongoDB: %v\n", err)
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Printf("Could not ping MongoDB: %v\n", err)
		return
	}

	fmt.Println("Successfully connected to MongoDB")

	repo := repository.NewMongoRepository(client)
	handler := &handlers.PlayerHandler{Repo: repo}

	http.Handle("/api/register", middleware.ValidateToken(http.HandlerFunc(handler.RegisterPlayer)))
	http.HandleFunc("/api/players", handler.GetPlayers)
	http.HandleFunc("/api/player/", handler.GetPlayer)

	port := os.Getenv("PORT")
	fmt.Println("Server started at http://localhost:" + port)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error starting the server: %v\n", err)
	}
}
