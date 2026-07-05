package api

import (
	"net/http"

	"halo-simulator/internal/domain"
)

type createEquipmentRequest struct {
	Type     domain.EquipmentType `json:"type"`
	Location string               `json:"location"`
}

func (h *Handler) createEquipment(w http.ResponseWriter, r *http.Request) {
	var req createEquipmentRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if req.Type == "" {
		http.Error(w, "type is required", http.StatusBadRequest)
		return
	}
	if req.Location == "" {
		req.Location = "Depot"
	}

	equip := h.equipment.Save(domain.GroundEquipment{
		Type:     req.Type,
		Status:   domain.EquipmentStatusAvailable,
		Location: req.Location,
	})

	writeJSON(w, http.StatusCreated, equip)
}
