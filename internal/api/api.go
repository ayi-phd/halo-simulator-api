package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler holds dependencies that will be injected in later steps
// (stores, event bus, etc). Empty for now.
type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Routes registers all API endpoints on a new chi router and returns it.
// Middleware (logging, recovery) is added by the caller in main.go.
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", h.health)

	r.Post("/flights", h.createFlight)
	r.Post("/crew", h.createCrew)
	r.Post("/equipment", h.createEquipment)
	r.Post("/events", h.publishEvent)

	r.Get("/assignments", h.listAssignments)
	r.Get("/dashboard", h.dashboard)
	r.Get("/metrics", h.metrics)

	return r
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "halo-simulator",
	})
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

func (h *Handler) metrics(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not yet implemented", http.StatusNotImplemented)
}
