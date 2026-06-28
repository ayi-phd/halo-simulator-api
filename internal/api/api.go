package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"halo-simulator/internal/metrics"
	"halo-simulator/internal/ws"
)

// Handler holds dependencies injected at startup.
type Handler struct {
	hub     *ws.Hub
	metrics *metrics.Collector
}

func NewHandler(hub *ws.Hub, metrics *metrics.Collector) *Handler {
	return &Handler{hub: hub, metrics: metrics}
}

// Routes registers all API endpoints on a new chi router and returns it.
// Middleware (logging, recovery) is added by the caller in main.go.
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", h.health)
	r.Get("/ws", h.hub.ServeWS)

	r.Post("/flights", h.createFlight)
	r.Post("/crew", h.createCrew)
	r.Post("/equipment", h.createEquipment)
	r.Post("/events", h.publishEvent)

	r.Get("/assignments", h.listAssignments)
	r.Get("/dashboard", h.dashboard)
	r.Get("/metrics", h.getMetrics)

	return r
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "halo-simulator",
	})
}

func (h *Handler) getMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.metrics.Snapshot())
}

func (h *Handler) createFlight(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not yet implemented", http.StatusNotImplemented)
}

func (h *Handler) createCrew(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not yet implemented", http.StatusNotImplemented)
}

func (h *Handler) createEquipment(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not yet implemented", http.StatusNotImplemented)
}

func (h *Handler) publishEvent(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not yet implemented", http.StatusNotImplemented)
}

func (h *Handler) listAssignments(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not yet implemented", http.StatusNotImplemented)
}

func (h *Handler) dashboard(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not yet implemented", http.StatusNotImplemented)
}
