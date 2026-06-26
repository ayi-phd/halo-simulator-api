package store

import (
	"sync"

	"github.com/google/uuid"
	"halo-simulator/internal/domain"
)

type FlightStore struct {
	mu      sync.RWMutex
	flights map[string]domain.Flight
	airlines map[string]domain.Airline
}

func NewFlightStore() *FlightStore {
	return &FlightStore{
		flights:  make(map[string]domain.Flight),
		airlines: make(map[string]domain.Airline),
	}
}

func (s *FlightStore) SaveFlight(f domain.Flight) domain.Flight {
	s.mu.Lock()
	defer s.mu.Unlock()
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	s.flights[f.ID] = f
	return f
}

func (s *FlightStore) GetFlight(id string) (domain.Flight, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.flights[id]
	return f, ok
}

func (s *FlightStore) AllFlights() []domain.Flight {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.Flight, 0, len(s.flights))
	for _, f := range s.flights {
		result = append(result, f)
	}
	return result
}

func (s *FlightStore) SaveAirline(a domain.Airline) domain.Airline {
	s.mu.Lock()
	defer s.mu.Unlock()
	if a.ID == "" {
		a.ID = uuid.NewString()
	}
	s.airlines[a.ID] = a
	return a
}

func (s *FlightStore) GetAirline(id string) (domain.Airline, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.airlines[id]
	return a, ok
}
