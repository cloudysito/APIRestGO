package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cloudysito/apirestgo/models"
	"github.com/cloudysito/apirestgo/repository"
)

type FactionHandler struct {
	repo *repository.FactionRepo
}

func NewFactionHandler(repo *repository.FactionRepo) *FactionHandler {
	return &FactionHandler{repo: repo}
}

func (h *FactionHandler) CreateFaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var faction models.Faction

	if err := json.NewDecoder(r.Body).Decode(&faction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateFaction(context.TODO(), &faction); err != nil {
		http.Error(w, "Failed to create faction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Faction created successfully!"})
}

func (h *FactionHandler) GetAllFactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	factions, err := h.repo.GetAllFactions(context.TODO())
	if err != nil {
		http.Error(w, "Failed to get factions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(factions)
}
