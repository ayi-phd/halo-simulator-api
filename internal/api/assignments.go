package api

import (
	"net/http"

	"halo-simulator/internal/domain"
)

func (h *Handler) listFlights(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.flights.AllFlights())
}

func (h *Handler) listEquipment(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.equipment.All())
}

func (h *Handler) listAssignments(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.assignments.All())
}

// dashboardFlight pairs a flight with its ground tasks.
type dashboardFlight struct {
	Flight domain.Flight `json:"flight"`
	Tasks  []domain.Task `json:"tasks"`
}

func (h *Handler) dashboard(w http.ResponseWriter, r *http.Request) {
	flights := h.flights.AllFlights()
	views := make([]dashboardFlight, 0, len(flights))
	for _, f := range flights {
		views = append(views, dashboardFlight{
			Flight: f,
			Tasks:  h.tasks.ByFlight(f.ID),
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"flights": views,
		"metrics": h.metrics.Snapshot(),
	})
}
