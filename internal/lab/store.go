package lab

import (
	"sort"
	"sync"
)

type Store interface {
	Create(Observation) error
	Get(string) (Observation, error)
	Update(Observation) error
	AddLabels(string, []string) error
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
func (s *MemoryStore) AddLabels(id string, labels []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.observations[id]
	if !ok {
		return ErrNotFound
	}
	item = cloneObservation(item)
	item.ApplyLabels(labels)
	s.observations[id] = cloneObservation(item)
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
func cloneObservation(item Observation) Observation {
	item.Labels = append([]string(nil), item.Labels...)
	if item.Review != nil {
		review := *item.Review
		item.Review = &review
	}
	return item
}
