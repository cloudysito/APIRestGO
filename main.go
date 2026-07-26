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
		log.Fatal("Error cargando el archivo .env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		fmt.Printf("Error al conectar a MongoDB: %v\n", err)
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Printf("No se pudo hacer ping a MongoDB: %v\n", err)
		return
	}

	fmt.Println("Conectado exitosamente a MongoDB")

	repo := repository.NewMongoRepository(client)
	handler := &handlers.JugadorHandler{Repo: repo}

	http.Handle("/api/registro", middleware.ValidarToken(http.HandlerFunc(handler.RegistrarJugadores)))
	http.HandleFunc("/api/jugadores", handler.ObtenerJugadores)
	http.HandleFunc("/api/jugador/", handler.ObtenerJugador)

	puerto := os.Getenv("PUERTO")
	fmt.Println("Servidor iniciado en http://localhost:" + puerto)

	err = http.ListenAndServe(":"+puerto, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
	}
}
