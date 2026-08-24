package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cloudysito/apirestgo/models"
	"github.com/cloudysito/apirestgo/repository"
	"github.com/go-chi/chi/v5"
)

type FactionHandler struct {
	repo *repository.FactionRepo
}

func NewFactionHandler(repo *repository.FactionRepo) *FactionHandler {
	return &FactionHandler{repo: repo}
}

func (h *FactionHandler) CreateFaction(w http.ResponseWriter, r *http.Request) {
	var faction models.Faction

	if err := json.NewDecoder(r.Body).Decode(&faction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateFaction(r.Context(), &faction); err != nil {
		http.Error(w, "Failed to create faction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Faction created successfully!"})
}

func (h *FactionHandler) GetAllFactions(w http.ResponseWriter, r *http.Request) {
	factions, err := h.repo.GetAllFactions(r.Context())
	if err != nil {
		http.Error(w, "Failed to get factions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(factions)
}

func (h *FactionHandler) GetFaction(w http.ResponseWriter, r *http.Request) {
	// chi.URLParam extracts the {name} segment defined in the route
	name := chi.URLParam(r, "name")

	faction, err := h.repo.GetFactionByName(r.Context(), name)
	if err != nil {
		http.Error(w, "Failed to get faction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(faction)
}

func (h *FactionHandler) UpdateFaction(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateFactionPartial(r.Context(), name, updates); err != nil {
		http.Error(w, "Failed to update faction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Faction updated successfully!"})
}

func (h *FactionHandler) DeleteFaction(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if err := h.repo.DeleteFaction(r.Context(), name); err != nil {
		http.Error(w, "Failed to delete faction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Faction deleted successfully!"})
}
