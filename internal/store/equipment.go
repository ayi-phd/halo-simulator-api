package store

import (
	"sync"

	"github.com/google/uuid"
	"halo-simulator/internal/domain"
)

type EquipmentStore struct {
	mu        sync.RWMutex
	equipment map[string]domain.GroundEquipment
}

func NewEquipmentStore() *EquipmentStore {
	return &EquipmentStore{
		equipment: make(map[string]domain.GroundEquipment),
	}
}

func (s *EquipmentStore) Save(e domain.GroundEquipment) domain.GroundEquipment {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	s.equipment[e.ID] = e
	return e
}

func (s *EquipmentStore) Get(id string) (domain.GroundEquipment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.equipment[id]
	return e, ok
}

func (s *EquipmentStore) All() []domain.GroundEquipment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.GroundEquipment, 0, len(s.equipment))
	for _, e := range s.equipment {
		result = append(result, e)
	}
	return result
}

// Available returns all equipment with status Available.
func (s *EquipmentStore) Available() []domain.GroundEquipment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.GroundEquipment
	for _, e := range s.equipment {
		if e.Status == domain.EquipmentStatusAvailable {
			result = append(result, e)
		}
	}
	return result
}
