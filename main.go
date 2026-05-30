package main

import (
	"fmt"
	"net/http"

	"github.com/cloudysito/apirestgo/handlers"
)

func main() {
	http.HandleFunc("/api/registro", handlers.RegistrarJugadores)
	http.HandleFunc("/api/jugadores", handlers.ObtenerJugadores)

	fmt.Println("Servidor iniciado en http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
	}
}
