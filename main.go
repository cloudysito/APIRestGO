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

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	client, err := mongo.Connect(ctx, clientOptions)
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
	db := client.Database("apirestgo")

	repo := repository.NewMongoRepository(client)
	handler := &handlers.PlayerHandler{Repo: repo}

	factionRepo := repository.NewFactionRepo(db)
	factionHandler := handlers.NewFactionHandler(factionRepo)

	r := chi.NewRouter()

	// Global middleware: runs on every request
	r.Use(chimiddleware.Logger)    // logs method, path, status and duration
	r.Use(chimiddleware.Recoverer) // recovers from panics and returns 500

	r.Route("/api", func(r chi.Router) {

		// --- Public routes (no auth required) ---
		r.Get("/players", handler.GetPlayers)
		r.Get("/players/{name}", handler.GetPlayer)
		r.Get("/factions", factionHandler.GetAllFactions)
		r.Get("/factions/{name}", factionHandler.GetFaction)

		// --- Protected routes (ValidateToken middleware applied to this group) ---
		r.Group(func(r chi.Router) {
			r.Use(middleware.ValidateToken)

			r.Post("/players", handler.RegisterPlayer)
			r.Delete("/players/{name}", handler.DeletePlayer)
			r.Put("/players/{name}/rank", handler.UpdateRank)
			r.Get("/stats", handler.GetStats)

			r.Post("/factions", factionHandler.CreateFaction)
			r.Put("/factions/{name}", factionHandler.UpdateFaction)
			r.Delete("/factions/{name}", factionHandler.DeleteFaction)
		})
	})

	port := os.Getenv("PORT")
	fmt.Println("Server started at http://localhost:" + port)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Printf("Error starting the server: %v\n", err)
	}
}
