package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
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
	go func(name string, rank string) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mmr := calculateMMR(rank)

		if err := h.Repo.UpdateMMR(ctx, name, mmr); err != nil {
			log.Printf("[GOROUTINE] Error updating MMR for %s: %v\n", name, err)
			return
		}
		log.Printf("[GOROUTINE] MMR set for %s: %d\n", name, mmr)
	}(newPlayer.Name, newPlayer.Rank)

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

func calculateMMR(rank string) int {
	baseMMR := map[string]int{
		"Bronze":   800,
		"Silver":   1100,
		"Gold":     1400,
		"Platinum": 1700,
		"Diamond":  2000,
	}
	mmr, ok := baseMMR[rank]
	if !ok {
		return 1000
	}
	return mmr + rand.Intn(101) - 50
}

func (h *PlayerHandler) RecalculateMMR(w http.ResponseWriter, r *http.Request) {
	players, err := h.Repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Error getting players", http.StatusInternalServerError)
		return
	}

	const numWorkers = 5
	jobs := make(chan models.Player, len(players))
	results := make(chan string, len(players))

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for player := range jobs {
				log.Printf("[Worker %d] Processing %s", workerID, player.Name)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				mmr := calculateMMR(player.Rank)

				err := h.Repo.UpdateMMR(ctx, player.Name, mmr)
				cancel()

				if err != nil {
					results <- fmt.Sprintf("FAILED: %s (%v)", player.Name, err)
				} else {
					results <- fmt.Sprintf("SUCCESS: %s MMR updated to %d", player.Name, mmr)
				}
			}
		}(i)
	}

	for _, player := range players {
		jobs <- player
	}
	close(jobs)

	wg.Wait()
	close(results)

	var summary []string
	for result := range results {
		summary = append(summary, result)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf("MMR recalculation completed for %d players", len(players)),
		"summary": summary,
	})
}
