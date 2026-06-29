package api

import (
	"net/http"

	"halo-simulator/internal/domain"
)

type createCrewRequest struct {
	Name     string          `json:"name"`
	Role     domain.CrewRole `json:"role"`
	Location string          `json:"location"`
}

func (h *Handler) createCrew(w http.ResponseWriter, r *http.Request) {
	var req createCrewRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if req.Name == "" || req.Role == "" {
		http.Error(w, "name and role are required", http.StatusBadRequest)
		return
	}
	if req.Location == "" {
		req.Location = "Base"
	}

	crew := h.crews.Save(domain.Crew{
		Name:     req.Name,
		Role:     req.Role,
		Status:   domain.CrewStatusAvailable,
		Location: req.Location,
	})

	writeJSON(w, http.StatusCreated, crew)
}
