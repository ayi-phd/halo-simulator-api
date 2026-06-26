package store

import (
	"sync"

	"github.com/google/uuid"
	"halo-simulator/internal/domain"
)

type CrewStore struct {
	mu    sync.RWMutex
	crews map[string]domain.Crew
}

func NewCrewStore() *CrewStore {
	return &CrewStore{
		crews: make(map[string]domain.Crew),
	}
}

func (s *CrewStore) Save(c domain.Crew) domain.Crew {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	s.crews[c.ID] = c
	return c
}

func (s *CrewStore) Get(id string) (domain.Crew, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.crews[id]
	return c, ok
}

func (s *CrewStore) All() []domain.Crew {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.Crew, 0, len(s.crews))
	for _, c := range s.crews {
		result = append(result, c)
	}
	return result
}

// Available returns all crew members with status Available.
func (s *CrewStore) Available() []domain.Crew {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.Crew
	for _, c := range s.crews {
		if c.Status == domain.CrewStatusAvailable {
			result = append(result, c)
		}
	}
	return result
}
