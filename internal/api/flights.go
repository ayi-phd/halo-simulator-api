package api

import (
	"net/http"
	"time"

	"halo-simulator/internal/domain"
	"halo-simulator/internal/events"
)

type createFlightRequest struct {
	FlightNumber       string    `json:"flight_number"`
	AirlineID          string    `json:"airline_id"`
	Gate               string    `json:"gate"`
	ScheduledDeparture time.Time `json:"scheduled_departure"`
}

func (h *Handler) createFlight(w http.ResponseWriter, r *http.Request) {
	var req createFlightRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if req.FlightNumber == "" || req.Gate == "" {
		http.Error(w, "flight_number and gate are required", http.StatusBadRequest)
		return
	}
	if req.ScheduledDeparture.IsZero() {
		req.ScheduledDeparture = time.Now().Add(2 * time.Hour)
	}

	flight := h.flights.SaveFlight(domain.Flight{
		FlightNumber:       req.FlightNumber,
		AirlineID:          req.AirlineID,
		Gate:               req.Gate,
		Status:             domain.FlightStatusArrived,
		ScheduledDeparture: req.ScheduledDeparture,
		ActualArrival:      time.Now(),
	})

	h.bus.Publish(events.Event{Type: events.FlightArrived, Payload: flight})
	writeJSON(w, http.StatusCreated, flight)
}
