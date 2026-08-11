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

func (h *FactionHandler) GetFaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing faction name", http.StatusBadRequest)
		return
	}

	faction, err := h.repo.GetFactionByName(context.TODO(), name)
	if err != nil {
		http.Error(w, "Failed to get faction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(faction)
}

func (h *FactionHandler) UpdateFaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing faction name", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateFactionPartial(context.TODO(), name, updates); err != nil {
		http.Error(w, "Failed to update faction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Faction updated successfully!"})
}

func (h *FactionHandler) DeleteFaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing faction name", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteFaction(context.TODO(), name); err != nil {
		http.Error(w, "Failed to delete faction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Faction deleted successfully!"})
}
