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

type PlayerHandler struct {
	Repo repository.PlayerRepository
}

func (h *PlayerHandler) RegisterPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	var newPlayer models.Player
	err := json.NewDecoder(r.Body).Decode(&newPlayer)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	err = h.Repo.Save(r.Context(), &newPlayer)
	if err != nil {
		http.Error(w, "Error saving to database", http.StatusInternalServerError)
		return
	}

	// GOROUTINE
	go func(name string) {
		fmt.Printf("\n[GOROUTINE] Calculating initial MMR for %s...\n", name)
		time.Sleep(5 * time.Second)
		fmt.Printf("[GOROUTINE] Calculation completed! Welcome %s\n", name)
	}(newPlayer.Name)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message": fmt.Sprintf("Player %s registered successfully", newPlayer.Name),
		"rank":    newPlayer.Rank,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *PlayerHandler) GetPlayers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed. Use GET.", http.StatusMethodNotAllowed)
		return
	}

	players, err := h.Repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Error getting players", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}

func (h *PlayerHandler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed. Use GET.", http.StatusMethodNotAllowed)
		return
	}

	name := strings.TrimPrefix(r.URL.Path, "/api/player/")

	if name == "" {
		http.Error(w, "Player name must be specified", http.StatusBadRequest)
		return
	}

	player, err := h.Repo.GetByName(r.Context(), name)
	if err != nil {
		http.Error(w, "Player not found in the database", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

func (h *PlayerHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed. Use GET.", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.Repo.GetRankStats(r.Context())
	if err != nil {
		http.Error(w, "Error getting stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}
