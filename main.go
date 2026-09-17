package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Get("/players", handler.GetPlayers)
		r.Get("/players/{name}", handler.GetPlayer)
		r.Get("/factions", factionHandler.GetAllFactions)
		r.Get("/factions/{name}", factionHandler.GetFaction)

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

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		fmt.Println("Server is running on port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port %s: %v\n", port, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("\nShutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exiting")
}
