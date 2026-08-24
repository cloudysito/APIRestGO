package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudysito/apirestgo/models"
	"github.com/cloudysito/apirestgo/repository"
	"github.com/go-chi/chi/v5"
)

type PlayerHandler struct {
	Repo repository.PlayerRepository
}

func (h *PlayerHandler) RegisterPlayer(w http.ResponseWriter, r *http.Request) {
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
		fmt.Printf("\n[GOROUTINE] Calculating initial elo for %s...\n", name)
		time.Sleep(5 * time.Second)
		fmt.Printf("[GOROUTINE] Calculation completed! Welcome %s\n", name)
	}(newPlayer.Name)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": fmt.Sprintf("Player %s registered successfully", newPlayer.Name),
		"rank":    newPlayer.Rank,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *PlayerHandler) GetPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := h.Repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Error getting players", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}

func (h *PlayerHandler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	// chi.URLParam extracts the {name} segment defined in the route
	name := chi.URLParam(r, "name")

	player, err := h.Repo.GetByName(r.Context(), name)
	if err != nil {
		http.Error(w, "Player not found in the database", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

func (h *PlayerHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.Repo.GetRankStats(r.Context())
	if err != nil {
		http.Error(w, "Error getting stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func (h *PlayerHandler) UpdateRank(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	newRank := data["rank"]
	if err := h.Repo.UpdateRank(r.Context(), name, newRank); err != nil {
		http.Error(w, "Error updating rank", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message": fmt.Sprintf("Player %s rank updated to %s", name, newRank),
	}
	json.NewEncoder(w).Encode(response)
}

func (h *PlayerHandler) DeletePlayer(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if err := h.Repo.DeletePlayer(r.Context(), name); err != nil {
		http.Error(w, "Error deleting player", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Player %s deleted successfully", name),
	})
}
