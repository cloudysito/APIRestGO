package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cloudysito/apirestgo/models"
	"github.com/cloudysito/apirestgo/repository"
)

type JugadorHandler struct {
	Repo repository.JugadorRepository
}

func (h *JugadorHandler) RegistrarJugadores(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido. Usa POST.", http.StatusMethodNotAllowed)
		return
	}

	var nuevoJugador models.Jugador
	err := json.NewDecoder(r.Body).Decode(&nuevoJugador)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	err = h.Repo.Guardar(nuevoJugador)
	if err != nil {
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	// GOROUTINE
	go func(nombre string) {
		fmt.Printf("\n[GOROUTINE] Calculando MMR inicial para %s...\n", nombre)
		time.Sleep(5 * time.Second)
		fmt.Printf("[GOROUTINE] ¡Cálculo terminado! Correo de bienvenida enviado a %s\n", nombre)
	}(nuevoJugador.Nombre)

	w.Header().Set("Content-Type", "application/json")
	respuesta := map[string]string{
		"mensaje": fmt.Sprintf("Jugador %s registrado exitosamente", nuevoJugador.Nombre),
		"rango":   nuevoJugador.Rango,
	}
	json.NewEncoder(w).Encode(respuesta)
}

func (h *JugadorHandler) ObtenerJugadores(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido. Usa GET.", http.StatusMethodNotAllowed)
		return
	}

	jugadores, err := h.Repo.ObtenerTodos()
	if err != nil {
		http.Error(w, "Error al obtener los jugadores", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jugadores)
}

func (h *JugadorHandler) ObtenerJugador(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido. Usa GET.", http.StatusMethodNotAllowed)
		return
	}

	nombre := strings.TrimPrefix(r.URL.Path, "/api/jugador/")

	if nombre == "" {
		http.Error(w, "Falta especificar el nombre del jugador", http.StatusBadRequest)
		return
	}

	jugador, err := h.Repo.ObtenerPorNombre(nombre)
	if err != nil {
		http.Error(w, "Jugador no encontrado en la base de datos", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jugador)
}
