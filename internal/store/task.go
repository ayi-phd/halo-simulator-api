package store

import (
	"sync"

	"github.com/google/uuid"
	"halo-simulator/internal/domain"
)

type TaskStore struct {
	mu    sync.RWMutex
	tasks map[string]domain.Task
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]domain.Task),
	}
}

func (s *TaskStore) Save(t domain.Task) domain.Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	s.tasks[t.ID] = t
	return t
}

func (s *TaskStore) Get(id string) (domain.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	return t, ok
}

func (s *TaskStore) All() []domain.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		result = append(result, t)
	}
	return result
}

// ByFlight returns all tasks for a given flight ID.
func (s *TaskStore) ByFlight(flightID string) []domain.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.Task
	for _, t := range s.tasks {
		if t.FlightID == flightID {
			result = append(result, t)
		}
	}
	return result
}
