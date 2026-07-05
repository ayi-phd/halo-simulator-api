package store

import (
	"sync"

	"github.com/google/uuid"
	"halo-simulator/internal/domain"
)

type AssignmentStore struct {
	mu          sync.RWMutex
	assignments map[string]domain.Assignment
}

func NewAssignmentStore() *AssignmentStore {
	return &AssignmentStore{
		assignments: make(map[string]domain.Assignment),
	}
}

func (s *AssignmentStore) Save(a domain.Assignment) domain.Assignment {
	s.mu.Lock()
	defer s.mu.Unlock()
	if a.ID == "" {
		a.ID = uuid.NewString()
	}
	s.assignments[a.ID] = a
	return a
}

func (s *AssignmentStore) Get(id string) (domain.Assignment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.assignments[id]
	return a, ok
}

func (s *AssignmentStore) All() []domain.Assignment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.Assignment, 0, len(s.assignments))
	for _, a := range s.assignments {
		result = append(result, a)
	}
	return result
}

// ByTask returns the active assignment for a given task ID, if any.
func (s *AssignmentStore) ByTask(taskID string) (domain.Assignment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.assignments {
		if a.TaskID == taskID && a.Status == domain.AssignmentStatusActive {
			return a, true
		}
	}
	return domain.Assignment{}, false
}
