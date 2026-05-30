package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudysito/apirestgo/models"
)

func RegistrarJugadores(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido, se espera un método POST.", http.StatusMethodNotAllowed)
		return
	}

	var nuevoJugador models.Jugador

	err := json.NewDecoder(r.Body).Decode(&nuevoJugador)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	models.BaseDatos = append(models.BaseDatos, nuevoJugador)

	w.Header().Set("Content-Type", "application/json")
	respuesta := map[string]string{
		"mensaje": fmt.Sprintf("Jugador %s registrado exitosamente", nuevoJugador.Nombre),
		"rango":   nuevoJugador.Rango,
	}

	json.NewEncoder(w).Encode(respuesta)
}

func ObtenerJugadores(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido, se espera un método GET.", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.BaseDatos)
}
