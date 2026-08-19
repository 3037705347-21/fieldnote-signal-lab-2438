package lab

import (
	"sort"
	"sync"
)

type Store interface {
	Create(Observation) error
	Get(string) (Observation, error)
	Update(Observation) error
	List() []Observation
}
type MemoryStore struct {
	mu           sync.RWMutex
	observations map[string]Observation
}

func NewMemoryStore() *MemoryStore { return &MemoryStore{observations: map[string]Observation{}} }
func (s *MemoryStore) Create(item Observation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.observations[item.ID]; ok {
		return ErrConflict
	}
	s.observations[item.ID] = cloneObservation(item)
	return nil
}
func (s *MemoryStore) Get(id string) (Observation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.observations[id]
	if !ok {
		return Observation{}, ErrNotFound
	}
	return cloneObservation(item), nil
}
func (s *MemoryStore) Update(item Observation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.observations[item.ID]; !ok {
		return ErrNotFound
	}
	s.observations[item.ID] = cloneObservation(item)
	return nil
}
func (s *MemoryStore) List() []Observation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Observation, 0, len(s.observations))
	for _, item := range s.observations {
		result = append(result, cloneObservation(item))
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
// cloneObservation returns a deep copy of item so that the slices and pointers
// reachable from an Observation returned to a caller cannot mutate the stored
// record. Labels is copied to its own backing array, and Review is copied to a
// fresh value with a new pointer. A nil Labels slice or Review stays nil so the
// (nil vs. empty) distinction used in JSON output is preserved.
func cloneObservation(item Observation) Observation {
	clone := item
	if item.Labels != nil {
		clone.Labels = make([]string, len(item.Labels))
		copy(clone.Labels, item.Labels)
	}
	if item.Review != nil {
		review := *item.Review
		clone.Review = &review
	}
	return clone
}
