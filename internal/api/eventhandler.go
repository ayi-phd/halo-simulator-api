package api

import (
	"net/http"
	"time"

	"halo-simulator/internal/domain"
	"halo-simulator/internal/events"
)

// publishEventRequest is a union-style payload: which fields are used
// depends on the event type.
type publishEventRequest struct {
	Type        events.EventType `json:"type"`
	FlightID    string           `json:"flight_id,omitempty"`
	EquipmentID string           `json:"equipment_id,omitempty"`
	CrewID      string           `json:"crew_id,omitempty"`
}

func (h *Handler) publishEvent(w http.ResponseWriter, r *http.Request) {
	var req publishEventRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	switch req.Type {
	case events.FlightArrived:
		h.handleFlightArrived(w, req.FlightID)
	case events.FlightDelayed:
		h.handleFlightDelayed(w, req.FlightID)
	case events.EquipmentBroken:
		h.handleEquipmentBroken(w, req.EquipmentID)
	case events.EquipmentAvailable:
		h.handleEquipmentRepaired(w, req.EquipmentID)
	case events.CrewUnavailable:
		h.handleCrewUnavailable(w, req.CrewID)
	default:
		http.Error(w, "unsupported event type", http.StatusBadRequest)
	}
}

func (h *Handler) handleFlightArrived(w http.ResponseWriter, flightID string) {
	flight, ok := h.flights.GetFlight(flightID)
	if !ok {
		http.Error(w, "flight not found", http.StatusNotFound)
		return
	}
	h.bus.Publish(events.Event{Type: events.FlightArrived, Payload: flight})
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) handleFlightDelayed(w http.ResponseWriter, flightID string) {
	flight, ok := h.flights.GetFlight(flightID)
	if !ok {
		http.Error(w, "flight not found", http.StatusNotFound)
		return
	}
	flight.Status = domain.FlightStatusDelayed
	flight.ScheduledDeparture = flight.ScheduledDeparture.Add(30 * time.Minute)
	h.flights.SaveFlight(flight)
	h.bus.Publish(events.Event{Type: events.FlightDelayed, Payload: flight})
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) handleEquipmentBroken(w http.ResponseWriter, equipmentID string) {
	equip, ok := h.equipment.Get(equipmentID)
	if !ok {
		http.Error(w, "equipment not found", http.StatusNotFound)
		return
	}
	equip.Status = domain.EquipmentStatusBroken
	h.equipment.Save(equip)
	h.bus.Publish(events.Event{Type: events.EquipmentBroken, Payload: equip})
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) handleEquipmentRepaired(w http.ResponseWriter, equipmentID string) {
	equip, ok := h.equipment.Get(equipmentID)
	if !ok {
		http.Error(w, "equipment not found", http.StatusNotFound)
		return
	}
	equip.Status = domain.EquipmentStatusAvailable
	h.equipment.Save(equip)
	h.bus.Publish(events.Event{Type: events.EquipmentAvailable, Payload: equip})
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) handleCrewUnavailable(w http.ResponseWriter, crewID string) {
	crew, ok := h.crews.Get(crewID)
	if !ok {
		http.Error(w, "crew not found", http.StatusNotFound)
		return
	}
	crew.Status = domain.CrewStatusOffDuty
	h.crews.Save(crew)
	h.bus.Publish(events.Event{Type: events.CrewUnavailable, Payload: crew})
	w.WriteHeader(http.StatusAccepted)
}
