package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"halo-simulator/internal/events"
	"halo-simulator/internal/metrics"
	"halo-simulator/internal/store"
	"halo-simulator/internal/ws"
)

// Config groups all Handler dependencies. Using a struct keeps NewHandler
// readable when the number of dependencies grows.
type Config struct {
	Hub         *ws.Hub
	Metrics     *metrics.Collector
	Flights     *store.FlightStore
	Tasks       *store.TaskStore
	Crews       *store.CrewStore
	Equipment   *store.EquipmentStore
	Assignments *store.AssignmentStore
	Bus         *events.Bus
}

type Handler struct {
	hub         *ws.Hub
	metrics     *metrics.Collector
	flights     *store.FlightStore
	tasks       *store.TaskStore
	crews       *store.CrewStore
	equipment   *store.EquipmentStore
	assignments *store.AssignmentStore
	bus         *events.Bus
}

func NewHandler(cfg Config) *Handler {
	return &Handler{
		hub:         cfg.Hub,
		metrics:     cfg.Metrics,
		flights:     cfg.Flights,
		tasks:       cfg.Tasks,
		crews:       cfg.Crews,
		equipment:   cfg.Equipment,
		assignments: cfg.Assignments,
		bus:         cfg.Bus,
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", h.serveDashboard)
	r.Get("/health", h.health)
	r.Get("/ws", h.hub.ServeWS)

	r.Post("/flights", h.createFlight)
	r.Post("/crew", h.createCrew)
	r.Post("/equipment", h.createEquipment)
	r.Post("/events", h.publishEvent)

	r.Get("/flights", h.listFlights)
	r.Get("/equipment", h.listEquipment)
	r.Get("/assignments", h.listAssignments)
	r.Get("/dashboard", h.dashboard)
	r.Get("/metrics", h.getMetrics)

	return r
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "halo-simulator",
	})
}

func (h *Handler) getMetrics(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.metrics.Snapshot())
}

// writeJSON sets Content-Type and encodes v as JSON with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// decodeJSON decodes the request body into v, returning false and writing
// a 400 if decoding fails.
func decodeJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}
