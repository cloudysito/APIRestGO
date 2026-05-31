package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudysito/apirestgo/handlers"
	"github.com/cloudysito/apirestgo/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
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

	http.HandleFunc("/api/registro", handler.RegistrarJugadores)
	http.HandleFunc("/api/jugadores", handler.ObtenerJugadores)

	fmt.Println("Servidor iniciado en http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
	}
}
